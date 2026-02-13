package runtime

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// Semaphore provides a system-wide concurrency limiter using filesystem locks.
// Lock files live in {tapeDir}/../locks/ and are shared across all processes
// in the tree (they all share QUINE_DATA_DIR).
type Semaphore struct {
	lockDir   string
	maxSlots  int
	sessionID string
	logWriter io.Writer // optional; operational log messages go here instead of stderr

	mu       sync.Mutex
	lockFile string // path of the currently held lock file, or "" if none
	seq      int    // monotonic counter for unique lock file names
}

// NewSemaphore creates a Semaphore.
// lockDir is typically filepath.Join(filepath.Dir(filepath.Clean(cfg.TapeDir)), "locks").
func NewSemaphore(lockDir string, maxSlots int, sessionID string) *Semaphore {
	return &Semaphore{
		lockDir:   lockDir,
		maxSlots:  maxSlots,
		sessionID: sessionID,
	}
}

// Acquire attempts to acquire a slot. It blocks until one is available.
// Creates a lock file named {sessionID}-{seq}.lock in the lock directory.
// If all slots are full, polls every 1 second.
// If blocked for > 60 seconds, logs a warning to stderr.
func (s *Semaphore) Acquire() error {
	s.mu.Lock()
	seq := s.seq
	s.seq++
	s.mu.Unlock()

	// Ensure lock directory exists.
	if err := os.MkdirAll(s.lockDir, 0o755); err != nil {
		return fmt.Errorf("semaphore: creating lock dir: %w", err)
	}

	lockName := fmt.Sprintf("%s-%d.lock", s.sessionID, seq)
	lockPath := filepath.Join(s.lockDir, lockName)

	start := time.Now()
	warned := false

	for {
		// Count existing lock files.
		count := s.countFiles()
		if count < s.maxSlots {
			// Try atomic create.
			f, err := os.OpenFile(lockPath, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o644)
			if err == nil {
				f.Close()
				s.mu.Lock()
				s.lockFile = lockPath
				s.mu.Unlock()
				return nil
			}
			// O_EXCL failed (race): another process grabbed it. Retry.
			if !os.IsExist(err) {
				return fmt.Errorf("semaphore: creating lock file: %w", err)
			}
		}

		// Slot not available — poll.
		if !warned && time.Since(start) > 60*time.Second {
			if w := s.logWriter; w != nil {
				fmt.Fprintf(w, "quine: semaphore blocked for >60s waiting for concurrency slot (%d/%d)\n",
					count, s.maxSlots)
			}
			warned = true
		}

		time.Sleep(1 * time.Second)
	}
}

// Release removes the lock file, freeing the slot.
func (s *Semaphore) Release() error {
	s.mu.Lock()
	path := s.lockFile
	s.lockFile = ""
	s.mu.Unlock()

	if path == "" {
		return nil
	}
	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("semaphore: removing lock file: %w", err)
	}
	return nil
}

// Count returns the current number of acquired slots (lock files in lockDir).
func (s *Semaphore) Count() int {
	return s.countFiles()
}

// countFiles returns the number of .lock files in the lock directory.
func (s *Semaphore) countFiles() int {
	entries, err := os.ReadDir(s.lockDir)
	if err != nil {
		return 0
	}
	count := 0
	for _, e := range entries {
		if !e.IsDir() && filepath.Ext(e.Name()) == ".lock" {
			count++
		}
	}
	return count
}
