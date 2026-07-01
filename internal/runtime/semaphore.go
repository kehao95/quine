package runtime

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"
)

// Semaphore provides a system-wide concurrency limiter using filesystem locks.
// Lock files live under QUINE_DATA_DIR/locks/ and are shared across all
// processes in the tree.
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
// lockDir is typically QUINE_DATA_DIR/locks/.
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
	if s.maxSlots <= 0 {
		return nil
	}

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
	if s.maxSlots <= 0 {
		return 0
	}
	return s.countFiles()
}

// IsFull returns true if all slots are currently occupied.
func (s *Semaphore) IsFull() bool {
	if s.maxSlots <= 0 {
		return false
	}
	return s.countFiles() >= s.maxSlots
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

// AgentRegistry tracks the total number of agents in the process tree.
// Each agent registers on startup and deregisters on shutdown.
//
// Internal registration authority lives in .agent files under the shared
// coordination directory; the runtime-owned liveness lock is now kept in
// "<QUINE_DATA_DIR>/locks/agents/<pid>.agent.lock".
// The public routing surface under runtimeRoot/pid/ is now authored only by
// the owning process.
type AgentRegistry struct {
	agentDir    string
	pidLockDir  string
	runtimeRoot string
	maxAgents   int
	sessionID   string
	runID       string
	logWriter   io.Writer

	mu          sync.Mutex
	agentFile   string // path of this agent's registration file
	pidLockFile *os.File
	pidLockPath string
	pid         int
}

type agentRegistration struct {
	SessionID string `json:"session_id"`
	RunID     string `json:"run_id"`
	PID       int    `json:"pid"`
}

const agentRegistrySharedFileMode os.FileMode = 0o664

// NewAgentRegistry creates an AgentRegistry.
// agentDir is the internal registration directory, typically the same as
// Semaphore's lockDir.
// maxAgents of 0 means unlimited.
func NewAgentRegistry(agentDir string, runtimeRoot string, maxAgents int, sessionID string, runID string) *AgentRegistry {
	if strings.TrimSpace(runID) == "" {
		runID = sessionID
	}
	return &AgentRegistry{
		agentDir:    agentDir,
		pidLockDir:  filepath.Join(agentDir, "agents"),
		runtimeRoot: runtimeRoot,
		maxAgents:   maxAgents,
		sessionID:   sessionID,
		runID:       runID,
		pid:         os.Getpid(),
	}
}

// Register creates an .agent file for this process and acquires the PID lock.
// Returns an error if the agent limit would be exceeded.
func (r *AgentRegistry) Register() error {
	// Ensure agent directory exists.
	if err := os.MkdirAll(r.agentDir, 0o755); err != nil {
		return fmt.Errorf("agent registry: creating dir: %w", err)
	}
	if err := os.MkdirAll(r.pidLockDir, 0o755); err != nil {
		return fmt.Errorf("agent registry: creating pid lock dir: %w", err)
	}
	if err := r.PruneStale(); err != nil {
		return fmt.Errorf("agent registry: prune stale registrations: %w", err)
	}

	pid := r.pid
	agentPath := filepath.Join(r.agentDir, r.runID+".agent")
	active, err := r.activeAgents()
	if err != nil {
		return fmt.Errorf("agent registry: scan active agents: %w", err)
	}
	for runID, reg := range active {
		if runID != r.runID && reg.SessionID == r.sessionID {
			return fmt.Errorf("session %q already active as run %q (pid %d)", r.sessionID, runID, reg.PID)
		}
	}
	alreadyRegistered := false
	if _, err := os.Stat(agentPath); err == nil {
		alreadyRegistered = true
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("agent registry: stat agent file: %w", err)
	}

	if r.maxAgents > 0 && !alreadyRegistered {
		if len(active) >= r.maxAgents {
			return fmt.Errorf("agent limit exceeded (%d/%d)", len(active), r.maxAgents)
		}
	}

	createdAgentFile := false
	if !alreadyRegistered {
		if f, err := os.OpenFile(agentPath, os.O_CREATE|os.O_EXCL|os.O_WRONLY, agentRegistrySharedFileMode); err != nil {
			if !os.IsExist(err) {
				return fmt.Errorf("agent registry: creating agent file: %w", err)
			}
		} else {
			if chmodErr := os.Chmod(agentPath, agentRegistrySharedFileMode); chmodErr != nil {
				f.Close()
				_ = os.Remove(agentPath)
				return fmt.Errorf("agent registry: chmod agent file: %w", chmodErr)
			}
			reg := agentRegistration{SessionID: r.sessionID, RunID: r.runID, PID: pid}
			data, marshalErr := json.Marshal(reg)
			if marshalErr != nil {
				f.Close()
				_ = os.Remove(agentPath)
				return fmt.Errorf("agent registry: marshal registration: %w", marshalErr)
			}
			if _, writeErr := fmt.Fprintf(f, "%s\n", data); writeErr != nil {
				f.Close()
				_ = os.Remove(agentPath)
				return fmt.Errorf("agent registry: writing agent file: %w", writeErr)
			}
			if closeErr := f.Close(); closeErr != nil {
				_ = os.Remove(agentPath)
				return fmt.Errorf("agent registry: closing agent file: %w", closeErr)
			}
			createdAgentFile = true
		}
	}

	if err := r.acquirePIDLock(pid); err != nil {
		if createdAgentFile {
			_ = os.Remove(agentPath)
		}
		return fmt.Errorf("agent registry: acquire pid lock: %w", err)
	}

	r.mu.Lock()
	r.agentFile = agentPath
	r.mu.Unlock()

	return nil
}

// Deregister removes this agent's .agent file and releases its PID lock.
func (r *AgentRegistry) Deregister() error {
	r.mu.Lock()
	path := r.agentFile
	r.agentFile = ""
	pathLock := r.pidLockPath
	r.mu.Unlock()

	if path == "" {
		if err := r.releasePIDLock(pathLock); err != nil {
			return err
		}
		return r.PruneStale()
	}
	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("agent registry: removing %s: %w", path, err)
	}
	if err := r.releasePIDLock(pathLock); err != nil {
		return err
	}
	return r.PruneStale()
}

// PublishSelfPID creates/repairs runtimeRoot/pid/<pid> for this session.
// This route must be owned by this process; no other process should create it.
func (r *AgentRegistry) PublishSelfPID() error {
	if strings.TrimSpace(r.runtimeRoot) == "" {
		return nil
	}
	pid := r.pid
	pidLink := filepath.Join(r.runtimeRoot, "pid", strconv.Itoa(pid))
	publicTarget := filepath.Join(r.runtimeRoot, "agent", r.sessionID, "public")
	if err := os.MkdirAll(filepath.Dir(pidLink), 0o755); err != nil {
		return fmt.Errorf("agent registry: creating pid dir: %w", err)
	}
	return replaceSymlink(pidLink, publicTarget)
}

// UnpublishSelfPID removes runtimeRoot/pid/<pid> when it belongs to this
// session. Missing entries are treated as success.
func (r *AgentRegistry) UnpublishSelfPID() error {
	if strings.TrimSpace(r.runtimeRoot) == "" {
		return nil
	}
	pid := r.pid
	pidLink := filepath.Join(r.runtimeRoot, "pid", strconv.Itoa(pid))
	target, err := os.Readlink(pidLink)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("agent registry: read pid link %s: %w", pidLink, err)
	}
	if !filepath.IsAbs(target) {
		target = filepath.Join(filepath.Dir(pidLink), target)
	}
	expected := filepath.Join(r.runtimeRoot, "agent", r.sessionID, "public")
	if filepath.Clean(target) != filepath.Clean(expected) {
		return nil
	}
	if err := os.Remove(pidLink); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("agent registry: remove pid link %s: %w", pidLink, err)
	}
	return nil
}

// Count returns the current number of registered agents.
func (r *AgentRegistry) Count() int {
	active, err := r.activeAgents()
	if err != nil {
		return 0
	}
	return len(active)
}

// IsFull returns true if the agent limit has been reached.
// Returns false if maxAgents is 0 (unlimited).
func (r *AgentRegistry) IsFull() bool {
	if r.maxAgents <= 0 {
		return false
	}
	return r.Count() >= r.maxAgents
}

// CanSpawn returns true if a new agent can be spawned (count < max).
// Returns true if maxAgents is 0 (unlimited).
func (r *AgentRegistry) CanSpawn() bool {
	if r.maxAgents <= 0 {
		return true
	}
	return r.Count() < r.maxAgents
}

// PruneStale removes dead registrations and prunes stale pid routes.
func (r *AgentRegistry) PruneStale() error {
	active, err := r.activeAgents()
	if err != nil {
		return fmt.Errorf("agent registry: scan active agents: %w", err)
	}
	if strings.TrimSpace(r.runtimeRoot) != "" {
		if err := r.removeLegacyLiveIndexes(); err != nil {
			return err
		}
	}
	if err := r.pruneStalePIDLocks(); err != nil {
		return err
	}
	if strings.TrimSpace(r.runtimeRoot) != "" {
		if err := r.prunePIDSurface(active); err != nil {
			return err
		}
	}
	return nil
}

func (r *AgentRegistry) activeAgents() (map[string]agentRegistration, error) {
	active := map[string]agentRegistration{}
	entries, err := os.ReadDir(r.agentDir)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return active, nil
		}
		return nil, err
	}
	for _, e := range entries {
		if e.IsDir() || filepath.Ext(e.Name()) != ".agent" {
			continue
		}
		fallbackID := strings.TrimSuffix(e.Name(), ".agent")
		path := filepath.Join(r.agentDir, e.Name())
		reg, ok := readAgentRegistration(path, fallbackID)
		if !ok {
			if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
				return nil, fmt.Errorf("remove stale agent file %s: %w", path, err)
			}
			continue
		}
		alive, err := r.isProcessAlive(reg.PID)
		if err != nil {
			return nil, fmt.Errorf("agent registry: check pid %d: %w", reg.PID, err)
		}
		if !alive {
			if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
				return nil, fmt.Errorf("remove stale agent file %s: %w", path, err)
			}
			continue
		}
		active[reg.RunID] = reg
	}
	return active, nil
}

func (r *AgentRegistry) prunePIDSurface(active map[string]agentRegistration) error {
	if strings.TrimSpace(r.runtimeRoot) == "" {
		return nil
	}
	pidRoot := filepath.Join(r.runtimeRoot, "pid")
	if err := os.MkdirAll(pidRoot, 0o755); err != nil {
		return fmt.Errorf("agent registry: creating pid dir: %w", err)
	}
	activePIDs := map[string]struct{}{}
	for _, reg := range active {
		activePIDs[strconv.Itoa(reg.PID)] = struct{}{}
	}
	entries, err := os.ReadDir(pidRoot)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return err
	}
	for _, entry := range entries {
		name := entry.Name()
		if _, want := activePIDs[name]; want {
			continue
		}
		if err := os.Remove(filepath.Join(pidRoot, name)); err != nil && !os.IsNotExist(err) {
			return fmt.Errorf("agent registry: remove stale pid route %s: %w", filepath.Join(pidRoot, name), err)
		}
	}
	return nil
}

func (r *AgentRegistry) removeLegacyLiveIndexes() error {
	if strings.TrimSpace(r.runtimeRoot) == "" {
		return nil
	}
	for _, legacyPath := range []string{
		filepath.Join(r.runtimeRoot, "agent", "live"),
		filepath.Join(r.runtimeRoot, "agent", "live_by_pid"),
	} {
		if err := os.RemoveAll(legacyPath); err != nil {
			return fmt.Errorf("agent registry: remove legacy live index %s: %w", legacyPath, err)
		}
	}
	return nil
}

func (r *AgentRegistry) acquirePIDLock(pid int) error {
	lockPath := filepath.Join(r.pidLockDir, fmt.Sprintf("%d.agent.lock", pid))
	f, err := os.OpenFile(lockPath, os.O_CREATE|os.O_RDWR, agentRegistrySharedFileMode)
	if err != nil {
		return fmt.Errorf("open pid lock file %s: %w", lockPath, err)
	}
	if err := os.Chmod(lockPath, agentRegistrySharedFileMode); err != nil && !errors.Is(err, os.ErrPermission) {
		f.Close()
		return fmt.Errorf("chmod pid lock file %s: %w", lockPath, err)
	}
	if err := syscall.Flock(int(f.Fd()), syscall.LOCK_EX|syscall.LOCK_NB); err != nil {
		f.Close()
		return fmt.Errorf("acquire pid lock: %w", err)
	}
	r.mu.Lock()
	r.pidLockFile = f
	r.pidLockPath = lockPath
	r.mu.Unlock()
	return nil
}

func (r *AgentRegistry) releasePIDLock(path string) error {
	r.mu.Lock()
	f := r.pidLockFile
	r.pidLockFile = nil
	r.pidLockPath = ""
	r.mu.Unlock()
	if f != nil {
		_ = syscall.Flock(int(f.Fd()), syscall.LOCK_UN)
		if err := f.Close(); err != nil {
			return fmt.Errorf("close pid lock file %s: %w", path, err)
		}
	}
	if path == "" {
		return nil
	}
	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("remove pid lock file %s: %w", path, err)
	}
	return nil
}

func (r *AgentRegistry) isProcessAlive(pid int) (bool, error) {
	if pid <= 0 {
		return false, nil
	}
	if pid == r.pid {
		return true, nil
	}

	lockPath := filepath.Join(r.pidLockDir, fmt.Sprintf("%d.agent.lock", pid))
	if _, err := os.Stat(lockPath); err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return false, err
		}
		return pidAlive(pid), nil
	}
	held, err := r.isPIDLockHeld(lockPath)
	if err != nil {
		return false, err
	}
	return held, nil
}

func (r *AgentRegistry) isPIDLockHeld(lockPath string) (bool, error) {
	f, err := os.OpenFile(lockPath, os.O_RDWR, agentRegistrySharedFileMode)
	if err != nil {
		return false, fmt.Errorf("open pid lock file %s: %w", lockPath, err)
	}
	defer f.Close()
	err = syscall.Flock(int(f.Fd()), syscall.LOCK_EX|syscall.LOCK_NB)
	if err == nil {
		// lock not held by any live process
		if err := syscall.Flock(int(f.Fd()), syscall.LOCK_UN); err != nil {
			return false, fmt.Errorf("release temp pid lock %s: %w", lockPath, err)
		}
		return false, nil
	}
	if errors.Is(err, syscall.EWOULDBLOCK) || errors.Is(err, syscall.EAGAIN) {
		return true, nil
	}
	return false, err
}

func (r *AgentRegistry) pruneStalePIDLocks() error {
	entries, err := os.ReadDir(r.pidLockDir)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return err
	}
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		pid, ok := parsePIDFromLockName(entry.Name())
		if !ok {
			continue
		}
		lockPath := filepath.Join(r.pidLockDir, entry.Name())
		if pid == r.pid {
			continue
		}
		f, acquired, err := r.tryAcquirePIDLockFile(lockPath)
		if err != nil {
			return fmt.Errorf("agent registry: check stale pid lock %s: %w", lockPath, err)
		}
		if !acquired {
			continue
		}
		if err := r.removeArtifactsForStalePID(pid); err != nil {
			if f != nil {
				_ = syscall.Flock(int(f.Fd()), syscall.LOCK_UN)
				_ = f.Close()
				f = nil
			}
			return err
		}
		if err := r.removeRegistrationsForPID(pid); err != nil {
			if f != nil {
				_ = syscall.Flock(int(f.Fd()), syscall.LOCK_UN)
				_ = f.Close()
				f = nil
			}
			return err
		}
		if f != nil {
			_ = syscall.Flock(int(f.Fd()), syscall.LOCK_UN)
			if err := f.Close(); err != nil {
				return fmt.Errorf("agent registry: close stale pid lock %s: %w", lockPath, err)
			}
			f = nil
		}
		if err := os.Remove(lockPath); err != nil && !os.IsNotExist(err) {
			return fmt.Errorf("agent registry: remove stale pid lock %s: %w", lockPath, err)
		}
	}
	return nil
}

func (r *AgentRegistry) removeRegistrationsForPID(pid int) error {
	if pid <= 0 {
		return nil
	}
	entries, err := os.ReadDir(r.agentDir)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return err
	}
	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".agent" {
			continue
		}
		path := filepath.Join(r.agentDir, entry.Name())
		reg, ok := readAgentRegistration(path, strings.TrimSuffix(entry.Name(), ".agent"))
		if !ok || reg.PID != pid {
			continue
		}
		if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
			return fmt.Errorf("agent registry: remove stale agent registration %s: %w", path, err)
		}
	}
	return nil
}

func (r *AgentRegistry) removeArtifactsForStalePID(pid int) error {
	if strings.TrimSpace(r.runtimeRoot) == "" {
		return nil
	}
	pidLink := filepath.Join(r.runtimeRoot, "pid", strconv.Itoa(pid))
	target, err := os.Readlink(pidLink)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("agent registry: read pid link %s: %w", pidLink, err)
	}
	if err := os.Remove(pidLink); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("agent registry: remove stale pid link %s: %w", pidLink, err)
	}
	if err := r.removeAgentRootByPIDTarget(pidLink, target, pid); err != nil {
		return err
	}
	return nil
}

func (r *AgentRegistry) removeAgentRootByPIDTarget(pidLinkPath, target string, pid int) error {
	if strings.TrimSpace(target) == "" {
		return nil
	}
	if !filepath.IsAbs(target) {
		target = filepath.Clean(filepath.Join(filepath.Dir(pidLinkPath), target))
	}
	publicRoot := filepath.Clean(target)
	if filepath.Base(publicRoot) != "public" {
		return nil
	}
	agentRoot := filepath.Clean(filepath.Dir(publicRoot))
	if err := r.removeAgentRootIfMatchesPID(agentRoot, pid); err != nil {
		return err
	}
	return nil
}

func (r *AgentRegistry) removeAgentRootIfMatchesPID(agentRoot string, pid int) error {
	sessionPath := filepath.Join(agentRoot, "status", "session.json")
	data, err := os.ReadFile(sessionPath)
	if err != nil {
		return nil
	}
	var status map[string]any
	if err := json.Unmarshal(data, &status); err != nil {
		return nil
	}
	recorded, ok := status["pid"].(float64)
	if !ok || int(recorded) != pid {
		return nil
	}
	if err := os.RemoveAll(agentRoot); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("agent registry: remove stale agent root %s: %w", agentRoot, err)
	}
	return nil
}

func parsePIDFromLockName(name string) (int, bool) {
	if !strings.HasSuffix(name, ".agent.lock") {
		return 0, false
	}
	pidText := strings.TrimSuffix(name, ".agent.lock")
	pid, err := strconv.Atoi(pidText)
	if err != nil || pid <= 0 {
		return 0, false
	}
	return pid, true
}

func readAgentRegistration(path string, fallbackID string) (agentRegistration, bool) {
	data, err := os.ReadFile(path)
	if err != nil {
		return agentRegistration{}, false
	}
	trimmed := strings.TrimSpace(string(data))
	var reg agentRegistration
	if strings.HasPrefix(trimmed, "{") {
		if err := json.Unmarshal([]byte(trimmed), &reg); err != nil {
			return agentRegistration{}, false
		}
	} else {
		pid, err := strconv.Atoi(trimmed)
		if err != nil {
			return agentRegistration{}, false
		}
		reg = agentRegistration{SessionID: fallbackID, RunID: fallbackID, PID: pid}
	}
	if reg.SessionID == "" {
		reg.SessionID = fallbackID
	}
	if reg.RunID == "" {
		reg.RunID = fallbackID
	}
	if reg.PID <= 0 {
		return agentRegistration{}, false
	}
	return reg, true
}

func pidAlive(pid int) bool {
	if pid <= 0 {
		return false
	}
	err := syscall.Kill(pid, 0)
	return err == nil || errors.Is(err, syscall.EPERM)
}
