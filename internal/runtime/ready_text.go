package runtime

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/kehao95/quine/internal/tape"
)

// This file owns experimental model-behavior heuristics behind
// QUINE_READY_TEXT_AUTO_IDLE and QUINE_EMPTY_ASSISTANT_SUCCESS. These phrase
// tables are tuning surfaces, not runtime physics, and may need retuning per
// model lane.

var readyExactPhrases = []string{
	"ready",
	"ready.",
	"ready!",
	"ready?",
	"done",
	"done.",
	"completed",
	"completed.",
	"standing by",
	"standing by.",
	"ack received",
	"ack received.",
}

var donePrefixes = []string{
	"done:",
	"done.",
	"completed:",
	"completed.",
}

var readyContinuationFragments = []string{
	"awaiting",
	"please provide",
	"provide the task",
	"what would you like",
	"how can i help",
}

func isReadyLikeTextOnly(msg tape.Message) bool {
	if len(msg.ToolCalls) != 0 ||
		len(bytes.TrimSpace(msg.StructuredContent)) != 0 ||
		strings.TrimSpace(msg.ReasoningContent) != "" ||
		len(msg.ReasoningItems) != 0 ||
		msg.Image != nil {
		return false
	}
	text := strings.ToLower(strings.TrimSpace(msg.Content))
	text = strings.Trim(text, " \t\r\n")
	for _, phrase := range readyExactPhrases {
		if text == phrase {
			return true
		}
	}
	for _, prefix := range donePrefixes {
		if strings.HasPrefix(text, prefix) {
			return true
		}
	}
	if !strings.HasPrefix(text, "ready.") && !strings.HasPrefix(text, "ready,") {
		return false
	}
	if len(strings.Fields(text)) > 8 {
		return false
	}
	for _, fragment := range readyContinuationFragments {
		if strings.Contains(text, fragment) {
			return true
		}
	}
	return false
}

func (r *Runtime) canTreatEmptyAssistantAsSuccess() bool {
	if r == nil || r.cfg == nil || !r.cfg.EmptyAssistantSuccess {
		return false
	}
	if r.tape != nil && r.tape.TurnCount > 0 {
		return true
	}
	if r.sh != nil {
		revision := strings.TrimSpace(r.sh.CurrentWorldRevision())
		return revision != "" && revision != "wr0"
	}
	return false
}

func (r *Runtime) finishEmptyAssistantSuccess() int {
	duration := time.Since(r.startTime)
	totalTokens := 0
	turnCount := 0
	if r.tape != nil {
		totalTokens = r.tape.TokensIn + r.tape.TokensOut
		turnCount = r.tape.TurnCount
	}
	r.log("empty assistant response after progress treated as success (%d turns, %.1fs, %d tokens)",
		turnCount, duration.Seconds(), totalTokens)
	return r.finalizeOutcome(0, "", tape.TermExit)
}

func (r *Runtime) shouldReadyTextAutoIdle(msg tape.Message) bool {
	return r != nil &&
		r.cfg != nil &&
		r.cfg.ReadyTextAutoIdle &&
		r.cfg.IdleToolEnabled() &&
		r.control != nil &&
		isReadyLikeTextOnly(msg)
}

func (r *Runtime) handleReadyTextAutoIdle() {
	r.log("ready-like text-only response auto-idled")
	delivery, pendingAtResume := r.waitForIdleResume()
	indicator, pendingAfterDelivery, delivered, interruptDelivered := r.consumeControlDelivery()
	fields := map[string]any{
		"event":             "ready_text_auto_idle_resumed",
		"delivery":          string(delivery),
		"pending_at_resume": pendingAtResume,
	}
	if indicator != "" {
		fields["inbox_indicator"] = indicator
	}
	if pendingAfterDelivery > 0 {
		fields["pending_count"] = pendingAfterDelivery
	}
	if len(delivered) > 0 {
		fields["incoming_messages"] = delivered
	}
	if interruptDelivered {
		fields["interrupt_notice"] = "Current operation was interrupted by peer control input."
	}
	body, err := json.MarshalIndent(fields, "", "  ")
	if err != nil {
		body = []byte(fmt.Sprintf("%#v", fields))
	}
	msg := tape.Message{
		Role: tape.RoleUser,
		Content: "Runtime control event resumed this process after a ready-like text-only response was auto-idled.\n" +
			string(body),
	}
	r.tape.Append(msg)
	r.writeTapeEntry(tape.MessageEntry(msg))
	r.appendContextMessage(msg)
}
