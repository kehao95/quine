//go:build !linux

package tools

import "fmt"

func (s *subjectiveFS) init(_, _ string) error {
	if !s.enabled || s.initialized {
		return nil
	}
	return fmt.Errorf("workspace physics are only supported on Linux")
}

func (s *subjectiveFS) commandEnv() []string {
	return nil
}

func (s *subjectiveFS) childEnvOverrides() []string {
	return nil
}

func (s *subjectiveFS) snapshot() (fsSnapshot, error) {
	return fsSnapshot{}, nil
}

func (s *subjectiveFS) formatMutations(before, after fsSnapshot) string {
	return ""
}

func (s *subjectiveFS) captureWorldRevision(kind string, turnID int) (worldRevision, error) {
	return worldRevision{}, nil
}

func (s *subjectiveFS) loadCurrentWorldRevision() (worldRevision, error) {
	return worldRevision{}, nil
}

func (s *subjectiveFS) restoreWorld(revision string) (string, string, error) {
	return "", "", nil
}

func (s *subjectiveFS) commit() error {
	return nil
}

func (s *subjectiveFS) rollback() error {
	return nil
}
