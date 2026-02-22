package tools

import "github.com/kehao95/quine/internal/llm"

// ShToolSchema returns the JSON Schema for the sh tool.
func ShToolSchema() llm.ToolSchema {
	return llm.ToolSchema{
		Name: "sh",
		Description: "Execute a POSIX shell command. Costs 1 execution.\n\n" +
			"Returns output immediately on normal completion, or [PAUSED] when a budget is exhausted.\n" +
			"Use `job` to read, resume, or kill paused jobs.",
		Parameters: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"command": map[string]any{
					"type":        "string",
					"description": "The shell command to execute.",
				},
				"timeout": map[string]any{
					"type":        "integer",
					"description": "Maximum wall-clock seconds. When exceeded, the job is SIGSTOP'd and [PAUSED] is returned. Default: no timeout.",
				},
				"output_limit": map[string]any{
					"type":        "integer",
					"description": "Maximum combined stdout+stderr bytes. When exceeded, the job is SIGSTOP'd and [PAUSED] is returned. Default: no limit.",
				},
			},
			"required": []string{"command"},
		},
	}
}

// JobToolSchema returns the JSON Schema for the job tool.
func JobToolSchema() llm.ToolSchema {
	return llm.ToolSchema{
		Name: "job",
		Description: "Manage a paused or running job.\n\n" +
			"- `job(id=N)` — Read accumulated output without resuming.\n" +
			"- `job(id=N, signal=\"cont\")` — Resume a paused job (optionally with a new budget).\n" +
			"- `job(id=N, signal=\"kill\")` — Terminate the job.",
		Parameters: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"id": map[string]any{
					"type":        "integer",
					"description": "Job ID (pgid) returned in the [PAUSED] header.",
				},
				"signal": map[string]any{
					"type":        "string",
					"enum":        []string{"cont", "kill"},
					"description": "\"cont\" resumes a paused job; \"kill\" terminates it.",
				},
				"timeout": map[string]any{
					"type":        "integer",
					"description": "New timeout budget in seconds (only meaningful with signal=\"cont\").",
				},
				"output_limit": map[string]any{
					"type":        "integer",
					"description": "New output budget in bytes (only meaningful with signal=\"cont\").",
				},
			},
			"required": []string{"id"},
		},
	}
}

// ForkToolSchema returns the JSON Schema for the fork tool.
func ForkToolSchema() llm.ToolSchema {
	return llm.ToolSchema{
		Name: "fork",
		Description: "Spawn a child agent with cloned context (horizontal scaling). " +
			"The child inherits your conversation history and starts with the given intent. " +
			"Use for parallel exploration, delegation, or breaking down complex tasks.",
		Parameters: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"intent": map[string]any{
					"type":        "string",
					"description": "The task or instruction for the child agent. Be specific about what you want the child to accomplish.",
				},
				"wait": map[string]any{
					"type":        "boolean",
					"description": "If true, block until child completes and return its output. If false (default), spawn child and continue immediately.",
				},
			},
			"required": []string{"intent"},
		},
	}
}

// ExecToolSchema returns the JSON Schema for the exec tool.
func ExecToolSchema() llm.ToolSchema {
	return llm.ToolSchema{
		Name: "exec",
		Description: "Metamorphosis: Replace yourself with a fresh instance while preserving the original mission. " +
			"Use this when your context is polluted with noise but the task isn't complete. " +
			"The new instance starts with: (1) Empty conversation history, (2) Same original intent from stdin, " +
			"(3) All wisdom preserved and merged with new wisdom you provide. This is vertical scaling — same mission, fresh brain.",
		Parameters: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"wisdom": map[string]any{
					"type":        "object",
					"description": "State to pass to your next incarnation as key-value pairs. Example: {\"files_checked\": \"15\", \"next_target\": \"shelf_02\", \"strategy\": \"try edges\"}. All values must be strings.",
					"additionalProperties": map[string]any{
						"type": "string",
					},
				},
				"persona": map[string]any{
					"type":        "string",
					"description": "Optional persona/system-prompt name to load (e.g. 'analyst', 'coder'). Looks for personas/{name}.md",
				},
			},
			"required": []string{},
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

// AllToolSchemas returns all tool schemas.
func AllToolSchemas() []llm.ToolSchema {
	return []llm.ToolSchema{
		ShToolSchema(),
		JobToolSchema(),
		ForkToolSchema(),
		ExecToolSchema(),
		ExitToolSchema(),
	}
}
