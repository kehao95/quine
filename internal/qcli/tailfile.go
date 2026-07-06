package qcli

import (
	"bytes"
	"context"
	"errors"
	"io"
	"os"
	"syscall"
	"time"
)

type fileIdentity struct {
	dev uint64
	ino uint64
}

type tailState struct {
	offset  int64
	partial []byte
	id      fileIdentity
	haveID  bool
}

type tailReadResult struct {
	lines []string
	reset bool
}

func pollIntervalFromEnv(defaultPoll time.Duration) time.Duration {
	if defaultPoll <= 0 {
		defaultPoll = 150 * time.Millisecond
	}
	raw := os.Getenv("QCLI_POLL_MS")
	if raw == "" {
		return defaultPoll
	}
	var ms int
	for _, b := range []byte(raw) {
		if b < '0' || b > '9' {
			return defaultPoll
		}
		ms = ms*10 + int(b-'0')
	}
	if ms <= 0 {
		return defaultPoll
	}
	return time.Duration(ms) * time.Millisecond
}

func readTailLines(path string, state *tailState) (tailReadResult, error) {
	info, err := os.Stat(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return tailReadResult{}, nil
		}
		return tailReadResult{}, err
	}
	id := identityOf(info)
	reset := false
	if state.haveID && id != state.id {
		reset = true
	}
	if info.Size() < state.offset {
		reset = true
	}
	if !state.haveID || reset {
		state.offset = 0
		state.partial = nil
		state.id = id
		state.haveID = true
	}

	f, err := os.Open(path)
	if err != nil {
		return tailReadResult{}, err
	}
	defer f.Close()
	if _, err := f.Seek(state.offset, io.SeekStart); err != nil {
		return tailReadResult{}, err
	}
	data, err := io.ReadAll(f)
	if err != nil {
		return tailReadResult{}, err
	}
	state.offset += int64(len(data))
	if len(data) == 0 {
		return tailReadResult{reset: reset}, nil
	}

	combined := append(append([]byte{}, state.partial...), data...)
	parts := bytes.Split(combined, []byte{'\n'})
	state.partial = append(state.partial[:0], parts[len(parts)-1]...)
	lines := make([]string, 0, len(parts)-1)
	for _, part := range parts[:len(parts)-1] {
		if len(bytes.TrimSpace(part)) == 0 {
			continue
		}
		lines = append(lines, string(part))
	}
	return tailReadResult{lines: lines, reset: reset}, nil
}

func readAllCompleteLines(path string) ([]string, tailState, error) {
	var st tailState
	res, err := readTailLines(path, &st)
	if err != nil {
		return nil, st, err
	}
	return res.lines, st, nil
}

func identityOf(info os.FileInfo) fileIdentity {
	if st, ok := info.Sys().(*syscall.Stat_t); ok {
		return fileIdentity{dev: uint64(st.Dev), ino: uint64(st.Ino)}
	}
	return fileIdentity{dev: uint64(info.ModTime().UnixNano()), ino: uint64(info.Size())}
}

func waitTick(ctx context.Context, d time.Duration) bool {
	timer := time.NewTimer(d)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return false
	case <-timer.C:
		return true
	}
}
