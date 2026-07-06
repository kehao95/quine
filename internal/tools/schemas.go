package tools

import (
	"strings"

	"github.com/kehao95/quine/internal/config"
	"github.com/kehao95/quine/internal/llm"
)

// ShToolSchema returns the JSON Schema for the sh tool.
func ShToolSchema(cfg *config.Config) llm.ToolSchema {
	interactiveEnabled := cfg == nil || cfg.ShInteractiveEnabled()
	timeoutOverrideEnabled := cfg == nil || cfg.ShTimeoutOverrideEnabled()
	stdinEnabled := cfg == nil || cfg.ShStdinEnabled()
	detachEnabled := cfg == nil || cfg.ShDetachEnabled()
	fsTelemetryEnabled := cfg == nil || cfg.FSMutationTelemetryEnabled()
	stdinDescription := "Verbatim stdin for the command. Quine materializes it as finite prewritten input, not a live stream. This supplies multi-line content without shell heredoc or quoting mechanics."
	noKillClause := "It does not kill the job."
	timeoutProcessClause := "the process group is stopped, not killed"
	if detachEnabled {
		noKillClause = "It does not kill or detach the job."
		timeoutProcessClause = "the process group is stopped, not killed or detached"
	}
	timeoutDescription := "Optional sync-shell protection timeout in seconds. In direct/host mode, if a synchronous `sh` call exceeds this bound, Quine stops the whole process group with `SIGSTOP` and returns `status=\"interrupted\"` plus `job.pid`, `job.path`, and `*_so_far` snapshots. " + noKillClause + " Under `overlay`, a sync timeout terminates the shell job, preserves workspace side effects that reached the shell boundary as a world revision, and returns snapshots; it is not resumable."
	switch {
	case detachEnabled && interactiveEnabled:
		timeoutDescription += " Detached and interactive jobs ignore this field."
	case detachEnabled:
		timeoutDescription += " Detached jobs ignore this field."
	case interactiveEnabled:
		timeoutDescription += " Interactive jobs ignore this field."
	}
	resultFields := "`exit_code`, `stdout`, `stderr`, and `world_revision` when available"
	timeoutFields := "`job.pid`, `job.path`, `world_revision` when available, and `*_so_far` snapshots"
	workspaceFields := "`world_revision`"
	if fsTelemetryEnabled {
		resultFields = "`exit_code`, `stdout`, `stderr`, `fs_mutations`, and `world_revision` when available"
		timeoutFields = "`job.pid`, `job.path`, `stdout_so_far`, `stderr_so_far`, `fs_mutations_so_far`, and `world_revision` when available"
		workspaceFields = "`fs_mutations` and `world_revision`"
	}
	timeoutClause := ""
	if timeoutOverrideEnabled {
		timeoutClause = " The optional `timeout` parameter sets a per-call sync-shell bound; in direct/host mode, an exceeded bound returns `status=\"interrupted\"` with " + timeoutFields + "; " + timeoutProcessClause + ". Under `overlay`, an exceeded bound terminates the shell job, preserves workspace side effects that reached the shell boundary as a world revision, and returns snapshots; it is not resumable."
	}
	interactiveImplClause := ""
	if cfg != nil && cfg.PromptImplDetails {
		interactiveImplClause = " Under `overlay`, interactive jobs run in a job-local workspace lineage and write `world_handle` after exit for optional `switch_world` adoption."
	}
	description := "Execute a POSIX shell command. Costs 1 execution. Default behavior waits for completion and returns bounded output. Tool results are compact JSON in `tool_result.content`; completed `sh` results expose fields such as " + resultFields + "." + timeoutClause
	jobClauses := []string{}
	if detachEnabled {
		jobClauses = append(jobClauses, "`detach=true` keeps a managed background job under `${QUINE_DATA_DIR}`")
	}
	if interactiveEnabled {
		jobClauses = append(jobClauses, "`interactive=true` keeps a PTY-backed job under `${QUINE_DATA_DIR}`")
	}
	if len(jobClauses) > 0 {
		description += " " + strings.Join(jobClauses, "; ") + ". `${QUINE_DATA_DIR}` is Quine runtime state, not the task workspace."
	}
	description += " When workspace physics are enabled, filesystem changes persist across `sh` calls in the provisional workspace and results include " + workspaceFields + "."
	if interactiveEnabled {
		if stdinEnabled {
			stdinDescription += " Mutually exclusive with `interactive=true`."
		}
		description += interactiveImplClause
	}
	if detachEnabled {
		description += " `detach=true` is unavailable under `overlay`."
	}
	detachDescription := "If true, keep the command as a managed background job under `${QUINE_DATA_DIR}` and return its absolute job directory immediately. Default: false."
	interactiveDescription := "If true, keep the command as a PTY-backed interactive POSIX job under `${QUINE_DATA_DIR}` and return its absolute job directory immediately."
	interactiveExclusions := []string{}
	if stdinEnabled {
		interactiveExclusions = append(interactiveExclusions, "`stdin`")
	}
	if detachEnabled {
		interactiveExclusions = append(interactiveExclusions, "`detach`")
	}
	if len(interactiveExclusions) > 0 {
		interactiveDescription += " Mutually exclusive with " + strings.Join(interactiveExclusions, " and ") + "."
	}
	interactiveDescription += " Default: false."
	if cfg != nil && cfg.PromptImplDetails {
		detachDescription = "If true, keep the command as a managed background job under `${QUINE_DATA_DIR}` and return its absolute job directory immediately. Detached job directories include `cmd`, `pid`, `started_at`, `out.log`, and `err.log`; after the job completes, its exit code appears at `<path>/exit` (or, if that write fails, the error is recorded at `<path>/exit_error`). Default: false."
		interactiveDescription = "If true, keep the command as a PTY-backed interactive POSIX job under `${QUINE_DATA_DIR}` and return its absolute job directory immediately. Interactive job directories include `cmd`, `pid`, `started_at`, screen snapshots, `in`, `winsize`, `events.log`, `events.hex`, and `input.log`; after the job completes, its exit code appears at `<path>/exit` (or, if that write fails, the error is recorded at `<path>/exit_error`). Under `overlay`, the job runs in a job-local workspace lineage; after exit, `workspace_session`, `fs_mutations` when enabled, `world_revision`, and `world_handle` appear so the world can be adopted with `switch_world`. Send POSIX signals with `kill` and the recorded pid/process group."
		if len(interactiveExclusions) > 0 {
			interactiveDescription += " Mutually exclusive with " + strings.Join(interactiveExclusions, " and ") + "."
		}
		interactiveDescription += " Default: false."
	}
	props := map[string]any{
		"command": map[string]any{
			"type":        "string",
			"description": "The shell command to execute.",
		},
	}
	if detachEnabled {
		props["detach"] = map[string]any{
			"type":        "boolean",
			"description": detachDescription,
		}
	}
	if stdinEnabled {
		props["stdin"] = map[string]any{
			"type":        "string",
			"description": stdinDescription,
		}
	}
	if timeoutOverrideEnabled {
		props["timeout"] = map[string]any{
			"type":        "integer",
			"description": timeoutDescription,
			"minimum":     0,
		}
	}
	if interactiveEnabled {
		props["interactive"] = map[string]any{
			"type":        "boolean",
			"description": interactiveDescription,
		}
	}
	required := []string{"command"}

	return llm.ToolSchema{
		Name:        "sh",
		Description: description,
		Parameters: map[string]any{
			"type":       "object",
			"properties": props,
			"required":   required,
		},
	}
}

// ForkToolSchema returns the JSON Schema for the fork tool.
func ForkToolSchema(cfg *config.Config) llm.ToolSchema {
	fsTelemetryEnabled := cfg == nil || cfg.FSMutationTelemetryEnabled()
	desc := "Spawn one or more child agents with the parent's current visible context (swarm fork). " +
		"Each child is another you under a different intent. " +
		"`fork` returns each child's process handles plus compact JSON fields such as `exit_code`, `stdout`, and `stderr`. " +
		"Child filesystem effects follow the configured workspace/world lineage; stdout/stderr are only captured process output. " +
		"Does not consume execution budget. "
	if cfg != nil && cfg.PromptImplDetails {
		desc += "Use fork when independent hypotheses, decoders, implementations, extractors, or verification strategies can be tried in parallel. " +
			"When no contract-shaped artifact or service is stable and 2-3 plausible approaches exist, prefer 2-3 labeled children over another long parent-only inspection. " +
			"Fork children preserve the parent mission as the active task contract; each child intent is a lane assignment, not a replacement mission. " +
			"Child intents should include lane-specific inputs only when they are not already in the parent mission or current visible context. " +
			"Do one cheap shared setup/probe if all lanes need it; then fork specialized heavyweight installs, downloads, transcription, OCR, builds, searches, or long-running probes. Do not make every child repeat the same setup. "
	}
	desc += "`mode=\"wait\"` blocks until all children finish and returns every result. `race` returns the first successful child and stops the rest. `forget` returns after spawning children without waiting for completion. "
	if fsTelemetryEnabled {
		desc += "When workspace physics are enabled and the fork completes before returning, results also include parent-side `fs_mutations` and `world_revision`; under subjective children these are normally empty or unchanged because child writes stay in child lineage. "
	}
	childrenDesc := "Array of child specs. Each child must provide `intent` and `scope`."
	scopeDesc := "Child scope path. `.` means inherit your current scope."
	intentDesc := "Concrete child intent. The parent mission remains active; state this child's lane-specific focus."
	if cfg != nil && cfg.PromptImplDetails {
		intentDesc = "Concrete lane assignment for this child. The parent mission remains active; name a distinct strategy, lane-specific inputs not already visible from the parent mission/context, the expected artifact/service when known, and the closest success check."
	}
	childProps := map[string]any{
		"intent": map[string]any{
			"type":        "string",
			"description": intentDesc,
		},
		"scope": map[string]any{
			"type":        "string",
			"description": scopeDesc,
		},
	}
	required := []string{"intent", "scope"}
	if cfg != nil && cfg.ForkWorldEnabled {
		subjectiveDesc := "`subjective` (default) starts a scoped subjective child; private lineage is guaranteed only under the `overlay` workspace backend. Under `direct`, writes remain host-visible and non-adoptable."
		protectionDesc := "`transactional` (default) requests subjective-world transactional protection. It provides private lineage only under the `overlay` workspace backend; under `direct`, the request is a scoped shared-workspace surface. `none` leaves the child on the host surface without transactional protection. Only `subjective + transactional` and `host + none` are currently supported."
		if cfg.WorkspaceTransactional() {
			subjectiveDesc = "`subjective` (default) starts a private subjective world under the `overlay` workspace backend."
			protectionDesc = "`transactional` (default) gives subjective-world transactional protection under the `overlay` workspace backend. `none` leaves the child on the host surface without transactional protection. Only `subjective + transactional` and `host + none` are currently supported."
		}
		desc += "Each child can choose `world=\"subjective\"` (default, backend-qualified subjective world) or `world=\"host\"` (host-side child). Each child can also choose `protection=\"transactional\"` (default) or `protection=\"none\"`. In this runtime, only `subjective + transactional` and `host + none` are legal pairs; private lineage requires the `overlay` workspace backend. "
		childrenDesc = "Array of child specs. Each child must provide `intent`. `world` defaults to `subjective`; `protection` defaults to `transactional`; `scope` defaults to `.` and is only meaningful for `world=\"subjective\"`."
		scopeDesc = "Child scope path, only meaningful for `world=\"subjective\"`. Defaults to `.`. Valid values are `.` or a relative subpath within the child's current root; absolute paths are invalid. For `world=\"host\"`, `scope` has no effect."
		childProps["world"] = map[string]any{
			"type":        "string",
			"enum":        []string{"host", "subjective"},
			"description": "Child world substrate. " + subjectiveDesc + " `host` starts a host-side child.",
		}
		childProps["protection"] = map[string]any{
			"type":        "string",
			"enum":        []string{"none", "transactional"},
			"description": "Child protection surface. " + protectionDesc,
		}
		childProps["scope"] = map[string]any{
			"type":        "string",
			"description": scopeDesc,
		}
		required = []string{"intent"}
	} else if cfg != nil && !cfg.WorkspaceEnabled {
		desc += "Children share the same filesystem surface as the parent and siblings. "
		scopeDesc = "Child scope placeholder. `.` is the only meaningful value when workspace physics are not enabled."
		childProps["scope"] = map[string]any{
			"type":        "string",
			"description": scopeDesc,
		}
	} else {
		desc += "Under Linux workspace physics, each child starts from your current world in its own process-local lineage. `overlay` is the workspace backend. "
		childrenDesc = "Array of child specs. Each child must provide `intent` and `scope`. `scope` must be `.` to inherit your current scope or a relative subpath under it."
		scopeDesc = "Child scope path. Must be `.` or a relative path under your current scope; absolute paths are invalid. `.` means inherit your current scope. Relative paths resolve from your current scope and must remain within the workspace root. The child's shell starts with cwd at that child scope. `scope` narrows working area inside the child lineage; it is not the lineage identity itself."
		childProps["scope"] = map[string]any{
			"type":        "string",
			"description": scopeDesc,
		}
	}
	return llm.ToolSchema{
		Name:        "fork",
		Description: desc + "Supports parallel exploration, delegation, and task decomposition.",
		Parameters: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"children": map[string]any{
					"type":        "array",
					"description": childrenDesc,
					"items": map[string]any{
						"type":       "object",
						"properties": childProps,
						"required":   required,
					},
					"minItems": 1,
				},
				"mode": map[string]any{
					"type":        "string",
					"enum":        []string{"wait", "race", "forget"},
					"description": "wait: block until all children finish and return every result, for comparison or merge. race (default): first child to exit 0 wins and the rest are killed; use when any one child can produce an acceptable artifact/service. forget: fire-and-forget, return PIDs immediately.",
				},
				"adopt_winner": map[string]any{
					"type":        "boolean",
					"description": "Only valid with `mode=\"race\"`. Set true for overlay subjective races when the winner's filesystem artifact should become parent state; if false or non-adoptable, inspect/copy/merge the best child result manually.",
				},
			},
			"required": []string{"children"},
		},
	}
}

// SpawnToolSchema returns the JSON Schema for the spawn tool.
func SpawnToolSchema(cfg *config.Config) llm.ToolSchema {
	childProps := map[string]any{
		"mission": map[string]any{
			"type":        "string",
			"description": "Mission string passed as argv to the fresh child Quine process.",
		},
	}
	childrenDesc := "Array of fresh child process specs. Each child must provide `mission`."
	desc := "Start one or more fresh-context Quine child processes from the configured binary. " +
		"Unlike fork, spawn does not import the parent's active context, seed, or anchor-memory surface; each child starts from normal runtime initialization with only its explicit mission. " +
		"Workspace, world, scope, protection, relation, capture, and wait/race/forget semantics match fork unless stated otherwise. " +
		"Use spawn when inherited assumptions would weaken independent critique while the current workspace artifact still matters. " +
		"The child process still belongs to the same runtime process tree for depth, agent-slot, retained relation, and environment governance. " +
		"Does not consume execution budget. `mode=\"wait\"` blocks until all children finish and returns every result. `race` returns the first successful child and stops the rest. `forget` returns after spawning children without waiting for completion."
	if cfg != nil && cfg.ForkWorldEnabled {
		childrenDesc = "Array of fresh child process specs. Each child must provide `mission`. `world` defaults to `subjective`; `protection` defaults to `transactional`; `scope` defaults to `.` and is only meaningful for `world=\"subjective\"`."
		childProps["world"] = map[string]any{
			"type":        "string",
			"enum":        []string{"host", "subjective"},
			"description": "Child world substrate. `subjective` (default) starts from the selected workspace lineage under the configured workspace backend. `host` starts a host-side fresh child.",
		}
		childProps["protection"] = map[string]any{
			"type":        "string",
			"enum":        []string{"none", "transactional"},
			"description": "Child protection surface. Only `subjective + transactional` and `host + none` are currently supported.",
		}
		childProps["scope"] = map[string]any{
			"type":        "string",
			"description": "Child scope path, only meaningful for `world=\"subjective\"`. Defaults to `.`. Valid values are `.` or a relative subpath within the child's current root; absolute paths are invalid.",
		}
	} else if cfg != nil && cfg.WorkspaceEnabled {
		desc += " Under Linux workspace physics, spawned children share the same workspace projection rules as fork and can use `scope` to start under a relative workspace subpath. "
		childProps["scope"] = map[string]any{
			"type":        "string",
			"description": "Child scope path. Defaults to `.`. Must be `.` or a relative path under your current scope; absolute paths are invalid.",
		}
	}
	return llm.ToolSchema{
		Name:        "spawn",
		Description: desc,
		Parameters: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"children": map[string]any{
					"type":        "array",
					"description": childrenDesc,
					"items": map[string]any{
						"type":       "object",
						"properties": childProps,
						"required":   []string{"mission"},
					},
					"minItems": 1,
				},
				"mode": map[string]any{
					"type":        "string",
					"enum":        []string{"wait", "race", "forget"},
					"description": "wait (default): block until all children finish and return every result. race: first child to exit 0 wins and the rest are killed. forget: fire-and-forget, return PIDs immediately.",
				},
			},
			"required": []string{"children"},
		},
	}
}

// ExecToolSchema returns the JSON Schema for the exec tool.
func ExecToolSchema(cfg *config.Config) llm.ToolSchema {
	desc := "Replace the current process image with a new executable. " +
		"Default behavior is sugar for re-execing quine through its configured self-reentry target; it preserves the current mission when one exists and uses no mission argv when none exists. " +
		"When both `target` and `argv` are omitted during quine self re-exec, the current mission state is preserved. " +
		"Providing `argv` explicitly replaces the mission passed to the new image. " +
		"A different target binary and explicit argv can be supplied. " +
		"If present, `config/next.env` is validated and applied to the successor's environment at this boundary; the runtime prompt owns the full staged-config mechanism."

	props := map[string]any{
		"target": map[string]any{
			"type":        "string",
			"description": "Optional executable path or name. If omitted, exec defaults to quine's configured self-reentry target. The runtime surface shows the launch path separately. Omission preserves the current mission state under default self re-exec. When the ephemeral body has been unlinked, a `target` of /proc/self/exe or /proc/<pid>/exe is rejected: re-executing the live process image recovers the original body instead of reconstructing a successor.",
		},
		"argv": map[string]any{
			"type":        "array",
			"description": "Optional full argv vector. If omitted: the self re-exec default is sugar for [configured self-reentry target, current mission] when a mission exists and [configured self-reentry target] when no mission exists, while an external target defaults to [target]. Supplying `argv` overrides that default and replaces the mission the new image receives.",
			"items": map[string]any{
				"type": "string",
			},
		},
	}
	return llm.ToolSchema{
		Name:        "exec",
		Description: desc,
		Parameters: map[string]any{
			"type":       "object",
			"properties": props,
			"required":   []string{},
		},
	}
}

func SwitchWorldToolSchema(cfg *config.Config) llm.ToolSchema {
	desc := "Switch the current provisional workspace view to a world target. Does not consume a turn. Results include `world_revision` for the switch turn."
	if cfg == nil || cfg.FSMutationTelemetryEnabled() {
		desc = "Switch the current provisional workspace view to a world target. Does not consume a turn. Results include `fs_mutations` and `world_revision` for the switch turn."
	}
	if cfg == nil || !cfg.CanRestoreWorld() {
		desc += " Requires Linux workspace physics with workspace revision mode `restore`."
	} else {
		desc += " Targets include local revision handles such as `wr0` or `wr3`, and subjective child world handles such as `world://<workspace-session>/<revision>`. `wr0` is the baseline world; later revisions are historical checkpoints captured after completed shell turns."
	}
	return llm.ToolSchema{
		Name:        "switch_world",
		Description: desc,
		Parameters: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"target": map[string]any{
					"type":        "string",
					"description": "Switch the workspace to a specific local revision or adopt a subjective child world. Examples: target=\"wr0\" restores the baseline world; target=\"wr3\" restores the world captured after the third completed shell turn; target=\"world://child-session/wr4\" adopts that child world into the current session and switches to it.",
				},
			},
			"required": []string{"target"},
		},
	}
}

// ExitToolSchema returns the JSON Schema for the exit tool.
func ExitToolSchema(cfg *config.Config) llm.ToolSchema {
	statusEnum := []string{"success", "failure"}
	description := "Finish your work and terminate. " +
		"Two modes: success (task complete), failure (task failed). " +
		"NOTE: This tool emits no stdout bytes. If stdout bytes are required, produce them before exit with `sh` writing to /dev/stdout."
	properties := map[string]any{
		"status": map[string]any{
			"type":        "string",
			"enum":        statusEnum,
			"description": "Task outcome. \"success\" = complete. \"failure\" = failed.",
		},
		"stderr": map[string]any{
			"type":        "string",
			"description": "Why the task failed. Required on failure. Must NOT be set on success.",
		},
	}
	if cfg != nil && !cfg.FailOnImpossible {
		statusEnum = []string{"success"}
		description = "Finish your work and terminate successfully. " +
			"In this runtime only `status=\"success\"` is available; blocked paths must continue recovering instead of explicitly exiting as failure. " +
			"NOTE: This tool emits no stdout bytes. If stdout bytes are required, produce them before exit with `sh` writing to /dev/stdout."
		properties = map[string]any{
			"status": map[string]any{
				"type":        "string",
				"enum":        statusEnum,
				"description": "Task outcome. Only \"success\" is available in this runtime.",
			},
		}
	}
	return llm.ToolSchema{
		Name:        "exit",
		Description: description,
		Parameters: map[string]any{
			"type":       "object",
			"properties": properties,
			"required":   []string{"status"},
		},
	}
}

// IdleToolSchema returns the JSON Schema for the idle tool.
func IdleToolSchema() llm.ToolSchema {
	return llm.ToolSchema{
		Name:        "idle",
		Description: "Suspend explicitly until an external control event resumes you. Does not consume execution budget. `idle` resumes on `poke`, `inject`, or `interrupt`; `poke` resumes without automatic context injection, while `inject` and payload-bearing `interrupt` can surface `incoming_messages` at the next safe point. The runtime prompt owns the full peer-control surface map.",
		Parameters: map[string]any{
			"type":                 "object",
			"properties":           map[string]any{},
			"required":             []string{},
			"additionalProperties": false,
		},
	}
}

// VisionToolSchema returns the JSON Schema for the vision tool.
func VisionToolSchema() llm.ToolSchema {
	return llm.ToolSchema{
		Name: "vision",
		Description: "Analyze an image file. Reads the file at the given path and delivers it " +
			"as a native image to the model's vision capabilities. " +
			"Supported formats: PNG, JPEG, GIF, WebP, PPM, and PGM. PPM/PGM inputs are converted to PNG before delivery. Does NOT consume a turn. " +
			"WARNING: Vision analysis can be imprecise, especially for fine-grained details. " +
			"Cross-verify extracted information with programmatic tools when accuracy is critical.",
		Parameters: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"path": map[string]any{
					"type":        "string",
					"description": "Path to the image file (absolute or relative to cwd).",
				},
				"instruction": map[string]any{
					"type":        "string",
					"description": "A text instruction describing what information to extract from the image.",
				},
			},
			"required": []string{"path"},
		},
	}
}

// MarkToolSchema returns the JSON Schema for the mark tool. When foldEnabled is
// false the `fold` consolidation move is omitted from both the description and
// the parameter set, leaving plain marks plus unfold over the same file
// substrate.
func MarkToolSchema(foldEnabled bool) llm.ToolSchema {
	properties := map[string]any{
		"resolution": map[string]any{
			"type":        "string",
			"description": "Low-entropy capture of what just became stably true in the current working set. Name the resolved result that just closed a subproblem or stabilized a pivot; keep it to the result, not plans.",
		},
	}
	description := "Crystallize the current memory frontier into an immutable anchor. Does not consume execution budget. Plain marks are low-cost working-memory checkpoints. Memory telemetry can point toward `mark`."
	if foldEnabled {
		description += " `fold=true` is the higher-order frontier reconfiguration move: it absorbs one or more earlier working anchors into remembered background when a newer resolution makes them background rather than active focus, including after `fork(mode=\"wait\")` returns several stable child findings that are now being compressed into one parent-level conclusion. The new anchor becomes the parent session's governing organizing point. Two parallel anchors alone do not require a fold. Memory telemetry can point toward `fold`; actual fold semantics depend on the newest resolution subsuming earlier anchors."
		properties["fold"] = map[string]any{
			"type":        "boolean",
			"description": "If true, consolidate one or more earlier frontier anchors into remembered background under this new higher-order resolution. Fold semantics apply when the current resolution makes those earlier anchors background rather than active focus, including when memory telemetry points toward `fold` or when several stable child findings have returned and are now being subsumed by one parent-level synthesis. The new anchor becomes the parent session's governing organizing point. Two active anchors by themselves are not enough reason to fold. Actual consolidation requires existing prior anchors; without them, the call records only the new anchor and reports `absorbed=0`.",
		}
	}
	return llm.ToolSchema{
		Name:        "mark",
		Description: description,
		Parameters: map[string]any{
			"type":       "object",
			"properties": properties,
			"required":   []string{"resolution"},
		},
	}
}

// UnfoldToolSchema returns the JSON Schema for the unfold tool.
func UnfoldToolSchema() llm.ToolSchema {
	return llm.ToolSchema{
		Name:        "unfold",
		Description: "Read one anchor and return its one-level structured view.",
		Parameters: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"anchor_id": map[string]any{
					"type":        "integer",
					"description": "Anchor identifier to inspect.",
				},
			},
			"required": []string{"anchor_id"},
		},
	}
}

// AllToolSchemas returns all tool schemas.
func AllToolSchemas(cfg *config.Config) []llm.ToolSchema {
	schemas := []llm.ToolSchema{
		ShToolSchema(cfg),
	}
	if cfg == nil || cfg.ForkEnabled() {
		schemas = append(schemas, ForkToolSchema(cfg))
	}
	if cfg != nil && cfg.SpawnEnabled() {
		schemas = append(schemas, SpawnToolSchema(cfg))
	}
	if cfg == nil || cfg.ExitEnabled() {
		schemas = append(schemas, ExitToolSchema(cfg))
	}
	if cfg != nil && cfg.IdleToolEnabled() {
		schemas = append(schemas, IdleToolSchema())
	}
	if cfg == nil || cfg.ExecEnabled {
		schemas = append(schemas, ExecToolSchema(cfg))
	}
	if cfg == nil || cfg.VisionEnabled {
		schemas = append(schemas, VisionToolSchema())
	}
	if cfg != nil && cfg.CanRestoreWorld() {
		schemas = append(schemas, SwitchWorldToolSchema(cfg))
	}
	if cfg != nil && cfg.AnchorMemoryEnabled {
		if cfg.AnchorMarkEnabled() {
			schemas = append(schemas, MarkToolSchema(cfg.AnchorFoldEnabled()))
		}
		schemas = append(schemas, UnfoldToolSchema())
	}
	return schemas
}
