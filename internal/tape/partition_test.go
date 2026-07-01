package tape

import (
	"encoding/json"
	"testing"
)

func TestPartitionOf(t *testing.T) {
	tp := NewTape("sess", "", 0, "model")

	meta := tp.MetaEntry()
	if got := PartitionOf(meta); got != PartitionMeta {
		t.Fatalf("PartitionOf(meta) = %q, want %q", got, PartitionMeta)
	}

	system := MessageEntry(Message{Role: RoleSystem, Content: "sys"})
	if got := PartitionOf(system); got != PartitionSystem {
		t.Fatalf("PartitionOf(system) = %q, want %q", got, PartitionSystem)
	}

	user := MessageEntry(Message{Role: RoleUser, Content: "user"})
	if got := PartitionOf(user); got != PartitionContext {
		t.Fatalf("PartitionOf(user) = %q, want %q", got, PartitionContext)
	}

	tool := ToolResultEntry(ToolResult{ToolID: "t1", Content: json.RawMessage(`{"status":"ok"}`)})
	if got := PartitionOf(tool); got != PartitionContext {
		t.Fatalf("PartitionOf(tool_result) = %q, want %q", got, PartitionContext)
	}
}

func TestContextEntries(t *testing.T) {
	tp := NewTape("sess", "", 0, "model")
	entries := []TapeEntry{
		tp.MetaEntry(),
		MessageEntry(Message{Role: RoleSystem, Content: "sys"}),
		MessageEntry(Message{Role: RoleUser, Content: "user"}),
		ToolResultEntry(ToolResult{ToolID: "t1", Content: json.RawMessage(`{"status":"ok"}`)}),
	}

	filtered := ContextEntries(entries)
	if len(filtered) != 2 {
		t.Fatalf("ContextEntries len = %d, want 2", len(filtered))
	}
	if got := PartitionOf(filtered[0]); got != PartitionContext {
		t.Fatalf("filtered[0] partition = %q, want %q", got, PartitionContext)
	}
	if got := PartitionOf(filtered[1]); got != PartitionContext {
		t.Fatalf("filtered[1] partition = %q, want %q", got, PartitionContext)
	}
}
