package runtime

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"github.com/kehao95/quine/internal/tape"
)

func TestFusePublicControlShellRedirectPostQueuesPayload(t *testing.T) {
	requireRuntimeSurfaceFUSESupport(t)

	cfg := testCfg(t)
	rt := NewWithProvider(cfg, &mockProvider{})
	silenceRuntime(rt)
	useRealPublicSurface(rt)
	rt.originalInput = "shell redirect over fuse"

	if err := rt.bootstrapAgentRoot(); err != nil {
		t.Fatalf("bootstrapAgentRoot failed: %v", err)
	}
	t.Cleanup(func() {
		if err := rt.cleanupAgentRoot(); err != nil {
			t.Fatalf("cleanupAgentRoot failed: %v", err)
		}
	})

	postPath := filepath.Join(cfg.AgentRoot(), "public", "ctl", "post")
	payload := "hello from shell redirect"
	cmd := exec.Command("sh", "-c", "printf '%s\n' \"$1\" > \"$2\"", "sh", payload, postPath)
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("shell redirect write failed: %v\n%s", err, output)
	}

	waitForCondition(t, 2*time.Second, func() bool {
		data, err := os.ReadFile(cfg.InboxPath())
		if err != nil {
			return false
		}
		var snapshot controlInboxSnapshot
		if err := json.Unmarshal(data, &snapshot); err != nil {
			return false
		}
		return snapshot.PendingCount == 1
	}, "fuse shell redirect inbox update")

	inboxData, err := os.ReadFile(cfg.InboxPath())
	if err != nil {
		t.Fatalf("read inbox snapshot: %v", err)
	}
	var inbox controlInboxSnapshot
	if err := json.Unmarshal(inboxData, &inbox); err != nil {
		t.Fatalf("unmarshal inbox snapshot: %v", err)
	}
	if len(inbox.Messages) != 1 || inbox.Messages[0].Payload != payload {
		t.Fatalf("inbox messages = %#v, want payload %q", inbox.Messages, payload)
	}
}

func TestApplyControlSurfaceAction(t *testing.T) {
	rt := &Runtime{
		cfg:      testCfg(t),
		control:  newControlState(),
		log:      func(string, ...any) {},
		logError: func(string, ...any) {},
	}
	silenceRuntime(rt)
	if err := rt.ensureControlSurface(); err != nil {
		t.Fatalf("ensureControlSurface() error: %v", err)
	}

	tests := []struct {
		name              string
		action            controlSurfaceAction
		payload           string
		wantPending       int
		wantPoke          bool
		wantInject        bool
		wantInterrupt     bool
		wantStoredPayload string
		wantErr           bool
	}{
		{
			name:              "post queues payload only",
			action:            controlActionPost,
			payload:           "hello\n",
			wantPending:       1,
			wantStoredPayload: "hello",
		},
		{
			name:              "poke queues payload and requests resume",
			action:            controlActionPoke,
			payload:           "poke payload\n",
			wantPending:       1,
			wantPoke:          true,
			wantStoredPayload: "poke payload",
		},
		{
			name:          "empty poke only requests resume",
			action:        controlActionPoke,
			payload:       "",
			wantPoke:      true,
			wantPending:   0,
			wantInterrupt: false,
		},
		{
			name:              "inject queues payload and requests delivery",
			action:            controlActionInject,
			payload:           "inject payload\n",
			wantPending:       1,
			wantInject:        true,
			wantStoredPayload: "inject payload",
		},
		{
			name:          "empty interrupt requests interrupt",
			action:        controlActionInterrupt,
			payload:       "",
			wantInterrupt: true,
			wantPending:   0,
		},
		{
			name:    "unknown action errors",
			action:  controlSurfaceAction("bogus"),
			payload: "hello",
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			rt.control = newControlState()
			if err := rt.writeControlInboxSnapshot(controlInboxSnapshot{}); err != nil {
				t.Fatalf("reset inbox snapshot: %v", err)
			}

			err := rt.applyControlSurfaceAction(tc.action, tc.payload)
			if tc.wantErr {
				if err == nil {
					t.Fatal("expected apply error")
				}
				return
			}
			if err != nil {
				t.Fatalf("applyControlSurfaceAction() error: %v", err)
			}

			rt.control.mu.Lock()
			defer rt.control.mu.Unlock()
			if len(rt.control.pending) != tc.wantPending {
				t.Fatalf("pending = %d, want %d", len(rt.control.pending), tc.wantPending)
			}
			if rt.control.pokeRequested != tc.wantPoke {
				t.Fatalf("pokeRequested = %v, want %v", rt.control.pokeRequested, tc.wantPoke)
			}
			if rt.control.injectRequested != tc.wantInject {
				t.Fatalf("injectRequested = %v, want %v", rt.control.injectRequested, tc.wantInject)
			}
			if rt.control.interruptRequested != tc.wantInterrupt {
				t.Fatalf("interruptRequested = %v, want %v", rt.control.interruptRequested, tc.wantInterrupt)
			}
			if tc.wantStoredPayload != "" {
				if got := rt.control.pending[0].Payload; got != tc.wantStoredPayload {
					t.Fatalf("stored payload = %q, want %q", got, tc.wantStoredPayload)
				}
			}
		})
	}
}

func TestFusePublicControlInjectDelivery(t *testing.T) {
	requireRuntimeSurfaceFUSESupport(t)

	mock := &mockProvider{
		responses: []tape.Message{
			{
				Role: tape.RoleAssistant,
				ToolCalls: []tape.ToolCall{
					{
						ID:   "call_1",
						Name: "sh",
						Arguments: map[string]any{
							"command": "sleep 1",
						},
					},
				},
			},
			{
				Role: tape.RoleAssistant,
				ToolCalls: []tape.ToolCall{
					{
						ID:   "call_2",
						Name: "exit",
						Arguments: map[string]any{
							"status": "success",
						},
					},
				},
			},
		},
	}

	cfg := testCfg(t)
	rt := NewWithProvider(cfg, mock)
	silenceRuntime(rt)
	useRealPublicSurface(rt)
	retainedRoot := cfg.SessionLogDir("")

	done := make(chan int, 1)
	go func() {
		done <- rt.Run("peer inject delivery over public ctl", "Begin.")
	}()

	waitForCondition(t, 2*time.Second, func() bool { return rt.activeProcess.Load() != nil }, "active shell")
	writeControlActionFile(t, filepath.Join(cfg.AgentRoot(), "public", "ctl", "inject"), "inject payload\n")

	if exitCode := <-done; exitCode != 0 {
		t.Fatalf("expected exit code 0, got %d", exitCode)
	}

	var shResult *tape.Message
	for _, msg := range rt.tape.Messages() {
		if msg.Role == tape.RoleToolResult {
			copyMsg := msg
			shResult = &copyMsg
			break
		}
	}
	if shResult == nil {
		t.Fatal("expected sh tool_result on tape")
	}

	payload := decodeToolContent(t, shResult.StructuredContent)
	runtimePayload := toolMap(t, payload, "runtime")
	rawIncoming, ok := runtimePayload["incoming_messages"].([]any)
	if !ok || len(rawIncoming) != 1 {
		t.Fatalf("incoming_messages = %#v, want single delivered message", runtimePayload["incoming_messages"])
	}
	incoming, ok := rawIncoming[0].(map[string]any)
	if !ok {
		t.Fatalf("incoming_messages[0] = %#v, want map", rawIncoming[0])
	}
	if toolString(t, incoming, "delivery") != string(controlDeliveryInject) {
		t.Fatalf("delivery = %q, want %q", toolString(t, incoming, "delivery"), controlDeliveryInject)
	}
	if toolString(t, incoming, "payload") != "inject payload" {
		t.Fatalf("payload = %q, want %q", toolString(t, incoming, "payload"), "inject payload")
	}

	inboxData, err := os.ReadFile(filepath.Join(retainedRoot, "status", "inbox.json"))
	if err != nil {
		t.Fatalf("read inbox snapshot: %v", err)
	}
	var inbox map[string]any
	if err := json.Unmarshal(inboxData, &inbox); err != nil {
		t.Fatalf("unmarshal inbox snapshot: %v", err)
	}
	if pending := toolInt(t, inbox, "pending_count"); pending != 0 {
		t.Fatalf("inbox pending_count = %d, want 0", pending)
	}
}

func TestFusePublicControlInterruptDelivery(t *testing.T) {
	requireRuntimeSurfaceFUSESupport(t)

	mock := &mockProvider{
		responses: []tape.Message{
			{
				Role: tape.RoleAssistant,
				ToolCalls: []tape.ToolCall{
					{
						ID:   "call_1",
						Name: "sh",
						Arguments: map[string]any{
							"command": "python3 -c 'import signal, sys, time; signal.signal(signal.SIGINT, lambda signum, frame: sys.exit(130)); time.sleep(5)'",
						},
					},
				},
			},
			{
				Role: tape.RoleAssistant,
				ToolCalls: []tape.ToolCall{
					{
						ID:   "call_2",
						Name: "exit",
						Arguments: map[string]any{
							"status": "success",
						},
					},
				},
			},
		},
	}

	cfg := testCfg(t)
	rt := NewWithProvider(cfg, mock)
	silenceRuntime(rt)
	useRealPublicSurface(rt)

	done := make(chan int, 1)
	go func() {
		done <- rt.Run("peer interrupt delivery over public ctl", "Begin.")
	}()

	waitForCondition(t, 2*time.Second, func() bool { return rt.activeProcess.Load() != nil }, "active shell")
	time.Sleep(150 * time.Millisecond)
	start := time.Now()
	writeControlActionFile(t, filepath.Join(cfg.AgentRoot(), "public", "ctl", "interrupt"), "interrupt payload\n")

	if exitCode := <-done; exitCode != 0 {
		t.Fatalf("expected exit code 0, got %d", exitCode)
	}
	if elapsed := time.Since(start); elapsed >= 2*time.Second {
		t.Fatalf("interrupt delivery took too long: %s", elapsed)
	}

	var shResult *tape.Message
	for _, msg := range rt.tape.Messages() {
		if msg.Role == tape.RoleToolResult {
			copyMsg := msg
			shResult = &copyMsg
			break
		}
	}
	if shResult == nil {
		t.Fatal("expected sh tool_result on tape")
	}

	payload := decodeToolContent(t, shResult.StructuredContent)
	if exitCode := toolInt(t, payload, "exit_code"); exitCode == 0 {
		t.Fatalf("expected interrupted tool exit_code to be non-zero, got %d", exitCode)
	}
	runtimePayload := toolMap(t, payload, "runtime")
	if got := toolString(t, runtimePayload, "interrupt_notice"); got != "Current operation was interrupted by peer control input." {
		t.Fatalf("interrupt_notice = %q", got)
	}
	rawIncoming, ok := runtimePayload["incoming_messages"].([]any)
	if !ok || len(rawIncoming) != 1 {
		t.Fatalf("incoming_messages = %#v, want single delivered message", runtimePayload["incoming_messages"])
	}
	incoming, ok := rawIncoming[0].(map[string]any)
	if !ok {
		t.Fatalf("incoming_messages[0] = %#v, want map", rawIncoming[0])
	}
	if toolString(t, incoming, "delivery") != string(controlDeliveryInterrupt) {
		t.Fatalf("delivery = %q, want %q", toolString(t, incoming, "delivery"), controlDeliveryInterrupt)
	}
	if toolString(t, incoming, "payload") != "interrupt payload" {
		t.Fatalf("payload = %q, want %q", toolString(t, incoming, "payload"), "interrupt payload")
	}
}

func TestFusePublicControlEmptyPokeRequestsResume(t *testing.T) {
	requireRuntimeSurfaceFUSESupport(t)

	cfg := testCfg(t)
	rt := NewWithProvider(cfg, &mockProvider{})
	silenceRuntime(rt)
	useRealPublicSurface(rt)
	rt.originalInput = "fuse empty poke"

	if err := rt.bootstrapAgentRoot(); err != nil {
		t.Fatalf("bootstrapAgentRoot failed: %v", err)
	}
	t.Cleanup(func() {
		if err := rt.cleanupAgentRoot(); err != nil {
			t.Fatalf("cleanupAgentRoot failed: %v", err)
		}
	})

	writeControlActionFile(t, filepath.Join(cfg.AgentRoot(), "public", "ctl", "poke"), "")

	waitForCondition(t, 2*time.Second, func() bool {
		rt.control.mu.Lock()
		defer rt.control.mu.Unlock()
		return rt.control.pokeRequested
	}, "empty fuse poke request")

	rt.control.mu.Lock()
	if got := len(rt.control.pending); got != 0 {
		rt.control.mu.Unlock()
		t.Fatalf("pending queue length = %d, want 0", got)
	}
	rt.control.mu.Unlock()
}

func TestFusePublicControlEmptyInterruptRequestsInterrupt(t *testing.T) {
	requireRuntimeSurfaceFUSESupport(t)

	cfg := testCfg(t)
	rt := NewWithProvider(cfg, &mockProvider{})
	silenceRuntime(rt)
	useRealPublicSurface(rt)
	rt.originalInput = "fuse empty interrupt"

	if err := rt.bootstrapAgentRoot(); err != nil {
		t.Fatalf("bootstrapAgentRoot failed: %v", err)
	}
	t.Cleanup(func() {
		if err := rt.cleanupAgentRoot(); err != nil {
			t.Fatalf("cleanupAgentRoot failed: %v", err)
		}
	})

	writeControlActionFile(t, filepath.Join(cfg.AgentRoot(), "public", "ctl", "interrupt"), "")

	waitForCondition(t, 2*time.Second, func() bool {
		rt.control.mu.Lock()
		defer rt.control.mu.Unlock()
		return rt.control.interruptRequested
	}, "empty fuse interrupt request")

	rt.control.mu.Lock()
	if got := len(rt.control.pending); got != 0 {
		rt.control.mu.Unlock()
		t.Fatalf("pending queue length = %d, want 0", got)
	}
	rt.control.mu.Unlock()
}

func TestControlInboxIndicator(t *testing.T) {
	mock := &mockProvider{
		responses: []tape.Message{
			{
				Role: tape.RoleAssistant,
				ToolCalls: []tape.ToolCall{
					{
						ID:   "call_1",
						Name: "sh",
						Arguments: map[string]any{
							"command": "sleep 1",
						},
					},
				},
			},
			{
				Role: tape.RoleAssistant,
				ToolCalls: []tape.ToolCall{
					{
						ID:   "call_2",
						Name: "exit",
						Arguments: map[string]any{
							"status": "success",
						},
					},
				},
			},
		},
	}

	cfg := testCfg(t)
	rt := NewWithProvider(cfg, mock)
	silenceRuntime(rt)
	retainedRoot := cfg.SessionLogDir("")

	done := make(chan int, 1)
	go func() {
		done <- rt.Run("peer inbox indicator", "Begin.")
	}()

	waitForCondition(t, 2*time.Second, func() bool { return rt.activeProcess.Load() != nil }, "active shell")
	if _, err := rt.enqueueControlPayload(controlActionPost, "hello from peer"); err != nil {
		t.Fatalf("enqueue control payload: %v", err)
	}

	if exitCode := <-done; exitCode != 0 {
		t.Fatalf("expected exit code 0, got %d", exitCode)
	}

	var shResult *tape.Message
	for _, msg := range rt.tape.Messages() {
		if msg.Role == tape.RoleToolResult {
			copyMsg := msg
			shResult = &copyMsg
			break
		}
	}
	if shResult == nil {
		t.Fatal("expected sh tool_result on tape")
	}

	payload := decodeToolContent(t, shResult.StructuredContent)
	runtimePayload := toolMap(t, payload, "runtime")
	if got := toolString(t, runtimePayload, "inbox_indicator"); got != controlInboxIndicator {
		t.Fatalf("inbox_indicator = %q, want %q", got, controlInboxIndicator)
	}
	if pending := toolInt(t, runtimePayload, "pending_count"); pending != 1 {
		t.Fatalf("pending_count = %d, want 1", pending)
	}

	inboxData, err := os.ReadFile(filepath.Join(retainedRoot, "status", "inbox.json"))
	if err != nil {
		t.Fatalf("read inbox snapshot: %v", err)
	}
	var inbox struct {
		PendingCount int `json:"pending_count"`
		Messages     []struct {
			Payload string `json:"payload"`
		} `json:"messages"`
	}
	if err := json.Unmarshal(inboxData, &inbox); err != nil {
		t.Fatalf("unmarshal inbox snapshot: %v", err)
	}
	if inbox.PendingCount != 1 || len(inbox.Messages) != 1 || inbox.Messages[0].Payload != "hello from peer" {
		t.Fatalf("unexpected inbox snapshot: %+v", inbox)
	}
}

func TestControlInjectDelivery(t *testing.T) {
	mock := &mockProvider{
		responses: []tape.Message{
			{
				Role: tape.RoleAssistant,
				ToolCalls: []tape.ToolCall{
					{
						ID:   "call_1",
						Name: "sh",
						Arguments: map[string]any{
							"command": "sleep 1",
						},
					},
				},
			},
			{
				Role: tape.RoleAssistant,
				ToolCalls: []tape.ToolCall{
					{
						ID:   "call_2",
						Name: "exit",
						Arguments: map[string]any{
							"status": "success",
						},
					},
				},
			},
		},
	}

	cfg := testCfg(t)
	rt := NewWithProvider(cfg, mock)
	silenceRuntime(rt)
	retainedRoot := cfg.SessionLogDir("")

	done := make(chan int, 1)
	go func() {
		done <- rt.Run("peer inject delivery", "Begin.")
	}()

	waitForCondition(t, 2*time.Second, func() bool { return rt.activeProcess.Load() != nil }, "active shell")
	if _, err := rt.enqueueControlPayload(controlActionInject, "inject payload"); err != nil {
		t.Fatalf("enqueue control payload: %v", err)
	}
	rt.requestControlInject()

	if exitCode := <-done; exitCode != 0 {
		t.Fatalf("expected exit code 0, got %d", exitCode)
	}

	var shResult *tape.Message
	for _, msg := range rt.tape.Messages() {
		if msg.Role == tape.RoleToolResult {
			copyMsg := msg
			shResult = &copyMsg
			break
		}
	}
	if shResult == nil {
		t.Fatal("expected sh tool_result on tape")
	}

	payload := decodeToolContent(t, shResult.StructuredContent)
	runtimePayload := toolMap(t, payload, "runtime")
	rawIncoming, ok := runtimePayload["incoming_messages"].([]any)
	if !ok || len(rawIncoming) != 1 {
		t.Fatalf("incoming_messages = %#v, want single delivered message", runtimePayload["incoming_messages"])
	}
	incoming, ok := rawIncoming[0].(map[string]any)
	if !ok {
		t.Fatalf("incoming_messages[0] = %#v, want map", rawIncoming[0])
	}
	if toolString(t, incoming, "delivery") != string(controlDeliveryInject) {
		t.Fatalf("delivery = %q, want %q", toolString(t, incoming, "delivery"), controlDeliveryInject)
	}
	if toolString(t, incoming, "payload") != "inject payload" {
		t.Fatalf("payload = %q, want %q", toolString(t, incoming, "payload"), "inject payload")
	}

	inboxData, err := os.ReadFile(filepath.Join(retainedRoot, "status", "inbox.json"))
	if err != nil {
		t.Fatalf("read inbox snapshot: %v", err)
	}
	var inbox map[string]any
	if err := json.Unmarshal(inboxData, &inbox); err != nil {
		t.Fatalf("unmarshal inbox snapshot: %v", err)
	}
	if pending := toolInt(t, inbox, "pending_count"); pending != 0 {
		t.Fatalf("inbox pending_count = %d, want 0", pending)
	}
}

func TestControlInterruptDelivery(t *testing.T) {
	mock := &mockProvider{
		responses: []tape.Message{
			{
				Role: tape.RoleAssistant,
				ToolCalls: []tape.ToolCall{
					{
						ID:   "call_1",
						Name: "sh",
						Arguments: map[string]any{
							"command": "python3 -c 'import signal, sys, time; signal.signal(signal.SIGINT, lambda signum, frame: sys.exit(130)); time.sleep(5)'",
						},
					},
				},
			},
			{
				Role: tape.RoleAssistant,
				ToolCalls: []tape.ToolCall{
					{
						ID:   "call_2",
						Name: "exit",
						Arguments: map[string]any{
							"status": "success",
						},
					},
				},
			},
		},
	}

	cfg := testCfg(t)
	rt := NewWithProvider(cfg, mock)
	silenceRuntime(rt)

	done := make(chan int, 1)
	go func() {
		done <- rt.Run("peer interrupt delivery", "Begin.")
	}()

	waitForCondition(t, 2*time.Second, func() bool { return rt.activeProcess.Load() != nil }, "active shell")
	time.Sleep(150 * time.Millisecond)
	start := time.Now()
	if _, err := rt.enqueueControlPayload(controlActionInterrupt, "interrupt payload"); err != nil {
		t.Fatalf("enqueue control payload: %v", err)
	}
	if forwarded := rt.requestControlInterrupt(true); !forwarded {
		t.Fatal("expected interrupt request to forward to active process")
	}

	if exitCode := <-done; exitCode != 0 {
		t.Fatalf("expected exit code 0, got %d", exitCode)
	}
	if elapsed := time.Since(start); elapsed >= 2*time.Second {
		t.Fatalf("interrupt delivery took too long: %s", elapsed)
	}

	var shResult *tape.Message
	for _, msg := range rt.tape.Messages() {
		if msg.Role == tape.RoleToolResult {
			copyMsg := msg
			shResult = &copyMsg
			break
		}
	}
	if shResult == nil {
		t.Fatal("expected sh tool_result on tape")
	}

	payload := decodeToolContent(t, shResult.StructuredContent)
	if exitCode := toolInt(t, payload, "exit_code"); exitCode == 0 {
		t.Fatalf("expected interrupted tool exit_code to be non-zero, got %d", exitCode)
	}
	runtimePayload := toolMap(t, payload, "runtime")
	if got := toolString(t, runtimePayload, "interrupt_notice"); got != "Current operation was interrupted by peer control input." {
		t.Fatalf("interrupt_notice = %q", got)
	}
	rawIncoming, ok := runtimePayload["incoming_messages"].([]any)
	if !ok || len(rawIncoming) != 1 {
		t.Fatalf("incoming_messages = %#v, want single delivered message", runtimePayload["incoming_messages"])
	}
	incoming, ok := rawIncoming[0].(map[string]any)
	if !ok {
		t.Fatalf("incoming_messages[0] = %#v, want map", rawIncoming[0])
	}
	if toolString(t, incoming, "delivery") != string(controlDeliveryInterrupt) {
		t.Fatalf("delivery = %q, want %q", toolString(t, incoming, "delivery"), controlDeliveryInterrupt)
	}
	if toolString(t, incoming, "payload") != "interrupt payload" {
		t.Fatalf("payload = %q, want %q", toolString(t, incoming, "payload"), "interrupt payload")
	}
}

func TestIdleResumesOnPokeWithoutContextInjection(t *testing.T) {
	mock := &mockProvider{
		responses: []tape.Message{
			{
				Role: tape.RoleAssistant,
				ToolCalls: []tape.ToolCall{
					{
						ID:        "call_idle",
						Name:      "idle",
						Arguments: map[string]any{},
					},
				},
			},
			{
				Role: tape.RoleAssistant,
				ToolCalls: []tape.ToolCall{
					{
						ID:   "call_exit",
						Name: "exit",
						Arguments: map[string]any{
							"status": "success",
						},
					},
				},
			},
		},
	}

	cfg := testCfg(t)
	cfg.IdleEnabled = true
	rt := NewWithProvider(cfg, mock)
	silenceRuntime(rt)

	done := make(chan int, 1)
	go func() {
		done <- rt.Run("idle poke delivery", "Begin.")
	}()

	waitForCondition(t, 2*time.Second, func() bool {
		return rt.control != nil && rt.control.started
	}, "control listener start")

	if _, err := rt.enqueueControlPayload(controlActionPoke, "poke from idle"); err != nil {
		t.Fatalf("enqueue control payload: %v", err)
	}
	select {
	case exitCode := <-done:
		t.Fatalf("idle resumed without poke signal, exit=%d", exitCode)
	case <-time.After(200 * time.Millisecond):
	}

	rt.requestControlPoke()
	if exitCode := <-done; exitCode != 0 {
		t.Fatalf("expected exit code 0, got %d", exitCode)
	}

	var idleResult *tape.Message
	for _, msg := range rt.tape.Messages() {
		if msg.Role == tape.RoleToolResult && msg.ToolID == "call_idle" {
			copyMsg := msg
			idleResult = &copyMsg
			break
		}
	}
	if idleResult == nil {
		t.Fatal("expected idle tool_result on tape")
	}

	payload := decodeToolContent(t, idleResult.StructuredContent)
	if got := toolString(t, payload, "delivery"); got != string(controlDeliveryPoke) {
		t.Fatalf("idle delivery = %q, want %q", got, controlDeliveryPoke)
	}
	if pending := toolInt(t, payload, "pending_count"); pending != 1 {
		t.Fatalf("idle pending_count = %d, want 1", pending)
	}
	if runtimePayload, _ := payload["runtime"].(map[string]any); runtimePayload != nil {
		if _, ok := runtimePayload["incoming_messages"]; ok {
			t.Fatalf("poke should not surface incoming_messages: %#v", runtimePayload)
		}
	}
}

func TestIdleInterruptDelivery(t *testing.T) {
	mock := &mockProvider{
		responses: []tape.Message{
			{
				Role: tape.RoleAssistant,
				ToolCalls: []tape.ToolCall{
					{
						ID:        "call_idle",
						Name:      "idle",
						Arguments: map[string]any{},
					},
				},
			},
			{
				Role: tape.RoleAssistant,
				ToolCalls: []tape.ToolCall{
					{
						ID:   "call_exit",
						Name: "exit",
						Arguments: map[string]any{
							"status": "success",
						},
					},
				},
			},
		},
	}

	cfg := testCfg(t)
	cfg.IdleEnabled = true
	rt := NewWithProvider(cfg, mock)
	silenceRuntime(rt)

	done := make(chan int, 1)
	go func() {
		done <- rt.Run("idle interrupt delivery", "Begin.")
	}()

	waitForCondition(t, 2*time.Second, func() bool {
		return rt.control != nil && rt.control.started
	}, "control listener start")
	if _, err := rt.enqueueControlPayload(controlActionInterrupt, "interrupt from idle"); err != nil {
		t.Fatalf("enqueue control payload: %v", err)
	}
	if forwarded := rt.requestControlInterrupt(true); forwarded {
		t.Fatal("idle interrupt should not forward to an active process")
	}

	if exitCode := <-done; exitCode != 0 {
		t.Fatalf("expected exit code 0, got %d", exitCode)
	}

	var idleResult *tape.Message
	for _, msg := range rt.tape.Messages() {
		if msg.Role == tape.RoleToolResult && msg.ToolID == "call_idle" {
			copyMsg := msg
			idleResult = &copyMsg
			break
		}
	}
	if idleResult == nil {
		t.Fatal("expected idle tool_result on tape")
	}

	payload := decodeToolContent(t, idleResult.StructuredContent)
	if got := toolString(t, payload, "delivery"); got != string(controlDeliveryInterrupt) {
		t.Fatalf("idle delivery = %q, want %q", got, controlDeliveryInterrupt)
	}
	runtimePayload := toolMap(t, payload, "runtime")
	if got := toolString(t, runtimePayload, "interrupt_notice"); got != "Current operation was interrupted by peer control input." {
		t.Fatalf("interrupt_notice = %q", got)
	}
	rawIncoming, ok := runtimePayload["incoming_messages"].([]any)
	if !ok || len(rawIncoming) != 1 {
		t.Fatalf("incoming_messages = %#v, want single delivered message", runtimePayload["incoming_messages"])
	}
	incoming, ok := rawIncoming[0].(map[string]any)
	if !ok {
		t.Fatalf("incoming_messages[0] = %#v, want map", rawIncoming[0])
	}
	if toolString(t, incoming, "delivery") != string(controlDeliveryInterrupt) {
		t.Fatalf("delivery = %q, want %q", toolString(t, incoming, "delivery"), controlDeliveryInterrupt)
	}
	if toolString(t, incoming, "payload") != "interrupt from idle" {
		t.Fatalf("payload = %q, want %q", toolString(t, incoming, "payload"), "interrupt from idle")
	}
}

func TestIdleResumesOnEmptyPoke(t *testing.T) {
	mock := &mockProvider{
		responses: []tape.Message{
			{
				Role: tape.RoleAssistant,
				ToolCalls: []tape.ToolCall{
					{
						ID:        "call_idle",
						Name:      "idle",
						Arguments: map[string]any{},
					},
				},
			},
			{
				Role: tape.RoleAssistant,
				ToolCalls: []tape.ToolCall{
					{
						ID:   "call_exit",
						Name: "exit",
						Arguments: map[string]any{
							"status": "success",
						},
					},
				},
			},
		},
	}

	cfg := testCfg(t)
	cfg.IdleEnabled = true
	rt := NewWithProvider(cfg, mock)
	silenceRuntime(rt)

	done := make(chan int, 1)
	go func() {
		done <- rt.Run("idle signal-only resume", "Begin.")
	}()

	waitForCondition(t, 2*time.Second, func() bool {
		return rt.control != nil && rt.control.started
	}, "control listener start")
	rt.requestControlPoke()

	if exitCode := <-done; exitCode != 0 {
		t.Fatalf("expected exit code 0, got %d", exitCode)
	}

	var idleResult *tape.Message
	for _, msg := range rt.tape.Messages() {
		if msg.Role == tape.RoleToolResult && msg.ToolID == "call_idle" {
			copyMsg := msg
			idleResult = &copyMsg
			break
		}
	}
	if idleResult == nil {
		t.Fatal("expected idle tool_result on tape")
	}

	payload := decodeToolContent(t, idleResult.StructuredContent)
	if got := toolString(t, payload, "delivery"); got != string(controlDeliveryPoke) {
		t.Fatalf("idle delivery = %q, want %q", got, controlDeliveryPoke)
	}
	if _, ok := payload["pending_count"]; ok {
		t.Fatalf("idle signal-only payload should not expose pending_count: %#v", payload)
	}
	runtimePayload, _ := payload["runtime"].(map[string]any)
	if runtimePayload != nil {
		if _, ok := runtimePayload["incoming_messages"]; ok {
			t.Fatalf("idle signal-only runtime should not expose incoming_messages: %#v", runtimePayload)
		}
	}
}
