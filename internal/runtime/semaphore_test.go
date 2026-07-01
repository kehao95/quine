package runtime

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"syscall"
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

func TestAgentRegistryRegistersUnlimitedAndSyncsPIDSurface(t *testing.T) {
	root := t.TempDir()
	lockDir := filepath.Join(root, "locks")
	registry := NewAgentRegistry(lockDir, root, 0, "session-a", "run-a")

	if err := registry.Register(); err != nil {
		t.Fatalf("Register failed: %v", err)
	}
	defer registry.Deregister()

	if got := registry.Count(); got != 1 {
		t.Fatalf("Count = %d, want 1", got)
	}

	agentFile := filepath.Join(lockDir, "run-a.agent")
	if _, err := os.Stat(agentFile); err != nil {
		t.Fatalf("expected agent file to exist: %v", err)
	}
	var reg agentRegistration
	data, err := os.ReadFile(agentFile)
	if err != nil {
		t.Fatalf("read agent file: %v", err)
	}
	if err := json.Unmarshal(data, &reg); err != nil {
		t.Fatalf("unmarshal agent file: %v", err)
	}
	if reg.SessionID != "session-a" || reg.RunID != "run-a" || reg.PID != os.Getpid() {
		t.Fatalf("agent registration = %#v, want session-a/run-a/current pid", reg)
	}

	// Registration no longer creates pid route directly;
	// it is published as a separate step.
	wantTarget := filepath.Join(root, "agent", "session-a", "public")
	pidLink := filepath.Join(root, "pid", strconv.Itoa(os.Getpid()))
	if err := registry.PublishSelfPID(); err != nil {
		t.Fatalf("PublishSelfPID failed: %v", err)
	}
	pidTarget, err := os.Readlink(pidLink)
	if err != nil {
		t.Fatalf("read pid link: %v", err)
	}
	if pidTarget != wantTarget {
		t.Fatalf("pid link = %q, want %q", pidTarget, wantTarget)
	}
	if err := registry.UnpublishSelfPID(); err != nil {
		t.Fatalf("UnpublishSelfPID failed: %v", err)
	}
	if _, err := os.Lstat(pidLink); !os.IsNotExist(err) {
		t.Fatalf("expected pid link removed after unpublish, got err=%v", err)
	}
	for _, legacyPath := range []string{
		filepath.Join(root, "agent", "live"),
		filepath.Join(root, "agent", "live_by_pid"),
	} {
		if _, err := os.Stat(legacyPath); !os.IsNotExist(err) {
			t.Fatalf("expected legacy live index %s removed, got err=%v", legacyPath, err)
		}
	}

	if err := registry.Deregister(); err != nil {
		t.Fatalf("Deregister failed: %v", err)
	}
	if _, err := os.Stat(agentFile); !os.IsNotExist(err) {
		t.Fatalf("expected agent file removed, got err=%v", err)
	}
}

func TestAgentRegistryRejectsConcurrentRunForSameSession(t *testing.T) {
	root := t.TempDir()
	lockDir := filepath.Join(root, "locks")

	first := NewAgentRegistry(lockDir, root, 0, "session-a", "run-a")
	if err := first.Register(); err != nil {
		t.Fatalf("first Register failed: %v", err)
	}
	defer first.Deregister()

	second := NewAgentRegistry(lockDir, root, 0, "session-a", "run-b")
	err := second.Register()
	if err == nil {
		second.Deregister()
		t.Fatal("expected second Register to fail for same live session")
	}
	if got := err.Error(); !strings.Contains(got, `session "session-a" already active`) {
		t.Fatalf("Register error = %q, want same-session active rejection", got)
	}
}

func TestAgentRegistryDeregisterAllowsSameSessionReentry(t *testing.T) {
	root := t.TempDir()
	lockDir := filepath.Join(root, "locks")

	first := NewAgentRegistry(lockDir, root, 0, "session-a", "run-a")
	if err := first.Register(); err != nil {
		t.Fatalf("first Register failed: %v", err)
	}
	if err := first.Deregister(); err != nil {
		t.Fatalf("first Deregister failed: %v", err)
	}

	second := NewAgentRegistry(lockDir, root, 0, "session-a", "run-b")
	if err := second.Register(); err != nil {
		t.Fatalf("second Register after deregister failed: %v", err)
	}
	defer second.Deregister()
}

func TestAgentRegistryPrunesStaleRunAndAllowsSessionResume(t *testing.T) {
	root := t.TempDir()
	lockDir := filepath.Join(root, "locks")
	if err := os.MkdirAll(lockDir, 0o755); err != nil {
		t.Fatalf("mkdir lockDir: %v", err)
	}

	oldReg := agentRegistration{SessionID: "session-a", RunID: "run-old", PID: 999999}
	oldData, err := json.Marshal(oldReg)
	if err != nil {
		t.Fatalf("marshal old registration: %v", err)
	}
	if err := os.WriteFile(filepath.Join(lockDir, "run-old.agent"), append(oldData, '\n'), 0o644); err != nil {
		t.Fatalf("write old agent file: %v", err)
	}
	if err := os.MkdirAll(filepath.Join(root, "pid"), 0o755); err != nil {
		t.Fatalf("mkdir pid dir: %v", err)
	}
	if err := replaceSymlink(filepath.Join(root, "pid", "999999"), filepath.Join(root, "agent", "session-a", "public")); err != nil {
		t.Fatalf("seed stale pid link: %v", err)
	}

	registry := NewAgentRegistry(lockDir, root, 0, "session-a", "run-new")
	if err := registry.Register(); err != nil {
		t.Fatalf("Register after stale same-session run failed: %v", err)
	}
	defer registry.Deregister()

	if _, err := os.Stat(filepath.Join(lockDir, "run-old.agent")); !os.IsNotExist(err) {
		t.Fatalf("expected old stale run registration removed, got err=%v", err)
	}
	if _, err := os.Lstat(filepath.Join(root, "pid", "999999")); !os.IsNotExist(err) {
		t.Fatalf("expected old stale pid link removed, got err=%v", err)
	}
	if _, err := os.Stat(filepath.Join(lockDir, "run-new.agent")); err != nil {
		t.Fatalf("expected new run registration: %v", err)
	}
}

func TestAgentRegistryPrunesStaleRegistrationsAndRepairsPIDSurface(t *testing.T) {
	root := t.TempDir()
	lockDir := filepath.Join(root, "locks")
	if err := os.MkdirAll(lockDir, 0o755); err != nil {
		t.Fatalf("mkdir lockDir: %v", err)
	}

	const liveSession = "session-a"
	livePID := os.Getpid()
	if err := os.WriteFile(filepath.Join(lockDir, liveSession+".agent"), []byte(strconv.Itoa(livePID)+"\n"), 0o644); err != nil {
		t.Fatalf("write live agent file: %v", err)
	}

	const staleSession = "stale-session"
	const stalePID = 999999
	if err := os.WriteFile(filepath.Join(lockDir, staleSession+".agent"), []byte(strconv.Itoa(stalePID)+"\n"), 0o644); err != nil {
		t.Fatalf("write stale agent file: %v", err)
	}

	legacyLiveDir := filepath.Join(root, "agent", "live")
	legacyLiveByPIDDir := filepath.Join(root, "agent", "live_by_pid")
	for _, dir := range []string{legacyLiveDir, legacyLiveByPIDDir, filepath.Join(root, "pid")} {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			t.Fatalf("mkdir %s: %v", dir, err)
		}
	}

	if err := replaceSymlink(filepath.Join(legacyLiveDir, liveSession), filepath.Join(root, "agent", "wrong-target")); err != nil {
		t.Fatalf("seed wrong legacy live link: %v", err)
	}
	if err := replaceSymlink(filepath.Join(legacyLiveDir, staleSession), filepath.Join(root, "agent", staleSession)); err != nil {
		t.Fatalf("seed stale legacy live link: %v", err)
	}
	if err := replaceSymlink(filepath.Join(legacyLiveByPIDDir, strconv.Itoa(stalePID)), filepath.Join(root, "agent", staleSession)); err != nil {
		t.Fatalf("seed stale legacy pid link: %v", err)
	}
	if err := replaceSymlink(filepath.Join(root, "pid", strconv.Itoa(livePID)), filepath.Join(root, "agent", "wrong-target")); err != nil {
		t.Fatalf("seed wrong pid link: %v", err)
	}
	if err := replaceSymlink(filepath.Join(root, "pid", strconv.Itoa(stalePID)), filepath.Join(root, "agent", staleSession)); err != nil {
		t.Fatalf("seed stale pid link: %v", err)
	}

	registry := NewAgentRegistry(lockDir, root, 0, liveSession, "run-a")
	if err := registry.PruneStale(); err != nil {
		t.Fatalf("PruneStale failed: %v", err)
	}

	if got := registry.Count(); got != 1 {
		t.Fatalf("Count = %d, want 1", got)
	}
	if _, err := os.Stat(filepath.Join(lockDir, staleSession+".agent")); !os.IsNotExist(err) {
		t.Fatalf("expected stale agent file removed, got err=%v", err)
	}
	for _, legacyPath := range []string{legacyLiveDir, legacyLiveByPIDDir} {
		if _, err := os.Stat(legacyPath); !os.IsNotExist(err) {
			t.Fatalf("expected legacy live index %s removed, got err=%v", legacyPath, err)
		}
	}

	if _, err := os.Lstat(filepath.Join(root, "pid", strconv.Itoa(stalePID))); !os.IsNotExist(err) {
		t.Fatalf("expected stale pid link removed, got err=%v", err)
	}

	pidTarget, err := os.Readlink(filepath.Join(root, "pid", strconv.Itoa(livePID)))
	if err != nil {
		t.Fatalf("read live pid link: %v", err)
	}
	// The runtime no longer repairs live pid targets during pruning.
	// Existing links are left untouched unless stale.
	if pidTarget != filepath.Join(root, "agent", "wrong-target") {
		t.Fatalf("live pid link unexpectedly changed to %q", pidTarget)
	}
}

func TestAgentRegistryEnforcesMaxAgents(t *testing.T) {
	root := t.TempDir()
	lockDir := filepath.Join(root, "locks")

	first := NewAgentRegistry(lockDir, root, 1, "session-a", "run-a")
	if err := first.Register(); err != nil {
		t.Fatalf("first Register failed: %v", err)
	}
	defer first.Deregister()

	second := NewAgentRegistry(lockDir, root, 1, "session-b", "run-b")
	err := second.Register()
	if err == nil {
		second.Deregister()
		t.Fatal("expected second Register to fail when maxAgents is reached")
	}
	if got, want := err.Error(), "agent limit exceeded (1/1)"; got != want {
		t.Fatalf("Register error = %q, want %q", got, want)
	}

	if count := first.Count(); count != 1 {
		t.Fatalf("Count = %d, want 1", count)
	}
}

func TestAgentRegistrySharedFilesAreGroupWritableDespiteUmask(t *testing.T) {
	oldUmask := syscall.Umask(0o022)
	defer syscall.Umask(oldUmask)

	root := t.TempDir()
	lockDir := filepath.Join(root, "locks")
	registry := NewAgentRegistry(lockDir, root, 0, "session-a", "run-a")

	if err := registry.Register(); err != nil {
		t.Fatalf("Register failed: %v", err)
	}
	defer registry.Deregister()

	for _, path := range []string{
		filepath.Join(lockDir, "run-a.agent"),
		filepath.Join(lockDir, "agents", fmt.Sprintf("%d.agent.lock", os.Getpid())),
	} {
		info, err := os.Stat(path)
		if err != nil {
			t.Fatalf("stat %s: %v", path, err)
		}
		if got := info.Mode().Perm(); got != agentRegistrySharedFileMode {
			t.Fatalf("mode %s = %o, want %o", path, got, agentRegistrySharedFileMode)
		}
	}
}

func TestAgentRegistryPublishesAndUnpublishesSelfPID(t *testing.T) {
	root := t.TempDir()
	lockDir := filepath.Join(root, "locks")

	registry := NewAgentRegistry(lockDir, root, 0, "session-a", "run-a")
	if err := registry.Register(); err != nil {
		t.Fatalf("Register failed: %v", err)
	}
	defer registry.Deregister()

	pid := os.Getpid()
	pidLink := filepath.Join(root, "pid", strconv.Itoa(pid))
	target := filepath.Join(root, "agent", "session-a", "public")
	if err := registry.PublishSelfPID(); err != nil {
		t.Fatalf("PublishSelfPID failed: %v", err)
	}
	pidTarget, err := os.Readlink(pidLink)
	if err != nil {
		t.Fatalf("readlink: %v", err)
	}
	if pidTarget != target {
		t.Fatalf("pid link = %q, want %q", pidTarget, target)
	}
	if err := registry.UnpublishSelfPID(); err != nil {
		t.Fatalf("UnpublishSelfPID failed: %v", err)
	}
	if _, err := os.Lstat(pidLink); !os.IsNotExist(err) {
		t.Fatalf("expected pid link removed after unpublish, got err=%v", err)
	}
}

func TestAgentRegistryPruneStalePIDLocksRemovesArtifacts(t *testing.T) {
	root := t.TempDir()
	lockDir := filepath.Join(root, "locks")
	const stalePID = 99999
	pidStr := strconv.Itoa(stalePID)
	stalePublic := filepath.Join(root, "agent", "stale-session", "public")
	staleStatus := filepath.Join(root, "agent", "stale-session", "status", "session.json")
	if err := os.MkdirAll(filepath.Dir(staleStatus), 0o755); err != nil {
		t.Fatalf("mkdir stale status dir: %v", err)
	}

	sessionPayload := []byte(`{"pid":` + pidStr + `}`)
	if err := os.WriteFile(staleStatus, sessionPayload, 0o644); err != nil {
		t.Fatalf("write stale session payload: %v", err)
	}

	pidLink := filepath.Join(root, "pid", pidStr)
	if err := os.MkdirAll(filepath.Dir(pidLink), 0o755); err != nil {
		t.Fatalf("mkdir pid root: %v", err)
	}
	if err := replaceSymlink(pidLink, stalePublic); err != nil {
		t.Fatalf("seed stale pid link: %v", err)
	}

	pidLockPath := filepath.Join(lockDir, "agents", fmt.Sprintf("%d.agent.lock", stalePID))
	if err := os.MkdirAll(filepath.Dir(pidLockPath), 0o755); err != nil {
		t.Fatalf("mkdir lock dir: %v", err)
	}
	lockFile, err := os.OpenFile(pidLockPath, os.O_CREATE|os.O_RDWR, 0o644)
	if err != nil {
		t.Fatalf("seed stale pid lock: %v", err)
	}
	lockFile.Close()
	staleRegPath := filepath.Join(lockDir, "stale-run.agent")
	staleReg, err := json.Marshal(agentRegistration{SessionID: "stale-session", RunID: "stale-run", PID: stalePID})
	if err != nil {
		t.Fatalf("marshal stale registration: %v", err)
	}
	if err := os.WriteFile(staleRegPath, append(staleReg, '\n'), 0o644); err != nil {
		t.Fatalf("write stale registration: %v", err)
	}

	registry := NewAgentRegistry(lockDir, root, 0, "session-live", "run-live")
	if err := registry.Register(); err != nil {
		t.Fatalf("Register failed: %v", err)
	}
	defer registry.Deregister()

	if _, err := os.Stat(pidLink); !os.IsNotExist(err) {
		t.Fatalf("expected stale pid link removed by stale pid lock sweep, got err=%v", err)
	}
	if _, err := os.Stat(stalePublic); !os.IsNotExist(err) {
		t.Fatalf("expected stale agent root removed, got err=%v", err)
	}
	if _, err := os.Stat(pidLockPath); !os.IsNotExist(err) {
		t.Fatalf("expected stale pid lock removed, got err=%v", err)
	}
	if _, err := os.Stat(staleRegPath); !os.IsNotExist(err) {
		t.Fatalf("expected stale agent registration removed, got err=%v", err)
	}
}
