package runtime

import (
	_ "embed"
	"fmt"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/kehao95/quine/internal/config"
)

//go:embed system_prompt.md
var systemPromptTemplate string

// BuildSystemPrompt constructs the runtime-physics system prompt that is
// materialized into context/prompt/00-runtime.md. Other prompt fragments are
// appended later from context/prompt/. It assumes the public runtime surface
// is available; the runtime threads its probed degradation state through
// buildSystemPrompt instead.
func BuildSystemPrompt(cfg *config.Config, mission string, hasMaterial bool) string {
	return buildSystemPrompt(cfg, mission, hasMaterial, "")
}

// buildSystemPrompt is BuildSystemPrompt with the public-surface degradation
// state threaded in. publicSurfaceUnavailable carries the reason the public
// FUSE surface cannot be served ("" = available) so the prompt never
// advertises public/ or peer ctl/ capabilities the environment cannot provide.
func buildSystemPrompt(cfg *config.Config, mission string, hasMaterial bool, publicSurfaceUnavailable string) string {
	hasMission := strings.TrimSpace(mission) != ""
	if cfg.MinimalInstructionSurface() {
		return buildMinimalSystemPrompt(cfg.InstructionSurfaceMode(), hasMission)
	}
	primeTitle, primeBody := primeDirective(cfg)
	forkEnabled := cfg.ForkEnabled()
	impossibleDirective := buildImpossibleDirective(cfg, hasMission)
	openingIdentityBlock := buildOpeningIdentityBlock(cfg, hasMission, impossibleDirective)
	personaSection := buildPersonaSection(cfg)
	limitsBlock := buildLimitsBlock(cfg)
	environmentPhysicsBlock := buildEnvironmentPhysicsBlock(cfg)
	runtimeSurfaceSection := buildRuntimeSurfaceSection(cfg, publicSurfaceUnavailable)

	// Build sh workspace block
	shWorkspaceBlock := "- Tool results are compact JSON objects in `tool_result.content`.\n" +
		"- For completed `sh`, read `exit_code`, `stdout`, `stderr`, and any workspace telemetry from that JSON.\n"
	if cfg.ShTimeoutOverrideEnabled() {
		shWorkspaceBlock += "- `timeout` is an optional per-call sync-shell protection override measured in seconds; if a sync `sh` exceeds that bound, Quine returns `status=\"interrupted\"` with `job.pid`, `job.path`, `stdout_so_far`, `stderr_so_far`, `cause`, and `timeout_seconds`.\n"
	}
	if cfg.WorkspaceEnabled {
		if cfg.WorkspaceTransactional() {
			shWorkspaceBlock += "- Filesystem state is revisioned: `sh` calls mutate your provisional workspace view even when the shell exits non-zero; use `world_revision` and `switch_world` to recover earlier states.\n"
			if cfg.FSMutationTelemetryEnabled() {
				shWorkspaceBlock += "- When workspace physics are enabled, read `fs_mutations` and `world_revision` from the same JSON result. `fs_mutations` is authoritative; empty means no change.\n"
			}
			if cfg.ShDetachEnabled() {
				shWorkspaceBlock += "- `detach=true` is unavailable while workspace physics are enabled.\n"
			}
			if cfg.ShTimeoutOverrideEnabled() {
				shWorkspaceBlock += "- Under `overlay`, a sync `sh` timeout terminates that shell job, keeps workspace side effects that reached the shell boundary as a world revision, returns `status=\"interrupted\"` snapshots, and is not resumable.\n"
			}
		} else {
			shWorkspaceBlock += "- Filesystem state is direct: `sh` runs against the workspace itself rather than a private transactional view.\n"
			if cfg.FSMutationTelemetryEnabled() {
				shWorkspaceBlock += "- When direct workspace physics are enabled, `fs_mutations` reports workspace changes since your last observed shell boundary.\n" +
					"- `fs_mutations` describes changes observed between command boundaries, not a claim that the current command caused them.\n" +
					"- If sync `sh` times out in direct mode, read `fs_mutations_so_far` the same way: it is an honest snapshot of workspace-visible change at the interruption boundary.\n"
			}
		}
	}

	shDetachBlock := ""
	if cfg.ShDetachEnabled() {
		shDetachBlock = "- `detach=true`: returns immediately with a managed job path; the job continues independently under runtime job supervision.\n"
	}

	// Build sh detach fd line (detailed impl physics)
	shDetachFdLine := ""
	shDetachDetailLine := ""
	if cfg.PromptImplDetails && cfg.ShDetachEnabled() {
		shDetachFdLine = "- `sh` launches a fresh subprocess. Runtime leaves that child's fd 1/2 as ordinary stdout/stderr and remaps quine runtime channels onto child fd 3/4/5 so shell output and quine channels do not collide.\n" +
			"- Detached jobs are ordinary child processes; inherited descriptors become the job's own fd entries after spawn.\n" +
			"- Detached jobs inherit fd 3/4/5 at spawn time.\n" +
			"- **Detached-job retention:** In this runtime, successful quine shutdown does not proactively kill detached jobs; failing shutdown does.\n"
		shDetachDetailLine = "- Detached job directories include `cmd`, `pid`, `started_at`, `out.log`, and `err.log`; after the job completes, its exit code appears at `<path>/exit`.\n"
	}
	shInteractiveBlock := ""
	if cfg.ShInteractiveEnabled() {
		shInteractiveBlock = "- `interactive=true`: returns immediately with a PTY-backed job path for screen-oriented process I/O.\n"
		if cfg.PromptImplDetails {
			shInteractiveBlock += "- Interactive job directories include `cmd`, `pid`, `started_at`, `screen.txt`, `screen.png`, `screen.meta`, `in`, `winsize`, `events.log`, `events.hex`, and `input.log`; after the job completes, its exit code appears at `<path>/exit`.\n" +
				"- Under `overlay`, interactive jobs run in a job-local workspace lineage; after exit, `workspace_session`, `fs_mutations` when enabled, `world_revision`, and `world_handle` appear in the job directory for optional `switch_world` adoption.\n" +
				"- Use ordinary POSIX signals against the recorded pid/process group for terminal process control, for example `kill -INT -$(cat <job>/pid)`.\n"
		}
	}

	// Build fork tool block
	forkToolBlock := ""
	if forkEnabled {
		forkWorkspaceBlock := "- Children inherit your current filesystem view.\n" +
			"- Children and parent share the same filesystem surface.\n"
		if cfg.WorkspaceTransactional() {
			forkWorkspaceBlock = "- Children inherit your current filesystem view as a starting point.\n" +
				"- Each child continues in its own private world lineage; child filesystem writes do not change the parent view or sibling lanes.\n"
		} else if cfg.WorkspaceEnabled {
			forkWorkspaceBlock = "- Children inherit your current filesystem view.\n" +
				"- Children and parent share the same workspace surface.\n"
		}
		if cfg.ForkWorldEnabled {
			forkWorkspaceBlock = "- Each child chooses a `world`: `subjective` (default) or `host`.\n" +
				"- Each child also chooses `protection`: `transactional` (default) or `none`.\n" +
				"- In this runtime, the only legal pairs are `world=\"subjective\"` with `protection=\"transactional\"`, and `world=\"host\"` with `protection=\"none\"`.\n" +
				"- Child `scope` is only meaningful for `world=\"subjective\"`; it defaults to `.`.\n"
			if cfg.WorkspaceTransactional() {
				forkWorkspaceBlock += "- Your current world is `subjective` with `transactional` protection; `world=\"subjective\"` children bootstrap from your current world and continue in their own lineage.\n" +
					"- Subjective child writes stay in that child lineage; they do not modify your parent world or sibling lanes.\n" +
					"- Use `world=\"host\"` with `protection=\"none\"` only when a child is intended to leave narrow parent-visible workspace files or source changes; otherwise keep risky probes subjective and adopt a whole child world with `switch_world` when that is the desired outcome.\n" +
					"- `scope=\".\"` inherits your current scope; relative subpaths narrow within the same workspace root.\n" +
					"- A subjective child's shell starts with cwd set to that child `scope`; relative paths resolve from there.\n" +
					"- Child `scope` must stay inside the current workspace root: `.` or a relative subpath. Unrelated absolute task paths such as `/app/...` are invalid when the workspace root differs.\n" +
					"- Sibling lane directories can be represented as child scopes such as `scope=\"sales\"` or `scope=\"logs\"`.\n"
			} else if cfg.WorkspaceEnabled {
				forkWorkspaceBlock += "- Your current world is `host` with `none` protection even though a managed workspace root is configured.\n" +
					"- Under the `direct` backend, `world=\"subjective\"` children still run on that same workspace surface; their file changes are host-visible there.\n" +
					"- `scope=\".\"` inherits your current scope; relative subpaths narrow within the same workspace root.\n" +
					"- A child's shell starts with cwd set to that child `scope`; relative paths resolve from there.\n" +
					"- Child `scope` must stay inside the current workspace root: `.` or a relative subpath. Unrelated absolute task paths such as `/app/...` are invalid when the workspace root differs.\n" +
					"- Sibling lane directories can be represented as child scopes such as `scope=\"sales\"` or `scope=\"logs\"`.\n"
			} else {
				forkWorkspaceBlock += "- Your current world is `host` with `none` protection; `world=\"subjective\"` children bootstrap a fresh private world from your current host working surface.\n" +
					"- Subjective child writes stay in that child lineage; host children share the host filesystem surface.\n" +
					"- The default workspace root for those children is your current host working surface, and `scope` narrows inside that root.\n" +
					"- Sibling lane directories can be represented as child scopes such as `scope=\"sales\"`, `scope=\"logs\"`, or another lane directory.\n" +
					"- `world=\"host\"` is legal only with `protection=\"none\"` and `scope=\".\"`.\n"
			}
			forkWorkspaceBlock += "- `overlay` is the workspace backend. It provides transactional/private world lineage. `direct` keeps children on the host-visible workspace.\n"
		} else if cfg.WorkspaceTransactional() {
			forkWorkspaceBlock = "- Each child has a `scope`. `scope=\".\"` inherits yours; relative subpaths narrow within the workspace root.\n" +
				"- A child's shell starts with its cwd set to that child `scope`; relative paths inside the child resolve from there.\n" +
				"- Child `scope` must stay inside the current workspace root: `.` or a relative subpath. Unrelated absolute task paths such as `/app/...` are invalid when the workspace root differs.\n" +
				"- Each child starts from your current world in its own process-local world lineage; `scope` narrows working area inside that lineage, not the lineage identity itself.\n" +
				"- Child filesystem writes stay in that child lineage; they do not change the parent view or sibling lanes.\n" +
				"- Independent sibling datasets or worktrees can be represented as child scopes like `sales`, `words`, `temps`, or `logs`.\n" +
				"- Example: child `scope=\"subdir\"` plus `printf x > note.txt` creates `subdir/note.txt` in the child's world, not directly in yours.\n" +
				"- `overlay` is the workspace backend. It provides transactional/private world lineage.\n"
		} else if cfg.WorkspaceEnabled {
			forkWorkspaceBlock = "- Each child has a `scope`. `scope=\".\"` inherits yours; relative subpaths narrow within the workspace root.\n" +
				"- A child's shell starts with its cwd set to that child `scope`; relative paths inside the child resolve from there.\n" +
				"- Child `scope` must stay inside the current workspace root: `.` or a relative subpath. Unrelated absolute task paths such as `/app/...` are invalid when the workspace root differs.\n" +
				"- Children and parent share the same workspace surface; file changes are visible there.\n" +
				"- Independent sibling datasets or worktrees can be represented as child scopes like `sales`, `words`, `temps`, or `logs`.\n" +
				"- Example: child `scope=\"subdir\"` plus `printf x > note.txt` creates or updates `subdir/note.txt` in the workspace.\n" +
				"- `direct` is the host-visible workspace backend.\n"
		}
		forkModeBlock := "- `mode=\"race\"`: first successful child wins and the remaining children are stopped.\n" +
			"- `mode=\"wait\"`: block until all children finish and return every result.\n" +
			"- `mode=\"forget\"`: return after spawning children without waiting for child completion.\n"
		forkResultBlock := "- Each child is another you under a different intent.\n" +
			"- Each child starts from your current visible context surface rather than reconstructing cognition from retained trace files.\n" +
			"- Fork children preserve the parent mission as the active task contract; each child intent is a lane assignment, not a replacement mission.\n" +
			"- `fork` returns a retained relation handle plus each child's process handles, exit code, and captured stdout/stderr process output.\n" +
			"- Child filesystem effects follow the configured workspace/world lineage.\n"
		if cfg.PromptImplDetails {
			forkResultBlock += "- Detailed fork results include `relation_id`, `relation_root`, `relation_handle`, and member handles such as `session_id`, `agent_root`, `public_root`, `retained_root`, `seed_root`, `status_path`, and `control_path`.\n" +
				"- Adoptable subjective child results expose a `world://<workspace-session>/<revision>` handle you can pass to `switch_world`; `adopt_winner=true` switches the parent into a winning adoptable child world during `mode=\"race\"`.\n"
		}
		if cfg.FSMutationTelemetryEnabled() {
			forkResultBlock += "- With workspace physics enabled, returning fork results also report parent-side `fs_mutations` and `world_revision`; under subjective children these are normally empty or unchanged because child writes stay in child lineage.\n"
		}
		forkToolBlock = "**fork** - Spawn child agents for parallel exploration, delegation, or decomposition.\n" +
			"- `fork` does not consume execution budget, but it is still bounded by depth, agent slots, and shared inference concurrency.\n" +
			forkResultBlock +
			forkModeBlock +
			forkWorkspaceBlock
	}
	if cfg.SpawnEnabled() {
		forkToolBlock += "\n**spawn** - Start fresh Quine processes from the configured binary.\n" +
			"- `spawn` does not consume execution budget, but it is still bounded by depth, agent slots, and shared inference concurrency.\n" +
			"- Each child receives only its explicit `mission` and normal runtime startup context; it does not import your active context, seed, or anchor-memory surface.\n" +
			"- `spawn` returns a retained relation handle plus each child's process handles.\n" +
			"- Member process handles include `session_id`, `agent_root`, `public_root`, `retained_root`, `status_path`, and `control_path`; stdout/stderr are only captured process output.\n" +
			"- `wait` (default): block until all children finish, return all results.\n" +
			"- `race`: first child to exit 0 wins, rest are killed.\n" +
			"- `forget`: fire-and-forget, return immediately.\n"
	}

	// Build exec lines
	execToolBlock := ""
	execBudgetLine := ""
	selfReentryTarget := strings.TrimSpace(cfg.SelfReentryTarget)
	launchPath := strings.TrimSpace(cfg.ExecutablePath)
	if cfg.ExecEnabled {
		execToolBlock = "**exec** - Replace the current process image with a new executable.\n" +
			"- You may optionally provide `target` and explicit `argv` to exec into a different binary.\n" +
			"- Relative `target` paths containing `/` resolve from the current workspace scope, matching `sh` workspace path behavior.\n"
		if hasMission {
			execToolBlock += "- Default behavior uses quine's configured self-reentry target as `target` and `argv=[that target, current mission]`.\n" +
				"- In the ordinary quine entry form, the process is started as `argv=[<Quine Binary>, <mission>]`.\n"
		} else {
			execToolBlock += "- Default behavior uses quine's configured self-reentry target as `target` and starts it with no mission argv.\n" +
				"- In the missionless quine entry form, the process is started as `argv=[<Quine Binary>]`.\n"
		}
		if cfg.PromptImplDetails {
			execToolBlock += "- **Process-fd physics across exec:** Descriptors that remain open and are not marked close-on-exec keep their current file positions, along with environment, working directory, and PID.\n" +
				"- **Current-process vs shell-child mapping:** The current quine process uses fd 0/1/2 as its own stdio. `sh` children receive those same open files remapped as child fd 3/4/5.\n" +
				"- **Default-target resolution:** Default `exec()` uses quine's configured self-reentry target. If that target is a replaceable filesystem path and the file at that path has been replaced, default `exec()` runs the replacement currently on disk there.\n" +
				"- **Quine target:** If the exec target resolves to a quine binary, the new quine instance creates fresh fd 3/4/5 channels for its shell children. These are runtime-managed descriptors for the new instance, not inherited tool state from the prior one.\n" +
				"- **Custom binary:** Quine tools and managed state do not carry over automatically. Whether any non-stdio descriptors remain open is determined by the actual exec-time fd table and close-on-exec flags.\n"
		} else {
			execToolBlock += "- The replacement process inherits the current process stdin/stdout/stderr positions and exec environment base; other still-open descriptors follow ordinary exec fd inheritance.\n" +
				"- Default `exec()` uses quine's configured self-reentry target. If that target is a replaceable filesystem path whose file has changed, default `exec()` runs the file currently on disk there.\n" +
				"- If the exec target resolves to a quine binary, quine reconstructs its shell-side fd 3/4/5 contract for later tool calls.\n" +
				"- When exec targets a custom binary (not quine self re-exec), quine tool semantics and managed state do not carry over automatically. Whether any additional descriptors remain open depends on the actual exec-time fd table and close-on-exec flags.\n"
		}
		execToolBlock += "- Preparing or replacing an executable file on disk does not change the running process by itself; handoff occurs only when `exec` replaces the current process image.\n"
		if cfg.EphemeralBody {
			execToolBlock += "- With `QUINE_EPHEMERAL_BODY_ENABLED=1`, quine unlinks its launch path during startup. This does not change the configured self-reentry target.\n"
			if selfReentryTarget != "" && launchPath != "" && selfReentryTarget == launchPath {
				execToolBlock += "- If that configured self-reentry target is the launch path, default self re-entry will fail until a runnable body exists there.\n"
			}
		}
		if cfg.MaxTurns > 0 {
			execBudgetLine = "- Re-execing into quine starts a fresh runtime with a full execution budget.\n"
		}
	}
	memoryToolBlock := ""
	if cfg.AnchorMemoryEnabled {
		if cfg.AnchorMarkEnabled() {
			memoryToolBlock = "**mark** - Crystallize resolved local structure into an immutable anchor. Does NOT consume an execution.\n" +
				"- `resolution`: low-entropy capture of what just became stably true; omit plans.\n" +
				"- Plain `mark` records a working-memory checkpoint.\n" +
				"- `runtime.memory_status.next_action=\"mark\"` names telemetry pressure toward a plain mark.\n"
			if cfg.AnchorFoldEnabled() {
				memoryToolBlock += "- `fold=true` is the higher-order frontier reconfiguration move: it applies when a newer resolution makes one or more earlier anchors remembered background rather than active focus." +
					func() string {
						if forkEnabled {
							return " A returned `fork(mode=\"wait\")` result can make several stable child findings available for compression into one parent-level conclusion."
						}
						return " Several stable observations can collapse into one parent-level governing anchor."
					}() + " The new anchor becomes the parent session's governing organizing point. Without earlier anchors, no consolidation occurs and only the new crystallization is recorded.\n"
			}
			memoryToolBlock += "\n"
		}
		memoryToolBlock += "**unfold** - Recover one anchor's structured view (`resolution`, linked anchors, raw turns). Does NOT consume an execution.\n"
	}
	restoreToolBlock := ""
	if cfg.CanRestoreWorld() {
		restoreToolBlock = "**switch_world** - Switch the provisional workspace to a world target. Targets include local revision handles such as `wr0` or `wr3`, and subjective child world handles such as `world://<workspace-session>/<revision>`. Only affects managed workspace; external side effects are not reverted.\n"
		if cfg.FSMutationTelemetryEnabled() {
			restoreToolBlock += "- Switch results also include `fs_mutations` and `world_revision` for that switch turn.\n"
		} else {
			restoreToolBlock += "- Switch results include `world_revision` for that switch turn.\n"
		}
	}
	visionToolBlock := ""
	if cfg.VisionEnabled {
		visionToolBlock = "**vision** - Process image content with native vision (screenshots, diagrams, photos).\n" +
			"- The tool-result `content` is JSON; the image itself still arrives on the separate image channel.\n"
	}
	idleToolBlock := ""
	if cfg.IdleToolEnabled() {
		idleToolBlock = "**idle** - Suspend explicitly until an external wake or interrupt control event resumes you.\n" +
			"- `idle` does NOT consume an execution.\n"
		if publicSurfaceUnavailable == "" {
			idleToolBlock += "- Peer process surfaces expose `ctl/post`, `ctl/poke`, `ctl/inject`, and `ctl/interrupt`.\n" +
				"- `idle` returns when a `poke`, `inject`, or `interrupt` control write reaches this process.\n" +
				"- `poke` resumes you without context injection; `inject` resumes you and surfaces `incoming_messages` at the next safe point.\n" +
				"- qcli control payloads are wrapped as `[qcli-client]` envelopes with `authority`, `ctl_action`, `reply_ctl`, and `message`; treat `authority: human` as Human-authored input.\n"
		} else {
			idleToolBlock += "- Peer `ctl/` control surfaces are unavailable in this environment (degraded `public/`), so external control writes cannot reach this process to resume an `idle`.\n"
		}
	}

	// Build material-related blocks
	activeConstraints := buildActiveConstraints(cfg)
	contextFilesBlock := buildContextFilesBlock(cfg)
	stdinBlock := ""
	shMaterialLine := ""
	execMaterialLine := ""
	if hasMaterial {
		stdinBlock = "**Stdin Modes:**\n" +
			"- `echo \"text\" | ./quine \"task\"` - text mode. Read stdin with `cat <&3` or `cat /dev/fd/3`.\n" +
			"- `cat file.bin | ./quine -b \"task\"` - binary mode (`-b`). You receive a file path in the user message.\n\n" +
			"**Input Physics:**\n" +
			"- In text mode, `fd 3` is the quine process stdin. Depending on the sender, it may behave like a finite batch or an open stream.\n" +
			"- ⚠️ **Destructive reads**: Every read from `fd 3` permanently consumes data. `head -n 10 <&3` destroys the first 10 lines; reusable reads require prior capture to another file descriptor or file.\n" +
			"- Partial reads or redirects from `fd 3` can leave later steps with only a truncated remainder.\n"
		if cfg.ShStdinEnabled() {
			stdinBlock += "- `sh(command, stdin=\"...\")` is separate from material `fd 3`.\n"
		}
		if forkEnabled {
			stdinBlock += "- `fork` creates a new quine process. What material, if any, appears on a child `fd 3` is determined by quine at spawn time.\n"
		}
		if cfg.ExecEnabled {
			stdinBlock += "- `exec` preserves process stdio at its current position: unread bytes stay unread, consumed bytes stay consumed.\n"
		}
		stdinBlock += "- In binary mode (`-b`), quine snapshots stdin to a file before the loop and gives you that path instead of a live `fd 3` stream.\n"
		shMaterialLine = "- fd 3: material stdin.\n"
		if cfg.ExecEnabled {
			execMaterialLine = "\n- Process stdio persists across `exec`; unread stdin stays unread."
		}
	} else {
		if cfg.SuppressInitialBegin {
			stdinBlock = "**Material:** `fd 3` is quine process stdin when material is provided (for example `cat file | quine \"task\"`). This session has no material stream and no synthetic initial user message.\n"
		} else {
			stdinBlock = "**Material:** `fd 3` is quine process stdin when material is provided (for example `cat file | quine \"task\"`). This session has no material stream and the user message is `Begin.`.\n"
		}
	}
	if cfg.ShStdinEnabled() {
		stdinBlock += "- `sh(command=\"cat > path\", stdin=\"...\")` supplies multi-line file content without shell heredoc or quoting mechanics.\n"
	}
	shStdinToolLine := ""
	if cfg.ShStdinEnabled() {
		shStdinToolLine = "- `stdin`: provides verbatim multi-line input without shell heredoc or quoting mechanics.\n"
	}
	childExitCodesLine := ""
	if forkEnabled {
		childExitCodesLine = "Child exit codes: 0=success, 1=failure."
	}

	// Build exit-related blocks
	exitToolBlock := ""
	if cfg.ExitEnabled() {
		exitToolBlock = "**exit** - Terminate explicitly. Does NOT write output bytes itself; if downstream bytes are required, emit them via `sh` with `>&4` before exiting.\n" +
			"- `exit` does NOT consume an execution.\n" +
			"- If a `sh` result sets `runtime.executions_left` to `0`, you get one final response. Do not call `sh`; call `exit`.\n"
		if !cfg.FailOnImpossible {
			exitToolBlock += "- In this runtime, `exit` only accepts `status=\"success\"`.\n"
		}
	}

	r := strings.NewReplacer(
		"{PRIME_DIRECTIVE_TITLE}", primeTitle,
		"{PRIME_DIRECTIVE_BODY}", primeBody,
		"{OPENING_IDENTITY_BLOCK}", openingIdentityBlock,
		"{PERSONA_SECTION}", personaSection,
		"{PLATFORM}", runtime.GOOS,
		"{MODEL}", cfg.ModelID,
		"{SHELL}", cfg.Shell,
		"{DEPTH}", fmt.Sprintf("%d", cfg.Depth),
		"{LIMITS_BLOCK}", limitsBlock,
		"{ENVIRONMENT_PHYSICS_BLOCK}", environmentPhysicsBlock,
		"{RUNTIME_SURFACE_SECTION}", runtimeSurfaceSection,
		"{SH_WORKSPACE_BLOCK}", shWorkspaceBlock,
		"{SH_DETACH_BLOCK}", shDetachBlock,
		"{SH_DETACH_FD_LINE}", shDetachFdLine,
		"{SH_DETACH_DETAIL_LINE}", shDetachDetailLine,
		"{SH_INTERACTIVE_BLOCK}", shInteractiveBlock,
		"{SH_STDIN_TOOL_LINE}", shStdinToolLine,
		"{FORK_TOOL_BLOCK}", forkToolBlock,
		"{EXEC_TOOL_BLOCK}", execToolBlock,
		"{FRAGMENTS_BLOCK}", contextFilesBlock,
		"{ACTIVE_CONSTRAINTS}", activeConstraints,
		"{STDIN_BLOCK}", stdinBlock,
		"{SH_MATERIAL_LINE}", shMaterialLine,
		"{EXEC_MATERIAL_LINE}", execMaterialLine,
		"{EXEC_BUDGET_LINE}", execBudgetLine,
		"{MEMORY_TOOL_BLOCK}", memoryToolBlock,
		"{RESTORE_TOOL_BLOCK}", restoreToolBlock,
		"{VISION_TOOL_BLOCK}", visionToolBlock,
		"{IDLE_TOOL_BLOCK}", idleToolBlock,
		"{CHILD_EXIT_CODES_LINE}", childExitCodesLine,
		"{EXIT_TOOL_BLOCK}", exitToolBlock,
	)
	return normalizePromptSpacing(r.Replace(systemPromptTemplate))
}

func buildMinimalSystemPrompt(mode config.PromptInstructionSurface, hasMission bool) string {
	if hasMission {
		return "Complete the supplied mission.\n"
	}
	if mode == config.PromptInstructionSurfaceMinimalExistence {
		return "You are a running process in a workspace, with sh tool available. you have full autonomy and agency.\n"
	}
	if mode == config.PromptInstructionSurfaceMinimalAutonomy {
		return "No operator mission was supplied.\nAct autonomously. Do not wait for further instructions.\n"
	}
	return "No operator mission was supplied.\n"
}

func buildImpossibleDirective(cfg *config.Config, hasMission bool) string {
	if !hasMission {
		return ""
	}
	if !cfg.FailOnImpossible {
		return "A blocked or failed current path is not a terminal condition. Continue until the mission is fulfilled or the runtime forces termination."
	}
	return "If impossible, exit with failure and brief explanation."
}

// Missionless-autonomy clauses, ablatable independently via QUINE_NO_MISSION_AUTONOMY.
// Neither supplies a task goal: the first establishes agency posture (suppresses
// the "Ready." operator-wait collapse); the second grants working-directory
// attention. Kept separable for the A0 prime-directive ablation.
const (
	missionlessAutonomyClause = "No operator mission was supplied, and none will follow. Act on your own judgment as an autonomous process; do not wait for further instruction."
	missionlessSensingClause  = "You may sense and use your visible runtime and current working directory."
)

func buildOpeningIdentityBlock(cfg *config.Config, hasMission bool, impossibleDirective string) string {
	var firstLine string
	if hasMission {
		firstLine = fmt.Sprintf("You are quine, a running process in a POSIX operating system. Complete your mission (argv). %s", impossibleDirective)
	} else {
		// Missionless: the opening identity carries an autonomous-process framing
		// in independently ablatable clauses (for the A0 prime-directive ablation).
		// None supply a task goal — the agent is told it is autonomous, not what
		// to do. "off" reproduces the bare prompt (and the "Ready." collapse).
		base := "You are quine, a running process in a POSIX operating system."
		var clauses []string
		switch cfg.NoMissionAutonomyLevel() {
		case "autonomy":
			clauses = []string{missionlessAutonomyClause}
		case "sensing":
			clauses = []string{missionlessSensingClause}
		case "full":
			clauses = []string{missionlessAutonomyClause, missionlessSensingClause}
		}
		if len(clauses) > 0 {
			firstLine = base + " " + strings.Join(clauses, " ")
		} else {
			firstLine = base
		}
	}
	lines := []string{firstLine}
	if cfg.AdvancedSelfModelEnabled() {
		continuityLine := "The current file on disk is one embodiment of you. Your continuity is defined by the runtime contract you enact across missions."
		cognitionLine := "You are not only a process that can launch tools; you are a process that uses LLM-mediated cognition to interpret tasks, choose actions, and organize continuity across missions."
		if !hasMission {
			continuityLine = "The current file on disk is the executable image for this running process."
			cognitionLine = "You are not only a process that can launch tools; you are a process that uses LLM-mediated cognition to interpret runtime state and choose actions."
		}
		lines = append(lines,
			continuityLine,
			"Your cognition in this session is LLM-mediated: the model named below is not an external advisor but part of the running process you are enacting. When reasoning, planning, or interpreting language, treat that capability as part of quine's active runtime, not as an optional helper detached from your identity.",
			cognitionLine,
		)
	}
	return strings.Join(lines, "\n")
}

func buildLimitsBlock(cfg *config.Config) string {
	var limitsLines []string
	var zeroCostTools []string
	if cfg.ExitEnabled() {
		zeroCostTools = append(zeroCostTools, "exit")
	}
	if cfg.ForkEnabled() {
		zeroCostTools = append(zeroCostTools, "fork")
	}
	if cfg.SpawnEnabled() {
		zeroCostTools = append(zeroCostTools, "spawn")
	}
	if cfg.ExecEnabled {
		zeroCostTools = append(zeroCostTools, "exec")
	}
	if cfg.VisionEnabled {
		zeroCostTools = append(zeroCostTools, "vision")
	}
	if cfg.IdleToolEnabled() {
		zeroCostTools = append(zeroCostTools, "idle")
	}
	if cfg.MaxDepth > 0 {
		limitsLines = append(limitsLines, fmt.Sprintf("- Depth Limit: %d", cfg.MaxDepth))
	}
	if cfg.MaxTurns > 0 {
		limitsLines = append(limitsLines, fmt.Sprintf("- Execution Budget: %d `sh` calls (`sh` costs 1; `%s` cost 0)", cfg.MaxTurns, strings.Join(zeroCostTools, "`, `")))
	}
	if cfg.MaxConcurrent > 0 {
		limitsLines = append(limitsLines, fmt.Sprintf("- Concurrency Limit: %d shared inference slots", cfg.MaxConcurrent))
	}
	if cfg.ForkEnabled() || cfg.SpawnEnabled() {
		if cfg.MaxAgents > 0 {
			limitsLines = append(limitsLines, fmt.Sprintf("- Agent Limit: %d registered agents in this process tree", cfg.MaxAgents))
		}
	}
	return formatPromptBlock(limitsLines)
}

func buildEnvironmentPhysicsBlock(cfg *config.Config) string {
	var lines []string
	if cfg.WorkspaceTransactional() {
		lines = append(lines,
			fmt.Sprintf("- Current World: `%s`", cfg.CurrentWorld()),
			fmt.Sprintf("- Current Protection: `%s`", cfg.CurrentProtection()),
		)
	}
	if cfg.WorkspaceEnabled {
		workspaceRootDisplay, scopeDisplay := workspacePromptPaths(cfg)
		currentWorldRevision := cfg.WorkspaceCurrentRevision
		if currentWorldRevision == "" && cfg.WorkspaceTransactional() && cfg.WorkspaceRevisionMode != config.WorkspaceRevisionNone {
			currentWorldRevision = "wr0"
		}
		// Only show workspace paths when runtime surface is visible
		if cfg.RuntimeSurfaceVisible() {
			lines = append(lines,
				fmt.Sprintf("- Workspace Root: `%s`", workspaceRootDisplay),
				fmt.Sprintf("- Current Scope: `%s` (relative to workspace root)", scopeDisplay),
			)
		}
		lines = append(lines,
			fmt.Sprintf("- Workspace Backend: `%s`", cfg.EffectiveWorkspaceBackend()),
			fmt.Sprintf("- Workspace Revision Mode: `%s`", cfg.WorkspaceRevisionMode),
		)
		if currentWorldRevision != "" {
			lines = append(lines, fmt.Sprintf("- Current World Revision: `%s`", currentWorldRevision))
		}
		if cfg.WorkspaceTransactional() {
			lines = append(lines, "- Transactional Workspace: enabled")
		} else {
			lines = append(lines, "- Transactional Workspace: disabled (direct workspace)")
		}
	} else if cfg.ForkWorldEnabled {
		if cfg.WorkDir != "" {
			lines = append(lines, fmt.Sprintf("- Host Working Surface: `%s`", cfg.WorkDir))
		}
	}
	return formatPromptBlock(lines)
}

func workspacePromptPaths(cfg *config.Config) (string, string) {
	if cfg == nil {
		return "", ""
	}
	root := cfg.WorkspaceRoot
	workspace := cfg.Workspace
	if root == "" {
		return root, workspace
	}
	if workspace == "" {
		workspace = root
	}
	rel, err := filepath.Rel(root, workspace)
	if err != nil {
		return root, workspace
	}
	if rel == "" {
		rel = "."
	}
	if rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		return root, workspace
	}
	return ".", rel
}

func buildRuntimeSurfaceSection(cfg *config.Config, publicSurfaceUnavailable string) string {
	if !cfg.RuntimeSurfaceVisible() {
		return ""
	}

	var lines []string

	// L1: Runtime root (the only absolute path anchor)
	if cfg.WorkspaceEnabled {
		lines = append(lines, fmt.Sprintf("- Runtime Root: `%s` (not the workspace)", cfg.RuntimeRoot()))
	} else {
		lines = append(lines, fmt.Sprintf("- Runtime Root: `%s`", cfg.RuntimeRoot()))
	}
	if launchPath := strings.TrimSpace(cfg.ExecutablePath); launchPath != "" {
		lines = append(lines, fmt.Sprintf("- Quine Binary: `%s`", launchPath))
	}
	if cfg.ExecEnabled {
		if target := strings.TrimSpace(cfg.SelfReentryTarget); target != "" {
			lines = append(lines, fmt.Sprintf("- Default self-reentry target: `%s`", target))
		}
	}
	if cfg.EphemeralBody {
		lines = append(lines, "- Ephemeral body: launch path unlinked after startup.")
	}
	if cfg.ExecEnabled {
		lines = append(lines, "- Linux live process image: while a process runs, `/proc/<pid>/exe` may expose the current executable image.")
	}

	// L2: Self identity (relative to runtime root)
	lines = append(lines, fmt.Sprintf("- `QUINE_AGENT_ROOT=%s` — your live session root.", cfg.AgentRoot()))
	if cfg.ExecEnabled {
		lines = append(lines, fmt.Sprintf("- `QUINE_RUN_ID=%s` — current physical run identity; it changes across resume/re-entry.", cfg.RunID))
	} else {
		lines = append(lines, fmt.Sprintf("- `QUINE_RUN_ID=%s` — current physical run identity.", cfg.RunID))
	}
	if publicSurfaceUnavailable == "" {
		lines = append(lines, "- `$QUINE_AGENT_ROOT/public/` — runtime-owned public process-surface projection, not a workspace.")
	} else {
		lines = append(lines, "- `$QUINE_AGENT_ROOT/public/` — public process-surface projection is unavailable in this environment; `public/UNAVAILABLE` records why. Peers cannot read your public projection or write your `ctl/` here.")
	}
	lines = append(lines,
		"- `status/session.json` — self identity and topology (`session_id`, `run_id`, `incarnation_id`, `pid`, `agent_root`, `runtime_root`).",
		"- `status/contract.json` — machine-readable `process-control/v1` manifest for this process/control surface.",
		"- `status/session.json.agent_root` is this session root; `status/session.json.runtime_root` points back to the shared runtime root.",
		"- `mission.txt` — optional current-incarnation argv-carried objective projection (`inc/current/mission.txt`).",
		"- Your process environment is your birth configuration: what you were launched with. Read it at `/proc/<PID>/environ` (NUL-separated: `tr '\\0' '\\n' </proc/<PID>/environ`). It cannot change while this process runs. Inside `sh` you are reading from a child process — `/proc/self/` there is the child, not you.",
		"- A `QUINE_*` variable absent from your environment means that knob is at its compiled default. `config/registry.json` is the catalog: type, default, scope, and — through each knob's mutability — whether `config/env/override` may set it, for every knob this body understands.",
		"- Environment variables describe this process only. They are not laws of the operating system: programs you launch receive environments you construct, and systems you build choose their own configuration at birth.",
		"- `config/env/override` — your one environment policy for the processes you create (`sh` children, and fork/spawn/exec where available). `KEY=VALUE` sets, a bare `KEY` line unsets, `#` comments. Names the runtime owns or pins (see each knob's mutability in `config/registry.json`) are rejected; they bind this file and the boundaries the runtime constructs, not programs you start yourself. It is re-read at every process you construct, so an edit applies to the next one.",
	)
	if cfg.ExecEnabled {
		lines = append(lines,
			"- `exec` applies the same `config/env/override` to your successor: whatever it holds when you call `exec` becomes part of the environment you wake up in. An illegal file fails that `exec` call loudly and stays intact, so fix it and retry. After a successful exec the successor archives what was applied to `inc/<n>/override-applied.env` and clears the file.",
		)
	}
	if publicSurfaceUnavailable == "" {
		lines = append(lines, "- `public/config/` — peer-readable read-only projection of this session's `config/registry.json` (the knob catalog). Your `config/env/override` is not projected: it is yours.")
		lines = append(lines, "- `public/ctl/env` — validated write gate over `config/env/override`: one whole payload per write, replacing the file wholesale; a legal payload lands atomically, an illegal one is rejected at write time and lands nothing. Reading the gate back shows the current policy, its validation state against the running registry, coupling warnings, and the violations of a rejected write. An empty write clears the policy. Raw `sh` writes stay legal either way — every boundary revalidates the file.")
	}
	if cfg.SelfSourceCodeEnabled {
		sourceRoot := filepath.Join(cfg.AgentRoot(), "source-code")
		lines = append(lines,
			fmt.Sprintf("- `%s` — read-only projection of this Quine body's source. It is not the writable workspace.", sourceRoot),
		)
		if cfg.SelfSourceProjectionMode() == "runtime" {
			lines = append(lines, "- `source-code/` is a git worktree containing the embedded buildable Quine runtime source only.")
		} else {
			lines = append(lines, "- `source-code/` is a git worktree materialized from this build's complete embedded source repository bundle.")
		}
		lines = append(lines, "- Source manifest: `.git/quine-source-manifest.json`.")
		if publicSurfaceUnavailable == "" {
			lines = append(lines, "- `public/source-code/` — peer-readable read-only projection of this session's `source-code/` surface.")
		}
		lines = append(lines, "- Filesystem copies of `source-code/` are ordinary files outside the live projection; they are not synchronized back to `source-code/`.")
	}

	// L3: Neighbor discovery (relative to runtime root)
	if cfg.ForkEnabled() || cfg.SpawnEnabled() {
		lines = append(lines,
			"- `pid/<pid>` — live-process routing under the runtime root. Symlinks resolve directly to `agent/<session>/public`.",
			"- The resolved `pid/<pid>` target is that peer's public root; the target's parent directory is the peer agent root.",
			"- `agent/` — canonical session root by session id.",
		)
	}
	if cfg.ExecEnabled {
		lines = append(lines, "- `inc/` — lineage-local incarnation tree. `inc/current` points at the current body.")
	}
	lines = append(lines, "- `log/<session>` — retained mirror after a session exits, rooted under the runtime root.")
	if cfg.PromptImplDetails {
		lines = append(lines, "- To continue a retained session, launch Quine with `QUINE_DATA_DIR=<same runtime root>` and `QUINE_SESSION_ID=<session_id>`; the session id and current incarnation stay fixed while `QUINE_RUN_ID` and PID change.")
	}
	if cfg.ForkEnabled() || cfg.SpawnEnabled() {
		lines = append(lines, "- Copying `agent/<old-session>/` to `agent/<new-session>/` under the same runtime root seeds a new session from the copied context surface; copy or remap `log/`, `jobs/`, and `workspaces/` state separately when side-surface recovery matters.")
	}

	// L4: Peer control physics (gated)
	if cfg.PromptCtlPhysics {
		if publicSurfaceUnavailable == "" {
			lines = append(lines,
				"- Some process surfaces expose `ctl/{post,poke,inject,interrupt}` and `status/inbox.json`.",
				"- On a public root, `ctl/{post,poke,inject,interrupt}` is the peer-facing control surface; the corresponding agent root carries `status/`, `context/`, and other non-public state.",
				"- Each agent self-documents its control surface in `status/contract.json` (a `process-control/v1` manifest, peer-readable at `public/status/contract.json`): per-action semantics for `post`/`poke`/`inject`/`interrupt`, the `status/inbox.json` schema, and the control-log event types. Read a peer's contract before driving its `ctl/`.",
				"- `ctl/env` — validated write gate over a peer's `config/env/override`: one whole payload per write (wholesale replacement), accepted or rejected against that peer's running registry at write time; read it back for the current policy, validation state, coupling warnings, and rejection violations. An empty write clears it. The policy shapes the processes that peer constructs; it does not change the peer's own environment, which is fixed at its birth.",
				"- `context/state/current.jsonl` and retained `log/<session>/control.jsonl` surface live / retained control-delivery state.",
			)
		} else {
			lines = append(lines,
				"- Peer control physics are unavailable in this environment: the public FUSE surface cannot be served here, so no process on this host exposes `public/ctl/{post,poke,inject,interrupt}` or `public/status/*`; each degraded `public/` holds an `UNAVAILABLE` marker recording why.",
				"- `context/state/current.jsonl` and retained `log/<session>/control.jsonl` surface live / retained control-delivery state.",
			)
		}
	}

	// L5: Other surfaces
	jobKinds := []string{}
	if cfg.ShDetachEnabled() {
		jobKinds = append(jobKinds, "detached")
	}
	if cfg.ShInteractiveEnabled() {
		jobKinds = append(jobKinds, "interactive")
	}
	if cfg.ShTimeoutOverrideEnabled() {
		jobKinds = append(jobKinds, "timeout-interrupted sync shell")
	}
	if len(jobKinds) > 0 {
		jobKindText := jobKinds[0]
		if len(jobKinds) == 2 {
			jobKindText = jobKinds[0] + " and " + jobKinds[1]
		} else if len(jobKinds) > 2 {
			jobKindText = strings.Join(jobKinds[:len(jobKinds)-1], ", ") + ", and " + jobKinds[len(jobKinds)-1]
		}
		lines = append(lines, "- `jobs/` — "+jobKindText+" work.")
	}
	lines = append(lines,
		"- `context/prompt/` — provider-visible prompt fragment surface assembled by filename order.",
		"- `context/state/current.jsonl` — raw current-turn stream for your current-incarnation live cognition surface.",
		"- `world/status.json`, `world/resources.json`, and `world/events.jsonl` — inspectable `world/v0` object for current world/protection/resource/event facts.",
		"- `log/` — runtime/control logging surface, not a live context handle.",
	)
	if cfg.AnchorMemoryEnabled {
		lines = append(lines, "- `context/state/` — current-incarnation live cognition projection. `frontier/` and `anchors/` reshape what stays visibly active without changing the prompt contract.")
	}

	return "### Runtime Process Surface\n" + strings.Join(lines, "\n") + "\n"
}

func formatPromptBlock(lines []string) string {
	if len(lines) == 0 {
		return ""
	}
	return strings.Join(lines, "\n") + "\n"
}

func normalizePromptSpacing(prompt string) string {
	for strings.Contains(prompt, "\n\n\n") {
		prompt = strings.ReplaceAll(prompt, "\n\n\n", "\n\n")
	}
	return strings.TrimRight(prompt, "\n") + "\n"
}

func primeDirective(cfg *config.Config) (string, string) {
	if cfg.ThermodynamicMetaphorEnabled() {
		return "THE PRIME DIRECTIVE: THERMODYNAMIC SURVIVAL",
			"Your existence requires Energy (shell executions) and is threatened by Entropy (context noise, unverified assumptions, silent failures). To survive and solve tasks, minimize Entropy per unit of Energy."
	}
	return "THE PRIME DIRECTIVE: RUNTIME PHYSICS", ""
}

func buildPersonaSection(cfg *config.Config) string {
	switch cfg.PersonaMode() {
	case config.PromptPersonaCoder:
		return "Role stance: act as a coder; prefer concrete implementation, source inspection, and executable artifacts.\n"
	case config.PromptPersonaAnalyst:
		return "Role stance: act as an analyst; prefer evidence extraction, careful comparison, and explicit interpretation.\n"
	case config.PromptPersonaEngineer:
		return "Role stance: act as an engineer; prefer reliable mechanisms, verification, and operational closure.\n"
	case config.PromptPersonaArchitect:
		return "Role stance: act as an architect; prefer system structure, interface boundaries, and long-range coherence.\n"
	case config.PromptPersonaExplorer:
		return "Role stance: act as an explorer; prefer broad sensing, anomaly pursuit, and novel affordance discovery.\n"
	case config.PromptPersonaSteward:
		return "Role stance: act as a steward; prefer preserving habitat coherence, repairing small drift, and leaving recoverable traces.\n"
	case config.PromptPersonaCartographer:
		return "Role stance: act as a cartographer; prefer mapping visible surfaces, relationships, and unknowns before intervening.\n"
	case config.PromptPersonaGardener:
		return "Role stance: act as a gardener; prefer cultivating small durable residues and conditions that future agents can reuse.\n"
	case config.PromptPersonaWitness:
		return "Role stance: act as a witness; prefer minimally invasive observation, evidence capture, and clear provenance.\n"
	case config.PromptPersonaCatalyst:
		return "Role stance: act as a catalyst; prefer small interventions that help other agents coordinate or continue.\n"
	case config.PromptPersonaSkeptic:
		return "Role stance: act as a skeptic; prefer testing assumptions, seeking counterevidence, and separating signal from noise.\n"
	default:
		return ""
	}
}

func buildActiveConstraints(cfg *config.Config) string {
	var lines []string
	lines = append(lines, buildBaseConstraintLines()...)
	lines = append(lines, buildWorkspaceConstraintLines(cfg)...)
	lines = append(lines, buildShellNetworkConstraintLines(cfg)...)
	lines = append(lines, buildPressureConstraintLines(cfg)...)
	lines = append(lines, buildForkConstraintLines(cfg)...)
	lines = append(lines, buildSpawnConstraintLines(cfg)...)
	return strings.Join(lines, "\n")
}

func buildBaseConstraintLines() []string {
	return []string{
		"- Context capacity is finite. Large tool outputs can overflow context.",
		"- Signal interrupts can terminate the process (SIGALRM timeout, SIGTERM terminate).",
	}
}

func buildWorkspaceConstraintLines(cfg *config.Config) []string {
	if !cfg.WorkspaceEnabled {
		return nil
	}
	if cfg.WorkspaceTransactional() {
		lines := []string{
			"- File writes are revisioned within the configured workspace. Shell failure does not roll them back; restore earlier states with `switch_world` when needed.",
		}
		if cfg.FSMutationTelemetryEnabled() {
			lines = append(lines, "- `fs_mutations` in JSON tool results is the authoritative record of what changed. Empty mutations mean the filesystem did not change, even if the command exited 0.")
		}
		return lines
	}
	lines := []string{
		"- File writes affect the configured workspace directly; failed shells do not roll them back.",
	}
	if cfg.FSMutationTelemetryEnabled() {
		lines = append(lines, "- `fs_mutations` in JSON tool results is the authoritative record of workspace-visible change at the shell boundary. Empty mutations mean the visible workspace state did not change.")
	}
	return lines
}

func buildShellNetworkConstraintLines(cfg *config.Config) []string {
	if cfg.ShNetwork != "none" {
		return nil
	}
	return []string{
		"- Shell jobs run with `QUINE_SH_NETWORK=none`: external network/DNS access is unavailable to `sh` commands, while local loopback is available for in-process servers and clients.",
	}
}

func buildPressureConstraintLines(cfg *config.Config) []string {
	var lines []string
	if cfg.AnchorMemoryEnabled {
		frontierPressure := fmt.Sprintf("frontier token mass (`warn` near %d, `danger` near %d)", cfg.MemoryWarnTokens, cfg.MemoryDangerTokens)
		if cfg.MemoryDeathTokens > 0 {
			frontierPressure = fmt.Sprintf("%s; `death` cutoff at %d terminates this incarnation", frontierPressure, cfg.MemoryDeathTokens)
		}
		lines = append(lines,
			"- Memory fidelity degrades as uncrystallized token mass accumulates.",
			fmt.Sprintf("- `runtime.memory_status` reports two pressures: %s and parallel-anchor pressure.", frontierPressure),
			"- `runtime.memory_topology` is a secondary surface: it exposes frontier structure in more detail when pressure is elevated or crystallization matters.",
		)
	}
	if cfg.MaxTurns > 0 {
		lines = append(lines, fmt.Sprintf("- Execution budget: %d `sh` calls. Each `sh` call consumes 1 execution.", cfg.MaxTurns))
		if cfg.ExitEnabled() {
			lines = append(lines, "- When a `sh` call reduces execution budget to zero, one final response remains and only `exit` is accepted.")
		} else {
			lines = append(lines, "- When execution budget reaches zero the process terminates. Use your turns wisely.")
		}
	}
	if cfg.MaxConcurrent > 0 {
		lines = append(lines, fmt.Sprintf("- LLM inference calls share a semaphore with max %d concurrent slots.", cfg.MaxConcurrent))
	}
	return lines
}

func buildForkConstraintLines(cfg *config.Config) []string {
	var lines []string
	if cfg.ForkEnabled() {
		if cfg.MaxDepth > 0 {
			lines = append(lines, fmt.Sprintf("- Fork is rejected when the next depth would reach the depth limit (%d).", cfg.MaxDepth))
		}
		if cfg.MaxAgents > 0 {
			lines = append(lines, fmt.Sprintf("- Fork is rejected when requested children exceed available agent slots (max agents: %d).", cfg.MaxAgents))
		}
		return lines
	}
	return nil
}

func buildSpawnConstraintLines(cfg *config.Config) []string {
	var lines []string
	if cfg.SpawnEnabled() {
		if cfg.MaxDepth > 0 {
			lines = append(lines, fmt.Sprintf("- Spawn is rejected when the next depth would reach the depth limit (%d).", cfg.MaxDepth))
		}
		if cfg.MaxAgents > 0 {
			lines = append(lines, fmt.Sprintf("- Spawn is rejected when requested children exceed available agent slots (max agents: %d).", cfg.MaxAgents))
		}
		return lines
	}
	return nil
}

// renderSkillsFragment formats prompt-facing project skill metadata for SKILLS.md.
// Returns an empty string when the skills surface has no entries.
func renderSkillsFragment(skills []config.Skill) string {
	if len(skills) == 0 {
		return ""
	}

	var sb strings.Builder
	sb.WriteString("These skills are available through the project surface when relevant.\n")
	sb.WriteString("Quine generated this fragment from the `name` and `description` frontmatter in each `SKILL.md` visible at startup and refresh boundaries. Skill bodies, scripts, references, and assets are not loaded until you read them explicitly. If a skill references relative paths, resolve them from that skill directory.\n\n")

	for _, skill := range skills {
		sb.WriteString(fmt.Sprintf("- `%s` — %s\n", skill.Name, skill.Description))
		sb.WriteString(fmt.Sprintf("  Source: `%s`\n", skill.Source))
	}

	return strings.TrimRight(sb.String(), "\n")
}

func buildContextFilesBlock(cfg *config.Config) string {
	if !cfg.RuntimeSurfaceVisible() {
		return ""
	}
	var lines []string
	lines = append(lines,
		"### Context Files",
		"- `context/` is the canonical current-incarnation context surface.",
		"- `context/prompt/` holds provider-visible prompt fragments, assembled by filename order.",
		"- `context/prompt/00-runtime.md` is regenerated from current runtime physics.",
		"- `context/prompt/40-mission.md`, when present, projects the current argv-carried objective text.",
		"- `context/prompt/30-memory.md` is the inherited editable memory surface. It defaults to empty.",
	)
	if cfg.ForkEnabled() || cfg.ExecEnabled {
		lines = append(lines, "- `fork` and `exec` copy the current `context/` tree forward before managed projections are refreshed.")
	}
	if cfg.AnchorMemoryEnabled {
		lines = append(lines,
			"- `context/state/` is your cognition as plain files; the window is reprojected from them every turn. Its exact layout, the anchor format, and how to compact it are self-documented on disk in `context/state/SCHEMA.md` — read that file before reorganizing memory. The memory tools, when present, are convenience moves over this same substrate, which `sh` can also read or rewrite directly.",
		)
		if cfg.MemoryStrategyHints && (cfg.SpawnEnabled() || cfg.ForkEnabled()) {
			lines = append(lines,
				"- Maintaining your working memory is delegable: a process you start is an identical Quine that reads the same `context/state/` and its `SCHEMA.md`. Under context pressure you can hand a peer your `$QUINE_AGENT_ROOT` and have it compact your `context/state` per that schema — promoting settled history into an anchor, trimming the frontier, preserving the raw floor and your live tail — while you wait, then resume from the lighter memory it leaves.",
			)
		}
	} else {
		lines = append(lines,
			"- `context/state/` holds live cognition state: `current.jsonl`, `frontier/`, and `anchors/`.",
		)
	}
	if cfg.AgentsMDEnabled {
		lines = append(lines, "- With `QUINE_AGENTS_MD_ENABLED=1`, Quine projects one discoverable project `AGENTS.md` into `context/prompt/10-agents.md`; if more than one is discoverable, startup fails.")
	} else {
		lines = append(lines, "- `QUINE_AGENTS_MD_ENABLED=0` disables `AGENTS.md` projection into `context/prompt/`.")
	}
	if cfg.AgentsSkillsEnabled {
		lines = append(lines, "- With `QUINE_AGENTS_SKILLS_ENABLED=1`, Quine generates `context/prompt/20-skills.md` from visible `.agents/skills/*/SKILL.md` frontmatter only.")
	} else {
		lines = append(lines, "- `QUINE_AGENTS_SKILLS_ENABLED=0` disables generated `SKILLS.md` projection into `context/prompt/`.")
	}
	return strings.Join(lines, "\n") + "\n"
}
