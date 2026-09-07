package app

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/kehao95/quine/internal/qcli"
)

const defaultSpawnMission = "idle and wait"

func Main() int {
	return Run(os.Args[1:], os.Stdout, os.Stderr)
}

func Run(args []string, stdout, stderr io.Writer) int {
	if len(args) == 0 {
		return runInteractive(nil, stdout, stderr)
	}
	switch args[0] {
	case "--help", "-h", "help":
		return runHelp(stdout)
	case "attach":
		return runAttach(args[1:], stdout, stderr)
	case "spawn":
		return runSpawn(args[1:], stdout, stderr)
	case "serve":
		return runServe(args[1:], stdout, stderr)
	case "web":
		return runWeb(args[1:], stdout, stderr)
	default:
		fmt.Fprintf(stderr, "qcli: unknown command %q\n", args[0])
		return runHelp(stderr)
	}
}

func runHelp(out io.Writer) int {
	fmt.Fprintln(out, "usage: qcli [attach|spawn|serve] [options]")
	fmt.Fprintln(out, "")
	fmt.Fprintln(out, "commands:")
	fmt.Fprintln(out, "  qcli                 attach the only live peer, or spawn a waiting peer, then open the TUI")
	fmt.Fprintln(out, "  attach               attach a selected live peer, serve qcli/1, then open the TUI")
	fmt.Fprintln(out, "  spawn [mission]      spawn quine, attach it, serve qcli/1, then open the TUI")
	fmt.Fprintln(out, "  serve                headless HTTP+SSE qcli/1 kernel")
	fmt.Fprintln(out, "  web                  attach/spawn, serve kernel + web client, open browser")
	fmt.Fprintln(out, "")
	fmt.Fprintln(out, "selectors:")
	fmt.Fprintln(out, "  --agent-root PATH    resolve directly from agent/<session> root")
	fmt.Fprintln(out, "  --pid PID            resolve from runtime-root pid/<pid>")
	fmt.Fprintln(out, "  --session ID         resolve from runtime-root agent/<session>")
	fmt.Fprintln(out, "  --runtime-root DIR   constrain discovery to one runtime root")
	fmt.Fprintln(out, "")
	fmt.Fprintln(out, "env:")
	fmt.Fprintln(out, "  QCLI_HTTP_ADDR, QCLI_CTL_TIMEOUT_MS, QCLI_POLL_MS, QCLI_QUINE_BIN, QCLI_RUNTIME_ROOT")
	return 0
}

func runAttach(args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("attach", flag.ContinueOnError)
	fs.SetOutput(stderr)
	resolve := addResolveFlags(fs)
	httpAddr := fs.String("http", envOrDefault("QCLI_HTTP_ADDR", "127.0.0.1:0"), "HTTP listen address for qcli kernel")
	clientID := fs.String("client-id", "", "client endpoint id")
	headless := fs.Bool("headless", false, "serve only; do not launch TUI")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	ctx, cancel := signalContext()
	defer cancel()
	root, _ := qcli.DiscoverRuntimeRoot(resolve().RuntimeRoot)
	kernel, err := qcli.NewKernel(ctx, qcli.KernelOptions{RuntimeRoot: root, ClientID: *clientID})
	if err != nil {
		fmt.Fprintf(stderr, "qcli: %v\n", err)
		return 1
	}
	defer kernel.Close()
	if _, err := kernel.Attach(resolve()); err != nil {
		fmt.Fprintf(stderr, "qcli: attach: %v\n", err)
		return 1
	}
	return serveAndMaybeTUI(ctx, kernel, *httpAddr, !*headless, stdout, stderr)
}

func runSpawn(args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("spawn", flag.ContinueOnError)
	fs.SetOutput(stderr)
	runtimeRoot := fs.String("runtime-root", "", "runtime root to use for the spawned quine")
	quineBinary := fs.String("quine", "", "path to quine binary")
	httpAddr := fs.String("http", envOrDefault("QCLI_HTTP_ADDR", "127.0.0.1:0"), "HTTP listen address for qcli kernel")
	clientID := fs.String("client-id", "", "client endpoint id")
	headless := fs.Bool("headless", false, "serve only; do not launch TUI")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	mission := strings.TrimSpace(strings.Join(fs.Args(), " "))
	if mission == "" {
		mission = defaultSpawnMission
	}
	ctx, cancel := signalContext()
	defer cancel()
	root, err := qcli.DiscoverRuntimeRoot(*runtimeRoot)
	if err != nil {
		fmt.Fprintf(stderr, "qcli: %v\n", err)
		return 1
	}
	kernel, err := qcli.NewKernel(ctx, qcli.KernelOptions{RuntimeRoot: root, ClientID: *clientID, QuineBinary: *quineBinary})
	if err != nil {
		fmt.Fprintf(stderr, "qcli: %v\n", err)
		return 1
	}
	defer kernel.Close()
	if _, err := kernel.Spawn(qcli.SpawnOptions{Mission: mission, RuntimeRoot: root, QuineBinary: *quineBinary, Stdout: toFile(stdout), Stderr: toFile(stderr), Stdin: os.Stdin}); err != nil {
		fmt.Fprintf(stderr, "qcli: spawn: %v\n", err)
		return 1
	}
	return serveAndMaybeTUI(ctx, kernel, *httpAddr, !*headless, stdout, stderr)
}

func runServe(args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("serve", flag.ContinueOnError)
	fs.SetOutput(stderr)
	resolve := addResolveFlags(fs)
	httpAddr := fs.String("http", envOrDefault("QCLI_HTTP_ADDR", "127.0.0.1:7777"), "HTTP listen address for qcli kernel")
	clientID := fs.String("client-id", "", "client endpoint id")
	quineBinary := fs.String("quine", "", "path to quine binary for /command spawn")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	ctx, cancel := signalContext()
	defer cancel()
	root, err := qcli.DiscoverRuntimeRoot(resolve().RuntimeRoot)
	if err != nil {
		fmt.Fprintf(stderr, "qcli: %v\n", err)
		return 1
	}
	kernel, err := qcli.NewKernel(ctx, qcli.KernelOptions{RuntimeRoot: root, ClientID: *clientID, QuineBinary: *quineBinary})
	if err != nil {
		fmt.Fprintf(stderr, "qcli: %v\n", err)
		return 1
	}
	defer kernel.Close()
	if hasSelector(resolve()) {
		if _, err := kernel.Attach(resolve()); err != nil {
			fmt.Fprintf(stderr, "qcli: attach: %v\n", err)
			return 1
		}
	}
	fmt.Fprintf(stdout, "qcli: serving qcli/1 on http://%s\n", *httpAddr)
	fmt.Fprintf(stdout, "qcli: runtime-root = %s\n", root)
	if err := qcli.ServeHTTP(ctx, kernel, qcli.HTTPOptions{Addr: *httpAddr}); err != nil && !errors.Is(err, context.Canceled) {
		fmt.Fprintf(stderr, "qcli: serve: %v\n", err)
		return 1
	}
	return 0
}

// runWeb attaches (or resolves/spawns the default peer, like the no-arg
// entry), serves the kernel + built web client, prints the URL, and makes a
// best-effort attempt to open a browser (WEB.md §3 entry point).
func runWeb(args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("web", flag.ContinueOnError)
	fs.SetOutput(stderr)
	resolve := addResolveFlags(fs)
	httpAddr := fs.String("http", envOrDefault("QCLI_HTTP_ADDR", "127.0.0.1:0"), "HTTP listen address for qcli kernel")
	clientID := fs.String("client-id", "", "client endpoint id")
	quineBinary := fs.String("quine", "", "path to quine binary")
	noOpen := fs.Bool("no-open", false, "do not try to open a browser")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	ctx, cancel := signalContext()
	defer cancel()
	root, err := qcli.DiscoverRuntimeRoot(resolve().RuntimeRoot)
	if err != nil {
		fmt.Fprintf(stderr, "qcli: %v\n", err)
		return 1
	}
	kernel, err := qcli.NewKernel(ctx, qcli.KernelOptions{RuntimeRoot: root, ClientID: *clientID, QuineBinary: *quineBinary})
	if err != nil {
		fmt.Fprintf(stderr, "qcli: %v\n", err)
		return 1
	}
	defer kernel.Close()
	switch {
	case hasSelector(resolve()):
		if _, err := kernel.Attach(resolve()); err != nil {
			fmt.Fprintf(stderr, "qcli: attach: %v\n", err)
			return 1
		}
	default:
		if agent, err := qcli.DefaultAgent(root); err == nil {
			if err := kernel.AttachAgent(agent); err != nil {
				fmt.Fprintf(stderr, "qcli: attach default peer: %v\n", err)
				return 1
			}
		} else if _, err := kernel.Spawn(qcli.SpawnOptions{Mission: defaultSpawnMission, RuntimeRoot: root, QuineBinary: *quineBinary, Stdout: toFile(stdout), Stderr: toFile(stderr), Stdin: os.Stdin}); err != nil {
			fmt.Fprintf(stderr, "qcli: no single live peer and spawn failed: %v\n", err)
			return 1
		}
	}
	serverURL, stop, err := startHTTPServer(ctx, kernel, *httpAddr)
	if err != nil {
		fmt.Fprintf(stderr, "qcli: serve: %v\n", err)
		return 1
	}
	defer stop()
	fmt.Fprintf(stdout, "qcli: serving qcli/1 + web client on %s\n", serverURL)
	if dist := qcli.DiscoverWebDist(); dist == "" {
		fmt.Fprintf(stdout, "qcli: no built web client found — run `make build-qcli-web` and reload\n")
	}
	if !*noOpen {
		openBrowser(ctx, serverURL, stdout)
	}
	<-ctx.Done()
	return 0
}

// openBrowser is best-effort and non-fatal: absence of xdg-open/open just
// leaves the printed URL as the entry point.
func openBrowser(ctx context.Context, url string, stdout io.Writer) {
	for _, opener := range []string{"xdg-open", "open"} {
		if path, err := exec.LookPath(opener); err == nil {
			if err := exec.CommandContext(ctx, path, url).Start(); err == nil {
				return
			}
		}
	}
	fmt.Fprintf(stdout, "qcli: open %s in a browser\n", url)
}

func runInteractive(args []string, stdout, stderr io.Writer) int {
	ctx, cancel := signalContext()
	defer cancel()
	root, err := qcli.DiscoverRuntimeRoot("")
	if err != nil {
		fmt.Fprintf(stderr, "qcli: %v\n", err)
		return 1
	}
	kernel, err := qcli.NewKernel(ctx, qcli.KernelOptions{RuntimeRoot: root})
	if err != nil {
		fmt.Fprintf(stderr, "qcli: %v\n", err)
		return 1
	}
	defer kernel.Close()
	if agent, err := qcli.DefaultAgent(root); err == nil {
		if err := kernel.AttachAgent(agent); err != nil {
			fmt.Fprintf(stderr, "qcli: attach default peer: %v\n", err)
			return 1
		}
	} else {
		if peers := qcli.ScanRoster(root, nil); len(peers) > 1 {
			fmt.Fprintf(stderr, "qcli: multiple live peers; use qcli attach --pid PID or qcli serve --pid PID\n")
			return 1
		}
		if _, err := kernel.Spawn(qcli.SpawnOptions{Mission: defaultSpawnMission, RuntimeRoot: root, Stdout: toFile(stdout), Stderr: toFile(stderr), Stdin: os.Stdin}); err != nil {
			fmt.Fprintf(stderr, "qcli: no single live peer and spawn failed: %v\n", err)
			return 1
		}
	}
	return serveAndMaybeTUI(ctx, kernel, "127.0.0.1:0", true, stdout, stderr)
}

func serveAndMaybeTUI(ctx context.Context, kernel *qcli.Kernel, addr string, launchTUI bool, stdout, stderr io.Writer) int {
	serverURL, stop, err := startHTTPServer(ctx, kernel, addr)
	if err != nil {
		fmt.Fprintf(stderr, "qcli: serve: %v\n", err)
		return 1
	}
	defer stop()
	fmt.Fprintf(stdout, "qcli: serving qcli/1 on %s\n", serverURL)
	if !launchTUI {
		<-ctx.Done()
		return 0
	}
	code, err := runTUI(ctx, serverURL, stdout, stderr)
	if err != nil {
		fmt.Fprintf(stderr, "qcli: %v\n", err)
		fmt.Fprintf(stderr, "qcli: kernel remains available while this process runs: %s\n", serverURL)
		fmt.Fprintf(stderr, "qcli: manual TUI: QCLI_SERVER=%s bun run operator/qcli/tui/src/index.tsx\n", serverURL)
		<-ctx.Done()
		return 1
	}
	return code
}

func startHTTPServer(ctx context.Context, kernel *qcli.Kernel, addr string) (string, func(), error) {
	if strings.TrimSpace(addr) == "" {
		addr = "127.0.0.1:0"
	}
	if err := qcli.ValidateHTTPAddr(addr); err != nil {
		return "", nil, err
	}
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return "", nil, err
	}
	server := &http.Server{Handler: kernel.Handler()}
	done := make(chan struct{})
	go func() {
		defer close(done)
		_ = server.Serve(ln)
	}()
	go func() {
		<-ctx.Done()
		shutCtx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		_ = server.Shutdown(shutCtx)
	}()
	stop := func() {
		shutCtx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		_ = server.Shutdown(shutCtx)
		<-done
	}
	return "http://" + ln.Addr().String(), stop, nil
}

func runTUI(ctx context.Context, serverURL string, stdout, stderr io.Writer) (int, error) {
	bun, err := exec.LookPath("bun")
	if err != nil {
		return 1, fmt.Errorf("bun not found")
	}
	tuiDir, tuiEntry, err := findTUIEntry()
	if err != nil {
		return 1, err
	}
	cmd := exec.CommandContext(ctx, bun, "run", tuiEntry)
	cmd.Dir = tuiDir
	cmd.Env = append(os.Environ(), "QCLI_SERVER="+serverURL)
	cmd.Stdout = toFileOr(stdout, os.Stdout)
	cmd.Stderr = toFileOr(stderr, os.Stderr)
	cmd.Stdin = os.Stdin
	err = cmd.Run()
	if err == nil {
		return 0, nil
	}
	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) {
		return exitErr.ExitCode(), nil
	}
	return 1, err
}

func findTUIEntry() (string, string, error) {
	candidates := []string{}
	if wd, err := os.Getwd(); err == nil {
		for dir := wd; ; dir = filepath.Dir(dir) {
			candidates = append(candidates, filepath.Join(dir, "operator", "qcli", "tui"))
			parent := filepath.Dir(dir)
			if parent == dir {
				break
			}
		}
	}
	if exe, err := os.Executable(); err == nil {
		dir := filepath.Dir(exe)
		candidates = append(candidates,
			filepath.Join(dir, "operator", "qcli", "tui"),
			filepath.Join(dir, "..", "operator", "qcli", "tui"),
		)
	}
	for _, dir := range candidates {
		entry := filepath.Join(dir, "src", "index.tsx")
		if info, err := os.Stat(entry); err == nil && !info.IsDir() {
			return dir, filepath.Join("src", "index.tsx"), nil
		}
	}
	return "", "", fmt.Errorf("cannot locate operator/qcli/tui/src/index.tsx")
}

func addResolveFlags(fs *flag.FlagSet) func() qcli.ResolveOptions {
	runtimeRoot := fs.String("runtime-root", "", "runtime root containing pid/ and agent/")
	agentRoot := fs.String("agent-root", "", "agent root path")
	pid := fs.Int("pid", 0, "target pid")
	session := fs.String("session", "", "target session id")
	return func() qcli.ResolveOptions {
		return qcli.ResolveOptions{RuntimeRoot: *runtimeRoot, AgentRoot: *agentRoot, PID: *pid, Session: *session}
	}
}

func hasSelector(opts qcli.ResolveOptions) bool {
	n := 0
	if strings.TrimSpace(opts.AgentRoot) != "" {
		n++
	}
	if opts.PID > 0 {
		n++
	}
	if strings.TrimSpace(opts.Session) != "" {
		n++
	}
	return n > 0
}

func signalContext() (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(context.Background())
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	go func() {
		select {
		case <-sigCh:
			cancel()
		case <-ctx.Done():
		}
		signal.Stop(sigCh)
	}()
	return ctx, cancel
}

func envOrDefault(key, fallback string) string {
	if v := strings.TrimSpace(os.Getenv(key)); v != "" {
		return v
	}
	return fallback
}

func toFile(w io.Writer) *os.File {
	if f, ok := w.(*os.File); ok {
		return f
	}
	return nil
}

func toFileOr(w io.Writer, fallback *os.File) *os.File {
	if f, ok := w.(*os.File); ok {
		return f
	}
	return fallback
}

func intString(v int) string {
	return strconv.Itoa(v)
}
