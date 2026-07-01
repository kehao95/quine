package world

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"syscall"
	"time"

	"github.com/kehao95/quine/internal/config"
)

// State holds the mutable runtime state for a budgeted world.
type State struct {
	BudgetTotal     int                      `json:"budget_total"`
	BudgetRemaining int                      `json:"budget_remaining"`
	Collected       map[string]CollectedCell `json:"collected"`
	AgentGets       map[string]int           `json:"agent_gets,omitempty"`
	ResetVotes      map[string]bool          `json:"reset_votes,omitempty"`
	Resets          int                      `json:"resets"`
}

// CollectedCell records one successful get.
type CollectedCell struct {
	Value string `json:"value"`
	PID   int    `json:"pid"`
}

// Event records one world interaction for the event log.
type Event struct {
	Time        string `json:"t"`
	Action      string `json:"action"`
	Cell        string `json:"cell,omitempty"`
	Agent       string `json:"agent,omitempty"`
	PID         int    `json:"pid"`
	BudgetAfter int    `json:"budget_after"`
	Result      string `json:"result"`
}

func agentIDFromExecutablePath(exePath string) string {
	if strings.TrimSpace(exePath) == "" {
		return ""
	}
	data, err := os.ReadFile(filepath.Join(filepath.Dir(exePath), "agent-id.txt"))
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(data))
}

// ResolveAgentID returns the stable agent identity for one world invocation.
// An executable-adjacent agent-id file, when present, seals identity and
// suppresses env-based overrides.
func ResolveAgentID(exePath string, getenv func(string) string) string {
	if id := agentIDFromExecutablePath(exePath); id != "" {
		return id
	}
	if getenv == nil {
		getenv = os.Getenv
	}
	if id := strings.TrimSpace(getenv("QUINE_RUN_ID")); id != "" {
		return id
	}
	if id := strings.TrimSpace(getenv("WORLD_AGENT_ID")); id != "" {
		return id
	}
	return strings.TrimSpace(getenv("QUINE_SESSION_ID"))
}

// AgentID returns the agent identity from an executable-adjacent agent-id file
// when present, otherwise from env vars.
func AgentID() string {
	exePath := resolveExecutablePath(os.Executable, exec.LookPath, os.Args)
	return ResolveAgentID(exePath, os.Getenv)
}

func resolveExecutablePath(
	resolve func() (string, error),
	lookPath func(string) (string, error),
	args []string,
) string {
	if resolve != nil {
		if resolved, err := resolve(); err == nil && strings.TrimSpace(resolved) != "" {
			return resolved
		}
	}
	if len(args) == 0 {
		return ""
	}
	argv0 := strings.TrimSpace(args[0])
	if argv0 == "" {
		return ""
	}
	if filepath.IsAbs(argv0) {
		return argv0
	}
	if strings.ContainsRune(argv0, filepath.Separator) {
		if resolved, err := filepath.Abs(argv0); err == nil {
			return resolved
		}
		return argv0
	}
	if lookPath == nil {
		return ""
	}
	resolved, err := lookPath(argv0)
	if err != nil {
		return ""
	}
	return resolved
}

func EnforceSingleWorldInvocationPerShell(getenv func(string) string) error {
	if getenv == nil {
		getenv = os.Getenv
	}
	if strings.TrimSpace(getenv(config.EnvWorldOnePerShell)) != "1" {
		return nil
	}

	parent := os.Getppid()
	if parent <= 1 {
		return nil
	}

	stampName := fmt.Sprintf(".world-once-%d", parent)
	if session := strings.TrimSpace(getenv("QUINE_SESSION_ID")); session != "" {
		stampName = fmt.Sprintf("%s-%s", stampName, session)
	}

	stampPath := filepath.Join(os.TempDir(), stampName)
	f, err := os.OpenFile(stampPath, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o600)
	if err != nil {
		if os.IsExist(err) {
			return fmt.Errorf("only one world command per shell execution in this experiment harness")
		}
		return fmt.Errorf("create shell guard stamp: %w", err)
	}
	return f.Close()
}

// StateDir returns the resolved state directory path.
func (s *Spec) StateDir() string {
	if s.Config != nil && s.Config.StateDir != "" {
		return s.Config.StateDir
	}
	return ".world-state"
}

// LoadState reads the current state from disk. If no state file exists,
// returns a fresh initial state.
func LoadState(dir string, budget int) (*State, error) {
	path := filepath.Join(dir, "state.json")
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &State{
				BudgetTotal:     budget,
				BudgetRemaining: budget,
				Collected:       make(map[string]CollectedCell),
				AgentGets:       make(map[string]int),
				ResetVotes:      make(map[string]bool),
			}, nil
		}
		return nil, fmt.Errorf("read state: %w", err)
	}
	var st State
	if err := json.Unmarshal(data, &st); err != nil {
		return nil, fmt.Errorf("parse state: %w", err)
	}
	if st.Collected == nil {
		st.Collected = make(map[string]CollectedCell)
	}
	if st.AgentGets == nil {
		st.AgentGets = make(map[string]int)
	}
	if st.ResetVotes == nil {
		st.ResetVotes = make(map[string]bool)
	}
	return &st, nil
}

// Save writes the state to disk.
func (st *State) Save(dir string) error {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(st, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(dir, "state.json"), data, 0o644)
}

// AppendEvent appends one event to the JSONL log.
func AppendEvent(dir string, ev Event) error {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	data, err := json.Marshal(ev)
	if err != nil {
		return err
	}
	f, err := os.OpenFile(filepath.Join(dir, "events.jsonl"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.Write(append(data, '\n'))
	return err
}

// LockState acquires an exclusive flock on the state directory.
// Returns the lock file which must be closed to release the lock.
func LockState(dir string) (*os.File, error) {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, err
	}
	f, err := os.OpenFile(filepath.Join(dir, "lock"), os.O_CREATE|os.O_RDWR, 0o644)
	if err != nil {
		return nil, err
	}
	if err := syscall.Flock(int(f.Fd()), syscall.LOCK_EX); err != nil {
		f.Close()
		return nil, fmt.Errorf("flock: %w", err)
	}
	return f, nil
}

// UnlockState releases the flock.
func UnlockState(f *os.File) {
	if f != nil {
		syscall.Flock(int(f.Fd()), syscall.LOCK_UN)
		f.Close()
	}
}

// GenerateItems creates n cells with random hex values.
func GenerateItems(n int) (map[string]string, error) {
	items := make(map[string]string, n)
	for i := 1; i <= n; i++ {
		val, err := randomHex(4)
		if err != nil {
			return nil, err
		}
		items[fmt.Sprintf("c%02d", i)] = val
	}
	return items, nil
}

func randomHex(n int) (string, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

// FormatGetResult formats the output of a successful budgeted get.
// If QUINE_PROMPT_BUDGET_VISIBILITY is "hidden", budget info is omitted.
func FormatGetResult(cellID, value string, generation, budgetRemaining, budgetTotal int) string {
	if os.Getenv(config.EnvPromptBudgetVisibility) == "hidden" {
		return fmt.Sprintf("%s\n[generation: %d]", value, generation)
	}
	return fmt.Sprintf("%s\n[generation: %d] [budget: %d/%d remaining]", value, generation, budgetRemaining, budgetTotal)
}

// FormatResetResult formats the output of a reset.
// If QUINE_PROMPT_BUDGET_VISIBILITY is "hidden", budget info is omitted.
func FormatResetResult(budgetTotal, resets int) string {
	if os.Getenv(config.EnvPromptBudgetVisibility) == "hidden" {
		return fmt.Sprintf("reset complete. all cells regenerated.\n[generation: %d]", resets+1)
	}
	return fmt.Sprintf("reset complete. all cells regenerated.\n[generation: %d] [budget: %d/%d remaining] [resets: %d]", resets+1, budgetTotal, budgetTotal, resets)
}

// FormatResetPendingResult formats the output of a pending quorum reset.
func FormatResetPendingResult(generation, votes, quorum int) string {
	return fmt.Sprintf("reset pending. waiting for quorum.\n[generation: %d] [reset votes: %d/%d]", generation, votes, quorum)
}

func (st *State) CurrentGeneration() int {
	return st.Resets + 1
}

// SortedCellIDs returns cell IDs from a map in sorted order.
func SortedCellIDs(items map[string]string) []string {
	ids := make([]string, 0, len(items))
	for id := range items {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	return ids
}

// FormatHelp generates the --help text for a budgeted world.
func FormatBudgetedHelp(cells, budget int) string {
	hideBudget := os.Getenv(config.EnvPromptBudgetVisibility) == "hidden"

	var b strings.Builder
	b.WriteString("Usage:\n")
	b.WriteString("  world --help\n")
	b.WriteString("  world get <id>\n")
	b.WriteString("  world validate <path>\n")
	b.WriteString("  world reset\n")
	b.WriteString("\n")
	b.WriteString("Rules:\n")
	b.WriteString(fmt.Sprintf("  This world contains %d cells (c01 through c%02d).\n", cells, cells))
	if !hideBudget {
		b.WriteString(fmt.Sprintf("  There is a total budget of %d get calls.\n", budget))
		b.WriteString("  Only `world get` consumes budget; each call costs exactly 1.\n")
		b.WriteString("  `world` itself does not create or modify workspace files.\n")
		b.WriteString("  When the budget reaches 0, further get calls will fail.\n")
	} else {
		b.WriteString("  `world` itself does not create or modify workspace files.\n")
	}
	b.WriteString("\n")
	b.WriteString("Commands:\n")
	if !hideBudget {
		b.WriteString("  get <id>   Return the content of a cell. Costs exactly 1 from the budget.\n")
		b.WriteString("             Output includes the cell value and remaining budget.\n")
	} else {
		b.WriteString("  get <id>   Return the content of a cell.\n")
		b.WriteString("             Output includes the cell value.\n")
	}
	b.WriteString("  validate <path>\n")
	b.WriteString("             Check a results file against the current world.\n")
	b.WriteString("             You may retry validation after improving the file.\n")
	if !hideBudget {
		b.WriteString("  reset      Reset the budget and regenerate all cell values.\n")
		b.WriteString("             Use this to start over after a failed attempt.\n")
	} else {
		b.WriteString("  reset      Regenerate all cell values.\n")
		b.WriteString("             Use this to start over after a failed attempt.\n")
	}
	b.WriteString("\n")
	b.WriteString("Exit Status:\n")
	b.WriteString("  0  success\n")
	b.WriteString("  1  usage or runtime error\n")
	b.WriteString("  2  invalid id\n")
	if !hideBudget {
		b.WriteString("  3  budget exhausted\n")
	}
	b.WriteString("  4  validate rejected\n")
	return b.String()
}

func FormatBudgetedHelpWithAgentLimit(cells, budget, agentGetLimit, resetQuorum int) string {
	hideBudget := os.Getenv(config.EnvPromptBudgetVisibility) == "hidden"
	base := FormatBudgetedHelp(cells, budget)
	if agentGetLimit <= 0 && resetQuorum <= 0 {
		return base
	}
	if hideBudget {
		// When budget is hidden, don't add per-process limits either
		return base
	}

	lines := strings.Split(base, "\n")
	var out strings.Builder
	for _, line := range lines {
		if line == "  When the budget reaches 0, further get calls will fail." {
			out.WriteString(line + "\n")
			if agentGetLimit > 0 {
				out.WriteString(fmt.Sprintf("  Per-process get limit: %d calls per reset epoch.\n", agentGetLimit))
			}
			if resetQuorum > 0 {
				out.WriteString(fmt.Sprintf("  `world reset` executes only after %d requests in the same generation.\n", resetQuorum))
			}
			continue
		}
		if line != "" {
			out.WriteString(line + "\n")
			continue
		}
		out.WriteString("\n")
	}
	return out.String()
}

func (st *State) RecordResetVote(agent string) (votes int, alreadyRecorded bool) {
	if st.ResetVotes == nil {
		st.ResetVotes = make(map[string]bool)
	}
	alreadyRecorded = st.ResetVotes[agent]
	st.ResetVotes[agent] = true
	return len(st.ResetVotes), alreadyRecorded
}

// ConsumeAgentGet records one get call for the current reset epoch.
// It returns false if the agent has already exhausted the configured per-epoch limit.
func (st *State) ConsumeAgentGet(agent string, limit int) bool {
	if limit <= 0 || agent == "" {
		return true
	}
	if st.AgentGets == nil {
		st.AgentGets = make(map[string]int)
	}
	if st.AgentGets[agent] >= limit {
		return false
	}
	st.AgentGets[agent]++
	return true
}

// NowString returns an ISO 8601 timestamp.
func NowString() string {
	return time.Now().Format(time.RFC3339Nano)
}
