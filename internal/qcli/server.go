package qcli

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

type Kernel struct {
	ctx    context.Context
	cancel context.CancelFunc

	mu          sync.Mutex
	session     *Session
	runtimeRoot string
	clientID    string
	quineBinary string
	poll        time.Duration
	ctlTimeout  time.Duration
	webDist     string
}

type KernelOptions struct {
	RuntimeRoot  string
	ClientID     string
	QuineBinary  string
	PollInterval time.Duration
	CtlTimeout   time.Duration
	// WebDist points at the built web client to serve under /. Empty means
	// discover (QCLI_WEB_DIST, then operator/qcli/web/dist); "-" disables.
	WebDist string
}

type CommandRequest struct {
	Op          string         `json:"op"`
	Selector    ResolveOptions `json:"selector,omitempty"`
	Mission     string         `json:"mission,omitempty"`
	RuntimeRoot string         `json:"runtime_root,omitempty"`
	Action      ControlAction  `json:"action,omitempty"`
	Text        string         `json:"text,omitempty"`
}

type HTTPOptions struct {
	Addr string
}

func NewKernel(ctx context.Context, opts KernelOptions) (*Kernel, error) {
	root := strings.TrimSpace(opts.RuntimeRoot)
	if root == "" {
		var err error
		root, err = DiscoverRuntimeRoot("")
		if err != nil {
			return nil, err
		}
	}
	webDist := strings.TrimSpace(opts.WebDist)
	switch webDist {
	case "":
		webDist = DiscoverWebDist()
	case "-":
		webDist = ""
	}
	kctx, cancel := context.WithCancel(ctx)
	return &Kernel{
		ctx:         kctx,
		cancel:      cancel,
		runtimeRoot: root,
		clientID:    opts.ClientID,
		quineBinary: opts.QuineBinary,
		poll:        opts.PollInterval,
		ctlTimeout:  opts.CtlTimeout,
		webDist:     webDist,
	}, nil
}

func (k *Kernel) Close() {
	k.cancel()
	k.mu.Lock()
	session := k.session
	k.session = nil
	k.mu.Unlock()
	if session != nil {
		session.Close()
	}
}

func (k *Kernel) Attach(opts ResolveOptions) (Agent, error) {
	if strings.TrimSpace(opts.RuntimeRoot) == "" {
		opts.RuntimeRoot = k.runtimeRoot
	}
	agent, err := Resolve(opts)
	if err != nil {
		return Agent{}, err
	}
	if err := k.setAgent(agent); err != nil {
		return Agent{}, err
	}
	return agent, nil
}

func (k *Kernel) AttachAgent(agent Agent) error {
	return k.setAgent(agent)
}

func (k *Kernel) Spawn(opts SpawnOptions) (Agent, error) {
	if strings.TrimSpace(opts.RuntimeRoot) == "" {
		opts.RuntimeRoot = k.runtimeRoot
	}
	if strings.TrimSpace(opts.QuineBinary) == "" {
		opts.QuineBinary = k.quineBinary
	}
	agent, err := Spawn(opts)
	if err != nil {
		return Agent{}, err
	}
	if err := k.setAgent(agent); err != nil {
		return Agent{}, err
	}
	return agent, nil
}

func (k *Kernel) setAgent(agent Agent) error {
	session, err := StartSession(k.ctx, agent, SessionOptions{
		ClientID:     k.clientID,
		PollInterval: k.poll,
		CtlTimeout:   k.ctlTimeout,
	})
	if err != nil {
		return err
	}
	k.mu.Lock()
	prev := k.session
	k.session = session
	k.runtimeRoot = agent.RuntimeRoot
	k.mu.Unlock()
	if prev != nil {
		prev.Close()
	}
	return nil
}

func (k *Kernel) currentSession() (*Session, bool) {
	k.mu.Lock()
	defer k.mu.Unlock()
	return k.session, k.session != nil
}

func (k *Kernel) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/events", k.handleEvents)
	mux.HandleFunc("/context", k.handleContext)
	mux.HandleFunc("/status", k.handleStatus)
	mux.HandleFunc("/roster", k.handleRoster)
	mux.HandleFunc("/peer-contract", k.handlePeerContract)
	mux.HandleFunc("/command", k.handleCommand)
	mux.HandleFunc("/healthz", k.handleHealthz)
	mux.HandleFunc("/history", k.handleHistory)
	mux.HandleFunc("/", k.handleRoot)
	return mux
}

func ServeHTTP(ctx context.Context, kernel *Kernel, opts HTTPOptions) error {
	addr := strings.TrimSpace(opts.Addr)
	if addr == "" {
		addr = strings.TrimSpace(os.Getenv("QCLI_HTTP_ADDR"))
	}
	if addr == "" {
		addr = "127.0.0.1:7777"
	}
	if err := ValidateHTTPAddr(addr); err != nil {
		return err
	}
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("qcli: listen %s: %w", addr, err)
	}
	server := &http.Server{Handler: kernel.Handler()}
	done := make(chan error, 1)
	go func() { done <- server.Serve(ln) }()
	select {
	case <-ctx.Done():
		shutCtx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		_ = server.Shutdown(shutCtx)
		return ctx.Err()
	case err := <-done:
		if err == http.ErrServerClosed {
			return nil
		}
		return err
	}
}

func ValidateHTTPAddr(addr string) error {
	if os.Getenv("QCLI_HTTP_UNSAFE") == "1" {
		return nil
	}
	host, _, err := net.SplitHostPort(addr)
	if err != nil {
		return fmt.Errorf("qcli: invalid http addr %q: %w", addr, err)
	}
	if host == "localhost" {
		return nil
	}
	ip := net.ParseIP(host)
	if ip != nil && ip.IsLoopback() {
		return nil
	}
	return fmt.Errorf("qcli: refusing non-loopback bind %q without QCLI_HTTP_UNSAFE=1", addr)
}

func (k *Kernel) handleEvents(w http.ResponseWriter, r *http.Request) {
	setJSONHeaders(w)
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	session, ok := k.currentSession()
	if !ok {
		writeError(w, http.StatusConflict, ErrorEvent{Type: "error", Scope: "server", Code: "no_peer", Message: "no peer attached", Recoverable: true})
		return
	}
	flusher, ok := w.(http.Flusher)
	if !ok {
		writeError(w, http.StatusInternalServerError, ErrorEvent{Type: "error", Scope: "server", Code: "streaming_unsupported", Message: "streaming unsupported", Recoverable: false})
		return
	}
	h := w.Header()
	h.Set("Content-Type", "text/event-stream")
	h.Set("Cache-Control", "no-cache")
	h.Set("Connection", "keep-alive")
	ch, unsub := session.Subscribe()
	defer unsub()
	for {
		select {
		case <-r.Context().Done():
			return
		case ev, ok := <-ch:
			if !ok {
				return
			}
			writeSSE(w, ev)
			flusher.Flush()
		}
	}
}

func (k *Kernel) handleContext(w http.ResponseWriter, r *http.Request) {
	setJSONHeaders(w)
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	session, ok := k.currentSession()
	if !ok {
		writeError(w, http.StatusConflict, ErrorEvent{Type: "error", Scope: "server", Code: "no_peer", Message: "no peer attached", Recoverable: true})
		return
	}
	_ = json.NewEncoder(w).Encode(session.Cells())
}

func (k *Kernel) handleStatus(w http.ResponseWriter, r *http.Request) {
	setJSONHeaders(w)
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	session, ok := k.currentSession()
	if !ok {
		writeError(w, http.StatusConflict, ErrorEvent{Type: "error", Scope: "server", Code: "no_peer", Message: "no peer attached", Recoverable: true})
		return
	}
	_ = json.NewEncoder(w).Encode(session.Status())
}

func (k *Kernel) handleRoster(w http.ResponseWriter, r *http.Request) {
	setJSONHeaders(w)
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var attached *Agent
	if session, ok := k.currentSession(); ok {
		agent := session.Agent()
		attached = &agent
	}
	_ = json.NewEncoder(w).Encode(RosterResponse{Type: "roster", Peers: ScanRoster(k.runtimeRoot, attached)})
}

func (k *Kernel) handlePeerContract(w http.ResponseWriter, r *http.Request) {
	setJSONHeaders(w)
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	session, ok := k.currentSession()
	if !ok {
		writeError(w, http.StatusConflict, ErrorEvent{Type: "error", Scope: "server", Code: "no_peer", Message: "no peer attached", Recoverable: true})
		return
	}
	contract, err := ReadPeerContract(session.Agent())
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			writeError(w, http.StatusNotFound, ErrorEvent{Type: "error", Scope: "server", Code: "not_found", Message: err.Error(), Recoverable: true})
			return
		}
		writeError(w, http.StatusInternalServerError, ErrorEvent{Type: "error", Scope: "server", Code: "read_failed", Message: err.Error(), Recoverable: true})
		return
	}
	_ = json.NewEncoder(w).Encode(PeerContractResponse{Type: "peer_contract", Contract: contract})
}

func (k *Kernel) handleCommand(w http.ResponseWriter, r *http.Request) {
	setJSONHeaders(w)
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var cmd CommandRequest
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		writeError(w, http.StatusBadRequest, ErrorEvent{Type: "error", Scope: "server", Code: "invalid_json", Message: err.Error(), Recoverable: true})
		return
	}
	switch strings.TrimSpace(cmd.Op) {
	case "attach":
		agent, err := k.Attach(cmd.Selector)
		if err != nil {
			writeError(w, statusForError(err), commandError("attach", err))
			return
		}
		_ = json.NewEncoder(w).Encode(AttachedResponse{Type: "attached", Peer: agent.PeerInfo()})
	case "spawn":
		root := strings.TrimSpace(cmd.RuntimeRoot)
		if root == "" {
			root = k.runtimeRoot
		}
		agent, err := k.Spawn(SpawnOptions{Mission: cmd.Mission, RuntimeRoot: root})
		if err != nil {
			writeError(w, http.StatusBadRequest, ErrorEvent{Type: "error", Scope: "spawn", Code: "spawn_failed", Message: err.Error(), Recoverable: false})
			return
		}
		_ = json.NewEncoder(w).Encode(AttachedResponse{Type: "attached", Peer: agent.PeerInfo()})
	case "send":
		session, ok := k.currentSession()
		if !ok {
			writeError(w, http.StatusConflict, ErrorEvent{Type: "error", Scope: "send", Code: "no_peer", Message: "no peer attached", Recoverable: true})
			return
		}
		action := cmd.Action
		if action == "" {
			action = ControlActionInject
		}
		resp, err := session.Send(action, cmd.Text)
		if err != nil {
			writeError(w, statusForError(err), commandError("send", err))
			return
		}
		_ = json.NewEncoder(w).Encode(resp)
	default:
		writeError(w, http.StatusBadRequest, ErrorEvent{Type: "error", Scope: "server", Code: "unknown_op", Message: fmt.Sprintf("unknown op %q", cmd.Op), Recoverable: true})
	}
}

func (k *Kernel) handleHealthz(w http.ResponseWriter, r *http.Request) {
	setJSONHeaders(w)
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	_, attached := k.currentSession()
	_ = json.NewEncoder(w).Encode(map[string]any{"ok": true, "contract": ContractVersion, "attached": attached})
}

func (k *Kernel) handleHistory(w http.ResponseWriter, r *http.Request) {
	setJSONHeaders(w)
	writeError(w, http.StatusNotImplemented, ErrorEvent{Type: "error", Scope: "server", Code: "reserved", Message: "/history is reserved in qcli/1", Recoverable: false})
}

func setJSONHeaders(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
}

func writeSSE(w http.ResponseWriter, ev any) {
	data, err := json.Marshal(ev)
	if err != nil {
		return
	}
	_, _ = fmt.Fprintf(w, "data: %s\n\n", data)
}

func writeError(w http.ResponseWriter, status int, ev ErrorEvent) {
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(ev)
}

func commandError(scope string, err error) ErrorEvent {
	code := "failed"
	recoverable := true
	switch {
	case errors.Is(err, ErrNoPeer):
		code = "no_peer"
	case errors.Is(err, ErrTargetRequired):
		code = "target_required"
	case errors.Is(err, ErrTargetNotFound):
		code = "not_found"
	case errors.Is(err, ErrRegisterTimeout):
		code = "register_timeout"
	case errors.Is(err, ErrWriteTimeout):
		code = "write_timeout"
	case errors.Is(err, ErrEmptyPayload):
		code = "empty_payload"
	case errors.Is(err, ErrUnreachable):
		code = "unreachable"
	}
	return ErrorEvent{Type: "error", Scope: scope, Code: code, Message: err.Error(), Recoverable: recoverable}
}

func statusForError(err error) int {
	switch {
	case errors.Is(err, ErrNoPeer):
		return http.StatusConflict
	case errors.Is(err, ErrTargetRequired), errors.Is(err, ErrEmptyPayload):
		return http.StatusBadRequest
	case errors.Is(err, ErrTargetNotFound):
		return http.StatusNotFound
	case errors.Is(err, ErrRegisterTimeout):
		return http.StatusServiceUnavailable
	case errors.Is(err, ErrWriteTimeout), errors.Is(err, ErrUnreachable):
		return http.StatusServiceUnavailable
	default:
		return http.StatusBadRequest
	}
}
