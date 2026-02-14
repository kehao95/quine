package runtime

import (
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"
)

func TestSemaphoreAcquireRelease(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "locks")
	sem := NewSemaphore(dir, 5, "test-session")

	if err := sem.Acquire(); err != nil {
		t.Fatalf("Acquire failed: %v", err)
	}

	if c := sem.Count(); c != 1 {
		t.Errorf("expected count=1 after acquire, got %d", c)
	}

	if err := sem.Release(); err != nil {
		t.Fatalf("Release failed: %v", err)
	}

	if c := sem.Count(); c != 0 {
		t.Errorf("expected count=0 after release, got %d", c)
	}
}

func TestSemaphoreMaxSlots(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "locks")
	maxSlots := 3

	// Use separate semaphores (different sessionIDs) to simulate different processes.
	sems := make([]*Semaphore, maxSlots)
	for i := 0; i < maxSlots; i++ {
		sems[i] = NewSemaphore(dir, maxSlots, "session-"+string(rune('A'+i)))
		if err := sems[i].Acquire(); err != nil {
			t.Fatalf("Acquire %d failed: %v", i, err)
		}
	}

	if c := sems[0].Count(); c != maxSlots {
		t.Errorf("expected count=%d, got %d", maxSlots, c)
	}

	// Start a goroutine that tries to acquire one more — should block.
	blocked := NewSemaphore(dir, maxSlots, "session-blocked")
	done := make(chan error, 1)
	var mu sync.Mutex
	acquired := false

	go func() {
		err := blocked.Acquire()
		mu.Lock()
		acquired = true
		mu.Unlock()
		done <- err
	}()

	// Give the goroutine time to attempt acquire and block.
	time.Sleep(200 * time.Millisecond)

	mu.Lock()
	if acquired {
		mu.Unlock()
		t.Fatal("goroutine acquired a slot when all slots should be full")
	}
	mu.Unlock()

	// Release one slot.
	if err := sems[0].Release(); err != nil {
		t.Fatalf("Release failed: %v", err)
	}

	// Wait for the blocked goroutine to proceed.
	select {
	case err := <-done:
		if err != nil {
			t.Fatalf("blocked Acquire returned error: %v", err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("blocked goroutine did not acquire slot within 5 seconds")
	}

	mu.Lock()
	if !acquired {
		mu.Unlock()
		t.Error("expected goroutine to have acquired the slot")
	} else {
		mu.Unlock()
	}

	// Clean up.
	if err := blocked.Release(); err != nil {
		t.Fatalf("Release blocked failed: %v", err)
	}
	for i := 1; i < maxSlots; i++ {
		if err := sems[i].Release(); err != nil {
			t.Fatalf("Release %d failed: %v", i, err)
		}
	}
}

func TestSemaphoreIsFull(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "locks")
	maxSlots := 2

	sem1 := NewSemaphore(dir, maxSlots, "session-1")
	sem2 := NewSemaphore(dir, maxSlots, "session-2")

	// Initially not full
	if sem1.IsFull() {
		t.Error("expected IsFull=false when no slots acquired")
	}

	// Acquire one slot
	if err := sem1.Acquire(); err != nil {
		t.Fatalf("Acquire 1 failed: %v", err)
	}
	if sem1.IsFull() {
		t.Errorf("expected IsFull=false with 1/%d slots", maxSlots)
	}

	// Acquire second slot - now full
	if err := sem2.Acquire(); err != nil {
		t.Fatalf("Acquire 2 failed: %v", err)
	}
	if !sem1.IsFull() {
		t.Error("expected IsFull=true when all slots acquired")
	}

	// Release one - no longer full
	if err := sem1.Release(); err != nil {
		t.Fatalf("Release 1 failed: %v", err)
	}
	if sem2.IsFull() {
		t.Error("expected IsFull=false after releasing one slot")
	}

	// Clean up
	sem2.Release()
}

func TestSemaphoreCreatesDirectory(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "nested", "deep", "locks")

	// Verify it doesn't exist yet.
	if _, err := os.Stat(dir); !os.IsNotExist(err) {
		t.Fatalf("expected directory to not exist, got err: %v", err)
	}

	sem := NewSemaphore(dir, 5, "test-session")

	if err := sem.Acquire(); err != nil {
		t.Fatalf("Acquire failed: %v", err)
	}
	defer sem.Release()

	// Verify directory was created.
	info, err := os.Stat(dir)
	if err != nil {
		t.Fatalf("expected directory to exist after Acquire, got: %v", err)
	}
	if !info.IsDir() {
		t.Fatal("expected lock path to be a directory")
	}
}
