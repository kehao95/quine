package app

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/kehao95/quine/internal/qcli"
)

func TestRunHelp(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := Run([]string{"--help"}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("code=%d", code)
	}
	if stderr.Len() != 0 {
		t.Fatalf("stderr=%q", stderr.String())
	}
	for _, want := range []string{"attach", "spawn", "serve", "QCLI_HTTP_ADDR"} {
		if !strings.Contains(stdout.String(), want) {
			t.Fatalf("help missing %q: %s", want, stdout.String())
		}
	}
}

func TestStartHTTPServerHealthz(t *testing.T) {
	agentRoot := makeAppTestAgent(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	root := filepath.Dir(filepath.Dir(agentRoot))
	kernel, err := qcli.NewKernel(ctx, qcli.KernelOptions{RuntimeRoot: root, PollInterval: 20 * time.Millisecond})
	if err != nil {
		t.Fatal(err)
	}
	defer kernel.Close()
	if _, err := kernel.Attach(qcli.ResolveOptions{RuntimeRoot: root, AgentRoot: agentRoot}); err != nil {
		t.Fatal(err)
	}

	url, stop, err := startHTTPServer(ctx, kernel, "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer stop()
	resp, err := http.Get(url + "/healthz")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	data, _ := io.ReadAll(resp.Body)
	if !strings.Contains(string(data), `"contract":"qcli/1"`) || !strings.Contains(string(data), `"attached":true`) {
		t.Fatalf("healthz=%s", data)
	}
}

func TestAttachUsesQCLIHTTPAddrDefault(t *testing.T) {
	agentRoot := makeAppTestAgent(t)
	t.Setenv("QCLI_HTTP_ADDR", "0.0.0.0:0")
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := Run([]string{"attach", "--headless", "--agent-root", agentRoot}, &stdout, &stderr)
	if code != 1 {
		t.Fatalf("code=%d stdout=%q stderr=%q", code, stdout.String(), stderr.String())
	}
	if !strings.Contains(stderr.String(), "QCLI_HTTP_UNSAFE=1") {
		t.Fatalf("stderr missing unsafe bind refusal: %q", stderr.String())
	}
}

func TestFindTUIEntryReturnsTUIDir(t *testing.T) {
	tuiDir, entry, err := findTUIEntry()
	if err != nil {
		t.Fatal(err)
	}
	if entry != filepath.Join("src", "index.tsx") {
		t.Fatalf("entry = %q", entry)
	}
	if _, err := os.Stat(filepath.Join(tuiDir, "bunfig.toml")); err != nil {
		t.Fatalf("tui dir %q does not contain bunfig.toml: %v", tuiDir, err)
	}
	if _, err := os.Stat(filepath.Join(tuiDir, entry)); err != nil {
		t.Fatalf("tui entry missing: %v", err)
	}
}

func makeAppTestAgent(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	session := "sess_app"
	agentRoot := filepath.Join(root, "agent", session)
	publicRoot := filepath.Join(agentRoot, "public")
	for _, dir := range []string{
		filepath.Join(agentRoot, "status"),
		filepath.Join(agentRoot, "context", "state"),
		filepath.Join(publicRoot, "status"),
		filepath.Join(publicRoot, "ctl"),
		filepath.Join(publicRoot, "log"),
		filepath.Join(root, "pid"),
	} {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			t.Fatal(err)
		}
	}
	status := `{"session_id":"sess_app","run_id":"run_app","incarnation_id":1,"pid":` + intString(os.Getpid()) + `,"ppid":1,"runtime_root":"` + root + `","agent_root":"` + agentRoot + `","model_id":"gpt-5.5","depth":1}` + "\n"
	for _, path := range []string{filepath.Join(agentRoot, "status", "session.json"), filepath.Join(publicRoot, "status", "session.json")} {
		if err := os.WriteFile(path, []byte(status), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	if err := os.WriteFile(filepath.Join(publicRoot, "status", "inbox.json"), []byte(`{"pending_count":0}`+"\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(publicRoot, "status", "contract.json"), []byte(`{"contract_version":"process-control/v0"}`+"\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(agentRoot, "context", "state", "current.jsonl"), nil, 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(agentRoot, "context", "state", "live.jsonl"), nil, 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(publicRoot, "log", "control.jsonl"), nil, 0o644); err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"post", "poke", "inject", "interrupt"} {
		if err := syscall.Mkfifo(filepath.Join(publicRoot, "ctl", name), 0o666); err != nil {
			t.Fatal(err)
		}
	}
	if err := os.Symlink(publicRoot, filepath.Join(root, "pid", intString(os.Getpid()))); err != nil {
		t.Fatal(err)
	}
	return agentRoot
}
