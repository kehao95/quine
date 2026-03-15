package tools

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/kehao95/quine/internal/config"
	"github.com/kehao95/quine/internal/tape"
)

const jobWrapperScript = `
job_pid="$$"
canon_dir="${QUINE_JOB_SESSION_DIR}/${job_pid}"

mkdir -p "$canon_dir"
printf '%s' "$QUINE_JOB_COMMAND" >"$canon_dir/cmd"

rm -f "$canon_dir/exit" "$canon_dir/exit.tmp"
mkfifo "$canon_dir/exit"
exec 9<>"$canon_dir/exit"

child_pid=""
term_requested=0
trap 'term_requested=1; if [ -n "$child_pid" ]; then kill -TERM "$child_pid" 2>/dev/null || true; fi' TERM INT HUP

cleanup_workspace_mounts() {
  if [ -z "${QUINE_WORKSPACE_ENABLED:-}" ] || [ "${QUINE_WORKSPACE_ENABLED:-}" != "1" ]; then
    return 0
  fi
  mount_root="${QUINE_WORKSPACE_MOUNT_BASE}/${job_pid}"
  merged_dir="${mount_root}/merged"
  umount "$QUINE_WORKSPACE_ROOT" 2>/dev/null || true
  umount "$merged_dir" 2>/dev/null || true
  rm -rf "$mount_root" 2>/dev/null || true
}

run_job_command() {
  if [ "${QUINE_WORKSPACE_ENABLED:-}" = "1" ]; then
    if [ "${QUINE_WORKSPACE_BACKEND:-overlay}" = "overlay" ]; then
      set -eu
      mount_root="${QUINE_WORKSPACE_MOUNT_BASE}/${job_pid}"
      merged_dir="${mount_root}/merged"
      work_dir="${mount_root}/work"
      mkdir -p "$merged_dir" "$work_dir"
      mount --make-rprivate / 2>/dev/null || mount --make-private /
      mount -t overlay overlay -o "lowerdir=${QUINE_WORKSPACE_ROOT},upperdir=${QUINE_WORKSPACE_UPPER},workdir=${work_dir}" "$merged_dir"
      mount --bind "$merged_dir" "$QUINE_WORKSPACE_ROOT"
    fi
    mkdir -p "$QUINE_WORKSPACE"
    cd "$QUINE_WORKSPACE"
    exec "$QUINE_JOB_SHELL" -c "$QUINE_JOB_COMMAND"
  fi
  exec "$QUINE_JOB_SHELL" -c "$QUINE_JOB_COMMAND"
}

if [ -n "${QUINE_JOB_STDIN_FILE:-}" ]; then
  run_job_command <"$QUINE_JOB_STDIN_FILE" >"$canon_dir/out.log" 2>"$canon_dir/err.log" &
else
  run_job_command >"$canon_dir/out.log" 2>"$canon_dir/err.log" &
fi
child_pid=$!

while :; do
  wait "$child_pid"
  status=$?
  if [ "$term_requested" -eq 0 ] || ! kill -0 "$child_pid" 2>/dev/null; then
    break
  fi
done

printf '%s\n' "$status" >"$canon_dir/exit.tmp"
printf '%s\n' "$status" >&9 || true
mv -f "$canon_dir/exit.tmp" "$canon_dir/exit"
exec 9>&-

if [ -n "${QUINE_JOB_STDIN_FILE:-}" ]; then
  rm -f "$QUINE_JOB_STDIN_FILE"
fi

cleanup_workspace_mounts

exit "$status"
`

// ShExecutor runs shell commands as ephemeral jobs backed by filesystem state.
//
// Every command gets a canonical job directory at:
//
//	${QUINE_DATA_DIR}/jobs/<session>/<pid>/
//
// Directory contents:
//   - cmd: original command string passed to sh
//   - out.log: child stdout
//   - err.log: child stderr
//   - exit: FIFO-then-regular-file completion status
//
// The special exit file starts life as a FIFO so `cat <job-path>/exit` blocks
// while the process runs; on completion the wrapper replaces it with a regular
// file containing the numeric exit status.
//
// Synchronous sh(command) calls remove their per-call job directory after the
// bounded result is constructed. Detached jobs retain their directory.
//
// TODO(quine): expose a stable subjective-world /job/<pid>/ namespace instead
// of returning raw runtime-root-backed absolute paths.
type ShExecutor struct {
	Shell     string
	MaxOutput int
	Env       []string
	Timeout   time.Duration

	// WorkDir is the working directory for shell commands.
	// If empty, defaults to the session's job directory.
	WorkDir string

	// Stdin is the material stdin file descriptor. Passed as fd 3
	// (ExtraFiles[0]) so the agent can read it via /dev/fd/3 or cat <&3.
	Stdin *os.File

	// Stdout is the deliverable stdout file descriptor. Passed as fd 4
	// (ExtraFiles[1]) so commands can write to >&4.
	Stdout *os.File

	// Stderr is the failure-signal file descriptor. Passed as fd 5
	// (ExtraFiles[2]) so commands can write to >&5.
	Stderr *os.File

	// ProcessStarted is called when a shell starts (for SIGINT forwarding).
	ProcessStarted func(*os.Process)
	// ProcessEnded is called when a shell exits.
	ProcessEnded func()

	DataDir    string // durable runtime-state root for jobs and related surfaces
	SessionID  string
	TurnID     int
	subjective *subjectiveFS

	mu       sync.Mutex
	detached map[int]*managedJob
}

type managedJob struct {
	ID           int
	cmd          *exec.Cmd
	detached     bool
	interactive  bool
	canonicalDir string
	displayDir   string
	outPath      string
	errPath      string
	exitPath     string
	doneCh       chan struct{}
	exitCode     int
	ptyState     *interactiveState
}

// NewShExecutor creates a ShExecutor from config with the given child
// environment. The childEnv slice should contain QUINE_* overrides (from
// Config.ChildEnv). These are merged with os.Environ() so that spawned
// commands inherit a complete environment with QUINE_* vars overlaid.
//
// QUINE_SESSION_ID and QUINE_TAPE_ID are stripped from both envs so each
// child quine process generates fresh per-process identity.
func NewShExecutor(cfg *config.Config, childEnv []string) *ShExecutor {
	// Filter out per-process identity from childEnv.
	filteredChildEnv := make([]string, 0, len(childEnv))
	for _, entry := range childEnv {
		if strings.HasPrefix(entry, "QUINE_SESSION_ID=") {
			continue
		}
		if strings.HasPrefix(entry, "QUINE_TAPE_ID=") {
			continue
		}
		filteredChildEnv = append(filteredChildEnv, entry)
	}

	// Filter out per-process identity from os.Environ() too.
	filteredOsEnv := make([]string, 0, len(os.Environ()))
	for _, entry := range os.Environ() {
		if strings.HasPrefix(entry, "QUINE_SESSION_ID=") {
			continue
		}
		if strings.HasPrefix(entry, "QUINE_TAPE_ID=") {
			continue
		}
		filteredOsEnv = append(filteredOsEnv, entry)
	}

	mergedEnv := MergeEnv(filteredOsEnv, filteredChildEnv)

	return &ShExecutor{
		Shell:      cfg.Shell,
		MaxOutput:  cfg.OutputTruncate,
		Env:        mergedEnv,
		Timeout:    time.Duration(cfg.ShTimeout) * time.Second,
		WorkDir:    cfg.WorkDir,
		DataDir:    cfg.DataDir,
		SessionID:  cfg.SessionID,
		subjective: newSubjectiveFS(cfg),
	}
}

// MergeEnv takes the OS environment and overlays child overrides.
// Keys from childEnv take precedence over osEnv.
func MergeEnv(osEnv []string, childEnv []string) []string {
	env := make(map[string]string, len(osEnv)+len(childEnv))
	order := make([]string, 0, len(osEnv)+len(childEnv))

	for _, entry := range osEnv {
		key, _, _ := strings.Cut(entry, "=")
		if _, exists := env[key]; !exists {
			order = append(order, key)
		}
		env[key] = entry
	}

	for _, entry := range childEnv {
		key, _, _ := strings.Cut(entry, "=")
		if _, exists := env[key]; !exists {
			order = append(order, key)
		}
		env[key] = entry
	}

	result := make([]string, 0, len(order))
	for _, key := range order {
		result = append(result, env[key])
	}
	return result
}

func (b *ShExecutor) initState() error {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.detached == nil {
		b.detached = make(map[int]*managedJob)
	}

	if b.Shell == "" {
		b.Shell = "/bin/sh"
	}
	if b.MaxOutput == 0 {
		b.MaxOutput = 20480
	}
	if b.DataDir == "" {
		tmpDir, err := os.MkdirTemp("", "quine-sh-data-*")
		if err != nil {
			return fmt.Errorf("creating temp data dir: %w", err)
		}
		b.DataDir = tmpDir
	}
	if b.SessionID == "" {
		b.SessionID = fmt.Sprintf("session-%d", time.Now().UnixNano())
	}
	if b.subjective != nil {
		if err := b.subjective.init(b.DataDir, b.SessionID); err != nil {
			return fmt.Errorf("initializing workspace physics: %w", err)
		}
	}
	return nil
}

func (b *ShExecutor) jobSessionRoot() (string, error) {
	dataDirAbs, err := filepath.Abs(b.DataDir)
	if err != nil {
		return "", fmt.Errorf("abs data dir: %w", err)
	}
	root := filepath.Join(dataDirAbs, "jobs", b.SessionID)
	if err := os.MkdirAll(root, 0o755); err != nil {
		return "", fmt.Errorf("creating session job root: %w", err)
	}
	return root, nil
}

func (b *ShExecutor) extraFiles() []*os.File {
	stdoutFD := b.Stdout
	if stdoutFD == nil {
		stdoutFD = os.Stdout
	}
	stdinFD := b.Stdin
	if stdinFD == nil {
		stdinFD = os.Stdin
	}
	stderrFD := b.Stderr
	if stderrFD == nil {
		stderrFD = os.Stderr
	}

	// Side channels follow fd+3 positional mapping:
	// runtime fd 0/1/2 -> child fd 3/4/5.
	return []*os.File{stdinFD, stdoutFD, stderrFD}
}

func (b *ShExecutor) startJob(command string, detach bool, stdin string) (*managedJob, error) {
	if err := b.initState(); err != nil {
		return nil, err
	}

	jobSessionRoot, err := b.jobSessionRoot()
	if err != nil {
		return nil, err
	}

	var stdinFile string
	if stdin != "" {
		f, err := os.CreateTemp(jobSessionRoot, "stdin-*")
		if err != nil {
			return nil, fmt.Errorf("creating stdin temp file: %w", err)
		}
		if _, err := f.WriteString(stdin); err != nil {
			f.Close()
			os.Remove(f.Name())
			return nil, fmt.Errorf("writing stdin temp file: %w", err)
		}
		if err := f.Close(); err != nil {
			os.Remove(f.Name())
			return nil, fmt.Errorf("closing stdin temp file: %w", err)
		}
		stdinFile = f.Name()
	}

	cmd := exec.Command(b.Shell, "-c", jobWrapperScript)
	useWorkspace := b.subjective != nil && b.subjective.enabled
	needsWorkspaceNamespace := useWorkspace && !b.subjective.directBackend()
	cmd.SysProcAttr = jobSysProcAttr(detach, needsWorkspaceNamespace)

	cmd.Env = MergeEnv(b.Env, []string{
		"QUINE_JOB_SHELL=" + b.Shell,
		"QUINE_JOB_COMMAND=" + command,
		"QUINE_JOB_SESSION_DIR=" + jobSessionRoot,
		"QUINE_JOB_STDIN_FILE=" + stdinFile,
	})
	// Use explicit WorkDir if set, otherwise fall back to the job session root
	if b.WorkDir != "" {
		cmd.Dir = b.WorkDir
	} else {
		cmd.Dir = jobSessionRoot
	}
	if useWorkspace {
		cmd.Env = MergeEnv(cmd.Env, b.subjective.commandEnv())
	}
	cmd.ExtraFiles = b.extraFiles()

	if err := cmd.Start(); err != nil {
		if stdinFile != "" {
			_ = os.Remove(stdinFile)
		}
		return nil, fmt.Errorf("starting process: %w", err)
	}

	pid := cmd.Process.Pid
	canonicalDir := filepath.Join(jobSessionRoot, strconv.Itoa(pid))
	job := &managedJob{
		ID:           pid,
		cmd:          cmd,
		detached:     detach,
		canonicalDir: canonicalDir,
		displayDir:   filepath.ToSlash(canonicalDir) + "/",
		outPath:      filepath.Join(canonicalDir, "out.log"),
		errPath:      filepath.Join(canonicalDir, "err.log"),
		exitPath:     filepath.Join(canonicalDir, "exit"),
		doneCh:       make(chan struct{}),
	}

	if err := b.waitForMaterialization(job); err != nil {
		_ = syscall.Kill(-job.ID, syscall.SIGKILL)
		_ = cmd.Wait()
		return nil, err
	}

	go b.awaitExit(job)
	return job, nil
}

func (b *ShExecutor) waitForMaterialization(job *managedJob) error {
	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		if _, err := os.Stat(job.exitPath); err == nil {
			return nil
		}
		if job.cmd.ProcessState != nil && job.cmd.ProcessState.Exited() {
			break
		}
		time.Sleep(10 * time.Millisecond)
	}
	return fmt.Errorf("job %d did not materialize job metadata", job.ID)
}

func (b *ShExecutor) awaitExit(job *managedJob) {
	code := exitCodeFromWait(job.cmd.Wait())
	job.exitCode = code
	close(job.doneCh)

	if job.detached {
		b.mu.Lock()
		delete(b.detached, job.ID)
		b.mu.Unlock()
	}
}

func exitCodeFromWait(err error) int {
	if err == nil {
		return 0
	}
	exitErr, ok := err.(*exec.ExitError)
	if !ok {
		return 1
	}
	if ws, ok := exitErr.Sys().(syscall.WaitStatus); ok {
		if ws.Exited() {
			return ws.ExitStatus()
		}
		if ws.Signaled() {
			return 128 + int(ws.Signal())
		}
	}
	code := exitErr.ExitCode()
	if code >= 0 {
		return code
	}
	return 1
}

func (b *ShExecutor) registerDetached(job *managedJob) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.detached[job.ID] = job
}

func (b *ShExecutor) readLog(path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	return string(data)
}

func (b *ShExecutor) Prepare() error {
	return b.initState()
}

// Execute dispatches a sh(command) call.
//
// timeout and outputLimit are currently ignored in the filesystem job model.
// interactive=true starts a PTY-backed interactive job and returns immediately
// with a filesystem control surface under the job directory.
func (b *ShExecutor) Execute(toolID string, command string, timeout time.Duration, outputLimit int, interactive bool, detach bool, stdin string) tape.ToolResult {
	_ = timeout
	_ = outputLimit

	if interactive {
		if stdin != "" {
			return tape.ToolResult{
				ToolID:  toolID,
				Content: "[SHELL ERROR] interactive=true cannot be combined with stdin",
				IsError: true,
			}
		}
		if detach {
			return tape.ToolResult{
				ToolID:  toolID,
				Content: "[SHELL ERROR] interactive=true already returns immediately; do not also set detach=true",
				IsError: true,
			}
		}
	}
	if (detach || interactive) && b.subjective != nil && b.subjective.enabled {
		return tape.ToolResult{
			ToolID:  toolID,
			Content: "[SHELL ERROR] interactive and detached jobs are not supported while Linux workspace physics are enabled",
			IsError: true,
		}
	}
	if b.subjective != nil && b.subjective.enabled {
		if err := b.initState(); err != nil {
			return tape.ToolResult{
				ToolID:  toolID,
				Content: fmt.Sprintf("[SHELL ERROR] %v", err),
				IsError: true,
			}
		}
	}

	beforeSnap, snapErr := fsSnapshot{}, error(nil)
	if b.subjective != nil && b.subjective.enabled {
		beforeSnap, snapErr = b.subjective.snapshot()
		if snapErr != nil {
			return tape.ToolResult{
				ToolID:  toolID,
				Content: fmt.Sprintf("[SHELL ERROR] snapshot before command: %v", snapErr),
				IsError: true,
			}
		}
	}

	var job *managedJob
	var err error
	if interactive {
		job, err = b.startInteractiveJob(command)
	} else {
		job, err = b.startJob(command, detach, stdin)
	}
	if err != nil {
		return tape.ToolResult{
			ToolID:  toolID,
			Content: fmt.Sprintf("[SHELL ERROR] %v", err),
			IsError: true,
		}
	}

	if interactive {
		b.registerDetached(job)
		if b.ProcessStarted != nil {
			b.ProcessStarted(job.cmd.Process)
		}
		if b.ProcessEnded != nil {
			b.ProcessEnded()
		}
		return tape.ToolResult{
			ToolID: toolID,
			Content: fmt.Sprintf(
				"[JOB] pid=%d path=%s (interactive)\nRead screen in `%sscreen.txt` or `%sscreen.png`.\nScreen metadata: `%sscreen.meta`\nInput keystrokes via `%sin`\nResize via `%swinsize`\nEvent log: `%sevents.log`\nWait for completion with `cat %sexit`.",
				job.ID, job.displayDir, job.displayDir, job.displayDir, job.displayDir, job.displayDir, job.displayDir, job.displayDir, job.displayDir,
			),
		}
	}

	if detach {
		b.registerDetached(job)
		if b.ProcessStarted != nil {
			b.ProcessStarted(job.cmd.Process)
		}
		if b.ProcessEnded != nil {
			b.ProcessEnded()
		}
		mutations := ""
		if b.subjective != nil && b.subjective.enabled {
			afterSnap, err := b.subjective.snapshot()
			if err == nil {
				mutations = b.subjective.formatMutations(beforeSnap, afterSnap)
			}
		}
		return tape.ToolResult{
			ToolID: toolID,
			Content: fmt.Sprintf(
				"[JOB] pid=%d path=%s (detached)\nYou can wait this job by `cat %sexit`\nSee output in `%sout.log` and `%serr.log`.%s",
				job.ID, job.displayDir, job.displayDir, job.displayDir, job.displayDir, optionalTrailingBlock(mutations),
			),
		}
	}

	if b.ProcessStarted != nil {
		b.ProcessStarted(job.cmd.Process)
	}
	<-job.doneCh
	if b.ProcessEnded != nil {
		b.ProcessEnded()
	}

	stdoutStr := b.applyOutputLimit(b.readLog(job.outPath))
	stderrStr := b.applyOutputLimit(b.readLog(job.errPath))
	mutations := ""
	worldRevisionBlock := ""
	if b.subjective != nil && b.subjective.enabled {
		afterSnap, err := b.subjective.snapshot()
		if err != nil {
			return tape.ToolResult{
				ToolID:  toolID,
				Content: fmt.Sprintf("[SHELL ERROR] snapshot after command: %v", err),
				IsError: true,
			}
		}
		mutations = b.subjective.formatMutations(beforeSnap, afterSnap)
		hasMutations := len(diffTree(beforeSnap.workspace, afterSnap.workspace)) > 0
		if hasMutations {
			if revision, err := b.subjective.captureWorldRevision("sh", b.TurnID); err != nil {
				return tape.ToolResult{
					ToolID:  toolID,
					Content: fmt.Sprintf("[SHELL ERROR] capture world revision after turn %d: %v", b.TurnID, err),
					IsError: true,
				}
			} else if revision.ID != "" {
				worldRevisionBlock = formatWorldRevisionCreated(revision, false)
			}
		} else if revision, err := b.subjective.loadCurrentWorldRevision(); err != nil {
			return tape.ToolResult{
				ToolID:  toolID,
				Content: fmt.Sprintf("[SHELL ERROR] load current world revision after turn %d: %v", b.TurnID, err),
				IsError: true,
			}
		} else if revision.ID != "" {
			worldRevisionBlock = formatWorldRevisionCreated(revision, true)
		}
	}

	// Cleanup job directory for synchronous (non-detached) jobs since
	// results have already been captured and returned inline.
	_ = os.RemoveAll(job.canonicalDir)

	content := fmt.Sprintf(
		"[EXIT CODE] %d\n[STDOUT]\n%s[STDERR]\n%s%s%s",
		job.exitCode, stdoutStr, stderrStr, optionalTrailingBlock(mutations), optionalTrailingBlock(worldRevisionBlock),
	)

	return tape.ToolResult{
		ToolID:  toolID,
		Content: content,
		IsError: job.exitCode != 0,
	}
}

func (b *ShExecutor) RestoreWorld(toolID string, revision string) tape.ToolResult {
	if strings.TrimSpace(revision) == "" {
		return tape.ToolResult{
			ToolID:  toolID,
			Content: "[RESTORE WORLD ERROR] revision is required",
			IsError: true,
		}
	}
	if b.subjective == nil || !b.subjective.enabled {
		return tape.ToolResult{
			ToolID:  toolID,
			Content: "[RESTORE WORLD ERROR] world restore requires Linux workspace physics",
			IsError: true,
		}
	}
	if !b.subjective.canRestoreWorld() {
		return tape.ToolResult{
			ToolID:  toolID,
			Content: "[RESTORE WORLD ERROR] restore_world requires workspace revision mode with restore support",
			IsError: true,
		}
	}
	if err := b.initState(); err != nil {
		return tape.ToolResult{
			ToolID:  toolID,
			Content: fmt.Sprintf("[RESTORE WORLD ERROR] %v", err),
			IsError: true,
		}
	}
	previous, current, err := b.subjective.restoreWorld(revision)
	if err != nil {
		return tape.ToolResult{
			ToolID:  toolID,
			Content: fmt.Sprintf("[RESTORE WORLD ERROR] %v", err),
			IsError: true,
		}
	}
	return tape.ToolResult{
		ToolID:  toolID,
		Content: fmt.Sprintf("[RESTORE WORLD] restored provisional workspace to revision %s%s", current, optionalTrailingBlock(formatWorldRevisionTransition(previous, current))),
	}
}

func (b *ShExecutor) CurrentWorldRevision() string {
	if b.subjective == nil {
		return ""
	}
	return b.subjective.currentWorldRevision()
}

func formatWorldRevisionCreated(revision worldRevision, unchanged bool) string {
	if revision.ID == "" {
		return ""
	}
	if unchanged {
		return fmt.Sprintf("[WORLD REVISION] current=%s (unchanged)", revision.ID)
	}
	if revision.Parent == "" {
		return fmt.Sprintf("[WORLD REVISION] current=%s", revision.ID)
	}
	return fmt.Sprintf("[WORLD REVISION] created=%s parent=%s current=%s", revision.ID, revision.Parent, revision.ID)
}

func formatWorldRevisionTransition(previous, current string) string {
	if strings.TrimSpace(current) == "" {
		return ""
	}
	if strings.TrimSpace(previous) == "" || previous == current {
		return fmt.Sprintf("[WORLD REVISION] current=%s", current)
	}
	return fmt.Sprintf("[WORLD REVISION] %s -> %s", previous, current)
}

func optionalTrailingBlock(block string) string {
	if strings.TrimSpace(block) == "" {
		return ""
	}
	return "\n" + block
}

// applyOutputLimit truncates with a visible notice for tool results.
func (b *ShExecutor) applyOutputLimit(s string) string {
	if b.MaxOutput <= 0 || len(s) <= b.MaxOutput {
		return s
	}
	total := len(s)
	return s[:b.MaxOutput] + fmt.Sprintf("\n...[Output Truncated, %d bytes total] Increase QUINE_OUTPUT_TRUNCATE to capture more in the tool result.", total)
}

// Close kills active detached jobs unless they should survive success.
func (b *ShExecutor) Close(keepDetached bool) error {
	if err := b.initState(); err != nil {
		return err
	}

	b.mu.Lock()
	jobs := make([]*managedJob, 0, len(b.detached))
	for _, job := range b.detached {
		jobs = append(jobs, job)
	}
	b.mu.Unlock()

	if !keepDetached {
		for _, job := range jobs {
			_ = syscall.Kill(-job.ID, syscall.SIGKILL)
		}
	}
	if keepDetached && b.subjective != nil && b.subjective.enabled {
		if err := b.subjective.commit(); err != nil {
			return err
		}
	}
	if !keepDetached && b.subjective != nil && b.subjective.enabled {
		if err := b.subjective.rollback(); err != nil {
			return err
		}
	}
	return nil
}

func (b *ShExecutor) ChildEnvOverrides() []string {
	if b.subjective == nil {
		return nil
	}
	return b.subjective.childEnvOverrides()
}

// ToInt converts various numeric types from JSON unmarshalling to int.
func ToInt(v any) int {
	switch n := v.(type) {
	case int:
		return n
	case int64:
		return int(n)
	case float64:
		return int(n)
	case float32:
		return int(n)
	}
	return 0
}

// toDuration converts an int (seconds) to time.Duration.
func toDuration(v any) time.Duration {
	secs := ToInt(v)
	if secs <= 0 {
		return 0
	}
	return time.Duration(secs) * time.Second
}
