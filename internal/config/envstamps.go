package config

// envstamps.go holds the stamp builders: the ONLY place resolved Config values
// are allowed to enter the environment of a process this runtime creates.
//
// Design authority:
//   Paper/theory/views/runtime-capability/env-process-boundary-brief.md
//
// A stamp is not a default. The distinction is the entire point of this file,
// and it has exactly one test:
//
//	Would the child derive a WRONG value without this line?
//
// If yes, the runtime owns the fact and must state it — the child's depth, the
// session it descends from, the runtime root it must join, the mount it will
// run inside. None of these are derivable by the child: it would guess, and
// guess wrong, and the process tree would silently split.
//
// If no — if the line merely restates a policy knob at its compiled default,
// something nobody chose — it does not belong here. Absence is how an unset
// knob is spelled, and config/registry.json is the catalog that says what
// absence means. The deleted synthesizer (Config.baseEnv) failed this test ~50
// times per child: it serialized every knob in the registry, defaults included,
// which is how a founder came to read `QUINE_SPAWN_ENABLED=0` — a line no
// operator authored — and conclude that constructing another Quine was
// physically impossible. Adding a line here that fails the test re-creates that
// failure, whatever it is named.
//
// Everything a stamp does not name is inherited (BuildChildEnv: stamps ⊕
// override ⊕ (environ − mask)). Stamps win over both, because a runtime-owned
// fact is not an opinion the agent gets to disagree with.

import (
	"path/filepath"
	"strconv"
	"strings"
)

// trimmed is the shared guard on every string stamp: a whitespace-only value is
// no value. Absence says "compiled default" honestly; a blank line says nothing
// and looks like a choice.
func trimmed(v string) string { return strings.TrimSpace(v) }

// joinedRoot is the guard on the two runtime-ROOT stamps (QUINE_DATA_DIR,
// QUINE_RETENTION_DIR). A root is stamped to say WHICH root this process joined,
// and a relative path does not say that — it says "resolve this against your own
// cwd", which is exactly the question the stamp exists to answer.
//
// It matters because the compiled default for DataDir is the RELATIVE ".quine/"
// (config.go: Load leaves it relative), and a child does not share our cwd: an sh
// child runs inside the workspace mount, a fork child runs in the executor's
// WorkDir. Stamping ".quine/" therefore reproduced the bug it was written to
// prevent — the child re-resolved the default against a different cwd, joined a
// different runtime root, and the process tree split with no error anywhere. The
// stamp only earns its place if it names a path that means the same thing from
// any cwd.
func joinedRoot(path string) string {
	path = trimmed(path)
	if path == "" {
		return ""
	}
	abs, err := filepath.Abs(path)
	if err != nil {
		// Abs only fails when the cwd is unreadable. The relative path is still
		// better than nothing (it is what we ourselves resolved against).
		return path
	}
	return abs
}

// ShellStamps are the facts an ordinary `sh` child cannot derive for itself.
//
// It gets no lineage marks: a program started from a shell is generic
// computation, not a member of this agent tree (brief E4). What it does get is
// where it is:
//
//   - QUINE_RUN_ID — the current physical run. The `world` binary and other
//     runtime-aware tools key off it; it changes on every activation, so an
//     inherited one would name a run that is over.
//   - QUINE_AGENT_ROOT — this session's live root, so `$QUINE_AGENT_ROOT/...`
//     resolves in a shell command. Never inherited: every process has its own.
//   - QUINE_DATA_DIR — the resolved runtime root this process actually joined,
//     ABSOLUTE (joinedRoot). The `world` binary resolves its spec from this name
//     (internal/world/spec.go), and the compiled default is the RELATIVE
//     ".quine/", so a child running from a different cwd — and an sh child, inside
//     the workspace mount, always is — would silently read a DIFFERENT world file:
//     no error, just divergence. The runtime knows which root it joined; it says
//     so, in terms that mean the same thing from anywhere.
//
// Workspace physics (QUINE_WORKSPACE_ENABLED and the overlay mount
// coordinates) and the QUINE_JOB_* wrapper plumbing are per-command facts, not
// per-process ones — the sh executor injects them at each job, where the mount
// and the job dir are actually known.
func (c *Config) ShellStamps() []string {
	if c == nil {
		return nil
	}
	var stamps []string
	if runID := trimmed(c.RunID); runID != "" {
		stamps = append(stamps, envKV(EnvRunID, runID))
	}
	if root := trimmed(c.AgentRoot()); root != "" {
		stamps = append(stamps, envKV(EnvAgentRoot, root))
	}
	if dataDir := joinedRoot(c.DataDir); dataDir != "" {
		stamps = append(stamps, envKV(EnvDataDir, dataDir))
	}
	return stamps
}

// ForkChildStamps are the facts a managed fork/spawn child cannot derive: it is
// a member of THIS agent tree, and membership is not something a process can
// work out by looking around. spawn shares fork's executor and therefore these
// stamps.
//
//   - QUINE_DEPTH = parent+1, QUINE_PARENT_SESSION = this session. Tree
//     membership. A child left to inherit them would claim its parent's
//     position in the tree.
//   - QUINE_DATA_DIR — the runtime root the child must JOIN, absolute. Same
//     criticality as the shell case and sharper: a child that fell back to
//     ".quine/" relative to its own cwd would join a different runtime root, and
//     the process tree would split without a single error.
//   - QUINE_RETENTION_DIR — only when this process has one, and absolute for the
//     same reason. Retained lineage state lives under a root the operator chose;
//     a child writing its record somewhere else loses it.
//   - workspace SCOPE and lineage, when physics are on. QUINE_WORKSPACE_ROOT and
//     QUINE_WORKSPACE are DERIVED from where the runtime started, not configured:
//     WORKSPACE defaults to the startup cwd, so a child re-deriving it from its
//     own cwd lands in the wrong scope entirely. Owner is stamped 0 because a
//     child never owns its parent's workspace — it borrows a view of it.
//     QUINE_WORKSPACE_BOOTSTRAP carries the parent's workspace session so the
//     child can adopt the lineage rather than start a fresh one.
//
// Deliberately NOT stamped — and this is the file's own test applied to itself:
// QUINE_WORKSPACE_BACKEND, QUINE_WORKSPACE_OVERLAY_DRIVER,
// QUINE_WORKSPACE_REVISION_MODE, QUINE_WORKSPACE_COMMIT_ON_SIGNAL. All four are
// exec-boundary (FREE) knobs, and all four are pure Load() defaults conditioned
// on the workspace root the child is already being handed: with physics on and
// nothing authored, a child derives backend=overlay, driver=kernel,
// revision-mode=restore and commit-on-signal=false by itself — identically. An
// operator who DID author one has it in environ, where the child inherits it. So
// stamping them answered a question nobody asked, and did two kinds of harm:
//
//	(1) it put QUINE_WORKSPACE_COMMIT_ON_SIGNAL=0 — a default-valued negative
//	    nobody authored — into every workspace-enabled fork child's environ and
//	    its retained birth record. That is the exact morphology of
//	    QUINE_SPAWN_ENABLED=0, the line that taught a founder its own
//	    reproduction was impossible. Synthesis does not stop being synthesis
//	    because it is spelled "stamp".
//	(2) stamps beat the override, so an agent writing
//	    QUINE_WORKSPACE_BACKEND=direct into config/env/override would have seen it
//	    accepted by the gate, forked — and got a child running overlay anyway,
//	    with no rejection and no log line. The policy said one thing and the
//	    boundary did another.
//
// The per-child values that ARE genuinely known (a subjective child's backend and
// revision mode) are stamped by the fork executor, which is where they are known
// (fork.go, world=subjective).
//
// Also not stamped: QUINE_SESSION_ID (per-child, injected by the fork executor
// which is the only thing that knows each child's id) and
// QUINE_WORKSPACE_CURRENT_REVISION (a revision handle in the parent's own
// workspace state; a child mounts its own view and computes its own).
//
// A world=host child is stripped of the workspace stamps by the fork executor's
// conditional filterWorkspacePhysics mask — it asked not to be in the
// workspace, and stamping is not a way around that.
func (c *Config) ForkChildStamps() []string {
	if c == nil {
		return nil
	}
	stamps := []string{
		envKV(EnvDepth, strconv.Itoa(c.Depth+1)),
	}
	if session := trimmed(c.SessionID); session != "" {
		stamps = append(stamps, envKV(EnvParentSession, session))
	}
	if dataDir := joinedRoot(c.DataDir); dataDir != "" {
		stamps = append(stamps, envKV(EnvDataDir, dataDir))
	}
	if retention := joinedRoot(c.RetentionDir); retention != "" {
		stamps = append(stamps, envKV(EnvRetentionDir, retention))
	}
	if c.WorkspaceEnabled {
		stamps = append(stamps,
			envKV(EnvWorkspaceRoot, c.WorkspaceRoot),
			envKV(EnvWorkspace, c.Workspace),
			envKV(EnvWorkspaceOwner, "0"),
		)
		if session := trimmed(c.WorkspaceSession); session != "" {
			stamps = append(stamps, envKV(EnvWorkspaceBootstrap, session))
		}
	}
	return stamps
}

// SelfReentryStamps are the facts that must survive an exec: the successor is
// the SAME quine wearing a new image, and it has to know that.
//
//   - QUINE_SESSION_ID, QUINE_TAPE_ID — lineage continuity. A successor that
//     generated fresh ones would be a different agent occupying the same
//     process.
//   - QUINE_PARENT_SESSION — preserved, not rewritten. exec does not create a
//     generation; the successor's parent is still whoever forked its
//     predecessor.
//   - QUINE_DEPTH — PRESERVED at the current value. This is a bugfix, not a
//     port: the deleted ExecEnv passed a literal 0, and since depth enforcement
//     reads the in-memory Depth (precheckProcessCreation), a fork→exec→fork
//     chain used to reset its own depth budget and walk straight out of
//     QUINE_MAX_DEPTH. exec is not a birth; it does not refill a budget.
//   - QUINE_DATA_DIR, QUINE_RETENTION_DIR — the roots this process is joined
//     to. exec preserves cwd, but not stamping the root the runtime actually
//     resolved would let the successor re-derive a different one.
//   - workspace RUNTIME-EMITTED state, when physics are on: the workspace
//     session it is already inside, whether it owns it, and the revision it is
//     currently at. These are masked from inheritance (they are runtime-owned),
//     so if the runtime does not restate them the successor wakes up believing
//     it owns a workspace it does not own, at no revision.
//
// Deliberately NOT stamped: every FREE knob, including the workspace knobs an
// operator or agent authors (ROOT, WORKSPACE, BACKEND, OVERLAY_DRIVER,
// REVISION_MODE, COMMIT_ON_SIGNAL). Stamps beat the override, and the override
// MUST be able to change free knobs at exec — that is the entire
// self-modification mechanism (brief E2). exec preserves cwd and inherits
// environ, so the successor re-derives identical values anyway; stamping them
// would buy nothing and silently disable staged self-reconfiguration.
func (c *Config) SelfReentryStamps() []string {
	if c == nil {
		return nil
	}
	stamps := []string{
		envKV(EnvDepth, strconv.Itoa(c.Depth)),
	}
	if session := trimmed(c.SessionID); session != "" {
		stamps = append(stamps, envKV(EnvSessionID, session))
	}
	if tapeID := trimmed(c.TapeID); tapeID != "" {
		stamps = append(stamps, envKV(EnvTapeID, tapeID))
	}
	if parent := trimmed(c.ParentSession); parent != "" {
		stamps = append(stamps, envKV(EnvParentSession, parent))
	}
	if dataDir := joinedRoot(c.DataDir); dataDir != "" {
		stamps = append(stamps, envKV(EnvDataDir, dataDir))
	}
	if retention := joinedRoot(c.RetentionDir); retention != "" {
		stamps = append(stamps, envKV(EnvRetentionDir, retention))
	}
	if c.WorkspaceEnabled {
		stamps = append(stamps, envKV(EnvWorkspaceOwner, bool01(c.WorkspaceOwner)))
		if session := trimmed(c.WorkspaceSession); session != "" {
			stamps = append(stamps, envKV(EnvWorkspaceSession, session))
		}
		if revision := trimmed(c.WorkspaceCurrentRevision); revision != "" {
			stamps = append(stamps, envKV(EnvWorkspaceCurrentRevision, revision))
		}
	}
	return stamps
}

// StampedEnvNames is the declared name set each stamp builder in this file may
// emit — their vocabulary, without their values.
//
// It exists as a drift guard and a review checkpoint on the stamp set. A stamp
// is the ONE place a resolved Config value may legitimately enter a child env,
// so every addition to it deserves the question "is this a runtime-owned fact
// the child cannot derive, or a policy default nobody chose?" — the second kind
// is resurrected synthesis (four workspace knobs were once stamped that way).
// TestStampedEnvNamesMatchBuilders pins this declared set against what the
// builders actually emit for a fully-populated Config, so a stamp added to a
// builder without being declared here fails the build, forcing that question to
// be answered at review time rather than discovered in a trace.
//
// Scope: the builders in THIS file. QUINE_AGENT_ROOT's exec stamp lives in the
// exec executor (it merges execProcessSurfaceEnv last), and the per-command sh
// physics and QUINE_JOB_* plumbing live in the sh executor.
func StampedEnvNames(b Boundary) []string {
	switch b {
	case BoundaryShell:
		return []string{EnvRunID, EnvAgentRoot, EnvDataDir}
	case BoundaryChild:
		return []string{
			EnvDepth, EnvParentSession, EnvDataDir, EnvRetentionDir,
			EnvWorkspaceRoot, EnvWorkspace, EnvWorkspaceOwner, EnvWorkspaceBootstrap,
		}
	case BoundarySelf:
		return []string{
			EnvDepth, EnvSessionID, EnvTapeID, EnvParentSession, EnvDataDir, EnvRetentionDir,
			EnvWorkspaceOwner, EnvWorkspaceSession, EnvWorkspaceCurrentRevision,
		}
	}
	return nil
}
