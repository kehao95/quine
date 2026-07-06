package runtime

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/kehao95/quine/internal/config"
	"github.com/kehao95/quine/internal/tools"
)

func bootstrapConfigSurfaceRuntime(t *testing.T, cfg *config.Config, mission string) *Runtime {
	t.Helper()
	rt := NewWithProvider(cfg, &mockProvider{})
	silenceRuntime(rt)
	rt.originalInput = mission
	t.Cleanup(func() {
		if err := rt.cleanupAgentRoot(); err != nil {
			t.Fatalf("cleanupAgentRoot failed: %v", err)
		}
	})
	if err := rt.bootstrapAgentRoot(); err != nil {
		t.Fatalf("bootstrapAgentRoot failed: %v", err)
	}
	return rt
}

func TestBootstrapAgentRootMaterializesConfigSurface(t *testing.T) {
	cfg := testCfg(t)
	cfg.MaxTurns = 7 // non-default knob: must show up in the resolved position
	bootstrapConfigSurfaceRuntime(t, cfg, "materialize config surface")

	configDir := filepath.Join(cfg.AgentRoot(), "config")

	regPath := filepath.Join(configDir, "registry.json")
	regData, err := os.ReadFile(regPath)
	if err != nil {
		t.Fatalf("read config/registry.json: %v", err)
	}
	wantReg, err := config.RegistryJSON()
	if err != nil {
		t.Fatalf("render compiled registry: %v", err)
	}
	wantReg = append(wantReg, '\n')
	if !bytes.Equal(regData, wantReg) {
		t.Fatalf("config/registry.json drifts from the compiled RegistryJSON payload")
	}
	regInfo, err := os.Stat(regPath)
	if err != nil {
		t.Fatalf("stat config/registry.json: %v", err)
	}
	if regInfo.Mode().Perm() != 0o444 {
		t.Fatalf("config/registry.json permissions = %o, want 444", regInfo.Mode().Perm())
	}

	resolvedPath := filepath.Join(configDir, "resolved.env")
	resolvedData, err := os.ReadFile(resolvedPath)
	if err != nil {
		t.Fatalf("read config/resolved.env: %v", err)
	}
	resolved := string(resolvedData)
	for _, want := range []string{
		"\n" + config.EnvMaxTurns + "=7\n", // non-default value
		"\n" + config.EnvModelID + "=" + cfg.ModelID + "\n",
		"\n" + config.EnvSessionID + "=" + cfg.SessionID + "\n",
		"\n" + config.EnvDepth + "=0\n",
	} {
		if !strings.Contains(resolved, want) {
			t.Fatalf("config/resolved.env missing %q:\n%s", strings.TrimSpace(want), resolved)
		}
	}
	if strings.Contains(resolved, cfg.APIKey) {
		t.Fatalf("config/resolved.env must not contain the raw API key value")
	}
	if !strings.Contains(resolved, config.EnvAPIKey+" is set; value redacted") {
		t.Fatalf("config/resolved.env should record API key presence without the value:\n%s", resolved)
	}
	resolvedInfo, err := os.Stat(resolvedPath)
	if err != nil {
		t.Fatalf("stat config/resolved.env: %v", err)
	}
	if resolvedInfo.Mode().Perm() != 0o644 {
		t.Fatalf("config/resolved.env permissions = %o, want 644 (runtime-owned cache, not read-only)", resolvedInfo.Mode().Perm())
	}

	snapPath := filepath.Join(cfg.AgentRoot(), "inc", "0", "config.env")
	snapData, err := os.ReadFile(snapPath)
	if err != nil {
		t.Fatalf("read inc/0/config.env birth snapshot: %v", err)
	}
	if !bytes.Equal(snapData, resolvedData) {
		t.Fatalf("inc/0/config.env should equal resolved.env at bootstrap")
	}
	snapInfo, err := os.Stat(snapPath)
	if err != nil {
		t.Fatalf("stat inc/0/config.env: %v", err)
	}
	if snapInfo.Mode().Perm() != 0o444 {
		t.Fatalf("inc/0/config.env permissions = %o, want 444", snapInfo.Mode().Perm())
	}
}

func TestBootstrapAgentRootRefreshesConfigSurfaceOnDrift(t *testing.T) {
	cfg := testCfg(t)
	rt := bootstrapConfigSurfaceRuntime(t, cfg, "refresh drifted config surface")

	configDir := filepath.Join(cfg.AgentRoot(), "config")
	regPath := filepath.Join(configDir, "registry.json")
	resolvedPath := filepath.Join(configDir, "resolved.env")

	if err := os.Chmod(regPath, 0o644); err != nil {
		t.Fatalf("chmod registry.json writable: %v", err)
	}
	if err := os.WriteFile(regPath, []byte("{\"stale\":true}\n"), 0o644); err != nil {
		t.Fatalf("corrupt registry.json: %v", err)
	}
	if err := os.Chmod(regPath, 0o444); err != nil {
		t.Fatalf("chmod registry.json read-only: %v", err)
	}
	if err := os.WriteFile(resolvedPath, []byte("QUINE_STALE=1\n"), 0o644); err != nil {
		t.Fatalf("corrupt resolved.env: %v", err)
	}

	if err := rt.bootstrapAgentRoot(); err != nil {
		t.Fatalf("second bootstrapAgentRoot failed: %v", err)
	}

	regData, err := os.ReadFile(regPath)
	if err != nil {
		t.Fatalf("read refreshed registry.json: %v", err)
	}
	wantReg, err := config.RegistryJSON()
	if err != nil {
		t.Fatalf("render compiled registry: %v", err)
	}
	wantReg = append(wantReg, '\n')
	if !bytes.Equal(regData, wantReg) {
		t.Fatalf("registry.json was not refreshed after drift")
	}
	regInfo, err := os.Stat(regPath)
	if err != nil {
		t.Fatalf("stat refreshed registry.json: %v", err)
	}
	if regInfo.Mode().Perm() != 0o444 {
		t.Fatalf("refreshed registry.json permissions = %o, want 444", regInfo.Mode().Perm())
	}

	resolvedData, err := os.ReadFile(resolvedPath)
	if err != nil {
		t.Fatalf("read refreshed resolved.env: %v", err)
	}
	if strings.Contains(string(resolvedData), "QUINE_STALE") {
		t.Fatalf("resolved.env was not re-rendered after drift")
	}
	if !bytes.Equal(resolvedData, cfg.ResolvedEnv()) {
		t.Fatalf("refreshed resolved.env should equal the current resolved position")
	}
}

func TestIncarnationConfigSnapshotWrittenOnce(t *testing.T) {
	cfg := testCfg(t)
	cfg.MaxTurns = 7
	rt := bootstrapConfigSurfaceRuntime(t, cfg, "write-once birth snapshot")

	snapPath := filepath.Join(cfg.AgentRoot(), "inc", "0", "config.env")
	birth, err := os.ReadFile(snapPath)
	if err != nil {
		t.Fatalf("read birth snapshot: %v", err)
	}

	// A later re-sync in the same incarnation must not rewrite the birth
	// snapshot even when the resolved position has changed.
	rt.cfg.MaxTurns = 9
	if err := rt.bootstrapAgentRoot(); err != nil {
		t.Fatalf("second bootstrapAgentRoot failed: %v", err)
	}

	resolvedData, err := os.ReadFile(filepath.Join(cfg.AgentRoot(), "config", "resolved.env"))
	if err != nil {
		t.Fatalf("read resolved.env: %v", err)
	}
	if !strings.Contains(string(resolvedData), config.EnvMaxTurns+"=9\n") {
		t.Fatalf("resolved.env should reflect the current position after re-sync")
	}
	after, err := os.ReadFile(snapPath)
	if err != nil {
		t.Fatalf("re-read birth snapshot: %v", err)
	}
	if !bytes.Equal(birth, after) {
		t.Fatalf("inc/0/config.env birth snapshot was rewritten; it must be immutable")
	}
	if !strings.Contains(string(after), config.EnvMaxTurns+"=7\n") {
		t.Fatalf("birth snapshot should keep the bootstrap-time position")
	}
}

func writeStagedNextEnvFile(t *testing.T, cfg *config.Config, content string) string {
	t.Helper()
	path := cfg.StagedNextEnvPath()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("mkdir staged config dir: %v", err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write staged config: %v", err)
	}
	return path
}

func TestBootstrapArchivesAndClearsStagedConfigOnExecReentry(t *testing.T) {
	t.Setenv(tools.ContextBootstrapEnv, "")
	cfg := testCfg(t)

	// Predecessor incarnation (inc/0) bootstraps and stages overrides for
	// its successor.
	bootstrapConfigSurfaceRuntime(t, cfg, "predecessor stages config")
	staged := "QUINE_OUTPUT_TRUNCATE=31337\n"
	stagedPath := writeStagedNextEnvFile(t, cfg, staged)

	// Successor: exec handover signaled by the staged bootstrap context.
	// (In the real flow the staged values are already baked into this
	// process's envp; the bootstrap step is pure bookkeeping.)
	t.Setenv(tools.ContextBootstrapEnv, t.TempDir())
	successor := NewWithProvider(cfg, &mockProvider{})
	silenceRuntime(successor)
	successor.originalInput = "successor consumes config"
	if err := successor.bootstrapAgentRoot(); err != nil {
		t.Fatalf("successor bootstrapAgentRoot failed: %v", err)
	}

	if cfg.IncarnationID != 1 {
		t.Fatalf("successor incarnation id = %d, want 1", cfg.IncarnationID)
	}
	if _, err := os.Lstat(stagedPath); !os.IsNotExist(err) {
		t.Fatalf("config/next.env should be cleared after a staged reentry (stat err = %v)", err)
	}

	// Archive lands in the PREDECESSOR's incarnation dir: inc/0, the
	// incarnation that staged the file.
	archivePath := filepath.Join(cfg.SessionIncarnationPath("", 0), "config-applied.env")
	archived, err := os.ReadFile(archivePath)
	if err != nil {
		t.Fatalf("read inc/0/config-applied.env archive: %v", err)
	}
	if string(archived) != staged {
		t.Fatalf("archive content = %q, want the staged content %q", archived, staged)
	}
	info, err := os.Stat(archivePath)
	if err != nil {
		t.Fatalf("stat archive: %v", err)
	}
	if info.Mode().Perm() != 0o444 {
		t.Fatalf("archive permissions = %o, want 444 (immutable applied-overrides archive)", info.Mode().Perm())
	}
	// The archive is also reachable through the agent-root inc/ symlink.
	if _, err := os.ReadFile(filepath.Join(cfg.AgentRoot(), "inc", "0", "config-applied.env")); err != nil {
		t.Fatalf("archive not visible through agent-root inc/ symlink: %v", err)
	}
	// The successor's own birth snapshot exists alongside.
	if _, err := os.Stat(filepath.Join(cfg.SessionIncarnationPath("", 1), "config.env")); err != nil {
		t.Fatalf("successor birth snapshot inc/1/config.env missing: %v", err)
	}

	// Idempotence: a re-run of the bootstrap (env still set, file gone)
	// is a clean no-op.
	if err := successor.bootstrapAgentRoot(); err != nil {
		t.Fatalf("re-bootstrap after consumption failed: %v", err)
	}
}

func TestBootstrapLeavesStagedConfigWithoutExecReentry(t *testing.T) {
	t.Setenv(tools.ContextBootstrapEnv, "")
	cfg := testCfg(t)
	rt := bootstrapConfigSurfaceRuntime(t, cfg, "non-staged bootstrap")

	staged := "QUINE_OUTPUT_TRUNCATE=31337\n"
	stagedPath := writeStagedNextEnvFile(t, cfg, staged)

	// A non-staged re-bootstrap (fresh start/resume shape) must not touch
	// the file: it belongs to THIS incarnation's own future exec.
	if err := rt.bootstrapAgentRoot(); err != nil {
		t.Fatalf("re-bootstrap failed: %v", err)
	}

	after, err := os.ReadFile(stagedPath)
	if err != nil {
		t.Fatalf("config/next.env should survive a non-staged bootstrap: %v", err)
	}
	if string(after) != staged {
		t.Fatalf("config/next.env content changed: %q -> %q", staged, after)
	}
	if _, err := os.Lstat(filepath.Join(cfg.SessionIncarnationPath("", 0), "config-applied.env")); !os.IsNotExist(err) {
		t.Fatalf("no config-applied.env archive may appear without a staged reentry (stat err = %v)", err)
	}
}

func TestConsumeStagedConfigWithoutPredecessorArchivesUnderCurrent(t *testing.T) {
	// A staged reentry at incarnation 0 (no predecessor dir — off-nominal
	// lineage state) keeps the evidence under the current incarnation.
	t.Setenv(tools.ContextBootstrapEnv, t.TempDir())
	cfg := testCfg(t)
	rt := NewWithProvider(cfg, &mockProvider{})
	silenceRuntime(rt)
	rt.originalInput = "orphan staged reentry"
	t.Cleanup(func() {
		if err := rt.cleanupAgentRoot(); err != nil {
			t.Fatalf("cleanupAgentRoot failed: %v", err)
		}
	})

	staged := "QUINE_OUTPUT_TRUNCATE=31337\n"
	stagedPath := writeStagedNextEnvFile(t, cfg, staged)

	if err := rt.bootstrapAgentRoot(); err != nil {
		t.Fatalf("bootstrapAgentRoot failed: %v", err)
	}
	if cfg.IncarnationID != 0 {
		t.Fatalf("incarnation id = %d, want 0 (no retained state)", cfg.IncarnationID)
	}
	if _, err := os.Lstat(stagedPath); !os.IsNotExist(err) {
		t.Fatalf("next.env should still be consumed (stat err = %v)", err)
	}
	archived, err := os.ReadFile(filepath.Join(cfg.SessionIncarnationPath("", 0), "config-applied.env"))
	if err != nil {
		t.Fatalf("fallback archive under the current incarnation missing: %v", err)
	}
	if string(archived) != staged {
		t.Fatalf("fallback archive content = %q, want %q", archived, staged)
	}
}

func TestLogCapabilityEnvDrift(t *testing.T) {
	cfg := testCfg(t)
	rt := NewWithProvider(cfg, &mockProvider{})
	silenceRuntime(rt)
	var logged []string
	rt.log = func(format string, args ...any) {
		logged = append(logged, fmt.Sprintf(format, args...))
	}

	// One genuine drift name (registry-unknown), plus every excluded
	// non-registry runtime namespace — the latter must stay silent.
	t.Setenv("QUINE_BOGUS_RENAMED_KNOB", "1")
	t.Setenv("QUINE_JOB_SHELL", "/bin/sh")
	t.Setenv("QUINE_TEST_HOOK", "1")
	t.Setenv("QUINE_AGENT_ROOT", "/somewhere")
	t.Setenv("QUINE_WORKSPACE_ENABLED", "1")
	t.Setenv(tools.ContextBootstrapEnv, "/staged")

	rt.logCapabilityEnvDrift()

	all := strings.Join(logged, "\n")
	if !strings.Contains(all, "QUINE_BOGUS_RENAMED_KNOB") {
		t.Fatalf("drift log should name the registry-unknown env, got %q", all)
	}
	for _, excluded := range []string{
		"QUINE_JOB_SHELL",
		"QUINE_TEST_HOOK",
		"QUINE_AGENT_ROOT",
		"QUINE_WORKSPACE_ENABLED",
		tools.ContextBootstrapEnv,
	} {
		if strings.Contains(all, excluded) {
			t.Fatalf("drift log must stay silent for excluded namespace %s, got %q", excluded, all)
		}
	}

	// After the unknown name is gone it must not be reported again. (Total
	// silence is not asserted: the ambient test environment may carry
	// legitimate drift such as a profile's QUINE_PROVIDER label.)
	t.Setenv("QUINE_BOGUS_RENAMED_KNOB", "")
	if err := os.Unsetenv("QUINE_BOGUS_RENAMED_KNOB"); err != nil {
		t.Fatalf("unset bogus env: %v", err)
	}
	logged = nil
	rt.logCapabilityEnvDrift()
	if strings.Contains(strings.Join(logged, "\n"), "QUINE_BOGUS_RENAMED_KNOB") {
		t.Fatalf("cleared env still reported as drift: %q", logged)
	}
}

func TestOnConfigMutatedReRendersResolvedEnvOnly(t *testing.T) {
	cfg := testCfg(t)
	cfg.WorkspaceEnabled = true
	cfg.WorkspaceRoot = t.TempDir()
	cfg.Workspace = cfg.WorkspaceRoot
	cfg.WorkspaceBackend = "direct"
	cfg.WorkspaceCurrentRevision = "rev-birth"
	rt := bootstrapConfigSurfaceRuntime(t, cfg, "config mutation re-render")

	configDir := filepath.Join(cfg.AgentRoot(), "config")
	regPath := filepath.Join(configDir, "registry.json")
	resolvedPath := filepath.Join(configDir, "resolved.env")
	snapPath := filepath.Join(cfg.AgentRoot(), "inc", "0", "config.env")

	regBefore, err := os.ReadFile(regPath)
	if err != nil {
		t.Fatalf("read registry.json: %v", err)
	}
	snapBefore, err := os.ReadFile(snapPath)
	if err != nil {
		t.Fatalf("read birth snapshot: %v", err)
	}
	if !strings.Contains(string(snapBefore), config.EnvWorkspaceCurrentRevision+"=rev-birth\n") {
		t.Fatalf("birth snapshot should carry the bootstrap-time revision")
	}

	// Simulate the workspace revision switch: post-D9 the only in-process
	// config mutation (syncWorldRevisionSurface calls onConfigMutated).
	rt.cfg.WorkspaceCurrentRevision = "rev-switched"
	rt.onConfigMutated()

	resolvedData, err := os.ReadFile(resolvedPath)
	if err != nil {
		t.Fatalf("read resolved.env after mutation: %v", err)
	}
	if !strings.Contains(string(resolvedData), config.EnvWorkspaceCurrentRevision+"=rev-switched\n") {
		t.Fatalf("resolved.env should reflect the switched workspace revision:\n%s", resolvedData)
	}

	regAfter, err := os.ReadFile(regPath)
	if err != nil {
		t.Fatalf("re-read registry.json: %v", err)
	}
	if !bytes.Equal(regBefore, regAfter) {
		t.Fatalf("onConfigMutated must not touch registry.json")
	}
	snapAfter, err := os.ReadFile(snapPath)
	if err != nil {
		t.Fatalf("re-read birth snapshot: %v", err)
	}
	if !bytes.Equal(snapBefore, snapAfter) {
		t.Fatalf("onConfigMutated must not touch the immutable birth snapshot")
	}
}
