package qcli

import (
	"encoding/json"
	"strings"

	"github.com/kehao95/quine/internal/tape"
)

const (
	CellMeta       = "meta"
	CellMessage    = "message"
	CellReasoning  = "reasoning"
	CellToolCall   = "tool_call"
	CellToolResult = "tool_result"
	CellOutcome    = "outcome"
	CellReply      = "reply"
	CellHumanPost  = "human_post"
)

func ReplyCell(text string, ts int64) Cell {
	return Cell{Kind: CellReply, Text: strptr(strings.TrimRight(text, "\r\n")), TS: int64ptr(ts)}
}

func HumanPostCell(author string, action ControlAction, text string) Cell {
	author = strings.TrimSpace(author)
	if author == "" {
		author = humanAuthor
	}
	a := string(action)
	return Cell{Kind: CellHumanPost, Author: strptr(author), Action: &a, Text: strptr(strings.TrimRight(text, "\r\n"))}
}

func ParseTapeLine(line string) []Cell {
	line = strings.TrimSpace(line)
	if line == "" {
		return nil
	}
	var entry tape.TapeEntry
	if err := json.Unmarshal([]byte(line), &entry); err != nil {
		return nil
	}
	switch entry.Type {
	case "meta":
		return parseMetaCells(entry.Data)
	case "message":
		return parseMessageCells(entry.Data)
	case "tool_result":
		return parseToolResultCells(entry.Data)
	case "outcome":
		return parseOutcomeCells(entry.Data)
	default:
		raw := line
		return []Cell{{Kind: entry.Type, Raw: &raw}}
	}
}

func parseMetaCells(data json.RawMessage) []Cell {
	var meta struct {
		SessionID string `json:"session_id"`
		ModelID   string `json:"model_id"`
		Depth     int    `json:"depth"`
	}
	if err := json.Unmarshal(data, &meta); err != nil {
		return nil
	}
	return []Cell{{
		Kind:    CellMeta,
		Text:    strptr(meta.ModelID),
		Session: strptr(meta.SessionID),
		Depth:   intptr(meta.Depth),
	}}
}

func parseMessageCells(data json.RawMessage) []Cell {
	var msg tape.Message
	if err := json.Unmarshal(data, &msg); err != nil {
		return nil
	}
	if author, action, body, ok := parseClientEnvelope(msg.Content); ok {
		if strings.TrimSpace(body) == "" {
			return nil
		}
		c := HumanPostCell(author, ControlAction(action), body)
		if msg.Timestamp != 0 {
			c.TS = int64ptr(msg.Timestamp)
		}
		return []Cell{c}
	}
	var cells []Cell
	role := string(msg.Role)
	if r := reasoningText(msg); r != "" {
		cells = append(cells, Cell{Kind: CellReasoning, Role: strptr(role), Text: strptr(r), TS: maybeTS(msg.Timestamp)})
	}
	if strings.TrimSpace(msg.Content) != "" {
		cells = append(cells, Cell{Kind: CellMessage, Role: strptr(role), Text: strptr(msg.Content), TS: maybeTS(msg.Timestamp)})
	}
	for _, call := range msg.ToolCalls {
		args := compactArgs(call.Arguments)
		cell := Cell{
			Kind:     CellToolCall,
			Role:     strptr(role),
			ToolName: strptr(call.Name),
			ToolID:   strptr(call.ID),
			Text:     strptr(args),
			TS:       maybeTS(msg.Timestamp),
		}
		cells = append(cells, cell)
	}
	return cells
}

func parseToolResultCells(data json.RawMessage) []Cell {
	var tr tape.ToolResult
	if err := json.Unmarshal(data, &tr); err != nil {
		return nil
	}
	if humans := extractIncomingHumanPosts(tr.Content); len(humans) > 0 {
		return humans
	}
	toolName, status, text := summarizeToolContent(tr.Content)
	cell := Cell{Kind: CellToolResult, ToolID: strptr(tr.ToolID), IsError: tr.IsError}
	if toolName != "" {
		cell.ToolName = strptr(toolName)
	}
	if status != "" {
		cell.Status = strptr(status)
	}
	if text != "" {
		cell.Text = strptr(text)
	}
	return []Cell{cell}
}

func extractIncomingHumanPosts(raw json.RawMessage) []Cell {
	var c struct {
		Runtime struct {
			IncomingMessages []struct {
				Payload string `json:"payload"`
			} `json:"incoming_messages"`
		} `json:"runtime"`
	}
	if err := json.Unmarshal(raw, &c); err != nil {
		return nil
	}
	var cells []Cell
	for _, m := range c.Runtime.IncomingMessages {
		author, action, body, ok := parseClientEnvelope(m.Payload)
		if !ok || strings.TrimSpace(body) == "" {
			continue
		}
		cells = append(cells, HumanPostCell(author, ControlAction(action), body))
	}
	return cells
}

func parseOutcomeCells(data json.RawMessage) []Cell {
	if blankRaw(data) {
		return nil
	}
	var oc tape.SessionOutcome
	if err := json.Unmarshal(data, &oc); err != nil {
		return nil
	}
	mode := string(oc.TerminationMode)
	text := strings.TrimSpace(mode)
	if text == "" {
		text = "outcome"
	}
	return []Cell{{
		Kind:            CellOutcome,
		Text:            strptr(text),
		ExitCode:        intptr(oc.ExitCode),
		TurnCount:       intptr(oc.TurnCount),
		TerminationMode: strptr(mode),
		TokensIn:        intptr(oc.TokensIn),
		TokensOut:       intptr(oc.TokensOut),
	}}
}

func maybeTS(ts int64) *int64 {
	if ts == 0 {
		return nil
	}
	return int64ptr(ts)
}

func reasoningText(msg tape.Message) string {
	if strings.TrimSpace(msg.ReasoningContent) != "" {
		return msg.ReasoningContent
	}
	var parts []string
	for _, item := range msg.ReasoningItems {
		for _, s := range item.Summary {
			if strings.TrimSpace(s) != "" {
				parts = append(parts, s)
			}
		}
	}
	return strings.Join(parts, "\n")
}

func compactArgs(args map[string]any) string {
	if len(args) == 0 {
		return ""
	}
	b, err := json.Marshal(args)
	if err != nil {
		return ""
	}
	return string(b)
}

func summarizeToolContent(raw json.RawMessage) (toolName, status, text string) {
	if blankRaw(raw) {
		return "", "", ""
	}
	var obj map[string]any
	if err := json.Unmarshal(raw, &obj); err == nil {
		if s, ok := obj["tool"].(string); ok {
			toolName = s
		}
		if s, ok := obj["status"].(string); ok {
			status = s
		}
		for _, k := range []string{"text", "output", "content", "stdout", "message"} {
			if v, ok := obj[k].(string); ok && strings.TrimSpace(v) != "" {
				return toolName, status, v
			}
		}
		return toolName, status, ""
	}
	var s string
	if err := json.Unmarshal(raw, &s); err == nil {
		return "", "", s
	}
	return "", "", strings.TrimSpace(string(raw))
}

func blankRaw(raw json.RawMessage) bool {
	t := strings.TrimSpace(string(raw))
	return t == "" || t == "null"
}
