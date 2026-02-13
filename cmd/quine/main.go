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
			os.Exit(126)
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
	// - TTY (no pipe): material = "Begin.", stdinReader = nil
	// - Piped data: requires --stdin-mode flag to specify handling
	material, stdinReader, err := handleStdin(cfg, mode)
	if err != nil {
		fmt.Fprintf(os.Stderr, "quine: reading stdin: %v\n", err)
		os.Exit(2)
	}

	rt, err := runtime.New(cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "quine: %v\n", err)
		os.Exit(1)
	}

	exitCode := rt.Run(mission, material, stdinReader)
	os.Exit(exitCode)
}

// handleStdin determines how to handle stdin and returns:
//   - material: the initial User Message content
//   - stdinReader: io.Reader for streaming (nil if no stdin pipe or binary mode)
//   - error: if handling fails
//
// The mode parameter controls behavior:
//   - stdinModeText: treat stdin as text, enable streaming (default)
//   - stdinModeBinary: treat stdin as binary, save to file (-b flag)
//
// When stdin is TTY (no pipe): material="Begin.", reader=nil (mode ignored)
func handleStdin(cfg *config.Config, mode stdinMode) (material string, stdinReader io.Reader, err error) {
	stat, err := os.Stdin.Stat()
	if err != nil {
		return "", nil, fmt.Errorf("stat stdin: %w", err)
	}

	// TTY (no piped data) — no mode needed
	if (stat.Mode() & os.ModeCharDevice) != 0 {
		return "Begin.", nil, nil
	}

	// Stdin is piped
	if mode == stdinModeText {
		// If resuming after exec and stdin is seekable, seek to the saved position
		// This handles the case where bufio had pre-read data that we need to re-read
		if cfg.StdinOffset > 0 {
			_, err := os.Stdin.Seek(cfg.StdinOffset, io.SeekStart)
			if err != nil {
				// Seek failed (probably a pipe) — this is expected for pipes
				// For pipes, exec preserves the position but loses bufio's buffer
				// The countingReader tracks what's been pulled from the kernel
			}
		}
		return "Streaming input available. Use the `read` tool to read lines from stdin. Read until EOF.", os.Stdin, nil
	}

	// Binary mode: read all and save to file
	allData, err := io.ReadAll(os.Stdin)
	if err != nil {
		return "", nil, fmt.Errorf("read binary stdin: %w", err)
	}
	if len(allData) == 0 {
		return "Begin.", nil, nil
	}

	os.MkdirAll(cfg.DataDir, 0o755)
	binaryPath := filepath.Join(cfg.DataDir, fmt.Sprintf("stdin-%s.bin", cfg.SessionID))
	if err := os.WriteFile(binaryPath, allData, 0o644); err != nil {
		return "", nil, fmt.Errorf("write binary stdin: %w", err)
	}

	return fmt.Sprintf("User sent a binary file at %s", binaryPath), nil, nil
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
