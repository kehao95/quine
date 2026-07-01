package tools

import (
	"strings"
	"testing"

	"github.com/kehao95/quine/internal/config"
)

func TestAllToolSchemas_Count(t *testing.T) {
	// Without config (nil), workspace revision tools should not be included.
	schemas := AllToolSchemas(nil)
	if len(schemas) != 5 {
		t.Fatalf("AllToolSchemas(nil) returned %d schemas, want 5", len(schemas))
	}
}

func TestMarkToolSchema_FoldGating(t *testing.T) {
	with := MarkToolSchema(true)
	props, _ := with.Parameters["properties"].(map[string]any)
	if _, ok := props["fold"]; !ok {
		t.Fatal("MarkToolSchema(true) should expose the fold parameter")
	}
	if !strings.Contains(with.Description, "fold=true") {
		t.Fatal("MarkToolSchema(true) description should mention fold=true")
	}

	without := MarkToolSchema(false)
	propsOff, _ := without.Parameters["properties"].(map[string]any)
	if _, ok := propsOff["fold"]; ok {
		t.Fatal("MarkToolSchema(false) should omit the fold parameter")
	}
	if strings.Contains(without.Description, "fold") {
		t.Fatalf("MarkToolSchema(false) description should not mention fold, got %q", without.Description)
	}
	if _, ok := propsOff["resolution"]; !ok {
		t.Fatal("MarkToolSchema(false) should still expose resolution")
	}
}

func TestAllToolSchemas_WithoutExec(t *testing.T) {
	cfg := &config.Config{ToolGates: config.ToolGates{ExecEnabled: false, VisionEnabled: true}}
	schemas := AllToolSchemas(cfg)
	for _, s := range schemas {
		if s.Name == "exec" {
			t.Fatal("AllToolSchemas() should omit exec when exec is disabled")
		}
	}
}

func TestAllToolSchemas_WithIdle(t *testing.T) {
	cfg := &config.Config{ToolGates: config.ToolGates{IdleEnabled: true, ExecEnabled: true, VisionEnabled: true}}
	schemas := AllToolSchemas(cfg)
	if len(schemas) != 6 {
		t.Fatalf("AllToolSchemas() with idle enabled should return 6 schemas, got %d", len(schemas))
	}

	found := false
	for _, s := range schemas {
		if s.Name != "idle" {
			continue
		}
		found = true
		if !strings.Contains(s.Description, "Suspend explicitly") {
			t.Fatalf("idle schema should describe explicit suspension, got %q", s.Description)
		}
		if !strings.Contains(s.Description, "`idle` resumes on `poke`, `inject`, or `interrupt`") {
			t.Fatalf("idle schema should describe resume classes, got %q", s.Description)
		}
		if !strings.Contains(s.Description, "runtime prompt owns the full peer-control surface map") {
			t.Fatalf("idle schema should defer full control-surface map to prompt, got %q", s.Description)
		}
		props := s.Parameters["properties"].(map[string]any)
		if len(props) != 0 {
			t.Fatalf("idle schema should not expose parameters, got %#v", props)
		}
		if additional, ok := s.Parameters["additionalProperties"].(bool); !ok || additional {
			t.Fatalf("idle schema should reject extra parameters, got %#v", s.Parameters["additionalProperties"])
		}
	}
	if !found {
		t.Fatal("AllToolSchemas() should include idle when QUINE_IDLE_ENABLED=1")
	}
}

func TestAllToolSchemas_WithoutFork(t *testing.T) {
	cfg := &config.Config{ToolGates: config.ToolGates{ExecEnabled: true, VisionEnabled: true}}
	cfg.ToolGates = config.DefaultToolGates()
	cfg.ToolGates.ForkEnabled = false
	schemas := AllToolSchemas(cfg)
	if len(schemas) != 4 {
		t.Fatalf("AllToolSchemas() with fork disabled should return 4 schemas, got %d", len(schemas))
	}
	for _, s := range schemas {
		if s.Name == "fork" {
			t.Fatal("AllToolSchemas() should omit fork when QUINE_FORK_ENABLED=0")
		}
	}
}

func TestAllToolSchemas_WithSpawn(t *testing.T) {
	cfg := &config.Config{ToolGates: config.ToolGates{ExecEnabled: true, VisionEnabled: true, SpawnEnabledFlag: true, ForkWorldEnabled: true}, WorkspaceConfig: config.WorkspaceConfig{WorkspaceEnabled: true}}
	schemas := AllToolSchemas(cfg)
	if len(schemas) != 6 {
		t.Fatalf("AllToolSchemas() with spawn enabled should return 6 schemas, got %d", len(schemas))
	}
	var found bool
	for _, s := range schemas {
		if s.Name != "spawn" {
			continue
		}
		found = true
		if !strings.Contains(s.Description, "does not import the parent's active context") {
			t.Fatalf("spawn schema should describe fresh context semantics, got %q", s.Description)
		}
		props := s.Parameters["properties"].(map[string]any)
		children := props["children"].(map[string]any)
		items := children["items"].(map[string]any)
		childProps := items["properties"].(map[string]any)
		if _, ok := childProps["mission"]; !ok {
			t.Fatalf("spawn child schema should expose mission, got %#v", childProps)
		}
		for _, key := range []string{"world", "protection", "scope"} {
			if _, ok := childProps[key]; !ok {
				t.Fatalf("spawn child schema should expose %s, got %#v", key, childProps)
			}
		}
		if !strings.Contains(s.Description, "Workspace, world, scope, protection") {
			t.Fatalf("spawn schema should describe fork-aligned workspace semantics, got %q", s.Description)
		}
	}
	if !found {
		t.Fatal("AllToolSchemas() should include spawn when QUINE_SPAWN_ENABLED=1")
	}
}

func TestAllToolSchemas_WithEscalation(t *testing.T) {
	// SmartModelID enables escalation.
	cfg := &config.Config{
		ToolGates: config.ToolGates{ExecEnabled: true, VisionEnabled: true},
		Escalation: config.Escalation{
			SmartModelID: "claude-opus",
			Escalated:    false,
		},
	}
	schemas := AllToolSchemas(cfg)
	if len(schemas) != 6 {
		t.Fatalf("AllToolSchemas() with escalation should return 6 schemas, got %d", len(schemas))
	}

	// Verify escalate is in the list
	found := false
	for _, s := range schemas {
		if s.Name == "escalate" {
			found = true
			break
		}
	}
	if !found {
		t.Error("AllToolSchemas() should include escalate tool when escalation is configured")
	}
}

func TestAllToolSchemas_PostEscalation(t *testing.T) {
	cfg := &config.Config{
		ToolGates: config.ToolGates{ExecEnabled: true, VisionEnabled: true},
		Escalation: config.Escalation{
			SmartModelID: "claude-opus",
			Escalated:    true,
		},
	}
	schemas := AllToolSchemas(cfg)
	if len(schemas) != 5 {
		t.Fatalf("AllToolSchemas() post-escalation should return 5 schemas, got %d", len(schemas))
	}

	// Verify escalate is NOT in the list
	for _, s := range schemas {
		if s.Name == "escalate" {
			t.Error("AllToolSchemas() should NOT include escalate tool after escalation")
		}
	}
}

func TestAllToolSchemas_WithAnchorMemory(t *testing.T) {
	cfg := &config.Config{ToolGates: config.ToolGates{ExecEnabled: true, AnchorMemoryEnabled: true, VisionEnabled: true}}
	schemas := AllToolSchemas(cfg)
	if len(schemas) != 7 {
		t.Fatalf("AllToolSchemas() with anchor memory should return 7 schemas, got %d", len(schemas))
	}
	var foundMark, foundUnfold bool
	for _, s := range schemas {
		if s.Name == "mark" {
			foundMark = true
			if !strings.Contains(s.Description, "Crystallize") {
				t.Fatalf("mark schema should describe crystallization, got %q", s.Description)
			}
			if !strings.Contains(s.Description, "working-memory checkpoints") {
				t.Fatalf("mark schema should describe working-memory checkpoints, got %q", s.Description)
			}
			if !strings.Contains(s.Description, "Does not consume execution budget") {
				t.Fatalf("mark schema should describe low cost, got %q", s.Description)
			}
			if !strings.Contains(s.Description, "Memory telemetry can point toward `mark`") {
				t.Fatalf("mark schema should describe telemetry pressure, got %q", s.Description)
			}
			props := s.Parameters["properties"].(map[string]any)
			resolution := props["resolution"].(map[string]any)
			fold := props["fold"].(map[string]any)
			if !strings.Contains(resolution["description"].(string), "current working set") {
				t.Fatalf("mark.resolution should describe working-set facts, got %q", resolution["description"])
			}
			if !strings.Contains(resolution["description"].(string), "stably true") {
				t.Fatalf("mark.resolution should describe stable facts, got %q", resolution["description"])
			}
			if !strings.Contains(resolution["description"].(string), "closed a subproblem") {
				t.Fatalf("mark.resolution should describe closure boundaries, got %q", resolution["description"])
			}
			if !strings.Contains(fold["description"].(string), "remembered background") {
				t.Fatalf("mark.fold should describe remembered-background folding, got %q", fold["description"])
			}
			if !strings.Contains(fold["description"].(string), "higher-order resolution") {
				t.Fatalf("mark.fold should describe higher-order resolution gating, got %q", fold["description"])
			}
			if !strings.Contains(fold["description"].(string), "Actual consolidation requires existing prior anchors") {
				t.Fatalf("mark.fold should require prior anchors for consolidation, got %q", fold["description"])
			}
			if !strings.Contains(fold["description"].(string), "memory telemetry points toward `fold`") {
				t.Fatalf("mark.fold should mention telemetry gating, got %q", fold["description"])
			}
			if !strings.Contains(fold["description"].(string), "stable child findings") {
				t.Fatalf("mark.fold should mention returned child findings as a valid fold case, got %q", fold["description"])
			}
			if !strings.Contains(s.Description, "higher-order frontier reconfiguration move") {
				t.Fatalf("mark schema should describe fold as frontier reconfiguration, got %q", s.Description)
			}
			if !strings.Contains(s.Description, "Memory telemetry can point toward `fold`") {
				t.Fatalf("mark schema should mention fold telemetry, got %q", s.Description)
			}
			if !strings.Contains(s.Description, "parent-level conclusion") {
				t.Fatalf("mark schema should mention parent-level synthesis after child findings, got %q", s.Description)
			}
			if !strings.Contains(s.Description, "governing organizing point") {
				t.Fatalf("mark schema should mention the new anchor as the governing organizing point, got %q", s.Description)
			}
		}
		if s.Name == "unfold" {
			foundUnfold = true
		}
	}
	if !foundMark || !foundUnfold {
		t.Fatalf("expected mark/unfold schemas, got mark=%v unfold=%v", foundMark, foundUnfold)
	}
}

func TestAllToolSchemas_WithAnchorMemoryAndEscalation(t *testing.T) {
	cfg := &config.Config{ToolGates: config.ToolGates{ExecEnabled: true, AnchorMemoryEnabled: true, VisionEnabled: true}, Escalation: config.Escalation{SmartModelID: "claude-opus"}}
	schemas := AllToolSchemas(cfg)
	if len(schemas) != 8 {
		t.Fatalf("AllToolSchemas() with anchor memory + escalation should return 8 schemas, got %d", len(schemas))
	}
}

func TestAllToolSchemas_WithIdleAndAnchorMemoryAndEscalation(t *testing.T) {
	cfg := &config.Config{ToolGates: config.ToolGates{IdleEnabled: true, ExecEnabled: true, AnchorMemoryEnabled: true, VisionEnabled: true}, Escalation: config.Escalation{SmartModelID: "claude-opus"}}
	schemas := AllToolSchemas(cfg)
	if len(schemas) != 9 {
		t.Fatalf("AllToolSchemas() with idle + anchor memory + escalation should return 9 schemas, got %d", len(schemas))
	}
}

func TestSwitchWorldToolSchema(t *testing.T) {
	schema := SwitchWorldToolSchema(&config.Config{WorkspaceConfig: config.WorkspaceConfig{WorkspaceEnabled: true, WorkspaceRevisionMode: config.WorkspaceRevisionRestore}})
	if schema.Name != "switch_world" {
		t.Fatalf("schema name = %q, want switch_world", schema.Name)
	}
	props := schema.Parameters["properties"].(map[string]any)
	if _, ok := props["target"]; !ok {
		t.Fatal("switch_world schema should expose target")
	}
	if !strings.Contains(schema.Description, "provisional workspace") {
		t.Fatalf("unexpected switch_world description: %q", schema.Description)
	}
	if !strings.Contains(schema.Description, "`fs_mutations` and `world_revision`") {
		t.Fatalf("switch_world schema should mention fs mutation telemetry by default, got %q", schema.Description)
	}
}

func TestForkToolSchema_ClarifiesParentSideWorldState(t *testing.T) {
	schema := ForkToolSchema(&config.Config{WorkspaceConfig: config.WorkspaceConfig{WorkspaceEnabled: true, WorkspaceRevisionMode: config.WorkspaceRevisionRestore}})
	if !strings.Contains(schema.Description, "stdout/stderr are only captured process output") {
		t.Fatalf("fork schema should not present stdout/stderr as the only work product, got %q", schema.Description)
	}
	if !strings.Contains(schema.Description, "parent-side `fs_mutations` and `world_revision`") {
		t.Fatalf("fork schema should describe parent-side world reporting, got %q", schema.Description)
	}
	if !strings.Contains(schema.Description, "child writes stay in child lineage") {
		t.Fatalf("fork schema should describe child lineage isolation, got %q", schema.Description)
	}
}

func TestForkToolSchema_HasNoArgvCompatibilityProperty(t *testing.T) {
	props := ForkToolSchema(&config.Config{}).Parameters["properties"].(map[string]any)
	if _, ok := props["argv"]; ok {
		t.Fatal("fork schema should not expose legacy argv compatibility property")
	}
}

func TestForkToolSchema_HostFilesystemWordingIsNeutral(t *testing.T) {
	schema := ForkToolSchema(&config.Config{})

	if !strings.Contains(schema.Description, "Children share the same filesystem surface as the parent and siblings.") {
		t.Fatalf("fork schema should describe shared host filesystem neutrally, got %q", schema.Description)
	}
	for _, forbidden := range []string{
		"conflict",
		"does not provide per-child filesystem isolation",
	} {
		if strings.Contains(schema.Description, forbidden) {
			t.Fatalf("fork schema should avoid negative shared-filesystem wording %q, got %q", forbidden, schema.Description)
		}
	}
}

func TestSchemasHideFSMutationTelemetryWhenDisabled(t *testing.T) {
	cfg := &config.Config{WorkspaceConfig: config.WorkspaceConfig{WorkspaceEnabled: true, WorkspaceRevisionMode: config.WorkspaceRevisionRestore}}
	cfg.ToolGates = config.DefaultToolGates()
	cfg.ToolGates.FSMutationTelemetry = false

	sh := ShToolSchema(cfg)
	if strings.Contains(sh.Description, "fs_mutations") {
		t.Fatalf("sh schema should hide fs mutation telemetry when disabled, got %q", sh.Description)
	}

	fork := ForkToolSchema(cfg)
	if strings.Contains(fork.Description, "parent-side `fs_mutations`") {
		t.Fatalf("fork schema should hide parent fs mutation telemetry when disabled, got %q", fork.Description)
	}

	switchWorld := SwitchWorldToolSchema(cfg)
	if strings.Contains(switchWorld.Description, "`fs_mutations` and `world_revision`") {
		t.Fatalf("switch_world schema should hide fs mutation telemetry when disabled, got %q", switchWorld.Description)
	}
	if !strings.Contains(switchWorld.Description, "Results include `world_revision` for the switch turn.") {
		t.Fatalf("switch_world schema should keep world revision reporting when disabled, got %q", switchWorld.Description)
	}
}

func TestExecToolSchema_DoesNotExposeResetWorld(t *testing.T) {
	schema := ExecToolSchema(&config.Config{WorkspaceConfig: config.WorkspaceConfig{WorkspaceEnabled: true, WorkspaceRevisionMode: config.WorkspaceRevisionRestore}})
	props := schema.Parameters["properties"].(map[string]any)
	if _, ok := props["reset_world"]; ok {
		t.Fatal("exec schema should not expose reset_world")
	}
	if _, ok := props["persona"]; ok {
		t.Fatal("exec schema should not expose persona")
	}
	if _, ok := props["target"]; !ok {
		t.Fatal("exec schema should expose target")
	}
	if _, ok := props["argv"]; !ok {
		t.Fatal("exec schema should expose argv")
	}
	if strings.Contains(schema.Description, "reset_world") {
		t.Fatalf("exec schema should not mention reset_world, got: %q", schema.Description)
	}
	if !strings.Contains(schema.Description, "sugar for re-execing quine through its configured self-reentry target; it preserves the current mission when one exists and uses no mission argv when none exists") {
		t.Fatalf("exec schema should explain the default self-exec sugar, got: %q", schema.Description)
	}
}

func TestExitToolSchema_DefaultAllowsFailure(t *testing.T) {
	schema := ExitToolSchema(&config.Config{PromptConfig: config.PromptConfig{FailOnImpossible: true}})
	props := schema.Parameters["properties"].(map[string]any)
	status := props["status"].(map[string]any)
	enum := status["enum"].([]string)

	if len(enum) != 2 || enum[0] != "success" || enum[1] != "failure" {
		t.Fatalf("exit status enum = %v, want [success failure]", enum)
	}
	if _, ok := props["stderr"]; !ok {
		t.Fatal("default exit schema should expose stderr")
	}
	if !strings.Contains(schema.Description, "Two modes: success") {
		t.Fatalf("default exit schema should describe success/failure modes, got %q", schema.Description)
	}
	if !strings.Contains(schema.Description, "This tool emits no stdout bytes") {
		t.Fatalf("default exit schema should narrow stdout behavior, got %q", schema.Description)
	}
	if strings.Contains(schema.Description, "stdout output is produced") {
		t.Fatalf("default exit schema should not frame stdout as the delivery surface, got %q", schema.Description)
	}
}

func TestExitToolSchema_SuccessOnlyWhenFailOnImpossibleDisabled(t *testing.T) {
	schema := ExitToolSchema(&config.Config{PromptConfig: config.PromptConfig{FailOnImpossible: false}})
	props := schema.Parameters["properties"].(map[string]any)
	status := props["status"].(map[string]any)
	enum := status["enum"].([]string)

	if len(enum) != 1 || enum[0] != "success" {
		t.Fatalf("exit status enum = %v, want [success]", enum)
	}
	if _, ok := props["stderr"]; ok {
		t.Fatal("success-only exit schema should not expose stderr")
	}
	if !strings.Contains(schema.Description, "only `status=\"success\"` is available") {
		t.Fatalf("success-only exit schema should describe restricted mode, got %q", schema.Description)
	}
	if !strings.Contains(schema.Description, "This tool emits no stdout bytes") {
		t.Fatalf("success-only exit schema should narrow stdout behavior, got %q", schema.Description)
	}
	if strings.Contains(schema.Description, "stdout output is produced") {
		t.Fatalf("success-only exit schema should not frame stdout as the delivery surface, got %q", schema.Description)
	}
}

func TestAllToolSchemas_WithWorkspaceRestoreMode(t *testing.T) {
	cfg := &config.Config{ToolGates: config.ToolGates{ExecEnabled: true, VisionEnabled: true}, WorkspaceConfig: config.WorkspaceConfig{WorkspaceEnabled: true, WorkspaceRevisionMode: config.WorkspaceRevisionRestore}}
	schemas := AllToolSchemas(cfg)
	var foundSwitch bool
	for _, s := range schemas {
		if s.Name == "switch_world" {
			foundSwitch = true
		}
	}
	if !foundSwitch {
		t.Fatal("restore mode should expose switch_world")
	}
}

func TestAllToolSchemas_WithWorkspaceNoneModeHidesRestore(t *testing.T) {
	cfg := &config.Config{ToolGates: config.ToolGates{ExecEnabled: true, VisionEnabled: true}, WorkspaceConfig: config.WorkspaceConfig{WorkspaceEnabled: true, WorkspaceRevisionMode: config.WorkspaceRevisionNone}}
	schemas := AllToolSchemas(cfg)
	for _, s := range schemas {
		if s.Name == "switch_world" {
			t.Fatal("none mode should not expose switch_world")
		}
	}
}

func TestShToolSchema_SingleModelMode(t *testing.T) {
	// Single-model mode: no SmartModelID configured
	schema := ShToolSchema(nil)

	props := schema.Parameters["properties"].(map[string]any)
	required := schema.Parameters["required"].([]string)

	// Should have command but NOT goal/strategy
	if _, ok := props["command"]; !ok {
		t.Error("sh schema should have 'command' property")
	}
	if _, ok := props["interactive"]; !ok {
		t.Error("sh schema should have 'interactive' property")
	}
	if _, ok := props["timeout"]; !ok {
		t.Error("sh schema should have 'timeout' property when timeout override is enabled")
	}
	if _, ok := props["goal"]; ok {
		t.Error("sh schema should NOT have 'goal' in single-model mode")
	}
	if _, ok := props["strategy"]; ok {
		t.Error("sh schema should NOT have 'strategy' in single-model mode")
	}

	// Only command should be required
	if len(required) != 1 || required[0] != "command" {
		t.Errorf("sh schema required should be ['command'] in single-model mode, got %v", required)
	}
}

func TestShToolSchema_HidesTimeoutWhenOverrideDisabled(t *testing.T) {
	cfg := &config.Config{}
	cfg.ToolGates = config.DefaultToolGates()
	cfg.ToolGates.ShTimeoutOverrideEnabled = false
	schema := ShToolSchema(cfg)

	props := schema.Parameters["properties"].(map[string]any)
	if _, ok := props["timeout"]; ok {
		t.Fatal("sh schema should not expose 'timeout' when timeout override is disabled")
	}
	for _, text := range []string{
		schema.Description,
		props["stdin"].(map[string]any)["description"].(string),
	} {
		forbidden := []string{
			"timeout parameter",
			"Optional sync-shell protection timeout",
			"sh(timeout",
			"runtime default shell timeout",
			"status=\"interrupted\"",
			"job.pid",
			"job.path",
			"stdout_so_far",
			"stderr_so_far",
			"*_so_far",
		}
		for _, term := range forbidden {
			if strings.Contains(text, term) {
				t.Fatalf("disabled sh timeout schema text should not contain %q: %s", term, text)
			}
		}
	}
}

func TestShToolSchema_EscalationMode(t *testing.T) {
	// Escalation mode: SmartModelID configured and not yet escalated
	cfg := &config.Config{ToolGates: config.ToolGates{ExecEnabled: true, VisionEnabled: true}, Escalation: config.Escalation{SmartModelID: "claude-opus", Escalated: false}}
	schema := ShToolSchema(cfg)

	props := schema.Parameters["properties"].(map[string]any)
	required := schema.Parameters["required"].([]string)

	// Should have command, goal, and strategy
	if _, ok := props["command"]; !ok {
		t.Error("sh schema should have 'command' property")
	}
	if _, ok := props["interactive"]; !ok {
		t.Error("sh schema should have 'interactive' in escalation mode")
	}
	if _, ok := props["goal"]; !ok {
		t.Error("sh schema should have 'goal' in escalation mode")
	}
	if _, ok := props["strategy"]; !ok {
		t.Error("sh schema should have 'strategy' in escalation mode")
	}

	// All three should be required
	if len(required) != 3 {
		t.Errorf("sh schema required should have 3 items in escalation mode, got %v", required)
	}
}

func TestShToolSchema_HidesInteractiveWhenDisabled(t *testing.T) {
	cfg := &config.Config{}
	cfg.ToolGates = config.DefaultToolGates()
	cfg.ToolGates.ShInteractiveEnabled = false
	schema := ShToolSchema(cfg)

	props := schema.Parameters["properties"].(map[string]any)
	if _, ok := props["interactive"]; ok {
		t.Fatal("sh schema should not expose 'interactive' when disabled")
	}
	for _, text := range []string{
		schema.Description,
		props["stdin"].(map[string]any)["description"].(string),
		props["timeout"].(map[string]any)["description"].(string),
	} {
		forbidden := []string{"interactive=true", "PTY", "screen snapshots"}
		for _, term := range forbidden {
			if strings.Contains(text, term) {
				t.Fatalf("disabled sh schema text should not contain %q: %s", term, text)
			}
		}
	}
}

func TestShToolSchema_HidesStdinWhenDisabled(t *testing.T) {
	cfg := &config.Config{}
	cfg.ToolGates = config.DefaultToolGates()
	cfg.ToolGates.ShStdinEnabled = false
	schema := ShToolSchema(cfg)

	props := schema.Parameters["properties"].(map[string]any)
	if _, ok := props["stdin"]; ok {
		t.Fatal("sh schema should not expose 'stdin' when disabled")
	}
	for _, text := range []string{
		schema.Description,
		props["interactive"].(map[string]any)["description"].(string),
	} {
		forbidden := []string{"stdin", "prewritten input", "heredoc"}
		for _, term := range forbidden {
			if strings.Contains(text, term) {
				t.Fatalf("disabled sh stdin schema text should not contain %q: %s", term, text)
			}
		}
	}
}

func TestShToolSchema_HidesDetachWhenDisabled(t *testing.T) {
	cfg := &config.Config{}
	cfg.ToolGates = config.DefaultToolGates()
	cfg.ToolGates.ShDetachEnabled = false
	schema := ShToolSchema(cfg)

	props := schema.Parameters["properties"].(map[string]any)
	if _, ok := props["detach"]; ok {
		t.Fatal("sh schema should not expose 'detach' when disabled")
	}
	for _, text := range []string{
		schema.Description,
		props["interactive"].(map[string]any)["description"].(string),
		props["timeout"].(map[string]any)["description"].(string),
	} {
		forbidden := []string{"detach=true", "Detached", "detached", "background job", "`detach`"}
		for _, term := range forbidden {
			if strings.Contains(text, term) {
				t.Fatalf("disabled sh detach schema text should not contain %q: %s", term, text)
			}
		}
	}
}

func TestShToolSchema_MinimalSurfaceOnlyExposesCommand(t *testing.T) {
	cfg := &config.Config{}
	cfg.ToolGates = config.DefaultToolGates()
	cfg.ToolGates.ShTimeoutOverrideEnabled = false
	cfg.ToolGates.ShStdinEnabled = false
	cfg.ToolGates.ShDetachEnabled = false
	cfg.ToolGates.ShInteractiveEnabled = false
	cfg.ToolGates.FSMutationTelemetry = false
	schema := ShToolSchema(cfg)

	props := schema.Parameters["properties"].(map[string]any)
	if len(props) != 1 {
		t.Fatalf("minimal sh schema should expose only command, got %v", props)
	}
	if _, ok := props["command"]; !ok {
		t.Fatal("minimal sh schema should expose command")
	}
	required := schema.Parameters["required"].([]string)
	if len(required) != 1 || required[0] != "command" {
		t.Fatalf("minimal sh required should be [command], got %v", required)
	}
	for _, forbidden := range []string{
		"stdin",
		"detach",
		"interactive",
		"timeout",
		"fs_mutations",
		"${QUINE_DATA_DIR}",
		"job.pid",
		"job.path",
		"*_so_far",
	} {
		if strings.Contains(schema.Description, forbidden) {
			t.Fatalf("minimal sh description should not contain %q: %s", forbidden, schema.Description)
		}
	}
}

func TestShToolSchema_DefaultHidesJobImplementationDetails(t *testing.T) {
	cfg := &config.Config{WorkspaceConfig: config.WorkspaceConfig{WorkspaceEnabled: true, WorkspaceBackend: "overlay", WorkspaceRevisionMode: config.WorkspaceRevisionRestore}}
	schema := ShToolSchema(cfg)
	props := schema.Parameters["properties"].(map[string]any)
	detach := props["detach"].(map[string]any)["description"].(string)
	interactive := props["interactive"].(map[string]any)["description"].(string)
	if !strings.Contains(interactive, "PTY-backed interactive POSIX job") {
		t.Fatalf("interactive schema should still expose the affordance, got %q", interactive)
	}
	for _, text := range []string{schema.Description, detach, interactive} {
		for _, forbidden := range []string{
			"Detached job directories include",
			"Interactive job directories include",
			"events.hex",
			"input.log",
			"job-local workspace lineage",
			"world_handle",
			"switch_world",
			"Send POSIX signals",
			"pid/process group",
		} {
			if strings.Contains(text, forbidden) {
				t.Fatalf("default sh schema should hide implementation detail %q: %s", forbidden, text)
			}
		}
	}
}

func TestShToolSchema_DetailedImplDescribesJobImplementationDetails(t *testing.T) {
	cfg := &config.Config{PromptConfig: config.PromptConfig{PromptImplDetails: true}, WorkspaceConfig: config.WorkspaceConfig{WorkspaceEnabled: true, WorkspaceBackend: "overlay", WorkspaceRevisionMode: config.WorkspaceRevisionRestore}}
	schema := ShToolSchema(cfg)
	props := schema.Parameters["properties"].(map[string]any)
	detach := props["detach"].(map[string]any)["description"].(string)
	interactive := props["interactive"].(map[string]any)["description"].(string)
	for _, want := range []string{
		"Interactive job directories include",
		"events.hex",
		"input.log",
		"job-local workspace lineage",
		"world_handle",
		"switch_world",
		"kill",
		"pid/process group",
	} {
		if !strings.Contains(interactive, want) {
			t.Fatalf("interactive schema description missing %q: %s", want, interactive)
		}
	}
	if !strings.Contains(detach, "Detached job directories include") {
		t.Fatalf("detach schema description missing detailed directory layout: %s", detach)
	}
	if !strings.Contains(schema.Description, "interactive jobs run in a job-local workspace lineage") {
		t.Fatalf("sh schema should describe overlay interactive lineage, got %q", schema.Description)
	}
}

func TestShToolSchema_PostEscalation(t *testing.T) {
	// Post-escalation: SmartModelID configured but already escalated
	cfg := &config.Config{ToolGates: config.ToolGates{ExecEnabled: true, VisionEnabled: true}, Escalation: config.Escalation{SmartModelID: "claude-opus", Escalated: true}}
	schema := ShToolSchema(cfg)

	props := schema.Parameters["properties"].(map[string]any)
	required := schema.Parameters["required"].([]string)

	// Post-escalation behaves like single-model mode (no STALL detection needed)
	if _, ok := props["goal"]; ok {
		t.Error("sh schema should NOT have 'goal' post-escalation")
	}
	if _, ok := props["strategy"]; ok {
		t.Error("sh schema should NOT have 'strategy' post-escalation")
	}

	// Only command should be required
	if len(required) != 1 || required[0] != "command" {
		t.Errorf("sh schema required should be ['command'] post-escalation, got %v", required)
	}
}

func TestForkToolSchema_WorkspaceDescriptionClarifiesChildCWD(t *testing.T) {
	cfg := &config.Config{PromptConfig: config.PromptConfig{PromptImplDetails: true}, WorkspaceConfig: config.WorkspaceConfig{WorkspaceEnabled: true}}
	schema := ForkToolSchema(cfg)

	if !strings.Contains(schema.Description, "Does not consume execution budget") {
		t.Fatalf("fork schema should describe low cost, got %q", schema.Description)
	}
	if !strings.Contains(schema.Description, "independent hypotheses, decoders, implementations, extractors, or verification strategies") {
		t.Fatalf("fork schema should describe parallel strategy use, got %q", schema.Description)
	}
	if !strings.Contains(schema.Description, "prefer 2-3 labeled children over another long parent-only inspection") {
		t.Fatalf("fork schema should describe when to fork instead of serial inspection, got %q", schema.Description)
	}
	if !strings.Contains(schema.Description, "Fork children preserve the parent mission as the active task contract") {
		t.Fatalf("fork schema should describe parent mission preservation, got %q", schema.Description)
	}
	if !strings.Contains(schema.Description, "each child intent is a lane assignment") {
		t.Fatalf("fork schema should describe child intents as lane assignments, got %q", schema.Description)
	}
	if !strings.Contains(schema.Description, "Do one cheap shared setup/probe if all lanes need it") {
		t.Fatalf("fork schema should describe shared setup before specialized forking, got %q", schema.Description)
	}
	if !strings.Contains(schema.Description, "`mode=\"wait\"` blocks until all children finish") {
		t.Fatalf("fork schema should describe wait mode mechanics, got %q", schema.Description)
	}
	if !strings.Contains(schema.Description, "`race` returns the first successful child") {
		t.Fatalf("fork schema should describe race mode mechanics, got %q", schema.Description)
	}

	props := schema.Parameters["properties"].(map[string]any)
	children := props["children"].(map[string]any)
	items := children["items"].(map[string]any)
	childProps := items["properties"].(map[string]any)
	intent := childProps["intent"].(map[string]any)
	scope := childProps["scope"].(map[string]any)
	desc := scope["description"].(string)
	mode := props["mode"].(map[string]any)
	modeDesc := mode["description"].(string)
	adopt := props["adopt_winner"].(map[string]any)
	adoptDesc := adopt["description"].(string)

	if !strings.Contains(intent["description"].(string), "distinct strategy") {
		t.Fatalf("fork intent description should request distinct child strategy, got %q", intent["description"])
	}
	if !strings.Contains(intent["description"].(string), "parent mission remains active") {
		t.Fatalf("fork intent description should describe parent mission preservation, got %q", intent["description"])
	}
	if !strings.Contains(intent["description"].(string), "lane-specific inputs") {
		t.Fatalf("fork intent description should request only lane-specific inputs, got %q", intent["description"])
	}
	if !strings.Contains(intent["description"].(string), "closest success check") {
		t.Fatalf("fork intent description should request success check, got %q", intent["description"])
	}
	if !strings.Contains(desc, "child's shell starts with cwd at that child scope") {
		t.Fatalf("fork scope description should clarify child cwd, got %q", desc)
	}
	if !strings.Contains(desc, "Must be `.` or a relative path under your current scope") {
		t.Fatalf("fork scope description should require relative paths, got %q", desc)
	}
	if !strings.Contains(desc, "absolute paths are invalid") {
		t.Fatalf("fork scope description should reject absolute paths, got %q", desc)
	}
	if !strings.Contains(desc, "`scope` narrows working area inside the child lineage") {
		t.Fatalf("fork scope description should clarify lineage-local scope, got %q", desc)
	}
	if strings.Contains(desc, "from the parent view") {
		t.Fatalf("fork workspace description should no longer imply parent-visible child writes, got %q", desc)
	}
	if !strings.Contains(modeDesc, "wait: block until all children finish and return every result") {
		t.Fatalf("fork mode description should describe wait mechanics, got %q", modeDesc)
	}
	if !strings.Contains(modeDesc, "for comparison or merge") {
		t.Fatalf("fork mode description should describe wait comparison/merge use, got %q", modeDesc)
	}
	if !strings.Contains(modeDesc, "race (default): first child to exit 0 wins and the rest are killed") {
		t.Fatalf("fork mode description should describe race mechanics, got %q", modeDesc)
	}
	if !strings.Contains(modeDesc, "any one child can produce an acceptable artifact/service") {
		t.Fatalf("fork mode description should describe race use, got %q", modeDesc)
	}
	if !strings.Contains(adoptDesc, "winner's filesystem artifact should become parent state") {
		t.Fatalf("fork adopt_winner description should describe artifact adoption, got %q", adoptDesc)
	}
}

func TestForkToolSchema_DefaultHidesStrategyCoaching(t *testing.T) {
	cfg := &config.Config{WorkspaceConfig: config.WorkspaceConfig{WorkspaceEnabled: true}}
	schema := ForkToolSchema(cfg)

	if !strings.Contains(schema.Description, "Spawn one or more child agents with the parent's current visible context") {
		t.Fatalf("fork schema should retain physical affordance, got %q", schema.Description)
	}
	if !strings.Contains(schema.Description, "`mode=\"wait\"` blocks until all children finish") {
		t.Fatalf("fork schema should retain mode semantics, got %q", schema.Description)
	}
	props := schema.Parameters["properties"].(map[string]any)
	children := props["children"].(map[string]any)
	items := children["items"].(map[string]any)
	childProps := items["properties"].(map[string]any)
	intent := childProps["intent"].(map[string]any)["description"].(string)
	if !strings.Contains(intent, "parent mission remains active") {
		t.Fatalf("default intent schema should preserve parent-mission semantics, got %q", intent)
	}
	for _, text := range []string{schema.Description, intent} {
		for _, forbidden := range []string{
			"independent hypotheses, decoders, implementations, extractors, or verification strategies",
			"prefer 2-3 labeled children",
			"Do one cheap shared setup/probe",
			"distinct strategy",
			"lane-specific inputs not already visible",
			"expected artifact/service",
			"closest success check",
		} {
			if strings.Contains(text, forbidden) {
				t.Fatalf("default fork schema should hide strategy coaching %q: %s", forbidden, text)
			}
		}
	}
}

func TestForkToolSchema_HidesWorldModesWhenDisabled(t *testing.T) {
	schema := ForkToolSchema(&config.Config{ToolGates: config.ToolGates{ForkWorldEnabled: false}})
	props := schema.Parameters["properties"].(map[string]any)
	children := props["children"].(map[string]any)
	items := children["items"].(map[string]any)
	childProps := items["properties"].(map[string]any)
	required := items["required"].([]string)

	if _, ok := childProps["world"]; ok {
		t.Fatal("fork schema should not expose world when feature flag is disabled")
	}
	if len(required) != 2 || required[0] != "intent" || required[1] != "scope" {
		t.Fatalf("required = %v, want [intent scope]", required)
	}
	if strings.Contains(schema.Description, "`world=\"subjective\"`") {
		t.Fatalf("fork schema should not teach world modes when disabled, got %q", schema.Description)
	}
}

func TestForkToolSchema_ExposesWorldModesWhenEnabled(t *testing.T) {
	schema := ForkToolSchema(&config.Config{ToolGates: config.ToolGates{ForkWorldEnabled: true}})
	props := schema.Parameters["properties"].(map[string]any)
	children := props["children"].(map[string]any)
	items := children["items"].(map[string]any)
	childProps := items["properties"].(map[string]any)
	required := items["required"].([]string)

	world, ok := childProps["world"].(map[string]any)
	if !ok {
		t.Fatal("fork schema should expose world when feature flag is enabled")
	}
	if got := world["enum"].([]string); len(got) != 2 || got[0] != "host" || got[1] != "subjective" {
		t.Fatalf("world enum = %v, want [host subjective]", got)
	}
	protection, ok := childProps["protection"].(map[string]any)
	if !ok {
		t.Fatal("fork schema should expose protection when feature flag is enabled")
	}
	if got := protection["enum"].([]string); len(got) != 2 || got[0] != "none" || got[1] != "transactional" {
		t.Fatalf("protection enum = %v, want [none transactional]", got)
	}
	if len(required) != 1 || required[0] != "intent" {
		t.Fatalf("required = %v, want [intent]", required)
	}
	scope := childProps["scope"].(map[string]any)
	desc := scope["description"].(string)
	if !strings.Contains(desc, "only meaningful for `world=\"subjective\"`") {
		t.Fatalf("scope description should scope itself to subjective children, got %q", desc)
	}
	if !strings.Contains(schema.Description, "private lineage requires the `overlay` workspace backend") {
		t.Fatalf("fork schema should qualify private lineage by backend, got %q", schema.Description)
	}
	if !strings.Contains(schema.Description, "only `subjective + transactional` and `host + none` are legal pairs") {
		t.Fatalf("fork schema should teach supported world/protection pairs, got %q", schema.Description)
	}
}

func TestForkToolSchema_WorldModesQualifyDirectBackend(t *testing.T) {
	schema := ForkToolSchema(&config.Config{ToolGates: config.ToolGates{ForkWorldEnabled: true}, WorkspaceConfig: config.WorkspaceConfig{WorkspaceEnabled: true, WorkspaceBackend: "direct"}})
	props := schema.Parameters["properties"].(map[string]any)
	children := props["children"].(map[string]any)
	items := children["items"].(map[string]any)
	childProps := items["properties"].(map[string]any)
	world := childProps["world"].(map[string]any)
	protection := childProps["protection"].(map[string]any)

	if !strings.Contains(schema.Description, "private lineage requires the `overlay` workspace backend") {
		t.Fatalf("fork schema should qualify lineage in top-level description, got %q", schema.Description)
	}
	if !strings.Contains(world["description"].(string), "Under `direct`, writes remain host-visible and non-adoptable") {
		t.Fatalf("world description should qualify direct backend, got %q", world["description"])
	}
	if !strings.Contains(protection["description"].(string), "private lineage only under the `overlay` workspace backend") {
		t.Fatalf("protection description should qualify overlay-only privacy, got %q", protection["description"])
	}
}

func TestShToolSchema_StdinDescriptionDescribesStructuredInput(t *testing.T) {
	schema := ShToolSchema(nil)
	props := schema.Parameters["properties"].(map[string]any)
	stdinProp := props["stdin"].(map[string]any)
	desc := stdinProp["description"].(string)

	if !strings.Contains(desc, "without shell heredoc or quoting mechanics") {
		t.Fatalf("stdin description should describe structured input mechanics, got %q", desc)
	}
}
