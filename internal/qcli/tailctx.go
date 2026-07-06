package qcli

import (
	"time"
)

func (s *Session) startContextTailer(initial tailState, poll time.Duration) {
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		s.tailContext(initial, poll)
	}()
}

func (s *Session) tailContext(st tailState, poll time.Duration) {
	poll = pollIntervalFromEnv(poll)
	for {
		select {
		case reason := <-s.resetReq:
			s.rebuildContext(reason, &st)
		default:
		}

		res, err := readTailLines(s.agent.ContextPath, &st)
		if err != nil {
			s.emit(ErrorEvent{Type: "error", Scope: "tail", Code: "context_read", Message: err.Error(), Recoverable: true})
		}
		if res.reset {
			s.applyContextReset("compaction", res.lines)
		} else if len(res.lines) > 0 {
			s.appendContextLines(res.lines, true)
		}

		select {
		case reason := <-s.resetReq:
			s.rebuildContext(reason, &st)
		case <-s.ctx.Done():
			return
		case <-time.After(poll):
		}
	}
}

func (s *Session) requestReset(reason string) {
	select {
	case s.resetReq <- reason:
	default:
	}
}

func (s *Session) rebuildContext(reason string, st *tailState) {
	lines, next, err := readAllCompleteLines(s.agent.ContextPath)
	if err != nil {
		s.emit(ErrorEvent{Type: "error", Scope: "tail", Code: "context_read", Message: err.Error(), Recoverable: true})
		return
	}
	*st = next
	s.applyContextReset(reason, lines)
}

func (s *Session) applyContextReset(reason string, lines []string) {
	s.mu.Lock()
	s.cells = nil
	s.pendingTools = map[string]bool{}
	s.lastContextAdvance = time.Now()
	s.mu.Unlock()
	s.emit(ContextResetEvent{Type: "context_reset", Reason: reason, TS: time.Now().UnixMilli()})
	s.appendContextLines(lines, true)
	s.emit(BackfillCompleteEvent{Type: "backfill_complete", Count: s.cellCount()})
}

func (s *Session) appendContextLines(lines []string, emit bool) {
	for _, line := range lines {
		for _, cell := range ParseTapeLine(line) {
			cell = s.stampCell(cell)
			s.updateWorkingState(cell)
			s.mu.Lock()
			s.cells = append(s.cells, cell)
			if emit {
				s.lastContextAdvance = time.Now()
			}
			s.mu.Unlock()
			if emit {
				s.emit(withCellType(cell))
			}
		}
	}
}

func (s *Session) stampCell(cell Cell) Cell {
	s.mu.Lock()
	s.nextSeq++
	cell.Seq = s.nextSeq
	s.mu.Unlock()
	return cell
}

func (s *Session) updateWorkingState(cell Cell) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.pendingTools == nil {
		s.pendingTools = map[string]bool{}
	}
	switch cell.Kind {
	case CellToolCall:
		if cell.ToolID != nil && *cell.ToolID != "" {
			s.pendingTools[*cell.ToolID] = true
		}
	case CellToolResult:
		if cell.ToolID != nil && *cell.ToolID != "" {
			delete(s.pendingTools, *cell.ToolID)
		}
	}
}

func (s *Session) cellCount() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return len(s.cells)
}
