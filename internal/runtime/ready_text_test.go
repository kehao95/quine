package runtime

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/kehao95/quine/internal/tape"
)

func TestReadyTextHeuristicRecognizesCurrentPhraseTables(t *testing.T) {
	var cases []struct {
		name string
		text string
	}
	for _, phrase := range readyExactPhrases {
		cases = append(cases, struct {
			name string
			text string
		}{
			name: fmt.Sprintf("exact/%s", phrase),
			text: phrase,
		})
	}
	for _, prefix := range donePrefixes {
		cases = append(cases, struct {
			name string
			text string
		}{
			name: fmt.Sprintf("done-prefix/%s", prefix),
			text: prefix + " summary text",
		})
	}
	for _, fragment := range readyContinuationFragments {
		cases = append(cases, struct {
			name string
			text string
		}{
			name: fmt.Sprintf("ready-continuation/%s", fragment),
			text: "ready, " + fragment,
		})
	}
	cases = append(cases,
		struct {
			name string
			text string
		}{name: "mixed-case-exact", text: "Done."},
		struct {
			name string
			text string
		}{name: "known-ready-comma-boundary", text: "ready, awaiting instructions"},
		struct {
			name string
			text string
		}{name: "known-done-colon-boundary", text: "done: summary text"},
	)

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if !isReadyLikeTextOnly(tape.Message{Role: tape.RoleAssistant, Content: tc.text}) {
				t.Fatalf("text %q should be treated as quiescent text-only", tc.text)
			}
		})
	}
}

func TestReadyTextHeuristicRejectsCurrentNegatives(t *testing.T) {
	cases := []struct {
		name string
		msg  tape.Message
	}{
		{
			name: "ordinary-text",
			msg:  tape.Message{Role: tape.RoleAssistant, Content: "I inspected the files and found the issue."},
		},
		{
			name: "ready-continuation-too-long",
			msg:  tape.Message{Role: tape.RoleAssistant, Content: "ready, awaiting detailed instructions about the next concrete implementation task"},
		},
		{
			name: "ready-with-unknown-continuation",
			msg:  tape.Message{Role: tape.RoleAssistant, Content: "ready, proceeding with the implementation"},
		},
		{
			name: "tool-call",
			msg: tape.Message{
				Role:    tape.RoleAssistant,
				Content: "ready",
				ToolCalls: []tape.ToolCall{{
					ID:   "call_1",
					Name: "idle",
				}},
			},
		},
		{
			name: "structured-content",
			msg: tape.Message{
				Role:              tape.RoleAssistant,
				Content:           "ready",
				StructuredContent: json.RawMessage(`{"status":"ready"}`),
			},
		},
		{
			name: "reasoning-content",
			msg: tape.Message{
				Role:             tape.RoleAssistant,
				Content:          "ready",
				ReasoningContent: "I am ready.",
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if isReadyLikeTextOnly(tc.msg) {
				t.Fatalf("message should not be treated as ready-like text-only: %#v", tc.msg)
			}
		})
	}
}
