package runtime

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/kehao95/quine/internal/config"
	"github.com/kehao95/quine/internal/tape"
)

func TestShTimeoutOverrideRejectsTimeoutArgument(t *testing.T) {
	sentinel := filepath.Join(t.TempDir(), "should-not-run.txt")
	mock := &mockProvider{
		responses: []tape.Message{
			{
				Role: tape.RoleAssistant,
				ToolCalls: []tape.ToolCall{
					{
						ID:   "call_1",
						Name: "sh",
						Arguments: map[string]any{
							"command": fmt.Sprintf("printf ran > %q", sentinel),
							"timeout": 1,
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
	cfg.ToolGates = config.DefaultToolGates()
	cfg.ToolGates.ShTimeoutOverrideEnabled = false
	rt := NewWithProvider(cfg, mock)
	silenceRuntime(rt)

	exitCode := rt.Run("try forbidden timeout override", "Begin.")
	if exitCode != 0 {
		t.Fatalf("expected exit code 0, got %d", exitCode)
	}
	if _, err := os.Stat(sentinel); !os.IsNotExist(err) {
		t.Fatalf("sh command should not have executed; stat err=%v", err)
	}

	var result *tape.Message
	for _, msg := range rt.tape.Messages() {
		if msg.Role == tape.RoleToolResult && msg.ToolID == "call_1" {
			msg := msg
			result = &msg
			break
		}
	}
	if result == nil {
		t.Fatal("expected rejected sh tool result")
	}
	payload := decodeToolContent(t, result.StructuredContent)
	if got := toolString(t, payload, "tool"); got != "sh" {
		t.Fatalf("tool = %q, want sh", got)
	}
	if got := toolString(t, payload, "status"); got != "rejected" {
		t.Fatalf("status = %q, want rejected", got)
	}
	if got := toolString(t, payload, "error"); !strings.Contains(got, "timeout override is disabled") {
		t.Fatalf("error = %q, want timeout override disabled", got)
	}
}

func TestShStdinRejectsArgumentWhenDisabled(t *testing.T) {
	sentinel := filepath.Join(t.TempDir(), "should-not-run.txt")
	mock := &mockProvider{
		responses: []tape.Message{
			{
				Role: tape.RoleAssistant,
				ToolCalls: []tape.ToolCall{
					{
						ID:   "call_1",
						Name: "sh",
						Arguments: map[string]any{
							"command": fmt.Sprintf("cat > %q", sentinel),
							"stdin":   "ran\n",
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
	cfg.ToolGates = config.DefaultToolGates()
	cfg.ToolGates.ShStdinEnabled = false
	rt := NewWithProvider(cfg, mock)
	silenceRuntime(rt)

	exitCode := rt.Run("try forbidden stdin", "Begin.")
	if exitCode != 0 {
		t.Fatalf("expected exit code 0, got %d", exitCode)
	}
	if _, err := os.Stat(sentinel); !os.IsNotExist(err) {
		t.Fatalf("sh command should not have executed; stat err=%v", err)
	}
	payload := rejectedToolPayload(t, rt, "call_1")
	if got := toolString(t, payload, "error"); !strings.Contains(got, "stdin is disabled") {
		t.Fatalf("error = %q, want stdin disabled", got)
	}
}

func TestShDetachRejectsArgumentWhenDisabled(t *testing.T) {
	sentinel := filepath.Join(t.TempDir(), "should-not-run.txt")
	mock := &mockProvider{
		responses: []tape.Message{
			{
				Role: tape.RoleAssistant,
				ToolCalls: []tape.ToolCall{
					{
						ID:   "call_1",
						Name: "sh",
						Arguments: map[string]any{
							"command": fmt.Sprintf("printf ran > %q", sentinel),
							"detach":  true,
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
	cfg.ToolGates = config.DefaultToolGates()
	cfg.ToolGates.ShDetachEnabled = false
	rt := NewWithProvider(cfg, mock)
	silenceRuntime(rt)

	exitCode := rt.Run("try forbidden detach", "Begin.")
	if exitCode != 0 {
		t.Fatalf("expected exit code 0, got %d", exitCode)
	}
	if _, err := os.Stat(sentinel); !os.IsNotExist(err) {
		t.Fatalf("sh command should not have executed; stat err=%v", err)
	}
	payload := rejectedToolPayload(t, rt, "call_1")
	if got := toolString(t, payload, "error"); !strings.Contains(got, "detach is disabled") {
		t.Fatalf("error = %q, want detach disabled", got)
	}
}

func TestShInteractiveRejectsArgumentWhenDisabled(t *testing.T) {
	sentinel := filepath.Join(t.TempDir(), "should-not-run.txt")
	mock := &mockProvider{
		responses: []tape.Message{
			{
				Role: tape.RoleAssistant,
				ToolCalls: []tape.ToolCall{
					{
						ID:   "call_1",
						Name: "sh",
						Arguments: map[string]any{
							"command":     fmt.Sprintf("printf ran > %q", sentinel),
							"interactive": true,
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
	cfg.ToolGates = config.DefaultToolGates()
	cfg.ToolGates.ShInteractiveEnabled = false
	rt := NewWithProvider(cfg, mock)
	silenceRuntime(rt)

	exitCode := rt.Run("try forbidden interactive mode", "Begin.")
	if exitCode != 0 {
		t.Fatalf("expected exit code 0, got %d", exitCode)
	}
	if _, err := os.Stat(sentinel); !os.IsNotExist(err) {
		t.Fatalf("sh command should not have executed; stat err=%v", err)
	}
	payload := rejectedToolPayload(t, rt, "call_1")
	if got := toolString(t, payload, "error"); !strings.Contains(got, "interactive mode is disabled") {
		t.Fatalf("error = %q, want interactive mode disabled", got)
	}
}

func rejectedToolPayload(t *testing.T, rt *Runtime, toolID string) map[string]any {
	t.Helper()
	for _, msg := range rt.tape.Messages() {
		if msg.Role == tape.RoleToolResult && msg.ToolID == toolID {
			payload := decodeToolContent(t, msg.StructuredContent)
			if got := toolString(t, payload, "tool"); got != "sh" {
				t.Fatalf("tool = %q, want sh", got)
			}
			if got := toolString(t, payload, "status"); got != "rejected" {
				t.Fatalf("status = %q, want rejected", got)
			}
			return payload
		}
	}
	t.Fatalf("expected rejected sh tool result for %s", toolID)
	return nil
}

func TestShStdinParameter(t *testing.T) {
	tmpDir := t.TempDir()
	outFile := tmpDir + "/result.txt"

	mock := &mockProvider{
		responses: []tape.Message{
			{
				Role: tape.RoleAssistant,
				ToolCalls: []tape.ToolCall{
					{
						ID:   "call_1",
						Name: "sh",
						Arguments: map[string]any{
							"command": "cat > " + outFile,
							"stdin":   "hello from llm stdin\n",
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

	exitCode := rt.Run("write a file via stdin", "Begin.")
	if exitCode != 0 {
		t.Errorf("expected exit code 0, got %d", exitCode)
	}

	data, err := os.ReadFile(outFile)
	if err != nil {
		t.Fatalf("output file not created: %v", err)
	}
	if !strings.Contains(string(data), "hello from llm stdin") {
		t.Errorf("file content = %q, expected 'hello from llm stdin'", string(data))
	}
}

// TestShStdinSpecialChars verifies that characters that would break shell
// quoting (quotes, backslashes, dollar signs) pass through unchanged.
func TestShStdinSpecialChars(t *testing.T) {
	tmpDir := t.TempDir()
	outFile := tmpDir + "/special.txt"

	tricky := "key = \"value\"\npath = C:\\data\nprice = $100\n"

	mock := &mockProvider{
		responses: []tape.Message{
			{
				Role: tape.RoleAssistant,
				ToolCalls: []tape.ToolCall{
					{
						ID:   "call_1",
						Name: "sh",
						Arguments: map[string]any{
							"command": "cat > " + outFile,
							"stdin":   tricky,
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

	exitCode := rt.Run("write a file with special chars", "Begin.")
	if exitCode != 0 {
		t.Errorf("expected exit code 0, got %d", exitCode)
	}

	data, err := os.ReadFile(outFile)
	if err != nil {
		t.Fatalf("output file not created: %v", err)
	}
	if string(data) != tricky {
		t.Errorf("content mismatch.\ngot:  %q\nwant: %q", string(data), tricky)
	}
}
