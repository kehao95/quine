package tools

import (
	"encoding/json"
	"errors"
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
wait_for_job_surface() {
  i=0
  while [ ! -d "$canon_dir" ]; do
    i=$((i + 1))
    if [ "$i" -ge 500 ]; then
      printf '%s\n' "job metadata not initialized at $canon_dir" >&2
      exit 97
    fi
    sleep 0.01
  done
}

wait_for_job_surface

child_pid=""
term_requested=0
trap 'term_requested=1; if [ -n "$child_pid" ]; then kill -TERM "$child_pid" 2>/dev/null || true; fi' TERM INT HUP

cleanup_workspace_mounts() {
  if [ -z "${QUINE_WORKSPACE_ENABLED:-}" ] || [ "${QUINE_WORKSPACE_ENABLED:-}" != "1" ]; then
    return 0
  fi
  mount_root="${QUINE_WORKSPACE_MOUNT_BASE}/${job_pid}"
  merged_dir="${mount_root}/merged"
  workspace_root="${QUINE_WORKSPACE_ROOT}"
  if [ "${QUINE_WORKSPACE_OVERLAY_DRIVER:-kernel}" = "fuse" ]; then
    cd / 2>/dev/null || true
    umount -l "$workspace_root" 2>/dev/null || true
    if command -v fusermount3 >/dev/null 2>&1; then
      fusermount3 -u "$merged_dir" 2>/dev/null || fusermount3 -uz "$merged_dir" 2>/dev/null || true
    elif command -v fusermount >/dev/null 2>&1; then
      fusermount -u "$merged_dir" 2>/dev/null || fusermount -uz "$merged_dir" 2>/dev/null || true
    fi
    umount -l "$merged_dir" 2>/dev/null || true
    rm -rf "$mount_root" 2>/dev/null || true
    return 0
  fi
  if command -v python3 >/dev/null 2>&1; then
    python3 - "$mount_root" "$workspace_root" "$merged_dir" <<'PY' >/dev/null 2>&1 || true
import os
import shutil
import subprocess
import sys

mount_root, workspace_root, merged_dir = sys.argv[1:4]

def safe_run(cmd):
    try:
        subprocess.run(cmd, stdout=subprocess.DEVNULL, stderr=subprocess.DEVNULL, timeout=5, check=False)
    except Exception:
        pass

pid = os.fork()
if pid:
    os.waitpid(pid, 0)
    raise SystemExit(0)

os.setsid()

pid = os.fork()
if pid:
    os._exit(0)

try:
    os.chdir("/")
except OSError:
    pass

devnull = os.open("/dev/null", os.O_RDWR)
for fd in (0, 1, 2):
    try:
        os.dup2(devnull, fd)
    except OSError:
        pass
if devnull > 2:
    os.close(devnull)

for target in (workspace_root, merged_dir, workspace_root):
    safe_run(["umount", "-l", target])

shutil.rmtree(mount_root, ignore_errors=True)
os._exit(0)
PY
  else
    /bin/sh -c '
      mount_root="$1"
      workspace_root="$2"
      merged_dir="$3"
      (
        trap "" HUP INT TERM
        cd / 2>/dev/null || true
        if command -v timeout >/dev/null 2>&1; then
          timeout 5 /bin/sh -c '"'"'
            workspace_root="$1"
            merged_dir="$2"
            umount -l "$workspace_root" 2>/dev/null || true
            umount -l "$merged_dir" 2>/dev/null || true
            umount -l "$workspace_root" 2>/dev/null || true
          '"'"' sh "$workspace_root" "$merged_dir" >/dev/null 2>&1 || true
        else
          umount -l "$workspace_root" 2>/dev/null || true
          umount -l "$merged_dir" 2>/dev/null || true
          umount -l "$workspace_root" 2>/dev/null || true
        fi
        rm -rf "$mount_root" 2>/dev/null || true
      ) </dev/null >/dev/null 2>&1 &
    ' sh "$mount_root" "$workspace_root" "$merged_dir" >/dev/null 2>&1 || true
  fi
}

enable_isolated_loopback() {
  if command -v ip >/dev/null 2>&1; then
    ip link set lo up
    return $?
  fi
  if command -v ifconfig >/dev/null 2>&1; then
    ifconfig lo up
    return $?
  fi
  if command -v python3 >/dev/null 2>&1; then
    python3 - <<'PY'
import os
import socket
import struct
import sys

RTM_NEWLINK = 16
NLM_F_REQUEST = 1
NLM_F_ACK = 4
IFF_UP = 1

try:
    index = socket.if_nametoindex("lo")
    sock = socket.socket(socket.AF_NETLINK, socket.SOCK_RAW, socket.NETLINK_ROUTE)
    ifi = struct.pack("BBHiII", socket.AF_UNSPEC, 0, 0, index, IFF_UP, IFF_UP)
    hdr = struct.pack("IHHII", 16 + len(ifi), RTM_NEWLINK, NLM_F_REQUEST | NLM_F_ACK, 1, os.getpid())
    sock.send(hdr + ifi)
    data = sock.recv(65535)
    if len(data) >= 20:
        error = struct.unpack("i", data[16:20])[0]
        if error != 0:
            raise OSError(-error, os.strerror(-error))
except Exception as exc:
    print(f"failed to enable isolated loopback: {exc}", file=sys.stderr)
    raise SystemExit(1)
PY
    return $?
  fi
  printf '%s\n' "QUINE_JOB_NETWORK=none requires ip, ifconfig, or python3 to enable loopback" >&2
  return 1
}

enter_workspace() {
  if [ "${QUINE_WORKSPACE_ENABLED:-}" = "1" ]; then
    if [ -n "${QUINE_WORKSPACE_INIT_ERROR:-}" ]; then
      printf '%s\n' "$QUINE_WORKSPACE_INIT_ERROR" >&2
      exit 1
    fi
    set -eu
    mount --make-rslave /
    mount_root="${QUINE_WORKSPACE_MOUNT_BASE}/${job_pid}"
    merged_dir="${mount_root}/merged"
    mkdir -p "$merged_dir"
    mount --bind "$QUINE_WORKSPACE_ROOT" "$QUINE_WORKSPACE_ROOT"
    mount --make-private "$QUINE_WORKSPACE_ROOT"
    overlay_options="lowerdir=${QUINE_WORKSPACE_LOWERDIR},upperdir=${QUINE_WORKSPACE_UPPER},workdir=${QUINE_WORKSPACE_WORKDIR}"
    if [ -n "${QUINE_WORKSPACE_OVERLAY_EXTRA_OPTS:-}" ]; then
      overlay_options="${overlay_options},${QUINE_WORKSPACE_OVERLAY_EXTRA_OPTS}"
    fi
    if [ "${QUINE_WORKSPACE_OVERLAY_DRIVER:-kernel}" = "fuse" ]; then
      if ! command -v fuse-overlayfs >/dev/null 2>&1; then
        printf '%s\n' "QUINE_WORKSPACE_OVERLAY_DRIVER=fuse requires fuse-overlayfs" >&2
        exit 1
      fi
      fuse_stderr="${mount_root}/fuse-overlayfs.stderr"
      fuse-overlayfs -f -o "$overlay_options" "$merged_dir" >/dev/null 2>"$fuse_stderr" &
      fuse_pid=$!
      merged_real="$(readlink -f "$merged_dir" 2>/dev/null || printf '%s' "$merged_dir")"
      i=0
      while ! mountpoint -q "$merged_dir" 2>/dev/null && ! grep -qs " $merged_dir " /proc/self/mounts && ! grep -qs " $merged_real " /proc/self/mounts; do
        if ! kill -0 "$fuse_pid" 2>/dev/null; then
          cat "$fuse_stderr" >&2 || true
          printf '%s\n' "fuse-overlayfs mount failed" >&2
          exit 1
        fi
        i=$((i + 1))
        if [ "$i" -ge 500 ]; then
          cat "$fuse_stderr" >&2 || true
          printf '%s\n' "fuse-overlayfs mount timed out" >&2
          exit 1
        fi
        sleep 0.01
      done
    else
      mount -t overlay overlay -o "$overlay_options" "$merged_dir"
    fi
    mount --bind "$merged_dir" "$QUINE_WORKSPACE_ROOT"
    mkdir -p "$QUINE_WORKSPACE"
    cd "$QUINE_WORKSPACE"
  fi
}

run_job_command() {
  if [ "${QUINE_JOB_NETWORK:-host}" = "none" ]; then
    enable_isolated_loopback
  fi
  enter_workspace || return "$?"
  exec "$QUINE_JOB_SHELL" -c "$QUINE_JOB_COMMAND"
}

run_interactive_job_command() {
  if [ "${QUINE_JOB_NETWORK:-host}" = "none" ]; then
    enable_isolated_loopback
  fi
  enter_workspace || return "$?"
  set +e
  "$QUINE_JOB_SHELL" -c "$QUINE_JOB_COMMAND"
  status=$?
  cleanup_workspace_mounts
  return "$status"
}

if [ "${QUINE_JOB_INTERACTIVE:-}" = "1" ]; then
  run_interactive_job_command
  exit "$?"
elif [ -n "${QUINE_JOB_STDIN_FILE:-}" ]; then
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
//   - pid: process-group leader PID
//   - started_at: UTC RFC3339Nano start timestamp
//   - out.log: child stdout
//   - err.log: child stderr
//   - exit: regular-file completion status, present only after the runtime
//     records a terminal outcome
//
// Synchronous sh(command) calls remove their per-call job directory after the
// bounded result is constructed. Detached jobs and timeout-interrupted sync
// jobs retain their directory.
//
// TODO(quine): expose a stable subjective-world /job/<pid>/ namespace instead
// of returning raw runtime-root-backed absolute paths.
type ShExecutor struct {
	Shell     string
	MaxOutput int
	// Env is the base environment of every sh child, as of construction.
	//
	// It is a snapshot, and a snapshot is not enough: the agent edits
	// config/env/override with an ordinary shell write, and the very next
	// command must see it. commandBaseEnv() rebuilds this from cfg at each
	// command, and every command-construction site calls that, not this field.
	// Env remains the value for executors built without a cfg (tests) and the
	// fallback when the override on disk is unreadable.
	Env     []string
	Timeout time.Duration
	Network string

	// cfg is retained solely so the boundary can be rebuilt per command. A
	// stale child-env policy is a design violation, not a performance
	// trade-off: reading one small file per sh call is the price of the
	// override meaning what it says.
	cfg *config.Config

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

	DataDir                    string // durable runtime-state root for jobs and related surfaces
	SessionID                  string
	RunID                      string
	TurnID                     int
	FSMutationTelemetryEnabled bool // expose fs_mutations telemetry in tool results
	subjective                 *subjectiveFS

	mu       sync.Mutex
	detached map[int]*managedJob
	pending  map[int]*managedJob
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
	workspace    *subjectiveFS
}

type shellStructuredJob struct {
	PID              int    `json:"pid"`
	Path             string `json:"path"`
	Interactive      bool   `json:"interactive,omitempty"`
	Detached         bool   `json:"detached,omitempty"`
	WorkspaceSession string `json:"workspace_session,omitempty"`
	Adoptable        bool   `json:"adoptable,omitempty"`
}

type shellStructuredResult struct {
	Tool             string              `json:"tool"`
	Mode             string              `json:"mode"`
	Status           string              `json:"status"`
	Job              *shellStructuredJob `json:"job,omitempty"`
	ExitCode         *int                `json:"exit_code,omitempty"`
	Stdout           string              `json:"stdout,omitempty"`
	Stderr           string              `json:"stderr,omitempty"`
	StdoutSoFar      string              `json:"stdout_so_far,omitempty"`
	StderrSoFar      string              `json:"stderr_so_far,omitempty"`
	FSMutations      string              `json:"fs_mutations,omitempty"`
	FSMutationsSoFar string              `json:"fs_mutations_so_far,omitempty"`
	WorldRevision    string              `json:"world_revision,omitempty"`
	Cause            string              `json:"cause,omitempty"`
	TimeoutSeconds   *int                `json:"timeout_seconds,omitempty"`
	Error            string              `json:"error,omitempty"`
}

type switchWorldStructuredResult struct {
	Tool          string `json:"tool"`
	Status        string `json:"status"`
	Target        string `json:"target,omitempty"`
	Revision      string `json:"revision,omitempty"`
	FSMutations   string `json:"fs_mutations,omitempty"`
	WorldRevision string `json:"world_revision,omitempty"`
	Error         string `json:"error,omitempty"`
}

func marshalStructuredContent(v any) json.RawMessage {
	data, err := json.Marshal(v)
	if err != nil {
		return nil
	}
	return json.RawMessage(data)
}

func shellErrorResult(toolID string, message string) tape.ToolResult {
	return tape.ToolResult{
		ToolID: toolID,
		Content: tape.MarshalToolResultContent(shellStructuredResult{
			Tool:   "sh",
			Mode:   "sync",
			Status: "error",
			Error:  message,
		}),
		IsError: true,
	}
}

// NewShExecutor creates a ShExecutor from config.
//
// An sh child inherits this process's environment minus the runtime-owned
// names, with the agent's config/env/override applied on top, plus the few
// facts it cannot derive for itself (config.ShellStamps: run id, agent root,
// resolved runtime root). It does NOT receive a synthesized copy of this
// process's policy: a knob nobody set is absent from a shell child exactly as
// it is absent here, and `./quine` launched from a shell is a new founder, not
// a marked descendant.
//
// The old code hand-rolled two strip loops over childEnv and os.Environ(); the
// shared mask in config.BoundaryBehavior subsumes both, and it is derived from
// the registry rather than kept by hand, so a new runtime-emitted knob cannot
// forget to be masked.
func NewShExecutor(cfg *config.Config) *ShExecutor {
	return &ShExecutor{
		Shell:                      cfg.Shell,
		MaxOutput:                  cfg.OutputTruncate,
		cfg:                        cfg,
		Env:                        shBoundaryEnv(cfg, nil),
		Timeout:                    time.Duration(cfg.ShTimeout) * time.Second,
		Network:                    cfg.ShNetwork,
		WorkDir:                    cfg.WorkDir,
		DataDir:                    cfg.DataDir,
		SessionID:                  cfg.SessionID,
		RunID:                      cfg.RunID,
		FSMutationTelemetryEnabled: cfg.FSMutationTelemetryEnabled(),
		subjective:                 newSubjectiveFS(cfg),
	}
}

// shBoundaryEnv builds the sh-boundary environment for cfg.
//
// The default git identity is a BASE layer, beneath everything else: Quine is
// the sole author of the work it does in its workspace, so every shell gets a
// usable git profile out of the box rather than having to run `git config`
// before its first commit — but a real GIT_* value in the inherited environment
// or the override still wins.
//
// A malformed override is reported through onError and then ignored: an
// unparseable policy file must not take the shell down with it. Enforcement is
// unaffected — an illegal line was never going to be applied.
func shBoundaryEnv(cfg *config.Config, onError func(error)) []string {
	if cfg == nil {
		return defaultGitIdentityEnv()
	}
	override, err := config.ReadEnvOverride(cfg.EnvOverridePath())
	if err != nil && onError != nil {
		onError(err)
	}
	env := config.BuildChildEnv(config.BoundaryShell, os.Environ(), override, cfg.ShellStamps())
	return MergeEnv(defaultGitIdentityEnv(), env)
}

// commandBaseEnv is the environment every sh command starts from. It re-reads
// config/env/override at each call: the agent's child-env policy is live, and a
// policy that took effect one command late would be a lie the surface tells
// about itself.
func (b *ShExecutor) commandBaseEnv() []string {
	if b == nil {
		return nil
	}
	if b.cfg == nil {
		return b.Env
	}
	return shBoundaryEnv(b.cfg, func(err error) {
		// Best-effort visibility; the sh path has no logger of its own. The
		// override surface's own read-back reports the violations in full.
		fmt.Fprintf(os.Stderr, "[quine] child-env override ignored: %v\n", err)
	})
}

// jobWrapperEnv is the QUINE_JOB_* plumbing the shell job wrapper reads back:
// the shell it runs under, the command it was given, its session dir, and its
// network regime. Runtime-emitted per job — masked from inheritance, stamped
// here, where the job is actually known. Extra names (stdin file, interactive
// marker) differ per job kind and are passed by the caller.
func jobWrapperEnv(shell, command, sessionDir, network string, extra ...string) []string {
	env := []string{
		"QUINE_JOB_SHELL=" + shell,
		"QUINE_JOB_COMMAND=" + command,
		"QUINE_JOB_SESSION_DIR=" + sessionDir,
		"QUINE_JOB_NETWORK=" + network,
	}
	return append(env, extra...)
}

// defaultGitIdentityEnv returns the default git author/committer identity for
// shell commands. Setting both author and committer covers commits made without
// any repo-local user.name/user.email configured. Callers layer this beneath the
// inherited and child environments so an explicit override always takes priority.
func defaultGitIdentityEnv() []string {
	const (
		name  = "Quine"
		email = "quine@kehao.me"
	)
	return []string{
		"GIT_AUTHOR_NAME=" + name,
		"GIT_AUTHOR_EMAIL=" + email,
		"GIT_COMMITTER_NAME=" + name,
		"GIT_COMMITTER_EMAIL=" + email,
	}
}

// ReadWorkspaceFile reads a path from the agent's current workspace view when
// workspace physics are enabled, falling back to the host-visible filesystem for
// paths outside that workspace.
func (b *ShExecutor) ReadWorkspaceFile(path string) ([]byte, error) {
	if b != nil && b.subjective != nil && b.subjective.enabled {
		data, err := b.subjective.readWorkspaceFile(path)
		if err == nil {
			return data, nil
		}
		if !errorsIsNotExist(err) {
			return nil, err
		}
	}
	return os.ReadFile(path)
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
	if b.pending == nil {
		b.pending = make(map[int]*managedJob)
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

func writeJobIdentityFiles(dir string, command string, pid int, startedAt time.Time) error {
	if err := os.WriteFile(filepath.Join(dir, "cmd"), []byte(command), 0o644); err != nil {
		return fmt.Errorf("writing cmd file: %w", err)
	}
	if err := os.WriteFile(filepath.Join(dir, "pid"), []byte(fmt.Sprintf("%d\n", pid)), 0o644); err != nil {
		return fmt.Errorf("writing pid file: %w", err)
	}
	startedAtText := startedAt.UTC().Format(time.RFC3339Nano) + "\n"
	if err := os.WriteFile(filepath.Join(dir, "started_at"), []byte(startedAtText), 0o644); err != nil {
		return fmt.Errorf("writing started_at file: %w", err)
	}
	return nil
}

func touchFile(path string) error {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	return f.Close()
}

func stageJobSurface(root string, pid int, init func(string) error) (string, error) {
	stageDir, err := os.MkdirTemp(root, fmt.Sprintf(".%d-", pid))
	if err != nil {
		return "", fmt.Errorf("creating staged job dir: %w", err)
	}
	cleanup := func() {
		_ = os.RemoveAll(stageDir)
	}
	if err := init(stageDir); err != nil {
		cleanup()
		return "", err
	}
	finalDir := filepath.Join(root, strconv.Itoa(pid))
	if err := os.Rename(stageDir, finalDir); err != nil {
		cleanup()
		return "", fmt.Errorf("publishing job dir: %w", err)
	}
	return finalDir, nil
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
	useOverlayWorkspace := b.subjective != nil && b.subjective.enabled && b.subjective.usesOverlayBackend()
	isolateNetwork := b.Network == "none"
	cmd.SysProcAttr = jobSysProcAttr(detach, useOverlayWorkspace, isolateNetwork)

	cmd.Env = MergeEnv(b.commandBaseEnv(), jobWrapperEnv(b.Shell, command, jobSessionRoot, b.Network,
		"QUINE_JOB_STDIN_FILE="+stdinFile,
	))
	// Under workspace physics, the wrapper is responsible for entering
	// QUINE_WORKSPACE after the subjective view is mounted. Starting the
	// outer process in a narrower workspace path can fail before the mount
	// exists (for example a child workspace that exists only in overlay state).
	if useOverlayWorkspace {
		cmd.Dir = jobSessionRoot
	} else if b.WorkDir != "" {
		// Use explicit WorkDir if set, otherwise fall back to the job session root
		cmd.Dir = b.WorkDir
	} else if b.subjective != nil && b.subjective.enabled && b.subjective.usesDirectBackend() && strings.TrimSpace(b.subjective.workspace) != "" {
		cmd.Dir = b.subjective.workspace
	} else {
		cmd.Dir = jobSessionRoot
	}
	if useOverlayWorkspace {
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
	startedAt := time.Now().UTC()
	cleanupStartFailure := func() {
		_ = syscall.Kill(-pid, syscall.SIGKILL)
		_ = cmd.Wait()
		if stdinFile != "" {
			_ = os.Remove(stdinFile)
		}
	}
	canonicalDir, err := stageJobSurface(jobSessionRoot, pid, func(dir string) error {
		if err := writeJobIdentityFiles(dir, command, pid, startedAt); err != nil {
			return err
		}
		if err := touchFile(filepath.Join(dir, "out.log")); err != nil {
			return fmt.Errorf("initializing out.log: %w", err)
		}
		if err := touchFile(filepath.Join(dir, "err.log")); err != nil {
			return fmt.Errorf("initializing err.log: %w", err)
		}
		return nil
	})
	if err != nil {
		cleanupStartFailure()
		return nil, err
	}
	exitPath := filepath.Join(canonicalDir, "exit")
	job := &managedJob{
		ID:           pid,
		cmd:          cmd,
		detached:     detach,
		canonicalDir: canonicalDir,
		displayDir:   filepath.ToSlash(canonicalDir) + "/",
		outPath:      filepath.Join(canonicalDir, "out.log"),
		errPath:      filepath.Join(canonicalDir, "err.log"),
		exitPath:     exitPath,
		doneCh:       make(chan struct{}),
	}

	go b.awaitExit(job)
	return job, nil
}

func (b *ShExecutor) awaitExit(job *managedJob) {
	code := exitCodeFromWait(job.cmd.Wait())
	job.exitCode = code
	if err := atomicWriteFile(job.exitPath, []byte(fmt.Sprintf("%d\n", code)), 0o644); err != nil {
		publishExitError(job, fmt.Sprintf("failed to publish exit code %d: %v", code, err))
	}
	close(job.doneCh)

	b.mu.Lock()
	delete(b.detached, job.ID)
	delete(b.pending, job.ID)
	b.mu.Unlock()
}

// publishExitError records that a managed job terminated but its terminal exit
// state could not be written to the canonical `exit` file (e.g. ENOSPC, or the
// job dir was removed). Pollers learn a job finished by reading `<dir>/exit`;
// silently dropping the write would leave a completed job looking perpetually
// in-flight with its status lost. Surface it via an `exit_error` sibling and a
// stderr breadcrumb instead.
func publishExitError(job *managedJob, msg string) {
	_ = atomicWriteFile(filepath.Join(job.canonicalDir, "exit_error"), []byte(msg+"\n"), 0o644)
	fmt.Fprintf(os.Stderr, "sh job %d: %s\n", job.ID, msg)
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

func (b *ShExecutor) registerPendingIfRunning(job *managedJob) bool {
	if job == nil {
		return false
	}
	select {
	case <-job.doneCh:
		return false
	default:
	}
	b.mu.Lock()
	defer b.mu.Unlock()
	select {
	case <-job.doneCh:
		return false
	default:
	}
	b.pending[job.ID] = job
	return true
}

func (b *ShExecutor) readLog(path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	return string(data)
}

func (b *ShExecutor) effectiveTimeout(requested time.Duration) time.Duration {
	if requested > 0 {
		return requested
	}
	return b.Timeout
}

func (b *ShExecutor) killSyncJobOnTimeout() bool {
	return b.subjective != nil && b.subjective.enabled && b.subjective.usesOverlayBackend()
}

func signalManagedJob(job *managedJob, sig syscall.Signal) error {
	if job == nil || job.cmd == nil || job.cmd.Process == nil || job.ID <= 0 {
		return nil
	}
	if err := syscall.Kill(-job.ID, sig); err == nil || err == syscall.ESRCH {
		return nil
	}
	if err := job.cmd.Process.Signal(sig); err != nil && !errors.Is(err, os.ErrProcessDone) {
		return err
	}
	return nil
}

type syncJobOutcome string

const (
	syncJobCompleted   syncJobOutcome = "completed"
	syncJobInterrupted syncJobOutcome = "interrupted"
)

func (b *ShExecutor) waitSyncJob(job *managedJob, timeout time.Duration, killOnTimeout bool) (syncJobOutcome, error) {
	if timeout <= 0 {
		<-job.doneCh
		return syncJobCompleted, nil
	}

	timer := time.NewTimer(timeout)
	defer timer.Stop()

	select {
	case <-job.doneCh:
		return syncJobCompleted, nil
	case <-timer.C:
		if killOnTimeout {
			if err := signalManagedJob(job, syscall.SIGTERM); err != nil {
				return "", err
			}
			select {
			case <-job.doneCh:
				return syncJobInterrupted, nil
			case <-time.After(time.Second):
				if err := signalManagedJob(job, syscall.SIGKILL); err != nil {
					return "", err
				}
			}
			select {
			case <-job.doneCh:
				return syncJobInterrupted, nil
			case <-time.After(5 * time.Second):
				return "", fmt.Errorf("timed out shell job %d did not exit after SIGKILL", job.ID)
			}
		}
		if err := signalManagedJob(job, syscall.SIGSTOP); err != nil {
			return "", err
		}
	}

	select {
	case <-job.doneCh:
		return syncJobCompleted, nil
	default:
	}
	if !b.registerPendingIfRunning(job) {
		return syncJobCompleted, nil
	}
	return syncJobInterrupted, nil
}

func (b *ShExecutor) Prepare() error {
	return b.initState()
}

// Execute dispatches a sh(command) call.
//
// outputLimit is currently ignored in the filesystem job model.
// interactive=true starts a PTY-backed interactive job and returns immediately
// with a filesystem control surface under the job directory.
func (b *ShExecutor) Execute(toolID string, command string, timeout time.Duration, outputLimit int, interactive bool, detach bool, stdin string) tape.ToolResult {
	_ = outputLimit

	if interactive {
		if stdin != "" {
			return shellErrorResult(toolID, "[SHELL ERROR] interactive=true cannot be combined with stdin")
		}
		if detach {
			return shellErrorResult(toolID, "[SHELL ERROR] interactive=true already returns immediately; do not also set detach=true")
		}
	}
	if detach && b.subjective != nil && b.subjective.enabled && b.subjective.usesOverlayBackend() {
		return shellErrorResult(toolID, "[SHELL ERROR] detached jobs are not supported while Linux workspace physics are enabled")
	}
	if b.subjective != nil && b.subjective.enabled {
		if err := b.initState(); err != nil {
			return shellErrorResult(toolID, fmt.Sprintf("[SHELL ERROR] %v", err))
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
		return shellErrorResult(toolID, fmt.Sprintf("[SHELL ERROR] %v", err))
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
			Content: tape.MarshalToolResultContent(shellStructuredResult{
				Tool:   "sh",
				Mode:   "interactive",
				Status: "spawned",
				Job: &shellStructuredJob{
					PID:              job.ID,
					Path:             job.displayDir,
					Interactive:      true,
					WorkspaceSession: jobWorkspaceSession(job),
					Adoptable:        jobWorkspaceSession(job) != "",
				},
			}),
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
		return tape.ToolResult{
			ToolID: toolID,
			Content: tape.MarshalToolResultContent(shellStructuredResult{
				Tool:   "sh",
				Mode:   "detached",
				Status: "spawned",
				Job: &shellStructuredJob{
					PID:      job.ID,
					Path:     job.displayDir,
					Detached: true,
				},
			}),
		}
	}

	if b.ProcessStarted != nil {
		b.ProcessStarted(job.cmd.Process)
	}
	effectiveTimeout := b.effectiveTimeout(timeout)
	outcome, err := b.waitSyncJob(job, effectiveTimeout, b.killSyncJobOnTimeout())
	if err != nil {
		return shellErrorResult(toolID, fmt.Sprintf("[SHELL ERROR] %v", err))
	}
	if b.ProcessEnded != nil {
		b.ProcessEnded()
	}

	stdoutStr := b.applyOutputLimit(b.readLog(job.outPath))
	stderrStr := b.applyOutputLimit(b.readLog(job.errPath))
	mutations := ""
	worldRevisionBlock := ""
	if b.subjective != nil && b.subjective.enabled {
		finalized, err := b.subjective.finalizeTurn("sh", b.TurnID)
		if err != nil {
			return shellErrorResult(toolID, fmt.Sprintf("[SHELL ERROR] finalize workspace turn %d: %v", b.TurnID, err))
		}
		if b.FSMutationTelemetryEnabled {
			mutations = finalized.Mutations
		}
		if finalized.Revision.ID != "" {
			worldRevisionBlock = formatWorldRevisionCreated(finalized.Revision, !finalized.Changed)
		}
	}

	if outcome == syncJobInterrupted {
		timeoutSeconds := int((effectiveTimeout + time.Second - 1) / time.Second)
		payload := shellStructuredResult{
			Tool:             "sh",
			Mode:             "sync",
			Status:           "interrupted",
			Job:              &shellStructuredJob{PID: job.ID, Path: job.displayDir},
			StdoutSoFar:      stdoutStr,
			StderrSoFar:      stderrStr,
			FSMutationsSoFar: mutations,
			WorldRevision:    worldRevisionBlock,
			Cause:            "timeout",
			TimeoutSeconds:   &timeoutSeconds,
		}
		return tape.ToolResult{
			ToolID:  toolID,
			Content: tape.MarshalToolResultContent(payload),
		}
	}

	// Cleanup job directory for completed synchronous jobs since
	// results have already been captured and returned inline.
	_ = os.RemoveAll(job.canonicalDir)

	payload := shellStructuredResult{
		Tool:          "sh",
		Mode:          "sync",
		Status:        "completed",
		ExitCode:      &job.exitCode,
		Stdout:        stdoutStr,
		Stderr:        stderrStr,
		FSMutations:   mutations,
		WorldRevision: worldRevisionBlock,
	}

	return tape.ToolResult{
		ToolID:  toolID,
		Content: tape.MarshalToolResultContent(payload),
		IsError: job.exitCode != 0,
	}
}

func (b *ShExecutor) SwitchWorld(toolID string, target string) tape.ToolResult {
	if strings.TrimSpace(target) == "" {
		return tape.ToolResult{
			ToolID: toolID,
			Content: tape.MarshalToolResultContent(switchWorldStructuredResult{
				Tool:   "switch_world",
				Status: "error",
				Target: target,
				Error:  "[SWITCH WORLD ERROR] target is required",
			}),
			IsError: true,
		}
	}
	if b.subjective == nil || !b.subjective.enabled {
		return tape.ToolResult{
			ToolID: toolID,
			Content: tape.MarshalToolResultContent(switchWorldStructuredResult{
				Tool:   "switch_world",
				Status: "error",
				Target: target,
				Error:  "[SWITCH WORLD ERROR] world switching requires Linux workspace physics",
			}),
			IsError: true,
		}
	}
	if !b.subjective.canRestoreWorld() {
		return tape.ToolResult{
			ToolID: toolID,
			Content: tape.MarshalToolResultContent(switchWorldStructuredResult{
				Tool:   "switch_world",
				Status: "error",
				Target: target,
				Error:  "[SWITCH WORLD ERROR] switch_world requires workspace revision mode with restore support",
			}),
			IsError: true,
		}
	}
	if err := b.initState(); err != nil {
		return tape.ToolResult{
			ToolID: toolID,
			Content: tape.MarshalToolResultContent(switchWorldStructuredResult{
				Tool:   "switch_world",
				Status: "error",
				Target: target,
				Error:  fmt.Sprintf("[SWITCH WORLD ERROR] %v", err),
			}),
			IsError: true,
		}
	}
	previous, current, err := b.subjective.switchWorld(target)
	if err != nil {
		return tape.ToolResult{
			ToolID: toolID,
			Content: tape.MarshalToolResultContent(switchWorldStructuredResult{
				Tool:   "switch_world",
				Status: "error",
				Target: target,
				Error:  fmt.Sprintf("[SWITCH WORLD ERROR] %v", err),
			}),
			IsError: true,
		}
	}
	mutations, err := b.subjective.restoreMutationBlock(previous, current)
	if err != nil {
		return tape.ToolResult{
			ToolID: toolID,
			Content: tape.MarshalToolResultContent(switchWorldStructuredResult{
				Tool:   "switch_world",
				Status: "error",
				Target: target,
				Error:  fmt.Sprintf("[SWITCH WORLD ERROR] switch mutation diff: %v", err),
			}),
			IsError: true,
		}
	}
	worldRevisionBlock := formatWorldRevisionTransition(previous, current)
	if !b.FSMutationTelemetryEnabled {
		mutations = ""
	}
	return tape.ToolResult{
		ToolID: toolID,
		Content: tape.MarshalToolResultContent(switchWorldStructuredResult{
			Tool:          "switch_world",
			Status:        "completed",
			Target:        target,
			Revision:      current,
			FSMutations:   mutations,
			WorldRevision: worldRevisionBlock,
		}),
	}
}

func (b *ShExecutor) CurrentWorldRevision() string {
	if b.subjective == nil {
		return ""
	}
	return b.subjective.currentWorldRevision()
}

func jobWorkspaceSession(job *managedJob) string {
	if job == nil || job.workspace == nil || !job.workspace.enabled || !job.workspace.usesOverlayBackend() {
		return ""
	}
	return job.workspace.workspaceSession
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
	return b.CloseWithOptions(keepDetached, keepDetached)
}

// CloseWithOptions separates detached-job lifetime from workspace finalization.
// This lets signal shutdown commit the latest completed workspace revision
// without sparing unfinished background processes.
func (b *ShExecutor) CloseWithOptions(keepDetached bool, commitWorkspace bool) error {
	if err := b.initState(); err != nil {
		return err
	}

	b.mu.Lock()
	jobs := make([]*managedJob, 0, len(b.detached))
	for _, job := range b.detached {
		jobs = append(jobs, job)
	}
	pending := make([]*managedJob, 0, len(b.pending))
	for _, job := range b.pending {
		pending = append(pending, job)
	}
	b.mu.Unlock()

	if !keepDetached {
		for _, job := range jobs {
			_ = syscall.Kill(-job.ID, syscall.SIGKILL)
		}
	}
	for _, job := range pending {
		_ = syscall.Kill(-job.ID, syscall.SIGKILL)
	}
	if commitWorkspace && b.subjective != nil && b.subjective.enabled {
		if err := b.subjective.commit(); err != nil {
			return err
		}
	}
	if !commitWorkspace && b.subjective != nil && b.subjective.enabled {
		if err := b.subjective.rollback(); err != nil {
			return err
		}
	}
	return nil
}

func RecoverWorkspaceCommit(cfg *config.Config) error {
	subjective := newSubjectiveFS(cfg)
	if subjective == nil || !subjective.enabled {
		return nil
	}
	if err := subjective.init(subjective.dataDir, subjective.sessionID); err != nil {
		return err
	}
	return subjective.commit()
}

func (b *ShExecutor) ResetWorkspaceAfterExternalCommit() {
	if b == nil || b.subjective == nil {
		return
	}
	b.mu.Lock()
	defer b.mu.Unlock()
	b.subjective.initialized = false
	b.subjective.currentRevision = ""
}

// Numeric tool-argument coercion lives in arg_coerce.go (intFromAny / IntArg),
// which distinguishes an absent argument from a present-but-uncoercible one
// instead of silently returning 0. The former ToInt/toDuration helpers masked
// that distinction and have been removed.
