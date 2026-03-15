package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/kehao95/quine/internal/config"
	"github.com/kehao95/quine/internal/runtime"
)

// stdinMode represents the expected stdin input type
type stdinMode int

const (
	stdinModeText   stdinMode = iota // default: text streaming
	stdinModeBinary                  // -b: binary input (save to file)
)

func main() {
	// Parse flags
	binaryMode := flag.Bool("b", false, "treat stdin as binary (save to file instead of streaming)")
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "usage: quine [-b] <mission>")
		fmt.Fprintln(os.Stderr, "       echo <text> | quine <mission>")
		fmt.Fprintln(os.Stderr, "       cat file.bin | quine -b <mission>")
		fmt.Fprintln(os.Stderr, "")
		fmt.Fprintln(os.Stderr, "flags:")
		flag.PrintDefaults()
		fmt.Fprintln(os.Stderr, "")
		fmt.Fprintln(os.Stderr, "workspace physics are enabled via QUINE_WORKSPACE / QUINE_WORKSPACE_ROOT env vars.")
	}
	flag.Parse()

	// Determine stdin mode
	mode := stdinModeText
	if *binaryMode {
		mode = stdinModeBinary
	}

	cfg, err := config.Load()
	if err != nil {
		if errors.Is(err, config.ErrDepthExceeded) {
			fmt.Fprintf(os.Stderr, "quine: max recursion depth exceeded (%d/%d)\n",
				depthFromEnv(), maxDepthFromEnv())
			os.Exit(1)
		}
		fmt.Fprintf(os.Stderr, "quine: %v\n", err)
		os.Exit(2)
	}

	// Determine mission: from remaining args or QUINE_ORIGINAL_INTENT (post-exec)
	//
	// The Quad-Channel Protocol (see Artifacts/implementation.md):
	//   - argv:   Intent/Mission (code segment, immutable)
	//   - stdin:  Material (data stream to process)
	//   - stdout: Deliverable (pure output)
	//   - stderr: Failure gradient (post-mortem)
	var mission string
	if cfg.OriginalIntent != "" {
		// Post-exec: mission preserved via environment
		mission = cfg.OriginalIntent
	} else if flag.NArg() > 0 {
		// Normal startup: mission from remaining args after flags
		mission = strings.Join(flag.Args(), " ")
	} else {
		flag.Usage()
		os.Exit(2)
	}

	if strings.TrimSpace(mission) == "" {
		fmt.Fprintln(os.Stderr, "quine: mission cannot be empty")
		os.Exit(2)
	}

	// Handle stdin:
	// - TTY (no pipe): material = "Begin."
	// - Piped text: material tells agent stdin is available via fd 3
	// - Piped binary (-b): material references saved file
	material, err := handleStdin(cfg, mode)
	if err != nil {
		fmt.Fprintf(os.Stderr, "quine: reading stdin: %v\n", err)
		os.Exit(2)
	}

	rt, err := runtime.New(cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "quine: %v\n", err)
		os.Exit(1)
	}

	exitCode := rt.Run(mission, material)
	os.Exit(exitCode)
}

// handleStdin determines how to handle stdin and returns the initial User
// Message content (material).
//
// The mode parameter controls behavior:
//   - stdinModeText: tell the agent that material is available via fd 3
//   - stdinModeBinary: read all stdin and save to a file (-b flag)
//
// When stdin is TTY (no pipe): material="Begin." (mode ignored)
func handleStdin(cfg *config.Config, mode stdinMode) (material string, err error) {
	stat, err := os.Stdin.Stat()
	if err != nil {
		return "", fmt.Errorf("stat stdin: %w", err)
	}

	// TTY (no piped data) — no mode needed
	if (stat.Mode() & os.ModeCharDevice) != 0 {
		return "Begin.", nil
	}

	// Stdin is piped
	if mode == stdinModeText {
		// Material is available to sh commands via runtime side channel fd 3.
		return "Input is piped as material on `fd 3`. This is the quine process stdin: it may behave like a finite batch or an open stream. Reads consume data, so later calls see only the unread remainder. Partial reads or redirects from `fd 3` can leave later steps with only a truncated remainder. It is available in `sh` via `cat <&3` or `cat /dev/fd/3`.", nil
	}

	// Binary mode: read all and save to file
	allData, err := io.ReadAll(os.Stdin)
	if err != nil {
		return "", fmt.Errorf("read binary stdin: %w", err)
	}
	if len(allData) == 0 {
		return "Begin.", nil
	}

	os.MkdirAll(cfg.DataDir, 0o755)
	binaryPath := filepath.Join(cfg.DataDir, fmt.Sprintf("stdin-%s.bin", cfg.SessionID))
	if err := os.WriteFile(binaryPath, allData, 0o644); err != nil {
		return "", fmt.Errorf("write binary stdin: %w", err)
	}

	return fmt.Sprintf("User sent a binary file at %s", binaryPath), nil
}

// depthFromEnv reads QUINE_DEPTH from environment for error reporting.
func depthFromEnv() int {
	v, err := strconv.Atoi(os.Getenv("QUINE_DEPTH"))
	if err != nil {
		return 0
	}
	return v
}

// maxDepthFromEnv reads QUINE_MAX_DEPTH from environment for error reporting.
func maxDepthFromEnv() int {
	v, err := strconv.Atoi(os.Getenv("QUINE_MAX_DEPTH"))
	if err != nil {
		return 5
	}
	return v
}
