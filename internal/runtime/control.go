package runtime

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/kehao95/quine/internal/tape"
)

const controlInboxIndicator = "You have pending inbox messages. Inspect QUINE_AGENT_ROOT/status/inbox.json."

type controlDelivery string

const (
	controlDeliveryPoke      controlDelivery = "poke"
	controlDeliveryInject    controlDelivery = "inject"
	controlDeliveryInterrupt controlDelivery = "interrupt"
)

type controlSurfaceAction string

const (
	controlActionPost      controlSurfaceAction = "post"
	controlActionPoke      controlSurfaceAction = "poke"
	controlActionInject    controlSurfaceAction = "inject"
	controlActionInterrupt controlSurfaceAction = "interrupt"
)

var controlSurfaceActions = []controlSurfaceAction{
	controlActionPost,
	controlActionPoke,
	controlActionInject,
	controlActionInterrupt,
}

type controlState struct {
	mu                 sync.Mutex
	nextID             int
	pending            []controlMessage
	pokeRequested      bool
	injectRequested    bool
	interruptRequested bool
	eventCh            chan struct{}
	// started records that the control surface has been set up for this run.
	// It is a readiness/idempotency flag, not a listener-lifecycle handle:
	// control writes arrive through the FUSE ctl node, not a goroutine.
	started bool
}

type controlMessage struct {
	ID         string `json:"id"`
	Payload    string `json:"payload"`
	ReceivedAt int64  `json:"received_at"`
	Notified   bool   `json:"-"`
}

type controlInboxMessage struct {
	ID         string `json:"id"`
	Payload    string `json:"payload"`
	ReceivedAt int64  `json:"received_at"`
}

type controlInboxSnapshot struct {
	PendingCount int                   `json:"pending_count"`
	Messages     []controlInboxMessage `json:"messages,omitempty"`
}

type controlRuntimeMessage struct {
	ID         string `json:"id"`
	Delivery   string `json:"delivery"`
	Payload    string `json:"payload"`
	ReceivedAt int64  `json:"received_at"`
}

type controlLogEntry struct {
	Kind         string                 `json:"kind"`
	Timestamp    int64                  `json:"timestamp"`
	Action       string                 `json:"action,omitempty"`
	Delivery     string                 `json:"delivery,omitempty"`
	PendingCount int                    `json:"pending_count,omitempty"`
	Message      *controlRuntimeMessage `json:"message,omitempty"`
}

func newControlState() *controlState {
	return &controlState{
		eventCh: make(chan struct{}, 1),
	}
}

func newControlStateFromSnapshot(snapshot controlInboxSnapshot) *controlState {
	state := newControlState()
	state.pending = make([]controlMessage, 0, len(snapshot.Messages))
	for _, msg := range snapshot.Messages {
		state.pending = append(state.pending, controlMessage{
			ID:         msg.ID,
			Payload:    msg.Payload,
			ReceivedAt: msg.ReceivedAt,
		})
		if strings.HasPrefix(msg.ID, "message-") {
			if id, err := strconv.Atoi(strings.TrimPrefix(msg.ID, "message-")); err == nil && id > state.nextID {
				state.nextID = id
			}
		}
	}
	return state
}

func (c *controlState) signalEvent() {
	if c == nil || c.eventCh == nil {
		return
	}
	select {
	case c.eventCh <- struct{}{}:
	default:
	}
}

func normalizeControlPayload(payload string) string {
	payload = strings.TrimSuffix(payload, "\n")
	payload = strings.TrimSuffix(payload, "\r")
	return payload
}

func (r *Runtime) publicControlSurfaceSummary(action controlSurfaceAction) []byte {
	snapshot := r.controlSnapshot()
	lines := []string{
		fmt.Sprintf("backend: %s", runtimeSurfaceBackendName),
		fmt.Sprintf("control_file: %s", action),
		fmt.Sprintf("pending_count: %d", snapshot.PendingCount),
	}
	switch action {
	case controlActionPost:
		lines = append(lines,
			"mode: queue-only",
			"usage: write payload to queue it without waking the target process",
			`example: echo "hello" > ctl/post`,
			"empty_write: no-op",
		)
	case controlActionPoke:
		lines = append(lines,
			"mode: queue-and-resume",
			"usage: write payload to queue it and resume idle without context injection",
			`example: echo "hello" > ctl/poke`,
			"empty_write: : > ctl/poke resumes idle without new payload",
		)
	case controlActionInject:
		lines = append(lines,
			"mode: queue-and-deliver",
			"usage: write payload to queue it and deliver it at the next safe point",
			`example: echo "hello" > ctl/inject`,
			"empty_write: no-op",
		)
	case controlActionInterrupt:
		lines = append(lines,
			"mode: interrupt-delivery",
			"usage: write optional payload to queue it and request urgent interrupt delivery",
			`example: echo "hello" > ctl/interrupt`,
			"empty_write: : > ctl/interrupt requests interrupt without new payload",
		)
	}
	return []byte(strings.Join(lines, "\n") + "\n")
}

func (r *Runtime) applyControlSurfaceAction(action controlSurfaceAction, payload string) error {
	hasPayload := strings.TrimSpace(normalizeControlPayload(payload)) != ""
	switch action {
	case controlActionPost:
		if !hasPayload {
			return nil
		}
		_, err := r.enqueueControlPayload(action, payload)
		return err
	case controlActionPoke:
		if hasPayload {
			if _, err := r.enqueueControlPayload(action, payload); err != nil {
				return err
			}
		}
		r.requestControlPoke()
		return nil
	case controlActionInject:
		if !hasPayload {
			return nil
		}
		if _, err := r.enqueueControlPayload(action, payload); err != nil {
			return err
		}
		r.requestControlInject()
		return nil
	case controlActionInterrupt:
		if hasPayload {
			if _, err := r.enqueueControlPayload(action, payload); err != nil {
				return err
			}
		}
		r.requestControlInterrupt(hasPayload)
		return nil
	default:
		return fmt.Errorf("unknown control action %q", action)
	}
}

func (r *Runtime) controlInboxFilePath() string {
	return filepath.Join(r.cfg.SessionRetainedDir(""), "status", "inbox.json")
}

func (r *Runtime) controlLogFilePath() string {
	return r.cfg.SessionControlLogPath("")
}

func (r *Runtime) loadRetainedControlState() error {
	if r == nil || r.cfg == nil {
		return nil
	}
	data, err := os.ReadFile(r.controlInboxFilePath())
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("read retained control inbox: %w", err)
	}
	var snapshot controlInboxSnapshot
	if err := json.Unmarshal(data, &snapshot); err != nil {
		return fmt.Errorf("parse retained control inbox: %w", err)
	}
	r.control = newControlStateFromSnapshot(snapshot)
	return nil
}

func ensureDirPath(path string, label string) error {
	if info, err := os.Lstat(path); err == nil {
		if info.IsDir() {
			return nil
		}
		if info.Mode()&os.ModeSymlink != 0 {
			if targetInfo, err := os.Stat(path); err == nil && targetInfo.IsDir() {
				return nil
			}
		}
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("stat %s: %w", label, err)
	}
	if err := os.MkdirAll(path, 0o755); err != nil {
		if info, statErr := os.Stat(path); statErr == nil && info.IsDir() {
			return nil
		}
		return fmt.Errorf("mkdir %s: %w", label, err)
	}
	return nil
}

// ensureControlSurface materializes the retained backing files for the control
// surface: the inbox SST and the control log. The peer-facing control files
// (ctl/{post,poke,inject,interrupt}) are virtual FUSE nodes — no real special
// files are created here. See Paper/core/decisions/2026-06-control-surface-fuse-only.md.
func (r *Runtime) ensureControlSurface() error {
	if err := ensureDirPath(filepath.Dir(r.controlInboxFilePath()), "inbox dir"); err != nil {
		return err
	}
	if err := ensureDirPath(filepath.Dir(r.controlLogFilePath()), "control log dir"); err != nil {
		return err
	}
	if err := r.writeControlInboxSnapshot(r.controlSnapshot()); err != nil {
		return err
	}
	f, err := os.OpenFile(r.controlLogFilePath(), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return fmt.Errorf("open control log: %w", err)
	}
	return f.Close()
}

func (c *controlState) snapshotLocked() controlInboxSnapshot {
	snapshot := controlInboxSnapshot{
		PendingCount: len(c.pending),
	}
	if len(c.pending) == 0 {
		return snapshot
	}
	snapshot.Messages = make([]controlInboxMessage, 0, len(c.pending))
	for _, msg := range c.pending {
		snapshot.Messages = append(snapshot.Messages, controlInboxMessage{
			ID:         msg.ID,
			Payload:    msg.Payload,
			ReceivedAt: msg.ReceivedAt,
		})
	}
	return snapshot
}

func (r *Runtime) controlSnapshot() controlInboxSnapshot {
	if r.control == nil {
		return controlInboxSnapshot{}
	}
	r.control.mu.Lock()
	defer r.control.mu.Unlock()
	return r.control.snapshotLocked()
}

func (r *Runtime) writeControlInboxSnapshot(snapshot controlInboxSnapshot) error {
	if err := ensureDirPath(filepath.Dir(r.controlInboxFilePath()), "inbox dir"); err != nil {
		return err
	}
	return writeJSONFile(r.controlInboxFilePath(), snapshot)
}

func (r *Runtime) appendControlLog(entry controlLogEntry) error {
	entry.Timestamp = time.Now().UnixMilli()
	data, err := json.Marshal(entry)
	if err != nil {
		return fmt.Errorf("marshal control log entry: %w", err)
	}
	if err := ensureDirPath(filepath.Dir(r.controlLogFilePath()), "control log dir"); err != nil {
		return err
	}
	f, err := os.OpenFile(r.controlLogFilePath(), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return fmt.Errorf("open control log: %w", err)
	}
	defer f.Close()
	if _, err := f.Write(append(data, '\n')); err != nil {
		return fmt.Errorf("write control log: %w", err)
	}
	return nil
}

func (r *Runtime) enqueueControlPayload(action controlSurfaceAction, payload string) (controlMessage, error) {
	var queued controlMessage
	if r.control == nil {
		return queued, nil
	}
	payload = normalizeControlPayload(payload)
	if strings.TrimSpace(payload) == "" {
		return queued, nil
	}

	var snapshot controlInboxSnapshot

	r.control.mu.Lock()
	r.control.nextID++
	queued = controlMessage{
		ID:         fmt.Sprintf("message-%d", r.control.nextID),
		Payload:    payload,
		ReceivedAt: time.Now().UnixMilli(),
	}
	r.control.pending = append(r.control.pending, queued)
	snapshot = r.control.snapshotLocked()
	r.control.mu.Unlock()

	if err := r.writeControlInboxSnapshot(snapshot); err != nil {
		return queued, err
	}
	msg := controlRuntimeMessage{
		ID:         queued.ID,
		Payload:    queued.Payload,
		ReceivedAt: queued.ReceivedAt,
	}
	if err := r.appendControlLog(controlLogEntry{
		Kind:    "received",
		Action:  string(action),
		Message: &msg,
	}); err != nil {
		return queued, err
	}
	return queued, nil
}

func (r *Runtime) requestControlPoke() {
	if r.control == nil {
		return
	}
	r.control.mu.Lock()
	r.control.pokeRequested = true
	pendingCount := len(r.control.pending)
	r.control.signalEvent()
	r.control.mu.Unlock()
	if err := r.appendControlLog(controlLogEntry{
		Kind:         "woke",
		Delivery:     string(controlDeliveryPoke),
		PendingCount: pendingCount,
	}); err != nil {
		r.log("control log append error: %v", err)
	}
}

func (r *Runtime) requestControlInject() {
	if r.control == nil {
		return
	}
	r.control.mu.Lock()
	r.control.injectRequested = true
	r.control.signalEvent()
	r.control.mu.Unlock()
}

func (r *Runtime) requestControlInterrupt(hasPayload bool) bool {
	if r.control != nil {
		r.control.mu.Lock()
		r.control.interruptRequested = true
		if hasPayload {
			r.control.injectRequested = true
		}
		r.control.signalEvent()
		r.control.mu.Unlock()
	}
	if proc := r.activeProcess.Load(); proc != nil {
		_ = syscall.Kill(-proc.Pid, syscall.SIGINT)
		return true
	}
	return false
}

func (r *Runtime) consumeControlDelivery() (string, int, []controlRuntimeMessage, bool) {
	if r.control == nil {
		return "", 0, nil, false
	}

	var indicator string
	var delivered []controlRuntimeMessage
	var interruptDelivered bool
	var delivery controlDelivery
	var snapshot controlInboxSnapshot

	r.control.mu.Lock()
	switch {
	case r.control.interruptRequested && r.control.injectRequested && len(r.control.pending) > 0:
		delivery = controlDeliveryInterrupt
		delivered = deliverPendingLocked(r.control.pending, delivery)
		r.control.pending = nil
		r.control.interruptRequested = false
		r.control.injectRequested = false
		r.control.pokeRequested = false
		interruptDelivered = true

	case r.control.injectRequested && len(r.control.pending) > 0:
		delivery = controlDeliveryInject
		delivered = deliverPendingLocked(r.control.pending, delivery)
		r.control.pending = nil
		r.control.injectRequested = false
		r.control.pokeRequested = false
		r.control.interruptRequested = false

	case r.control.interruptRequested:
		r.control.interruptRequested = false
		r.control.injectRequested = false
		r.control.pokeRequested = false
		interruptDelivered = true

	case r.control.pokeRequested:
		r.control.pokeRequested = false

	case len(r.control.pending) > 0:
		newlyNotified := false
		for i := range r.control.pending {
			if r.control.pending[i].Notified {
				continue
			}
			r.control.pending[i].Notified = true
			newlyNotified = true
		}
		if newlyNotified {
			indicator = controlInboxIndicator
		}
	}

	snapshot = r.control.snapshotLocked()
	r.control.mu.Unlock()

	if err := r.writeControlInboxSnapshot(snapshot); err != nil {
		r.log("control inbox sync error: %v", err)
	}
	for _, msg := range delivered {
		msgCopy := msg
		if err := r.appendControlLog(controlLogEntry{
			Kind:     "delivered",
			Delivery: string(delivery),
			Message:  &msgCopy,
		}); err != nil {
			r.log("control log append error: %v", err)
		}
	}
	return indicator, snapshot.PendingCount, delivered, interruptDelivered
}

func deliverPendingLocked(pending []controlMessage, delivery controlDelivery) []controlRuntimeMessage {
	delivered := make([]controlRuntimeMessage, 0, len(pending))
	for _, msg := range pending {
		delivered = append(delivered, controlRuntimeMessage{
			ID:         msg.ID,
			Delivery:   string(delivery),
			Payload:    msg.Payload,
			ReceivedAt: msg.ReceivedAt,
		})
	}
	return delivered
}

func (r *Runtime) consumeIdleResume() (controlDelivery, int, bool) {
	if r.control == nil {
		return "", 0, false
	}

	var delivery controlDelivery
	var pendingCount int
	var snapshot controlInboxSnapshot
	var syncSnapshot bool

	r.control.mu.Lock()
	switch {
	case r.control.interruptRequested:
		delivery = controlDeliveryInterrupt
		pendingCount = len(r.control.pending)
		if pendingCount == 0 {
			r.control.interruptRequested = false
			r.control.injectRequested = false
			r.control.pokeRequested = false
			snapshot = r.control.snapshotLocked()
			syncSnapshot = true
		}
	case r.control.injectRequested:
		delivery = controlDeliveryInject
		pendingCount = len(r.control.pending)
		if pendingCount == 0 {
			r.control.injectRequested = false
			r.control.pokeRequested = false
			snapshot = r.control.snapshotLocked()
			syncSnapshot = true
		}
	case r.control.pokeRequested:
		delivery = controlDeliveryPoke
		pendingCount = len(r.control.pending)
		r.control.pokeRequested = false
		if pendingCount == 0 {
			snapshot = r.control.snapshotLocked()
			syncSnapshot = true
		}
	}
	r.control.mu.Unlock()

	if delivery == "" {
		return "", 0, false
	}
	if syncSnapshot {
		if err := r.writeControlInboxSnapshot(snapshot); err != nil {
			r.log("control inbox sync error: %v", err)
		}
	}
	return delivery, pendingCount, true
}

func (r *Runtime) waitForIdleResume() (controlDelivery, int) {
	if r.control == nil {
		return "", 0
	}
	for {
		if delivery, pendingCount, ok := r.consumeIdleResume(); ok {
			return delivery, pendingCount
		}
		<-r.control.eventCh
	}
}

func (r *Runtime) appendControlStatus(msg *tape.Message) bool {
	if msg == nil || msg.Role != tape.RoleToolResult {
		return false
	}
	indicator, pendingCount, delivered, interruptDelivered := r.consumeControlDelivery()
	if indicator == "" && pendingCount == 0 && len(delivered) == 0 && !interruptDelivered {
		return false
	}

	changed := false
	if indicator != "" {
		updated, err := setRuntimeField(msg.StructuredContent, "inbox_indicator", indicator)
		if err != nil {
			r.log("tool result inbox-indicator update error: %v", err)
		} else {
			msg.StructuredContent = updated
			changed = true
		}
	}
	if pendingCount > 0 {
		updated, err := setRuntimeField(msg.StructuredContent, "pending_count", pendingCount)
		if err != nil {
			r.log("tool result pending-count update error: %v", err)
		} else {
			msg.StructuredContent = updated
			changed = true
		}
	}
	if len(delivered) > 0 {
		updated, err := setRuntimeField(msg.StructuredContent, "incoming_messages", delivered)
		if err != nil {
			r.log("tool result incoming-messages update error: %v", err)
		} else {
			msg.StructuredContent = updated
			changed = true
		}
	}
	if interruptDelivered {
		updated, err := setRuntimeField(msg.StructuredContent, "interrupt_notice", "Current operation was interrupted by peer control input.")
		if err != nil {
			r.log("tool result interrupt-notice update error: %v", err)
		} else {
			msg.StructuredContent = updated
			changed = true
		}
	}
	if changed {
		syncToolResultMessageContent(msg)
	}
	return changed
}

// setupControlSurface initializes control state and materializes the retained
// backing files. Peer control writes arrive through the unified FUSE ctl node
// (public/ctl/<action>), which calls applyControlSurfaceAction directly; there
// is no listener goroutine to start or stop.
func (r *Runtime) setupControlSurface() error {
	if r.control == nil {
		r.control = newControlState()
	}
	if err := r.ensureControlSurface(); err != nil {
		return err
	}
	r.control.mu.Lock()
	r.control.started = true
	r.control.mu.Unlock()
	return nil
}
