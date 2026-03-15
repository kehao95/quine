package tools

import (
	"github.com/kehao95/quine/internal/config"
	"github.com/kehao95/quine/internal/llm"
)

// ShToolSchema returns the JSON Schema for the sh tool.
// When cfg.CanEscalate() is true, goal/strategy params are added for STALL detection.
func ShToolSchema(cfg *config.Config) llm.ToolSchema {
	props := map[string]any{
		"command": map[string]any{
			"type":        "string",
			"description": "The shell command to execute.",
		},
		"detach": map[string]any{
			"type":        "boolean",
			"description": "If true, start the command in a managed background session and return immediately with an absolute filesystem-backed job directory path under `${QUINE_DATA_DIR}`. `${QUINE_DATA_DIR}` is Quine's runtime-state root, not the task workspace. Use this for daemons, overlap, or anything that must survive beyond the current `sh` call. `<path>/cmd` stores the original command, `<path>/out.log` and `<path>/err.log` store output, and `cat <path>/exit` waits for completion. In every sh environment, `$QUINE_JOB_SESSION_DIR` points to `${QUINE_DATA_DIR}/jobs/<session>`. Plain shell backgrounding with `&`, `nohup`, `disown`, or `setsid ... &` inside a normal `sh` call is not part of this managed persistence model. Default: false.",
		},
		"interactive": map[string]any{
			"type":        "boolean",
			"description": "If true, start the command under a PTY-backed interactive job and return immediately with an absolute filesystem-backed job directory path under `${QUINE_DATA_DIR}`. The job directory exposes `screen.txt`, `screen.png`, `screen.meta`, `in`, `winsize`, `events.log`, and `exit`. Use this for REPLs, shells, editors, or any program whose semantics depend on a terminal screen rather than append-only stdout/stderr. Mutually exclusive with `stdin` and `detach`. Default: false.",
		},
		"stdin": map[string]any{
			"type":        "string",
			"description": "Data to provide verbatim on the command's stdin with no shell escaping required. Quine materializes this as finite prewritten input for the command, so do not assume live stream or pipe-specific behavior. Ideal for writing text files with special characters, heredocs, or multi-line scripts. Example: sh(command=\"cat > file.py\", stdin=\"print('hello')\\n\"). Mutually exclusive with `interactive=true`.",
		},
	}
	required := []string{"command"}

	// Add goal/strategy only when escalation is available (for STALL detection)
	if cfg != nil && cfg.CanEscalate() {
		props["goal"] = map[string]any{
			"type":        "string",
			"description": "The OVERALL task objective in 2-5 abstract words (e.g. 'Decode hidden message', 'Fix auth bug', 'Build REST API'). Choose ONE goal for your entire session and reuse it VERBATIM on every sh call. Do NOT describe individual commands.",
		}
		props["strategy"] = map[string]any{
			"type":        "string",
			"description": "Your current approach to achieving the goal (e.g. 'parse with regex', 'geometric analysis'). Update when you pivot to a different approach.",
		}
		required = []string{"command", "goal", "strategy"}
	}

	return llm.ToolSchema{
		Name: "sh",
		Description: "Execute a POSIX shell command. Costs 1 execution.\n\n" +
			"Each `sh` call runs as a managed job with stdout/stderr captured to runtime-managed logs.\n" +
			"By default `sh` waits for completion, returns a bounded summary, and then removes the per-call job directory.\n" +
			"With `detach=true`, it returns immediately with an absolute job path under `${QUINE_DATA_DIR}` and keeps that job directory so you can inspect logs or wait via `cat <path>/exit`.\n" +
			"With `interactive=true`, it returns immediately with an absolute PTY-backed job path under `${QUINE_DATA_DIR}` whose control surface includes `screen.txt`, `screen.png`, `screen.meta`, `in`, `winsize`, `events.log`, and `exit`.\n" +
			"`${QUINE_DATA_DIR}` is Quine's runtime-state root for jobs, tapes, logs, and coordination files; it is separate from the task workspace.\n" +
			"Use `detach=true` for durable background work; do not assume plain shell `&` or `nohup` inside a normal `sh` call will persist across calls or after quine exits.\n" +
			"Use `interactive=true` when the process is operating a screen instead of emitting meaningful append-only stdout/stderr.\n" +
			"When Linux workspace physics are enabled, file changes persist across `sh` calls inside a provisional overlay-backed view and each completed result ends with `[FS MUTATIONS]` showing created, modified, and deleted paths. `detach=true` is unavailable in that mode.\n" +
			"Shell commands also get `$QUINE_JOB_SESSION_DIR` for concise path construction.",
		Parameters: map[string]any{
			"type":       "object",
			"properties": props,
			"required":   required,
		},
	}
}

// ForkToolSchema returns the JSON Schema for the fork tool.
func ForkToolSchema(cfg *config.Config) llm.ToolSchema {
	desc := "Spawn one or more child agents with cloned context (swarm fork). " +
		"Each child inherits your conversation history and starts with its own intent. "
	childrenDesc := "Array of child specs. Each child must provide `intent` and `workspace`."
	workspaceDesc := "Child workspace path. `.` means inherit your current workspace."
	if cfg != nil && !cfg.WorkspaceEnabled {
		desc += "Children share the filesystem. "
		workspaceDesc = "Child workspace placeholder. Use `.` when workspace physics are not enabled."
	} else {
		desc += "Under Linux workspace physics, each child should be pointed at an explicit workspace. "
		childrenDesc = "Array of child specs. Each child must provide `intent` and `workspace`. `workspace` may be `.` to inherit your current workspace or a narrower path within the workspace root."
		workspaceDesc = "Child workspace path. `.` means inherit your current workspace. Relative paths resolve from your current workspace and must remain within the workspace root. Children share the same mounted workspace surface."
	}
	return llm.ToolSchema{
		Name:        "fork",
		Description: desc + "Use for parallel exploration, delegation, or breaking down complex tasks.",
		Parameters: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"children": map[string]any{
					"type":        "array",
					"description": childrenDesc,
					"items": map[string]any{
						"type": "object",
						"properties": map[string]any{
							"intent": map[string]any{
								"type":        "string",
								"description": "Mission string for this child process.",
							},
							"workspace": map[string]any{
								"type":        "string",
								"description": workspaceDesc,
							},
						},
						"required": []string{"intent", "workspace"},
					},
					"minItems": 1,
				},
				"argv": map[string]any{
					"type":        "array",
					"description": "Legacy compatibility form. Prefer `children`. Each string becomes a child mission and implicitly uses workspace `.`.",
					"items":       map[string]any{"type": "string"},
					"minItems":    1,
				},
				"mode": map[string]any{
					"type":        "string",
					"enum":        []string{"wait", "race", "forget"},
					"description": "race (default): first child to exit 0 wins, rest are killed. wait: block until all children finish, return all results. forget: fire-and-forget, return PIDs immediately.",
				},
			},
			"required": []string{"children"},
		},
	}
}

// ExecToolSchema returns the JSON Schema for the exec tool.
func ExecToolSchema(cfg *config.Config) llm.ToolSchema {
	desc := "Replace yourself with a fresh instance while preserving the original mission. " +
		"Use this when your context is noisy but the task is not complete. " +
		"The new instance starts with: (1) empty conversation history, (2) same original intent from stdin, " +
		"(3) all wisdom preserved and merged with new wisdom you provide."

	props := map[string]any{
		"wisdom": map[string]any{
			"type":        "object",
			"description": "State to pass to your next incarnation. This is the ONLY information that survives exec — conversation history, tool outputs, and all other context are discarded. Values must be strings.",
			"additionalProperties": map[string]any{
				"type": "string",
			},
		},
		"persona": map[string]any{
			"type":        "string",
			"description": "Optional persona/system-prompt name to load (e.g. 'analyst', 'coder'). Looks for personas/{name}.md",
		},
	}
	return llm.ToolSchema{
		Name:        "exec",
		Description: desc,
		Parameters: map[string]any{
			"type": "object",
			"properties": props,
			"required": []string{},
		},
	}
}

func RestoreWorldToolSchema(cfg *config.Config) llm.ToolSchema {
	desc := "Restore the current provisional workspace view to a world revision handle. Does not consume a turn."
	if cfg == nil || !cfg.CanRestoreWorld() {
		desc += " Requires Linux workspace physics with workspace revision mode `restore`."
	} else {
		desc += " Use revision handles such as `wr0` or `wr3`. `wr0` is the baseline world; later revisions are historical checkpoints captured after completed shell turns."
	}
	return llm.ToolSchema{
		Name:        "restore_world",
		Description: desc,
		Parameters: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"revision": map[string]any{
					"type":        "string",
					"description": "Restore the workspace to a specific world revision handle. Example: revision=\"wr0\" restores the provisional world to the session baseline; revision=\"wr3\" restores the world captured after the third completed shell turn.",
				},
			},
			"required": []string{"revision"},
		},
	}
}

// ExitToolSchema returns the JSON Schema for the exit tool.
func ExitToolSchema() llm.ToolSchema {
	return llm.ToolSchema{
		Name: "exit",
		Description: "Finish your work and terminate. " +
			"Two modes: success (task complete), failure (task failed). " +
			"NOTE: This tool does NOT output to stdout. Use sh to write output to /dev/stdout.",
		Parameters: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"status": map[string]any{
					"type":        "string",
					"enum":        []string{"success", "failure"},
					"description": "Task outcome. \"success\" = complete. \"failure\" = failed.",
				},
				"stderr": map[string]any{
					"type":        "string",
					"description": "Why the task failed. Required on failure. Must NOT be set on success.",
				},
			},
			"required": []string{"status"},
		},
	}
}

// VisionToolSchema returns the JSON Schema for the vision tool.
func VisionToolSchema() llm.ToolSchema {
	return llm.ToolSchema{
		Name: "vision",
		Description: "Analyze an image file. Reads the file at the given path and delivers it " +
			"as a native image to the model's vision capabilities. " +
			"Supported formats: PNG, JPEG, GIF, WebP. Does NOT consume a turn. " +
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
					"type": "string",
					"description": "A text instruction describing what information to extract from the image. " +
						"Be specific about what to identify — vague instructions yield vague results.",
				},
			},
			"required": []string{"path"},
		},
	}
}

// MarkToolSchema returns the JSON Schema for the mark tool.
func MarkToolSchema() llm.ToolSchema {
	return llm.ToolSchema{
		Name: "mark",
		Description: "Compress current memory context into an immutable anchor. " +
			"Use fold=true to absorb existing frontier anchors into the new anchor.",
		Parameters: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"summary": map[string]any{
					"type":        "string",
					"description": "Human-readable summary for the anchor.",
				},
				"fold": map[string]any{
					"type":        "boolean",
					"description": "If true, absorb existing frontier anchors into the new anchor.",
				},
			},
			"required": []string{"summary"},
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

// EscalateToolSchema returns the JSON Schema for the escalate tool.
func EscalateToolSchema() llm.ToolSchema {
	return llm.ToolSchema{
		Name: "escalate",
		Description: "Request a more capable model. Use when you are stuck: repeated errors, " +
			"complex reasoning beyond your capacity, or cryptic failures you cannot resolve. " +
			"This is a one-way upgrade — use only when genuinely needed.",
		Parameters: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"reason": map[string]any{
					"type":        "string",
					"description": "Why you need a smarter model (what you tried, why it's not working).",
				},
			},
			"required": []string{"reason"},
		},
	}
}

// AllToolSchemas returns all tool schemas.
func AllToolSchemas(cfg *config.Config) []llm.ToolSchema {
	schemas := []llm.ToolSchema{
		ShToolSchema(cfg),
		ForkToolSchema(cfg),
		ExecToolSchema(cfg),
		ExitToolSchema(),
		VisionToolSchema(),
	}
	if cfg != nil && cfg.CanRestoreWorld() {
		schemas = append(schemas, RestoreWorldToolSchema(cfg))
	}
	if cfg != nil && cfg.AnchorMemoryEnabled {
		schemas = append(schemas, MarkToolSchema(), UnfoldToolSchema())
	}
	if cfg != nil && cfg.CanEscalate() {
		schemas = append(schemas, EscalateToolSchema())
	}
	return schemas
}
