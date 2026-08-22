//go:build !linux

package tools

import "fmt"

func (s *subjectiveFS) init(_, sessionID string) error {
	if !s.enabled || s.initialized {
		return nil
	}
	s.sessionID = sessionID
	if s.usesDirectBackend() {
		if s.workspaceSession == "" {
			s.workspaceSession = sessionID
		}
		if err := s.initDirectState(); err != nil {
			return err
		}
		s.initialized = true
		return nil
	}
	return fmt.Errorf("workspace physics are only supported on Linux")
}

func (s *subjectiveFS) commandEnv() []string {
	return nil
}

func (s *subjectiveFS) bootstrapWorkspaceState() error {
	return nil
}

func (s *subjectiveFS) exportCurrentTree(string) error {
	return fmt.Errorf("workspace physics are only supported on Linux")
}

func (s *subjectiveFS) captureWorldRevision(kind string, turnID int) (worldRevision, error) {
	return worldRevision{}, nil
}

func (s *subjectiveFS) loadCurrentWorldRevision() (worldRevision, error) {
	return worldRevision{}, nil
}

func (s *subjectiveFS) finalizeTurn(kind string, turnID int) (turnFinalizeResult, error) {
	if s.usesDirectBackend() {
		return s.finalizeDirectTurn()
	}
	return turnFinalizeResult{}, nil
}

func (s *subjectiveFS) importHostWorkspaceChanges(kind string, turnID int) (turnFinalizeResult, error) {
	if s.usesDirectBackend() {
		return s.finalizeDirectTurn()
	}
	return turnFinalizeResult{}, nil
}

func (s *subjectiveFS) observeHostWorkspaceChanges() error {
	if s.usesDirectBackend() {
		current, err := snapshotTree(s.workspace)
		if err != nil {
			return err
		}
		return s.saveObservedTree(current)
	}
	return nil
}

func (s *subjectiveFS) switchWorld(target string) (string, string, error) {
	return "", "", nil
}

func (s *subjectiveFS) restoreMutationBlock(previous, current string) (string, error) {
	return "", nil
}

func (s *subjectiveFS) readOverlayWorkspaceFile(string) ([]byte, error) {
	return nil, fmt.Errorf("workspace physics are only supported on Linux")
}

func (s *subjectiveFS) commit() error {
	return nil
}

func (s *subjectiveFS) rollback() error {
	return nil
}
