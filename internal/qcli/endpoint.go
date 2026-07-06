package qcli

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"time"
)

type ClientEndpoint struct {
	ID             string `json:"id"`
	RuntimeRoot    string `json:"runtime_root"`
	EndpointRoot   string `json:"endpoint_root"`
	PublicRoot     string `json:"public_root"`
	ControlPath    string `json:"control_path"`
	StatusPath     string `json:"status_path"`
	InboxPath      string `json:"inbox_path"`
	ControlLogPath string `json:"control_log_path"`
}

type endpointStatus struct {
	ID           string `json:"id"`
	Kind         string `json:"kind"`
	RuntimeRoot  string `json:"runtime_root"`
	EndpointRoot string `json:"endpoint_root"`
	PublicRoot   string `json:"public_root"`
	ControlPath  string `json:"control_path"`
	CreatedAt    int64  `json:"created_at"`
}

type endpointInbox struct {
	PendingCount int `json:"pending_count"`
	Messages     []struct {
		ID         string `json:"id"`
		Payload    string `json:"payload"`
		ReceivedAt int64  `json:"received_at"`
	} `json:"messages,omitempty"`
}

type endpointControlMessage struct {
	Payload    string
	ReceivedAt int64
}

func OpenClientEndpoint(runtimeRoot, id string) (ClientEndpoint, error) {
	root := strings.TrimSpace(runtimeRoot)
	if root == "" {
		return ClientEndpoint{}, fmt.Errorf("qcli: client endpoint requires runtime root")
	}
	if strings.TrimSpace(id) == "" {
		id = fmt.Sprintf("qcli-%d-%d", time.Now().UnixMilli(), os.Getpid())
	}
	endpointRoot := filepath.Join(root, "client", id)
	publicRoot := filepath.Join(endpointRoot, "public")
	statusDir := filepath.Join(publicRoot, "status")
	logDir := filepath.Join(publicRoot, "log")
	controlPath := filepath.Join(publicRoot, "ctl")
	for _, dir := range []struct {
		path  string
		label string
	}{
		{statusDir, "client endpoint status dir"},
		{logDir, "client endpoint log dir"},
		{controlPath, "client endpoint ctl dir"},
	} {
		if err := os.MkdirAll(dir.path, 0o755); err != nil {
			return ClientEndpoint{}, fmt.Errorf("qcli: create %s: %w", dir.label, err)
		}
	}
	if err := ensureClientControlFIFO(filepath.Join(controlPath, string(ControlActionPost))); err != nil {
		return ClientEndpoint{}, err
	}

	endpoint := ClientEndpoint{
		ID:             id,
		RuntimeRoot:    root,
		EndpointRoot:   endpointRoot,
		PublicRoot:     publicRoot,
		ControlPath:    controlPath,
		StatusPath:     filepath.Join(statusDir, "client.json"),
		InboxPath:      filepath.Join(statusDir, "inbox.json"),
		ControlLogPath: filepath.Join(logDir, "control.jsonl"),
	}
	status := endpointStatus{
		ID:           endpoint.ID,
		Kind:         "qcli-client",
		RuntimeRoot:  endpoint.RuntimeRoot,
		EndpointRoot: endpoint.EndpointRoot,
		PublicRoot:   endpoint.PublicRoot,
		ControlPath:  endpoint.ControlPath,
		CreatedAt:    time.Now().UnixMilli(),
	}
	if err := writeJSONFile(endpoint.StatusPath, status); err != nil {
		return ClientEndpoint{}, err
	}
	if err := writeJSONFile(endpoint.InboxPath, endpointInbox{PendingCount: 0}); err != nil {
		return ClientEndpoint{}, err
	}
	if _, err := os.Stat(endpoint.ControlLogPath); err != nil {
		if os.IsNotExist(err) {
			if err := os.WriteFile(endpoint.ControlLogPath, nil, 0o644); err != nil {
				return ClientEndpoint{}, fmt.Errorf("qcli: init client control log: %w", err)
			}
		} else {
			return ClientEndpoint{}, fmt.Errorf("qcli: stat client control log: %w", err)
		}
	}
	return endpoint, nil
}

func ServeClientEndpoint(ctx context.Context, endpoint ClientEndpoint, onPayload func(string, int64)) error {
	var seq int
	messages := make(chan endpointControlMessage)
	errs := make(chan error, 1)
	go func() {
		errs <- serveClientEndpointPost(ctx, filepath.Join(endpoint.ControlPath, string(ControlActionPost)), messages)
	}()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case err := <-errs:
			if err != nil {
				return err
			}
			return nil
		case msg := <-messages:
			seq++
			messageID := fmt.Sprintf("reply-%d", seq)
			if err := appendEndpointControlLog(endpoint.ControlLogPath, "received", ControlActionPost, "direct", messageID, msg.Payload, msg.ReceivedAt); err != nil {
				return err
			}
			if err := writeJSONFile(endpoint.InboxPath, endpointInbox{PendingCount: 0}); err != nil {
				return err
			}
			if onPayload != nil {
				onPayload(msg.Payload, msg.ReceivedAt)
			}
		}
	}
}

func ensureClientControlFIFO(path string) error {
	if info, err := os.Lstat(path); err == nil {
		if info.Mode()&os.ModeNamedPipe != 0 {
			return nil
		}
		if err := os.RemoveAll(path); err != nil {
			return fmt.Errorf("qcli: remove stale client ctl surface %s: %w", path, err)
		}
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("qcli: stat client ctl surface %s: %w", path, err)
	}
	if err := syscall.Mkfifo(path, 0o666); err != nil && !os.IsExist(err) {
		return fmt.Errorf("qcli: create client ctl fifo %s: %w", path, err)
	}
	return nil
}

func serveClientEndpointPost(ctx context.Context, path string, out chan<- endpointControlMessage) error {
	for {
		payload, err := readFIFOOnce(ctx, path)
		if err != nil {
			return err
		}
		select {
		case out <- endpointControlMessage{Payload: payload, ReceivedAt: time.Now().UnixMilli()}:
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func readFIFOOnce(ctx context.Context, path string) (string, error) {
	type result struct {
		payload string
		err     error
	}
	done := make(chan result, 1)
	go func() {
		f, err := os.OpenFile(path, os.O_RDONLY, os.ModeNamedPipe)
		if err != nil {
			done <- result{err: fmt.Errorf("qcli: open client ctl fifo: %w", err)}
			return
		}
		defer f.Close()
		data, err := io.ReadAll(f)
		if err != nil {
			done <- result{err: fmt.Errorf("qcli: read client ctl fifo: %w", err)}
			return
		}
		done <- result{payload: string(data)}
	}()

	select {
	case <-ctx.Done():
		deadline := time.NewTimer(2 * time.Second)
		defer deadline.Stop()
		tick := time.NewTicker(10 * time.Millisecond)
		defer tick.Stop()
		for {
			if f, err := os.OpenFile(path, os.O_RDWR|syscall.O_NONBLOCK, os.ModeNamedPipe); err == nil {
				_ = f.Close()
			}
			select {
			case <-done:
				return "", ctx.Err()
			case <-deadline.C:
				return "", ctx.Err()
			case <-tick.C:
			}
		}
	case res := <-done:
		return res.payload, res.err
	}
}

func removeClientEndpoint(endpoint ClientEndpoint) {
	if endpoint.EndpointRoot == "" {
		return
	}
	_ = os.RemoveAll(endpoint.EndpointRoot)
}

func writeJSONFile(path string, value interface{}) error {
	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return fmt.Errorf("qcli: marshal %s: %w", path, err)
	}
	data = append(data, '\n')
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("qcli: write %s: %w", path, err)
	}
	return nil
}

func appendEndpointControlLog(path, kind string, action ControlAction, delivery, messageID, payload string, receivedAt int64) error {
	entry := struct {
		Kind      string `json:"kind"`
		Timestamp int64  `json:"timestamp"`
		Action    string `json:"action,omitempty"`
		Delivery  string `json:"delivery,omitempty"`
		Message   struct {
			ID         string `json:"id"`
			Action     string `json:"action,omitempty"`
			Delivery   string `json:"delivery,omitempty"`
			Payload    string `json:"payload"`
			ReceivedAt int64  `json:"received_at"`
		} `json:"message"`
	}{
		Kind:      kind,
		Timestamp: time.Now().UnixMilli(),
		Action:    string(action),
		Delivery:  delivery,
	}
	entry.Message.ID = messageID
	entry.Message.Action = string(action)
	entry.Message.Delivery = delivery
	entry.Message.Payload = payload
	entry.Message.ReceivedAt = receivedAt

	data, err := json.Marshal(entry)
	if err != nil {
		return fmt.Errorf("qcli: marshal control log entry: %w", err)
	}
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return fmt.Errorf("qcli: open control log: %w", err)
	}
	defer f.Close()
	if _, err := f.Write(append(data, '\n')); err != nil {
		return fmt.Errorf("qcli: append control log: %w", err)
	}
	return nil
}

func ignoreCanceled(err error) error {
	if errors.Is(err, context.Canceled) {
		return nil
	}
	return err
}
