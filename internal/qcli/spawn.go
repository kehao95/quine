package qcli

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

var (
	currentExecutable = os.Executable
	lookPath          = exec.LookPath
)

type SpawnOptions struct {
	Mission     string
	QuineBinary string
	RuntimeRoot string
	WorkDir     string
	Env         []string
	Stdout      *os.File
	Stderr      *os.File
	Stdin       *os.File
	WaitTimeout time.Duration
}

func Spawn(opts SpawnOptions) (Agent, error) {
	quineBinary, err := resolveQuineBinary(strings.TrimSpace(opts.QuineBinary))
	if err != nil {
		return Agent{}, err
	}
	runtimeRoot, err := defaultRuntimeRoot(strings.TrimSpace(opts.RuntimeRoot), opts.WorkDir)
	if err != nil {
		return Agent{}, err
	}

	var args []string
	if mission := strings.TrimSpace(opts.Mission); mission != "" {
		args = append(args, mission)
	}
	cmd := exec.Command(quineBinary, args...)
	cmd.Dir = opts.WorkDir
	cmd.Env = setEnvKeyDefault(setEnvKey(opts.EnvOrDefault(), "QUINE_DATA_DIR", runtimeRoot), "QUINE_IDLE_ENABLED", "1")
	if opts.Stdout != nil {
		cmd.Stdout = opts.Stdout
	}
	if opts.Stderr != nil {
		cmd.Stderr = opts.Stderr
	}
	if opts.Stdin != nil {
		cmd.Stdin = opts.Stdin
	}
	if err := cmd.Start(); err != nil {
		return Agent{}, fmt.Errorf("qcli: spawn quine: %w", err)
	}

	exited := make(chan error, 1)
	go func() { exited <- cmd.Wait() }()

	wait := opts.WaitTimeout
	if wait <= 0 {
		wait = 5 * time.Second
	}
	deadline := time.NewTimer(wait)
	defer deadline.Stop()
	tick := time.NewTicker(100 * time.Millisecond)
	defer tick.Stop()

	link := filepath.Join(runtimeRoot, "pid", strconv.Itoa(cmd.Process.Pid))
	var lastErr error
	for {
		select {
		case err := <-exited:
			if err != nil {
				return Agent{}, fmt.Errorf("qcli: spawned quine exited before registration: %w", err)
			}
			return Agent{}, fmt.Errorf("qcli: spawned quine exited before registration")
		case <-deadline.C:
			if lastErr != nil {
				return Agent{}, fmt.Errorf("%w: timed out waiting for spawned quine to become ready under %s: %v", ErrRegisterTimeout, link, lastErr)
			}
			return Agent{}, fmt.Errorf("%w: timed out waiting for spawned quine to register under %s", ErrRegisterTimeout, link)
		case <-tick.C:
			if agentRoot, err := filepath.EvalSymlinks(link); err == nil {
				agent, loadErr := loadAgentFromRoot(agentRoot)
				if loadErr == nil {
					return agent, nil
				}
				lastErr = loadErr
			}
		}
	}
}

func (opts SpawnOptions) EnvOrDefault() []string {
	if opts.Env != nil {
		return opts.Env
	}
	return os.Environ()
}

func resolveQuineBinary(explicit string) (string, error) {
	if strings.TrimSpace(explicit) != "" {
		return filepath.Abs(strings.TrimSpace(explicit))
	}
	if env := strings.TrimSpace(os.Getenv("QCLI_QUINE_BIN")); env != "" {
		return filepath.Abs(env)
	}
	if self, err := currentExecutable(); err == nil {
		sibling := filepath.Join(filepath.Dir(self), "quine")
		if info, err := os.Stat(sibling); err == nil && !info.IsDir() {
			return sibling, nil
		}
	}
	path, err := lookPath("quine")
	if err != nil {
		return "", fmt.Errorf("qcli: cannot resolve quine binary; set QCLI_QUINE_BIN, pass --quine, place sibling quine next to qcli, or put quine on PATH")
	}
	return path, nil
}

func setEnvKey(env []string, key, value string) []string {
	prefix := key + "="
	out := make([]string, 0, len(env)+1)
	replaced := false
	for _, entry := range env {
		if strings.HasPrefix(entry, prefix) {
			if !replaced {
				out = append(out, prefix+value)
				replaced = true
			}
			continue
		}
		out = append(out, entry)
	}
	if !replaced {
		out = append(out, prefix+value)
	}
	return out
}

func setEnvKeyDefault(env []string, key, value string) []string {
	prefix := key + "="
	for _, entry := range env {
		if strings.HasPrefix(entry, prefix) {
			return env
		}
	}
	return append(env, prefix+value)
}
