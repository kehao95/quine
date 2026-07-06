package qcli

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"syscall"
	"time"
)

var ErrWriteTimeout = errors.New("qcli: ctl write timeout")

const defaultCtlTimeout = 3 * time.Second
const maxAbandonedWriters = 8

type CtlWriter struct {
	slots chan struct{}
}

func NewCtlWriter() *CtlWriter {
	return &CtlWriter{slots: make(chan struct{}, maxAbandonedWriters)}
}

func CtlTimeoutFromEnv() time.Duration {
	raw := strings.TrimSpace(os.Getenv("QCLI_CTL_TIMEOUT_MS"))
	if raw == "" {
		return defaultCtlTimeout
	}
	ms, err := strconv.Atoi(raw)
	if err != nil || ms <= 0 {
		return defaultCtlTimeout
	}
	return time.Duration(ms) * time.Millisecond
}

func (w *CtlWriter) Write(path string, payload []byte, timeout time.Duration) error {
	if w == nil {
		w = NewCtlWriter()
	}
	if timeout <= 0 {
		timeout = defaultCtlTimeout
	}
	select {
	case w.slots <- struct{}{}:
	default:
		return ErrUnreachable
	}

	done := make(chan error, 1)
	go func() {
		defer func() { <-w.slots }()
		done <- writeCtl(path, payload)
	}()

	timer := time.NewTimer(timeout)
	defer timer.Stop()
	select {
	case err := <-done:
		return classifyCtlWriteError(err)
	case <-timer.C:
		return ErrWriteTimeout
	}
}

func writeCtl(path string, payload []byte) error {
	flags := os.O_WRONLY | os.O_TRUNC
	f, err := os.OpenFile(path, flags, os.ModeNamedPipe)
	if err != nil {
		return err
	}
	if len(payload) > 0 {
		if _, err := f.Write(payload); err != nil {
			_ = f.Close()
			return err
		}
	}
	return f.Close()
}

func classifyCtlWriteError(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, syscall.ENOENT) {
		return fmt.Errorf("%w: %v", ErrNoPeer, err)
	}
	if errors.Is(err, syscall.ENOTCONN) || errors.Is(err, syscall.ENXIO) || errors.Is(err, syscall.EPIPE) {
		return fmt.Errorf("%w: %v", ErrUnreachable, err)
	}
	return err
}
