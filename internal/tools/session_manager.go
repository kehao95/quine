package tools

import (
	"fmt"
	"os"
	"sync"
)

// SessionManager manages named shell sessions.
// The owned map IS ownership — presence in the map means you own the session.
type SessionManager struct {
	owned map[string]*Session
	mu    sync.RWMutex

	// Configuration for spawning new sessions
	shell      string
	env        []string
	extraFiles []*os.File // deliverable stdout, material stdin
}

// NewSessionManager creates a new SessionManager.
func NewSessionManager(shell string, env []string, extraFiles []*os.File) *SessionManager {
	return &SessionManager{
		owned:      make(map[string]*Session),
		shell:      shell,
		env:        env,
		extraFiles: extraFiles,
	}
}

// Get returns a named session, or nil if it doesn't exist.
func (m *SessionManager) Get(name string) *Session {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.owned[name]
}

// GetOrCreate returns an existing session or creates a new one.
func (m *SessionManager) GetOrCreate(name string) (*Session, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if s, ok := m.owned[name]; ok {
		return s, nil
	}

	s, err := newSession(name, m.shell, m.env, m.extraFiles)
	if err != nil {
		return nil, err
	}

	// Start control pump for exit code detection
	go s.pumpControl()

	m.owned[name] = s
	return s, nil
}

// List returns the names of all owned sessions.
func (m *SessionManager) List() []string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	names := make([]string, 0, len(m.owned))
	for name := range m.owned {
		names = append(names, name)
	}
	return names
}

// CloseAll kills all owned sessions and clears the map.
func (m *SessionManager) CloseAll() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	var firstErr error
	for name, s := range m.owned {
		if err := s.Close(); err != nil && firstErr == nil {
			firstErr = fmt.Errorf("closing session %q: %w", name, err)
		}
		delete(m.owned, name)
	}
	return firstErr
}

// CloseSession kills a specific named session.
func (m *SessionManager) CloseSession(name string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	s, ok := m.owned[name]
	if !ok {
		return fmt.Errorf("session %q not found", name)
	}

	err := s.Close()
	delete(m.owned, name)
	return err
}

// Count returns the number of owned sessions.
func (m *SessionManager) Count() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.owned)
}
