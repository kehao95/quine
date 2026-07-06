package qcli

import (
	"encoding/json"
	"time"
)

type controlLogEntry struct {
	Kind         string                 `json:"kind"`
	Timestamp    int64                  `json:"timestamp"`
	Action       string                 `json:"action,omitempty"`
	Delivery     string                 `json:"delivery,omitempty"`
	PendingCount int                    `json:"pending_count,omitempty"`
	Message      *controlRuntimeMessage `json:"message,omitempty"`
}

type controlRuntimeMessage struct {
	ID         string `json:"id"`
	Action     string `json:"action,omitempty"`
	Delivery   string `json:"delivery,omitempty"`
	Payload    string `json:"payload"`
	ReceivedAt int64  `json:"received_at"`
}

func (s *Session) startControlTailer(poll time.Duration) {
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		s.tailControl(poll)
	}()
}

func (s *Session) tailControl(poll time.Duration) {
	poll = pollIntervalFromEnv(poll)
	var st tailState
	for {
		res, err := readTailLines(s.agent.ControlLogPath, &st)
		if err != nil {
			s.emit(ErrorEvent{Type: "error", Scope: "tail", Code: "control_log_read", Message: err.Error(), Recoverable: true})
		}
		for _, line := range res.lines {
			var entry controlLogEntry
			if err := json.Unmarshal([]byte(line), &entry); err != nil {
				continue
			}
			ev, ok := s.receiptFor(entry)
			if ok {
				s.emit(ev)
			}
		}
		if !waitTick(s.ctx, poll) {
			return
		}
	}
}

func (s *Session) receiptFor(entry controlLogEntry) (ReceiptEvent, bool) {
	stage := ""
	switch entry.Kind {
	case "received", "woke", "delivered":
		stage = entry.Kind
	default:
		return ReceiptEvent{}, false
	}
	ts := entry.Timestamp
	if ts == 0 {
		ts = time.Now().UnixMilli()
	}
	ev := ReceiptEvent{Type: "receipt", Stage: stage, TS: ts}
	if entry.Action != "" {
		ev.Action = strptr(entry.Action)
	}
	if entry.Delivery != "" {
		ev.Delivery = strptr(entry.Delivery)
	}
	if entry.PendingCount != 0 || entry.Kind == "woke" {
		ev.Pending = intptr(entry.PendingCount)
	}
	if entry.Message != nil {
		if entry.Message.ID != "" {
			ev.MessageID = strptr(entry.Message.ID)
		}
		if ev.Action == nil && entry.Message.Action != "" {
			ev.Action = strptr(entry.Message.Action)
		}
		if ev.Delivery == nil && entry.Message.Delivery != "" {
			ev.Delivery = strptr(entry.Message.Delivery)
		}
		ref := s.matchClientRef(entry.Kind, entry.Message.ID, entry.Message.Payload)
		if ref != "" {
			ev.ClientRef = strptr(ref)
		}
	}
	return ev, true
}
