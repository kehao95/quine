package runtime

import (
	"fmt"
	"strings"

	"github.com/kehao95/quine/internal/config"
	"github.com/kehao95/quine/internal/tape"
)

type toolOutcome struct {
	terminate       bool
	exitCode        int
	budgetExhausted bool
}

type toolSpec struct {
	name        string
	enabled     func(*config.Config) bool
	disabledMsg string
	handle      func(*Runtime, tape.Message, tape.ToolCall) toolOutcome
}

func newToolRegistry() map[string]toolSpec {
	return map[string]toolSpec{
		"exit": {
			name: "exit",
			enabled: func(cfg *config.Config) bool {
				return cfg.ExitEnabled()
			},
			disabledMsg: "Rejected: exit is disabled in this runtime. The process terminates when execution budget is exhausted or by signal.",
			handle: func(r *Runtime, _ tape.Message, tc tape.ToolCall) toolOutcome {
				code, ok := r.handleExit(tc)
				if ok {
					return toolOutcome{terminate: true, exitCode: code}
				}
				return toolOutcome{}
			},
		},
		"idle": {
			name: "idle",
			enabled: func(cfg *config.Config) bool {
				return cfg.IdleToolEnabled()
			},
			disabledMsg: "Rejected: idle is disabled in this runtime.",
			handle: func(r *Runtime, _ tape.Message, tc tape.ToolCall) toolOutcome {
				r.handleIdleEnabled(tc)
				return toolOutcome{}
			},
		},
		"sh": {
			name: "sh",
			handle: func(r *Runtime, _ tape.Message, tc tape.ToolCall) toolOutcome {
				if r.handleSh(tc) {
					return toolOutcome{budgetExhausted: true}
				}
				return toolOutcome{}
			},
		},
		"fork": {
			name: "fork",
			enabled: func(cfg *config.Config) bool {
				return cfg.ForkEnabled()
			},
			disabledMsg: "Rejected: fork is disabled in this runtime.",
			handle: func(r *Runtime, _ tape.Message, tc tape.ToolCall) toolOutcome {
				r.handleForkEnabled(tc)
				return toolOutcome{}
			},
		},
		"spawn": {
			name: "spawn",
			enabled: func(cfg *config.Config) bool {
				return cfg.SpawnEnabled()
			},
			disabledMsg: "Rejected: spawn is disabled in this runtime.",
			handle: func(r *Runtime, _ tape.Message, tc tape.ToolCall) toolOutcome {
				r.handleSpawnEnabled(tc)
				return toolOutcome{}
			},
		},
		"exec": {
			name: "exec",
			enabled: func(cfg *config.Config) bool {
				return cfg.ExecEnabled
			},
			disabledMsg: "Rejected: exec is disabled in this runtime.",
			handle: func(r *Runtime, _ tape.Message, tc tape.ToolCall) toolOutcome {
				r.handleExec(tc)
				return toolOutcome{}
			},
		},
		"switch_world": {
			name: "switch_world",
			handle: func(r *Runtime, _ tape.Message, tc tape.ToolCall) toolOutcome {
				r.handleSwitchWorld(tc)
				return toolOutcome{}
			},
		},
		"vision": {
			name: "vision",
			enabled: func(cfg *config.Config) bool {
				return cfg.VisionEnabled
			},
			disabledMsg: "Rejected: vision is disabled in this runtime.",
			handle: func(r *Runtime, _ tape.Message, tc tape.ToolCall) toolOutcome {
				r.handleVision(tc)
				return toolOutcome{}
			},
		},
		"mark": {
			name: "mark",
			handle: func(r *Runtime, source tape.Message, tc tape.ToolCall) toolOutcome {
				r.handleMark(source, tc)
				return toolOutcome{}
			},
		},
		"unfold": {
			name: "unfold",
			handle: func(r *Runtime, _ tape.Message, tc tape.ToolCall) toolOutcome {
				r.handleUnfold(tc)
				return toolOutcome{}
			},
		},
		"escalate": {
			name: "escalate",
			handle: func(r *Runtime, _ tape.Message, tc tape.ToolCall) toolOutcome {
				r.handleEscalate(tc)
				return toolOutcome{}
			},
		},
	}
}

func (r *Runtime) rejectDisabledTool(tc tape.ToolCall, spec toolSpec) {
	rejectMsg := runtimeToolResultMessage(tc.ID, spec.name, "rejected", map[string]any{
		"error": spec.disabledMsg,
	})
	r.appendRuntimeToolMessage(rejectMsg, true)
	r.log("%s rejected: disabled by config", spec.name)
}

func (r *Runtime) toolSpec(name string) (toolSpec, bool) {
	if r.toolRegistry != nil {
		spec, ok := r.toolRegistry[name]
		return spec, ok
	}
	spec, ok := newToolRegistry()[name]
	return spec, ok
}

func (r *Runtime) precheckProcessCreation(toolID, kind string, requested int) *tape.Message {
	kindLabel := strings.ToUpper(kind)
	// Check depth limit before process creation (disabled when MaxDepth <= 0).
	if r.cfg.MaxDepth > 0 && r.cfg.Depth+1 >= r.cfg.MaxDepth {
		r.log("%s rejected - depth limit exceeded (%d/%d)", kind, r.cfg.Depth+1, r.cfg.MaxDepth)
		errMsg := runtimeToolResultMessage(toolID, kind, "error", map[string]any{
			"error": fmt.Sprintf("[%s ERROR] Max recursion depth exceeded (%d/%d). Cannot spawn child.", kindLabel, r.cfg.Depth+1, r.cfg.MaxDepth),
		})
		return &errMsg
	}

	// All-or-Nothing slot allocation: if the agent requests N children but
	// fewer than N slots are available, the entire process creation fails.
	// No partial swarms - the agent must make a conscious decision.
	if r.cfg.MaxAgents > 0 {
		current := r.agentRegistry.Count()
		available := r.cfg.MaxAgents - current
		if available < requested {
			r.log("%s rejected - insufficient slots (requested %d, available %d, current %d/%d)",
				kind, requested, available, current, r.cfg.MaxAgents)
			errMsg := runtimeToolResultMessage(toolID, kind, "error", map[string]any{
				"error": fmt.Sprintf("[%s ERROR] Insufficient slots. Requested %d, available %d. Release agents or request fewer.", kindLabel, requested, available),
			})
			return &errMsg
		}
	}
	return nil
}
