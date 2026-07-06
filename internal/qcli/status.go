package qcli

import (
	"encoding/json"
	"os"
	"reflect"
	"time"
)

func ReadStatus(agent Agent, working bool) (StatusEvent, error) {
	status := sessionStatus{}
	if data, err := os.ReadFile(agent.SessionPath); err == nil {
		if err := json.Unmarshal(data, &status); err != nil {
			return StatusEvent{}, err
		}
	}
	pending := 0
	if data, err := os.ReadFile(agent.InboxPath); err == nil {
		var inbox inboxStatus
		if err := json.Unmarshal(data, &inbox); err == nil {
			pending = inbox.PendingCount
		}
	}
	ev := StatusEvent{
		Type:          "status",
		Live:          pidLive(agent.PID),
		Working:       working,
		PID:           agent.PID,
		PPID:          agent.PPID,
		Session:       agent.Session,
		RunID:         agent.RunID,
		Incarnation:   agent.Incarnation,
		ParentSession: agent.ParentSession,
		Model:         agent.Model,
		Depth:         agent.Depth,
		Pending:       pending,
		TokensIn:      nil,
		TokensOut:     nil,
	}
	if status.PID > 0 {
		ev.PID = status.PID
	}
	if status.PPID != 0 {
		ev.PPID = status.PPID
	}
	if status.SessionID != "" {
		ev.Session = status.SessionID
	}
	if status.RunID != "" {
		ev.RunID = status.RunID
	}
	if status.IncarnationID != 0 {
		ev.Incarnation = status.IncarnationID
	}
	if status.ParentSession != "" {
		ev.ParentSession = status.ParentSession
	}
	if status.ModelID != "" {
		ev.Model = status.ModelID
	}
	if status.Depth != 0 {
		ev.Depth = status.Depth
	}
	return ev, nil
}

func (s *Session) startStatusPoller() {
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		ticker := time.NewTicker(500 * time.Millisecond)
		defer ticker.Stop()
		s.pollStatus()
		for {
			select {
			case <-s.ctx.Done():
				return
			case <-ticker.C:
				s.pollStatus()
			}
		}
	}()
}

func (s *Session) pollStatus() {
	working := s.working()
	status, err := ReadStatus(s.agent, working)
	if err != nil {
		s.emit(ErrorEvent{Type: "error", Scope: "tail", Code: "status_read", Message: err.Error(), Recoverable: true})
		return
	}

	s.mu.Lock()
	prev := s.status
	first := s.status.Session == "" && s.status.PID == 0
	incarnationChanged := !first && prev.Incarnation != 0 && status.Incarnation != 0 && prev.Incarnation != status.Incarnation
	s.status = status
	s.mu.Unlock()

	if incarnationChanged {
		s.requestReset("incarnation")
	}
	if first || !reflect.DeepEqual(prev, status) {
		s.emit(status)
	}
}

func (s *Session) working() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	now := time.Now()
	if !s.lastContextAdvance.IsZero() && now.Sub(s.lastContextAdvance) < 2*time.Second {
		return true
	}
	if !s.lastLiveAdvance.IsZero() && now.Sub(s.lastLiveAdvance) < 2*time.Second {
		return true
	}
	if s.liveOpen {
		return true
	}
	return len(s.pendingTools) > 0
}
