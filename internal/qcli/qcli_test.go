package qcli

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"syscall"
	"testing"
	"time"
)

func TestFormatClientSignalPayloadGolden(t *testing.T) {
	endpoint := ClientEndpoint{ControlPath: "/tmp/qcli/client/qcli-1/public/ctl"}
	got := FormatClientSignalPayload(endpoint, ControlActionInject, "hello\n")
	want := "[qcli-client]\n" +
		"authority: human\n" +
		"ctl_action: inject\n" +
		"reply_ctl: /tmp/qcli/client/qcli-1/public/ctl/post\n" +
		"reply_required: false\n\n" +
		"message:\n" +
		"hello"
	if got != want {
		t.Fatalf("payload mismatch\n got: %q\nwant: %q", got, want)
	}
}

func TestCtlWriteTimeoutAndCap(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "inject")
	if err := syscall.Mkfifo(path, 0o666); err != nil {
		t.Fatalf("mkfifo: %v", err)
	}
	writer := NewCtlWriter()
	start := time.Now()
	err := writer.Write(path, []byte("payload"), 20*time.Millisecond)
	if !errors.Is(err, ErrWriteTimeout) {
		t.Fatalf("Write err = %v, want ErrWriteTimeout", err)
	}
	if elapsed := time.Since(start); elapsed > time.Second {
		t.Fatalf("timeout took too long: %s", elapsed)
	}
	for i := 0; i < maxAbandonedWriters-1; i++ {
		_ = writer.Write(path, []byte("payload"), time.Millisecond)
	}
	if err := writer.Write(path, []byte("payload"), time.Millisecond); !errors.Is(err, ErrUnreachable) {
		t.Fatalf("cap err = %v, want ErrUnreachable", err)
	}
}

func TestParseTapeLineContractCells(t *testing.T) {
	tests := []struct {
		name string
		line string
		want []Cell
	}{
		{
			name: "meta entry",
			line: `{"type":"meta","data":{"session_id":"sess_abc","model_id":"gpt-5.5","depth":2}}`,
			want: []Cell{{Kind: CellMeta, Text: strptr("gpt-5.5"), Session: strptr("sess_abc"), Depth: intptr(2)}},
		},
		{
			name: "message ts",
			line: `{"type":"message","data":{"role":"assistant","content":"hello there","timestamp":1719300000123}}`,
			want: []Cell{{Kind: CellMessage, Role: strptr("assistant"), Text: strptr("hello there"), TS: int64ptr(1719300000123)}},
		},
		{
			name: "reasoning and tool call order",
			line: `{"type":"message","data":{"role":"assistant","content":"let me check","reasoning_content":"why","tool_calls":[{"id":"c1","name":"bash","arguments":{"cmd":"ls"}}]}}`,
			want: []Cell{
				{Kind: CellReasoning, Role: strptr("assistant"), Text: strptr("why")},
				{Kind: CellMessage, Role: strptr("assistant"), Text: strptr("let me check")},
				{Kind: CellToolCall, Role: strptr("assistant"), ToolName: strptr("bash"), ToolID: strptr("c1"), Text: strptr(`{"cmd":"ls"}`)},
			},
		},
		{
			name: "tool result",
			line: `{"type":"tool_result","data":{"tool_id":"c1","content":{"tool":"bash","status":"completed","output":"3 files"},"is_error":false}}`,
			want: []Cell{{Kind: CellToolResult, ToolID: strptr("c1"), ToolName: strptr("bash"), Status: strptr("completed"), Text: strptr("3 files")}},
		},
		{
			name: "outcome structured",
			line: `{"type":"outcome","data":{"exit_code":0,"turn_count":7,"termination_mode":"exit","tokens_in":10,"tokens_out":3}}`,
			want: []Cell{{Kind: CellOutcome, Text: strptr("exit"), ExitCode: intptr(0), TurnCount: intptr(7), TerminationMode: strptr("exit"), TokensIn: intptr(10), TokensOut: intptr(3)}},
		},
		{
			name: "unknown raw",
			line: `{"type":"telemetry","data":{"x":1}}`,
			want: []Cell{{Kind: "telemetry", Raw: strptr(`{"type":"telemetry","data":{"x":1}}`)}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ParseTapeLine(tt.line)
			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("got %#v\nwant %#v", got, tt.want)
			}
		})
	}
}

func TestParseTapeLineHumanEnvelope(t *testing.T) {
	env := "[qcli-client]\\nauthority: human\\nctl_action: inject\\nreply_ctl: /x\\nreply_required: false\\n\\nmessage:\\nhello agent"
	line := `{"type":"message","data":{"role":"user","content":"` + env + `","timestamp":5}}`
	got := ParseTapeLine(line)
	if len(got) != 1 {
		t.Fatalf("got %d cells, want 1: %#v", len(got), got)
	}
	c := got[0]
	if c.Kind != CellHumanPost || c.Author == nil || *c.Author != humanAuthor || c.Action == nil || *c.Action != "inject" {
		t.Fatalf("cell = %#v", c)
	}
	if c.Text == nil || *c.Text != "hello agent" {
		t.Fatalf("text = %#v", c.Text)
	}
}

func TestParseTapeLineLegacyHumanEnvelopeDefaultsAuthority(t *testing.T) {
	env := "[qcli-client]\\nctl_action: inject\\nreply_ctl: /x\\nreply_required: false\\n\\nmessage:\\nhello legacy"
	line := `{"type":"message","data":{"role":"user","content":"` + env + `"}}`
	got := ParseTapeLine(line)
	if len(got) != 1 {
		t.Fatalf("got %d cells, want 1: %#v", len(got), got)
	}
	c := got[0]
	if c.Kind != CellHumanPost || c.Author == nil || *c.Author != humanAuthor {
		t.Fatalf("cell = %#v", c)
	}
	if c.Text == nil || *c.Text != "hello legacy" {
		t.Fatalf("text = %#v", c.Text)
	}
}

func TestKernelHTTPContractAndSend(t *testing.T) {
	agent := makeTestAgent(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	kernel, err := NewKernel(ctx, KernelOptions{RuntimeRoot: agent.RuntimeRoot, CtlTimeout: time.Second, PollInterval: 20 * time.Millisecond})
	if err != nil {
		t.Fatal(err)
	}
	defer kernel.Close()
	if err := kernel.AttachAgent(agent); err != nil {
		t.Fatal(err)
	}
	server := httptest.NewServer(kernel.Handler())
	defer server.Close()

	var health struct {
		OK       bool   `json:"ok"`
		Contract string `json:"contract"`
		Attached bool   `json:"attached"`
	}
	getJSON(t, server.URL+"/healthz", &health)
	if !health.OK || health.Contract != ContractVersion || !health.Attached {
		t.Fatalf("health = %#v", health)
	}

	frames := make(chan map[string]any, 8)
	sseCtx, sseCancel := context.WithCancel(context.Background())
	defer sseCancel()
	go readSSEFrames(t, sseCtx, server.URL+"/events", frames)
	for _, want := range []string{"hello", "status", "cell", "backfill_complete"} {
		got := waitFrameType(t, frames, want)
		if got["type"] != want {
			t.Fatalf("frame type = %v, want %s", got["type"], want)
		}
	}

	payloadCh := make(chan string, 1)
	go readFIFO(t, filepath.Join(agent.PublicRoot, "ctl", "inject"), payloadCh)
	resp, err := http.Post(server.URL+"/command", "application/json", strings.NewReader(`{"op":"send","action":"inject","text":"hello peer"}`))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		t.Fatalf("send status = %d body=%s", resp.StatusCode, body)
	}
	var queued QueuedResponse
	if err := json.NewDecoder(resp.Body).Decode(&queued); err != nil {
		t.Fatal(err)
	}
	if queued.Type != "queued" || queued.Action != "inject" || queued.ClientRef == "" {
		t.Fatalf("queued = %#v", queued)
	}
	select {
	case payload := <-payloadCh:
		if !strings.Contains(payload, "[qcli-client]\nauthority: human\nctl_action: inject") || !strings.Contains(payload, "message:\nhello peer") {
			t.Fatalf("payload = %q", payload)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for ctl payload")
	}
}

func TestCommandSendErrorPaths(t *testing.T) {
	agent := makeTestAgent(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	kernel, err := NewKernel(ctx, KernelOptions{RuntimeRoot: agent.RuntimeRoot, CtlTimeout: 10 * time.Millisecond})
	if err != nil {
		t.Fatal(err)
	}
	defer kernel.Close()
	server := httptest.NewServer(kernel.Handler())
	defer server.Close()

	resp, err := http.Post(server.URL+"/command", "application/json", strings.NewReader(`{"op":"send","action":"inject","text":"x"}`))
	if err != nil {
		t.Fatal(err)
	}
	var noPeer ErrorEvent
	if err := json.NewDecoder(resp.Body).Decode(&noPeer); err != nil {
		t.Fatal(err)
	}
	_ = resp.Body.Close()
	if resp.StatusCode != http.StatusConflict || noPeer.Code != "no_peer" {
		t.Fatalf("no peer status=%d body=%#v", resp.StatusCode, noPeer)
	}

	if err := kernel.AttachAgent(agent); err != nil {
		t.Fatal(err)
	}
	resp, err = http.Post(server.URL+"/command", "application/json", strings.NewReader(`{"op":"send","action":"inject","text":"   "}`))
	if err != nil {
		t.Fatal(err)
	}
	var empty ErrorEvent
	if err := json.NewDecoder(resp.Body).Decode(&empty); err != nil {
		t.Fatal(err)
	}
	_ = resp.Body.Close()
	if resp.StatusCode != http.StatusBadRequest || empty.Code != "empty_payload" {
		t.Fatalf("empty status=%d body=%#v", resp.StatusCode, empty)
	}

	payloadCh := make(chan string, 1)
	go readFIFO(t, filepath.Join(agent.PublicRoot, "ctl", "poke"), payloadCh)
	resp, err = http.Post(server.URL+"/command", "application/json", strings.NewReader(`{"op":"send","action":"poke","text":""}`))
	if err != nil {
		t.Fatal(err)
	}
	var poke QueuedResponse
	if err := json.NewDecoder(resp.Body).Decode(&poke); err != nil {
		t.Fatal(err)
	}
	_ = resp.Body.Close()
	if resp.StatusCode != http.StatusOK || poke.Type != "queued" || poke.Action != "poke" || poke.ClientRef == "" {
		t.Fatalf("empty poke status=%d body=%#v", resp.StatusCode, poke)
	}
	select {
	case payload := <-payloadCh:
		if payload != "" {
			t.Fatalf("empty poke payload = %q", payload)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for empty poke write")
	}
}

func TestPeerContractEndpoint(t *testing.T) {
	agent := makeTestAgent(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	kernel, err := NewKernel(ctx, KernelOptions{RuntimeRoot: agent.RuntimeRoot, PollInterval: 20 * time.Millisecond})
	if err != nil {
		t.Fatal(err)
	}
	defer kernel.Close()
	if err := kernel.AttachAgent(agent); err != nil {
		t.Fatal(err)
	}
	server := httptest.NewServer(kernel.Handler())
	defer server.Close()

	var got PeerContractResponse
	getJSON(t, server.URL+"/peer-contract", &got)
	if got.Type != "peer_contract" {
		t.Fatalf("type = %q", got.Type)
	}
	m, ok := got.Contract.(map[string]any)
	if !ok || m["contract_version"] != "process-control/v1" {
		t.Fatalf("contract = %#v", got.Contract)
	}
}

func TestReservedAndPreAttachEndpoints(t *testing.T) {
	agent := makeTestAgent(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	kernel, err := NewKernel(ctx, KernelOptions{RuntimeRoot: agent.RuntimeRoot})
	if err != nil {
		t.Fatal(err)
	}
	defer kernel.Close()
	server := httptest.NewServer(kernel.Handler())
	defer server.Close()

	resp, err := http.Get(server.URL + "/events")
	if err != nil {
		t.Fatal(err)
	}
	var noPeer ErrorEvent
	if err := json.NewDecoder(resp.Body).Decode(&noPeer); err != nil {
		t.Fatal(err)
	}
	_ = resp.Body.Close()
	if resp.StatusCode != http.StatusConflict || noPeer.Code != "no_peer" {
		t.Fatalf("/events before attach status=%d body=%#v", resp.StatusCode, noPeer)
	}

	resp, err = http.Get(server.URL + "/history")
	if err != nil {
		t.Fatal(err)
	}
	var reserved ErrorEvent
	if err := json.NewDecoder(resp.Body).Decode(&reserved); err != nil {
		t.Fatal(err)
	}
	_ = resp.Body.Close()
	if resp.StatusCode != http.StatusNotImplemented || reserved.Code != "reserved" {
		t.Fatalf("/history status=%d body=%#v", resp.StatusCode, reserved)
	}
}

func TestRosterProjectionMarksAttached(t *testing.T) {
	agent := makeTestAgent(t)
	roster := ScanRoster(agent.RuntimeRoot, &agent)
	if len(roster) != 1 {
		t.Fatalf("roster len=%d: %#v", len(roster), roster)
	}
	got := roster[0]
	if got.PID != agent.PID || got.Session != agent.Session || !got.Live || !got.Attached {
		t.Fatalf("roster entry = %#v, agent=%#v", got, agent)
	}
}

func TestSessionRestartBackfillIsFixtureDeterministic(t *testing.T) {
	agent := makeTestAgent(t)

	run := func() []Cell {
		t.Helper()
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		session, err := StartSession(ctx, agent, SessionOptions{PollInterval: 20 * time.Millisecond})
		if err != nil {
			t.Fatal(err)
		}
		defer session.Close()
		return session.Cells()
	}

	first := run()
	second := run()
	if !reflect.DeepEqual(first, second) {
		t.Fatalf("restart backfill mismatch\nfirst=%#v\nsecond=%#v", first, second)
	}
}

func TestControlReceiptMatchesClientRef(t *testing.T) {
	agent := makeTestAgent(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	session, err := StartSession(ctx, agent, SessionOptions{PollInterval: 20 * time.Millisecond, CtlTimeout: time.Second})
	if err != nil {
		t.Fatal(err)
	}
	defer session.Close()
	ch, unsub := session.Subscribe()
	defer unsub()
	drainInitial(t, ch)

	payloadCh := make(chan string, 1)
	go readFIFO(t, filepath.Join(agent.PublicRoot, "ctl", "inject"), payloadCh)
	queued, err := session.Send(ControlActionInject, "receipt me")
	if err != nil {
		t.Fatal(err)
	}
	payload := <-payloadCh
	appendControlFixture(t, agent.ControlLogPath, `{"kind":"received","timestamp":1,"action":"inject","message":{"id":"message-1","payload":`+quoteJSON(payload)+`,"received_at":1}}`)
	appendControlFixture(t, agent.ControlLogPath, `{"kind":"delivered","timestamp":2,"delivery":"inject","message":{"id":"message-1","delivery":"inject","payload":`+quoteJSON(payload)+`,"received_at":1}}`)

	received := waitSessionEvent[ReceiptEvent](t, ch, func(ev ReceiptEvent) bool {
		return ev.Stage == "received" && ev.ClientRef != nil && *ev.ClientRef == queued.ClientRef
	})
	if received.MessageID == nil || *received.MessageID != "message-1" {
		t.Fatalf("received = %#v", received)
	}
	delivered := waitSessionEvent[ReceiptEvent](t, ch, func(ev ReceiptEvent) bool {
		return ev.Stage == "delivered" && ev.ClientRef != nil && *ev.ClientRef == queued.ClientRef
	})
	if delivered.Delivery == nil || *delivered.Delivery != "inject" {
		t.Fatalf("delivered = %#v", delivered)
	}
}

func TestLiveTailAndContextReset(t *testing.T) {
	agent := makeTestAgent(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	session, err := StartSession(ctx, agent, SessionOptions{PollInterval: 20 * time.Millisecond})
	if err != nil {
		t.Fatal(err)
	}
	defer session.Close()
	ch, unsub := session.Subscribe()
	defer unsub()
	drainInitial(t, ch)

	livePath := filepath.Join(agent.AgentRoot, "context", "state", "live.jsonl")
	if err := os.WriteFile(livePath, []byte(`{"seq":1,"kind":"text_delta","text":"hi","ts":3}`+"\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	delta := waitSessionEvent[StreamDeltaEvent](t, ch, func(ev StreamDeltaEvent) bool {
		return ev.Kind == "text_delta" && ev.Text != nil && *ev.Text == "hi"
	})
	if delta.Generation != 1 {
		t.Fatalf("generation=%d, want 1", delta.Generation)
	}

	contextPath := filepath.Join(agent.AgentRoot, "context", "state", "current.jsonl")
	if err := os.WriteFile(contextPath, []byte(`{"type":"message","data":{"role":"assistant","content":"after reset","timestamp":4}}`+"\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	reset := waitSessionEvent[ContextResetEvent](t, ch, func(ev ContextResetEvent) bool {
		return ev.Reason == "compaction"
	})
	if reset.Type != "context_reset" {
		t.Fatalf("reset = %#v", reset)
	}
	cell := waitSessionEvent[Cell](t, ch, func(ev Cell) bool {
		return ev.Kind == CellMessage && ev.Text != nil && *ev.Text == "after reset"
	})
	if cell.Seq == 0 {
		t.Fatalf("cell seq not stamped: %#v", cell)
	}
}

func makeTestAgent(t *testing.T) Agent {
	t.Helper()
	root := t.TempDir()
	session := "sess_test"
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
	status := `{"session_id":"sess_test","run_id":"run_test","incarnation_id":2,"pid":` + intToString(os.Getpid()) + `,"ppid":1,"parent_session":"","runtime_root":"` + root + `","agent_root":"` + agentRoot + `","model_id":"gpt-5.5","depth":2}` + "\n"
	for _, path := range []string{filepath.Join(agentRoot, "status", "session.json"), filepath.Join(publicRoot, "status", "session.json")} {
		if err := os.WriteFile(path, []byte(status), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	if err := os.WriteFile(filepath.Join(publicRoot, "status", "inbox.json"), []byte(`{"pending_count":0}`+"\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(publicRoot, "status", "contract.json"), []byte(`{"contract_version":"process-control/v1","control_actions":{"inject":{"description":"queue-and-deliver"}}}`+"\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	contextLine := `{"type":"message","data":{"role":"assistant","content":"hello from context","timestamp":1}}` + "\n"
	if err := os.WriteFile(filepath.Join(agentRoot, "context", "state", "current.jsonl"), []byte(contextLine), 0o644); err != nil {
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
			t.Fatalf("mkfifo %s: %v", name, err)
		}
	}
	if err := os.Symlink(publicRoot, filepath.Join(root, "pid", intToString(os.Getpid()))); err != nil {
		t.Fatal(err)
	}
	agent, err := loadAgentFromRoot(agentRoot)
	if err != nil {
		t.Fatal(err)
	}
	return agent
}

func readFIFO(t *testing.T, path string, out chan<- string) {
	t.Helper()
	f, err := os.OpenFile(path, os.O_RDONLY, os.ModeNamedPipe)
	if err != nil {
		return
	}
	defer f.Close()
	data, _ := io.ReadAll(f)
	out <- string(data)
}

func appendControlFixture(t *testing.T, path, line string) {
	t.Helper()
	f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	if _, err := f.WriteString(line + "\n"); err != nil {
		t.Fatal(err)
	}
}

func quoteJSON(s string) string {
	b, _ := json.Marshal(s)
	return string(b)
}

func drainInitial(t *testing.T, ch <-chan any) {
	t.Helper()
	want := map[string]bool{"hello": false, "status": false, "backfill_complete": false}
	deadline := time.After(2 * time.Second)
	for {
		if allSeen(want) {
			return
		}
		select {
		case ev := <-ch:
			switch e := ev.(type) {
			case HelloEvent:
				want[e.Type] = true
			case StatusEvent:
				want[e.Type] = true
			case BackfillCompleteEvent:
				want[e.Type] = true
			}
		case <-deadline:
			t.Fatalf("timed out draining initial events: %#v", want)
		}
	}
}

func allSeen(m map[string]bool) bool {
	for _, seen := range m {
		if !seen {
			return false
		}
	}
	return true
}

func waitSessionEvent[T any](t *testing.T, ch <-chan any, match func(T) bool) T {
	t.Helper()
	deadline := time.After(2 * time.Second)
	for {
		select {
		case ev := <-ch:
			if typed, ok := ev.(T); ok && match(typed) {
				return typed
			}
		case <-deadline:
			var zero T
			t.Fatalf("timed out waiting for event %T", zero)
		}
	}
}

func readSSEFrames(t *testing.T, ctx context.Context, url string, out chan<- map[string]any) {
	t.Helper()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := scanner.Text()
		if !strings.HasPrefix(line, "data: ") {
			continue
		}
		var frame map[string]any
		if err := json.Unmarshal([]byte(strings.TrimPrefix(line, "data: ")), &frame); err == nil {
			out <- frame
		}
	}
}

func waitFrameType(t *testing.T, frames <-chan map[string]any, want string) map[string]any {
	t.Helper()
	deadline := time.After(2 * time.Second)
	for {
		select {
		case frame := <-frames:
			if frame["type"] == want {
				return frame
			}
		case <-deadline:
			t.Fatalf("timed out waiting for %s", want)
		}
	}
}

func getJSON(t *testing.T, url string, out any) {
	t.Helper()
	resp, err := http.Get(url)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		t.Fatalf("GET %s status=%d body=%s", url, resp.StatusCode, body)
	}
	if err := json.NewDecoder(resp.Body).Decode(out); err != nil {
		t.Fatal(err)
	}
}

func intToString(v int) string {
	return strconv.Itoa(v)
}
