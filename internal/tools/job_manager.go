package tools

import (
	"fmt"
	"os"
	"sync"
	"time"
)

// JobManager tracks all active jobs by pgid.
// The map IS ownership — presence means the job is managed here.
type JobManager struct {
	mu   sync.RWMutex
	jobs map[int]*Job

	// Configuration for spawning new jobs
	shell      string
	env        []string
	extraFiles []*os.File // [fd3 deliverable stdout, fd4 material stdin]
}

// NewJobManager creates a new JobManager.
func NewJobManager(shell string, env []string, extraFiles []*os.File) *JobManager {
	return &JobManager{
		jobs:       make(map[int]*Job),
		shell:      shell,
		env:        env,
		extraFiles: extraFiles,
	}
}

// Launch spawns a new job and registers it.
func (m *JobManager) Launch(command string, timeout time.Duration, outputLimit int, interactive bool) (*Job, error) {
	j, err := newJob(m.shell, command, m.env, m.extraFiles, timeout, outputLimit, interactive)
	if err != nil {
		return nil, err
	}

	m.mu.Lock()
	m.jobs[j.ID] = j
	m.mu.Unlock()

	return j, nil
}

// Get returns a job by pgid, or nil if not found.
func (m *JobManager) Get(pgid int) *Job {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.jobs[pgid]
}

// Remove removes a job from the registry (does not kill it).
func (m *JobManager) Remove(pgid int) {
	m.mu.Lock()
	delete(m.jobs, pgid)
	m.mu.Unlock()
}

// KillAll kills all registered jobs and clears the registry.
func (m *JobManager) KillAll() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	var firstErr error
	for pgid, j := range m.jobs {
		j.Kill()
		delete(m.jobs, pgid)
	}
	return firstErr
}

// Count returns the number of registered jobs.
func (m *JobManager) Count() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.jobs)
}

// --- High-level dispatch helpers used by tools.go ---

// RunResult holds the result of a completed or paused job.
type RunResult struct {
	Paused      bool
	JobID       int
	ExitCode    int
	Stdout      string
	Stderr      string
	StdoutBytes int // total bytes in stdout buffer at pause time
}

// RunSync launches a job and blocks until it completes or pauses.
// On pause it leaves the job registered; on completion it removes it.
func (m *JobManager) RunSync(command string, timeout time.Duration, outputLimit int, interactive bool) (RunResult, error) {
	j, err := m.Launch(command, timeout, outputLimit, interactive)
	if err != nil {
		return RunResult{}, fmt.Errorf("launching job: %w", err)
	}

	paused := j.Wait()

	if paused {
		// Job is paused — leave in registry for agent to resume/kill.
		return RunResult{
			Paused:      true,
			JobID:       j.ID,
			Stdout:      j.ReadStdout(),
			Stderr:      j.ReadStderr(),
			StdoutBytes: j.stdout.Len(),
		}, nil
	}

	// Natural completion — remove from registry.
	m.Remove(j.ID)
	return RunResult{
		Paused:   false,
		JobID:    j.ID,
		ExitCode: j.ExitCode(),
		Stdout:   j.ReadStdout(),
		Stderr:   j.ReadStderr(),
	}, nil
}
