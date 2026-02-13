package tools

import "github.com/kehao95/quine/internal/llm"

// ShToolSchema returns the JSON Schema for the sh tool.
func ShToolSchema() llm.ToolSchema {
	return llm.ToolSchema{
		Name:        "sh",
		Description: "Execute a POSIX shell command.",
		Parameters: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"command": map[string]any{
					"type":        "string",
					"description": "The shell command to execute",
				},
				"timeout": map[string]any{
					"type":        "integer",
					"description": "Optional timeout in seconds",
				},
				"stdout": map[string]any{
					"type":        "boolean",
					"description": "If true, pipe command stdout directly to the process stdout (passthrough). Use for delivering binary output or large artifacts. Command output will NOT appear in the tool result — it goes straight to the parent process.",
				},
			},
			"required": []string{"command"},
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
					"description": "Key-value pairs to pass to your next incarnation. Use this to transfer critical state like 'found_count', 'current_position', 'partial_result'. Values must be strings.",
					"additionalProperties": map[string]any{
						"type": "string",
					},
				},
				"reason": map[string]any{
					"type":        "string",
					"description": "Why you're exec'ing — logged for debugging (e.g. 'context too noisy after reading 50K tokens')",
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

// ReadToolSchema returns the JSON Schema for the read tool.
func ReadToolSchema() llm.ToolSchema {
	return llm.ToolSchema{
		Name: "read",
		Description: "Read lines from stdin (streaming input). Use this to process input incrementally. " +
			"Each call reads up to 500 lines (capped to prevent context overflow). " +
			"Call repeatedly until EOF. If you see [TRUNCATED], more data is available — " +
			"consider using exec with wisdom to reset context while preserving progress.",
		Parameters: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"lines": map[string]any{
					"type":        "integer",
					"description": "Number of lines to read. Default 1. Set to 0 to read a batch (up to 500 lines). Capped at 500 per call.",
				},
				"timeout": map[string]any{
					"type":        "integer",
					"description": "Timeout in seconds for blocking read. Default is 60 seconds.",
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
			"Three modes: success (task complete), failure (task failed), progress (partial work done). " +
			"NOTE: This tool does NOT output to stdout. All stdout must go through sh with stdout:true.",
		Parameters: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"status": map[string]any{
					"type":        "string",
					"enum":        []string{"success", "failure", "progress"},
					"description": "Task outcome. \"success\" = complete. \"failure\" = failed. \"progress\" = partial work done but couldn't finish.",
				},
				"stderr": map[string]any{
					"type":        "string",
					"description": "Why the task failed or is incomplete. Required on failure and progress. Must NOT be set on success.",
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
		ReadToolSchema(),
		ForkToolSchema(),
		ExecToolSchema(),
		ExitToolSchema(),
	}
}
