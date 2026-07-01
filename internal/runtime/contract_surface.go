package runtime

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	processControlContractVersion = "process-control/v0"
	worldSurfaceContractVersion   = "world/v0"
)

type runtimeContractManifest struct {
	ContractVersion  string                           `json:"contract_version"`
	Backend          string                           `json:"backend"`
	SessionID        string                           `json:"session_id"`
	RunID            string                           `json:"run_id"`
	IncarnationID    int                              `json:"incarnation_id"`
	AgentRoot        string                           `json:"agent_root"`
	PublicRoot       string                           `json:"public_root"`
	RuntimeRoot      string                           `json:"runtime_root"`
	Usage            string                           `json:"usage"`
	Surfaces         runtimeContractSurfaces          `json:"surfaces"`
	ControlActions   map[string]controlActionContract `json:"control_actions"`
	InboxSchema      map[string]string                `json:"inbox_schema"`
	ControlLogEvents map[string]string                `json:"control_log_events"`
	NonClaims        []string                         `json:"non_claims"`
}

type runtimeContractSurfaces struct {
	Contract        string `json:"contract"`
	Identity        string `json:"identity"`
	Inbox           string `json:"inbox"`
	Control         string `json:"control"`
	ControlLog      string `json:"control_log"`
	LiveContext     string `json:"live_context"`
	LiveGeneration  string `json:"live_generation"`
	PromptFragments string `json:"prompt_fragments"`
	WorldStatus     string `json:"world_status"`
}

type controlActionContract struct {
	Description    string `json:"description"`
	Queues         bool   `json:"queues,omitempty"`
	QueuesNonempty bool   `json:"queues_nonempty,omitempty"`
	WakesIdle      bool   `json:"wakes_idle,omitempty"`
	InjectsContext bool   `json:"injects_context,omitempty"`
	Interrupts     bool   `json:"interrupts,omitempty"`
	EmptyWrite     string `json:"empty_write,omitempty"`
}

func (r *Runtime) runtimeContractManifest() runtimeContractManifest {
	publicRoot := filepath.Join(r.cfg.AgentRoot(), "public")
	return runtimeContractManifest{
		ContractVersion: processControlContractVersion,
		Backend:         runtimeSurfaceBackendName,
		SessionID:       r.cfg.SessionID,
		RunID:           r.cfg.RunID,
		IncarnationID:   r.cfg.IncarnationID,
		AgentRoot:       r.cfg.AgentRoot(),
		PublicRoot:      publicRoot,
		RuntimeRoot:     r.cfg.RuntimeRoot(),
		Usage:           "To drive a peer process, read its public/status/contract.json, then write the named file under its public ctl/. Reads of status/* and this contract are scan-safe; writes to ctl/* are the only effectful actions. Retrieve queued payloads from status/inbox.json. The live_context surface is the canonical complete-cell context (provider input); live_generation is a transient display-only stream of in-flight generation deltas and is never provider input or recovery state.",
		Surfaces: runtimeContractSurfaces{
			Contract:        "status/contract.json",
			Identity:        "status/session.json",
			Inbox:           "status/inbox.json",
			Control:         "ctl",
			ControlLog:      "log/control.jsonl",
			LiveContext:     "../context/state/current.jsonl", // complete cells, provider input, recovered
			LiveGeneration:  "../context/state/live.jsonl",    // transient generation deltas, display-only, never input, never recovered
			PromptFragments: "../context/prompt",
			WorldStatus:     "../world/status.json",
		},
		ControlActions: map[string]controlActionContract{
			"post": {
				Description: "queue-only: a non-empty write enqueues one payload (after trailing-newline normalization); does not wake idle and does not inject context.",
				Queues:      true,
				EmptyWrite:  "no-op",
			},
			"poke": {
				Description:    "queue-and-resume: a write resumes idle without automatic context injection; a non-empty write also queues a payload for later inbox retrieval.",
				QueuesNonempty: true,
				WakesIdle:      true,
				EmptyWrite:     "resume idle without new payload",
			},
			"inject": {
				Description:    "queue-and-deliver: a non-empty write queues one payload, resumes idle, and surfaces it via incoming_messages at the next safe point.",
				Queues:         true,
				WakesIdle:      true,
				InjectsContext: true,
				EmptyWrite:     "no-op",
			},
			"interrupt": {
				Description:    "interrupt-delivery: an empty write interrupts without payload; a non-empty write also queues one payload for urgent delivery. If a tool subprocess is active, the runtime first forwards a real SIGINT into that process group.",
				QueuesNonempty: true,
				InjectsContext: true,
				Interrupts:     true,
				EmptyWrite:     "request interrupt without new payload",
			},
		},
		InboxSchema: map[string]string{
			"pending_count":          "number of queued payloads not yet retrieved",
			"messages[].id":          "stable id of a queued payload",
			"messages[].payload":     "the payload string written via a ctl action",
			"messages[].received_at": "receipt timestamp of the payload",
		},
		ControlLogEvents: map[string]string{
			"received.action":    "a control write was received (receipt is distinct from wake-up and delivery)",
			"woke.delivery":      "an idle process was resumed by a control write",
			"delivered.delivery": "a queued payload was surfaced into the process context",
		},
		NonClaims: []string{
			"no broad ctl/* tree",
			"no callback/correlation/routing protocol",
			"no full inference-preemption guarantee for interrupt",
			"no total private world contract",
			"no claim that copying a live process image via procfs is behavioral reconstruction",
		},
	}
}

func (r *Runtime) writeRuntimeContract(statusDir string) error {
	return writeJSONFile(filepath.Join(statusDir, "contract.json"), r.runtimeContractManifest())
}

type runtimeWorldStatus struct {
	ContractVersion        string `json:"contract_version"`
	World                  string `json:"world"`
	Protection             string `json:"protection"`
	Backend                string `json:"backend,omitempty"`
	WorkspaceEnabled       bool   `json:"workspace_enabled,omitempty"`
	WorkspaceRoot          string `json:"workspace_root,omitempty"`
	Workspace              string `json:"workspace,omitempty"`
	Scope                  string `json:"scope,omitempty"`
	WorkspaceSession       string `json:"workspace_session,omitempty"`
	CurrentRevision        string `json:"current_revision,omitempty"`
	RevisionMode           string `json:"revision_mode,omitempty"`
	CanRollback            bool   `json:"can_rollback,omitempty"`
	CanSwitchRevision      bool   `json:"can_switch_revision,omitempty"`
	Owner                  bool   `json:"owner,omitempty"`
	FSMutationTelemetry    bool   `json:"fs_mutation_telemetry,omitempty"`
	LastFSMutationsSurface string `json:"last_fs_mutations_surface,omitempty"`
	EventsSurface          string `json:"events_surface"`
	ResourcesSurface       string `json:"resources_surface"`
}

type runtimeWorldResources struct {
	ContractVersion string                 `json:"contract_version"`
	Workspace       *runtimeWorldWorkspace `json:"workspace,omitempty"`
	Tools           []string               `json:"tools"`
	Control         map[string]string      `json:"control"`
	NonClaims       []string               `json:"non_claims"`
}

type runtimeWorldWorkspace struct {
	Root       string `json:"root,omitempty"`
	Current    string `json:"current,omitempty"`
	Scope      string `json:"scope,omitempty"`
	Visibility string `json:"visibility"`
	Mutable    bool   `json:"mutable"`
}

type runtimeWorldEvent struct {
	Kind        string `json:"kind"`
	Timestamp   int64  `json:"timestamp"`
	Revision    string `json:"revision,omitempty"`
	FSMutations string `json:"fs_mutations,omitempty"`
}

func (r *Runtime) runtimeWorldStatus() runtimeWorldStatus {
	backend := ""
	fsMutationTelemetry := false
	if r.cfg.WorkspaceEnabled {
		backend = r.cfg.EffectiveWorkspaceBackend()
		fsMutationTelemetry = r.cfg.FSMutationTelemetryEnabled()
	}
	lastMutationSurface := ""
	if strings.TrimSpace(r.lastFSMutations) != "" {
		lastMutationSurface = "fs-mutations.latest"
	}
	revisionMode := ""
	if mode := strings.TrimSpace(string(r.cfg.WorkspaceRevisionMode)); mode != "" && mode != "none" {
		revisionMode = mode
	}
	return runtimeWorldStatus{
		ContractVersion:        worldSurfaceContractVersion,
		World:                  string(r.cfg.CurrentWorld()),
		Protection:             string(r.cfg.CurrentProtection()),
		Backend:                backend,
		WorkspaceEnabled:       r.cfg.WorkspaceEnabled,
		WorkspaceRoot:          r.cfg.WorkspaceRoot,
		Workspace:              r.cfg.Workspace,
		Scope:                  runtimeWorkspaceScope(r.cfg.WorkspaceRoot, r.cfg.Workspace),
		WorkspaceSession:       r.cfg.WorkspaceSession,
		CurrentRevision:        r.cfg.WorkspaceCurrentRevision,
		RevisionMode:           revisionMode,
		CanRollback:            r.cfg.WorkspaceTransactional(),
		CanSwitchRevision:      r.cfg.CanRestoreWorld(),
		Owner:                  r.cfg.WorkspaceOwner,
		FSMutationTelemetry:    fsMutationTelemetry,
		LastFSMutationsSurface: lastMutationSurface,
		EventsSurface:          "events.jsonl",
		ResourcesSurface:       "resources.json",
	}
}

func (r *Runtime) runtimeWorldResources() runtimeWorldResources {
	var workspace *runtimeWorldWorkspace
	visibility := "host-visible"
	if r.cfg.WorkspaceTransactional() {
		visibility = "subjective"
	}
	if r.cfg.WorkspaceEnabled {
		workspace = &runtimeWorldWorkspace{
			Root:       r.cfg.WorkspaceRoot,
			Current:    r.cfg.Workspace,
			Scope:      runtimeWorkspaceScope(r.cfg.WorkspaceRoot, r.cfg.Workspace),
			Visibility: visibility,
			Mutable:    true,
		}
	}
	nonClaims := []string{
		"world/v0 is an inspectable workspace-world slice, not a total private world",
		"resources list present runtime affordances only; absent names are not disabled affordances",
		"resources are descriptive runtime affordances, not a complete capability ontology",
	}
	if r.cfg.CanRestoreWorld() {
		nonClaims = append(nonClaims, "revision switching is available only when the active workspace backend supports restore")
	}
	return runtimeWorldResources{
		ContractVersion: worldSurfaceContractVersion,
		Workspace:       workspace,
		Tools:           r.runtimeWorldToolNames(),
		Control: map[string]string{
			"process_contract": "../status/contract.json",
			"ctl":              "../public/ctl",
			"inbox":            "../public/status/inbox.json",
			"control_log":      "../public/log/control.jsonl",
		},
		NonClaims: nonClaims,
	}
}

func (r *Runtime) runtimeWorldToolNames() []string {
	names := []string{"sh"}
	if r.cfg.ForkEnabled() {
		names = append(names, "fork")
	}
	if r.cfg.SpawnEnabled() {
		names = append(names, "spawn")
	}
	if r.cfg.ExitEnabled() {
		names = append(names, "exit")
	}
	if r.cfg.IdleToolEnabled() {
		names = append(names, "idle")
	}
	if r.cfg.ExecEnabled {
		names = append(names, "exec")
	}
	if r.cfg.VisionEnabled {
		names = append(names, "vision")
	}
	if r.cfg.CanRestoreWorld() {
		names = append(names, "switch_world")
	}
	if r.cfg.AnchorMemoryEnabled {
		names = append(names, "mark", "unfold")
	}
	if r.cfg.CanEscalate() {
		names = append(names, "escalate")
	}
	return names
}

func runtimeWorkspaceScope(root, current string) string {
	root = strings.TrimSpace(root)
	current = strings.TrimSpace(current)
	if root == "" || current == "" {
		return ""
	}
	rel, err := filepath.Rel(root, current)
	if err != nil || rel == "." {
		return "."
	}
	return filepath.ToSlash(rel)
}

func (r *Runtime) syncInspectableWorldSurface(retainedWorldDir string) error {
	if err := os.MkdirAll(retainedWorldDir, 0o755); err != nil {
		return fmt.Errorf("mkdir world surface: %w", err)
	}
	if err := writeTextFile(filepath.Join(retainedWorldDir, "workspace_root"), r.cfg.WorkspaceRoot); err != nil {
		return err
	}
	if err := writeTextFile(filepath.Join(retainedWorldDir, "workspace"), r.cfg.Workspace); err != nil {
		return err
	}
	if err := writeJSONFile(filepath.Join(retainedWorldDir, "status.json"), r.runtimeWorldStatus()); err != nil {
		return err
	}
	if err := writeJSONFile(filepath.Join(retainedWorldDir, "resources.json"), r.runtimeWorldResources()); err != nil {
		return err
	}
	eventsPath := filepath.Join(retainedWorldDir, "events.jsonl")
	if _, err := os.Stat(eventsPath); os.IsNotExist(err) {
		if err := os.WriteFile(eventsPath, nil, 0o644); err != nil {
			return fmt.Errorf("init world events: %w", err)
		}
	} else if err != nil {
		return fmt.Errorf("stat world events: %w", err)
	}

	fsMutationsPath := filepath.Join(retainedWorldDir, "fs-mutations.latest")
	mutationBlock := strings.TrimSpace(r.lastFSMutations)
	if mutationBlock == "" {
		if err := os.Remove(fsMutationsPath); err != nil && !os.IsNotExist(err) {
			return fmt.Errorf("remove fs mutations snapshot: %w", err)
		}
		return nil
	}
	if err := writeTextFile(fsMutationsPath, r.lastFSMutations); err != nil {
		return err
	}
	return r.appendWorldMutationEventOnce(retainedWorldDir, r.lastFSMutations)
}

func (r *Runtime) appendWorldMutationEventOnce(retainedWorldDir string, mutationBlock string) error {
	markerPath := filepath.Join(retainedWorldDir, ".fs-mutations.evented")
	if data, err := os.ReadFile(markerPath); err == nil && string(data) == mutationBlock {
		return nil
	}
	entry := runtimeWorldEvent{
		Kind:        "mutation_observed",
		Timestamp:   time.Now().UnixMilli(),
		Revision:    r.cfg.WorkspaceCurrentRevision,
		FSMutations: mutationBlock,
	}
	data, err := json.Marshal(entry)
	if err != nil {
		return fmt.Errorf("marshal world event: %w", err)
	}
	f, err := os.OpenFile(filepath.Join(retainedWorldDir, "events.jsonl"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return fmt.Errorf("open world events: %w", err)
	}
	defer f.Close()
	if _, err := f.Write(append(data, '\n')); err != nil {
		return fmt.Errorf("write world event: %w", err)
	}
	return os.WriteFile(markerPath, []byte(mutationBlock), 0o644)
}
