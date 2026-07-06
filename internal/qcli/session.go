package qcli

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"
)

type SessionOptions struct {
	ClientID     string
	PollInterval time.Duration
	EventBuffer  int
	CtlTimeout   time.Duration
}

type Session struct {
	agent    Agent
	endpoint ClientEndpoint
	writer   *CtlWriter
	timeout  time.Duration

	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
	once   sync.Once

	resetReq chan string

	mu                 sync.Mutex
	nextSeq            int64
	cells              []Cell
	status             StatusEvent
	pendingTools       map[string]bool
	lastContextAdvance time.Time
	lastLiveAdvance    time.Time
	liveOpen           bool
	nextClientRef      int64
	recentSends        []sendRecord
	messageRefs        map[string]string
	subs               map[int]chan any
	nextSub            int
}

type sendRecord struct {
	ClientRef string
	Payload   string
}

func StartSession(ctx context.Context, agent Agent, opts SessionOptions) (*Session, error) {
	endpoint, err := OpenClientEndpoint(agent.RuntimeRoot, opts.ClientID)
	if err != nil {
		return nil, err
	}
	sessionCtx, cancel := context.WithCancel(ctx)
	timeout := opts.CtlTimeout
	if timeout <= 0 {
		timeout = CtlTimeoutFromEnv()
	}
	s := &Session{
		agent:        agent,
		endpoint:     endpoint,
		writer:       NewCtlWriter(),
		timeout:      timeout,
		ctx:          sessionCtx,
		cancel:       cancel,
		resetReq:     make(chan string, 1),
		pendingTools: map[string]bool{},
		messageRefs:  map[string]string{},
		subs:         map[int]chan any{},
	}
	lines, initialState, err := readAllCompleteLines(agent.ContextPath)
	if err != nil {
		cancel()
		removeClientEndpoint(endpoint)
		return nil, err
	}
	s.appendContextLines(lines, false)
	if status, err := ReadStatus(agent, false); err == nil {
		s.status = status
	}
	s.startEndpointWatcher()
	s.startContextTailer(initialState, opts.PollInterval)
	s.startLiveTailer(opts.PollInterval)
	s.startControlTailer(opts.PollInterval)
	s.startStatusPoller()
	return s, nil
}

func (s *Session) Agent() Agent {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.agent
}

func (s *Session) Endpoint() ClientEndpoint {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.endpoint
}

func (s *Session) Status() StatusEvent {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.status
}

func (s *Session) Cells() []Cell {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make([]Cell, len(s.cells))
	copy(out, s.cells)
	return out
}

func (s *Session) Subscribe() (<-chan any, func()) {
	s.mu.Lock()
	id := s.nextSub
	s.nextSub++
	cells := make([]Cell, len(s.cells))
	copy(cells, s.cells)
	status := s.status
	if status.Type == "" {
		status.Type = "status"
	}
	ch := make(chan any, len(cells)+1024)
	ch <- HelloEvent{Type: "hello", Contract: ContractVersion, EndpointID: s.endpoint.ID, Peer: s.agent.PeerInfo()}
	ch <- status
	for _, cell := range cells {
		ch <- withCellType(cell)
	}
	ch <- BackfillCompleteEvent{Type: "backfill_complete", Count: len(cells)}
	s.subs[id] = ch
	s.mu.Unlock()

	return ch, func() {
		s.mu.Lock()
		defer s.mu.Unlock()
		if c, ok := s.subs[id]; ok {
			delete(s.subs, id)
			close(c)
		}
	}
}

func (s *Session) emit(ev any) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for id, ch := range s.subs {
		select {
		case ch <- ev:
		default:
			delete(s.subs, id)
			close(ch)
		}
	}
}

func (s *Session) Send(action ControlAction, text string) (QueuedResponse, error) {
	if !validControlAction(action) {
		return QueuedResponse{}, fmt.Errorf("qcli: unknown control action %q", action)
	}
	path, allowEmpty, err := controlPathForAction(s.agent, action)
	if err != nil {
		return QueuedResponse{}, err
	}
	trimmed := strings.TrimRight(text, "\r\n")
	empty := strings.TrimSpace(trimmed) == ""
	if empty && !allowEmpty {
		return QueuedResponse{}, ErrEmptyPayload
	}

	ref := s.nextRef()
	var payload []byte
	envelope := ""
	if !empty {
		envelope = FormatClientSignalPayload(s.endpoint, action, trimmed)
		payload = []byte(envelope)
		s.recordSend(ref, envelope)
	}
	if err := s.writer.Write(path, payload, s.timeout); err != nil {
		return QueuedResponse{}, err
	}
	return QueuedResponse{Type: "queued", Action: string(action), ClientRef: ref}, nil
}

func (s *Session) nextRef() string {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.nextClientRef++
	return fmt.Sprintf("c-%d", s.nextClientRef)
}

func (s *Session) recordSend(ref, payload string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.recentSends = append(s.recentSends, sendRecord{ClientRef: ref, Payload: payload})
	if len(s.recentSends) > 64 {
		copy(s.recentSends, s.recentSends[len(s.recentSends)-64:])
		s.recentSends = s.recentSends[:64]
	}
}

func (s *Session) matchClientRef(kind, messageID, payload string) string {
	s.mu.Lock()
	defer s.mu.Unlock()
	switch kind {
	case "received":
		for i := len(s.recentSends) - 1; i >= 0; i-- {
			rec := s.recentSends[i]
			if rec.Payload != payload {
				continue
			}
			if messageID != "" {
				s.messageRefs[messageID] = rec.ClientRef
			}
			return rec.ClientRef
		}
	case "delivered":
		if ref := s.messageRefs[messageID]; ref != "" {
			return ref
		}
	}
	return ""
}

func (s *Session) startEndpointWatcher() {
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		err := ServeClientEndpoint(s.ctx, s.endpoint, func(payload string, ts int64) {
			cell := s.stampCell(ReplyCell(payload, ts))
			s.mu.Lock()
			s.cells = append(s.cells, cell)
			s.mu.Unlock()
			s.emit(withCellType(cell))
		})
		if err != nil && !errors.Is(err, context.Canceled) {
			s.emit(ErrorEvent{Type: "error", Scope: "tail", Code: "client_endpoint", Message: err.Error(), Recoverable: false})
		}
	}()
}

func (s *Session) Close() {
	s.once.Do(func() {
		s.cancel()
		s.wg.Wait()
		removeClientEndpoint(s.endpoint)
		s.mu.Lock()
		for id, ch := range s.subs {
			delete(s.subs, id)
			close(ch)
		}
		s.mu.Unlock()
	})
}
