package qcli

import (
	"context"
	"encoding/json"
	"time"
)

type liveEntry struct {
	Seq      int    `json:"seq"`
	Kind     string `json:"kind"`
	Text     string `json:"text,omitempty"`
	ToolID   string `json:"tool_id,omitempty"`
	ToolName string `json:"tool_name,omitempty"`
	TS       int64  `json:"ts"`
}

func (s *Session) startLiveTailer(poll time.Duration) {
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		s.tailLive(poll)
	}()
}

func (s *Session) tailLive(poll time.Duration) {
	poll = pollIntervalFromEnv(poll)
	var st tailState
	generation := 1
	for {
		res, err := readTailLines(s.agent.LivePath, &st)
		if err != nil {
			s.emit(ErrorEvent{Type: "error", Scope: "tail", Code: "live_read", Message: err.Error(), Recoverable: true})
		}
		if res.reset && generation > 0 {
			generation++
			s.setLiveOpen(false)
		}
		for _, line := range res.lines {
			var entry liveEntry
			if err := json.Unmarshal([]byte(line), &entry); err != nil {
				continue
			}
			if entry.TS == 0 {
				entry.TS = time.Now().UnixMilli()
			}
			ev := StreamDeltaEvent{
				Type:       "stream_delta",
				Generation: generation,
				Seq:        entry.Seq,
				Kind:       entry.Kind,
				TS:         entry.TS,
			}
			if entry.Text != "" {
				ev.Text = strptr(entry.Text)
			}
			if entry.ToolID != "" {
				ev.ToolID = strptr(entry.ToolID)
			}
			if entry.ToolName != "" {
				ev.ToolName = strptr(entry.ToolName)
			}
			s.markLiveAdvance(entry.Kind != "completed")
			s.emit(ev)
		}
		if !waitTick(s.ctx, poll) {
			return
		}
	}
}

func (s *Session) markLiveAdvance(open bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.lastLiveAdvance = time.Now()
	s.liveOpen = open
}

func (s *Session) setLiveOpen(open bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.liveOpen = open
}

var _ = context.Canceled
