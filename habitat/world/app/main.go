package app

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	worldpkg "github.com/kehao95/quine/internal/world"
)

func Main() int {
	return Run(os.Args[1:], os.Stdout, os.Stderr)
}

func Run(args []string, stdout, stderr io.Writer) int {
	return runWithIO(args, stdout, stderr)
}

func run(args []string, stdout, stderr io.Writer) int {
	if len(args) == 0 {
		return runHelp(stdout, stderr)
	}
	return runWithIO(args, stdout, stderr)
}

func runWithIO(args []string, stdout, stderr io.Writer) int {
	if len(args) == 0 {
		return runHelp(stdout, stderr)
	}
	if err := worldpkg.EnforceSingleWorldInvocationPerShell(os.Getenv); err != nil {
		fmt.Fprintf(stderr, "world: %v\n", err)
		return 98
	}

	switch args[0] {
	case "--help", "-h", "help":
		return runHelp(stdout, stderr)
	case "get":
		return runGet(args[1:], stdout, stderr)
	case "validate":
		return runValidate(args[1:], stdout, stderr)
	case "reset":
		return runReset(stdout, stderr)
	default:
		fmt.Fprintf(stderr, "world: unknown command %q\n", args[0])
		return runHelp(stderr, stderr)
	}
}

func runHelp(target, stderr io.Writer) int {
	spec, _, err := worldpkg.LoadDefault()
	if err != nil {
		// No spec loaded; show generic usage.
		writeUsage(target)
		return 0
	}
	if spec.Budgeted() {
		fmt.Fprint(target, worldpkg.FormatBudgetedHelpWithAgentLimit(spec.Config.Cells, spec.Config.Budget, spec.Config.AgentGetLimit, spec.Config.ResetQuorum))
		return 0
	}
	writeUsage(target)
	return 0
}

func resolveBudgetedAgentID(spec *worldpkg.Spec, pid int) (string, error) {
	agent := strings.TrimSpace(worldpkg.AgentID())
	if agent != "" {
		return agent, nil
	}
	if spec != nil && spec.Config != nil && (spec.Config.AgentGetLimit > 0 || spec.Config.ResetQuorum > 0) {
		return "", fmt.Errorf("agent identity unavailable; use the runtime-provided world executable")
	}
	return fmt.Sprintf("pid-%d", pid), nil
}

func runGet(args []string, stdout, stderr io.Writer) int {
	if len(args) != 1 || strings.TrimSpace(args[0]) == "" {
		fmt.Fprintln(stderr, "world: usage: world get <id>")
		return 1
	}

	spec, _, err := worldpkg.LoadDefault()
	if err != nil {
		fmt.Fprintf(stderr, "world: %v\n", err)
		return 1
	}

	if spec.Budgeted() {
		return runBudgetedGet(spec, args[0], stdout, stderr)
	}

	payload, err := spec.Payload(args[0])
	if err != nil {
		if errors.Is(err, worldpkg.ErrUnknownID) {
			fmt.Fprintf(stderr, "world: %v\n", err)
			return 2
		}
		fmt.Fprintf(stderr, "world: %v\n", err)
		return 1
	}

	if _, err := io.WriteString(stdout, payload); err != nil {
		fmt.Fprintf(stderr, "world: write stdout: %v\n", err)
		return 1
	}
	if !strings.HasSuffix(payload, "\n") {
		if _, err := io.WriteString(stdout, "\n"); err != nil {
			fmt.Fprintf(stderr, "world: write stdout: %v\n", err)
			return 1
		}
	}

	return 0
}

func runValidate(args []string, stdout, stderr io.Writer) int {
	if len(args) != 1 || strings.TrimSpace(args[0]) == "" {
		fmt.Fprintln(stderr, "world: usage: world validate <path>")
		return 1
	}

	spec, _, err := worldpkg.LoadDefault()
	if err != nil {
		fmt.Fprintf(stderr, "world: %v\n", err)
		return 1
	}
	if !spec.Budgeted() {
		fmt.Fprintln(stderr, "world: validate is only available in budgeted mode")
		return 1
	}

	return runBudgetedValidate(strings.TrimSpace(args[0]), stdout, stderr)
}

func runBudgetedGet(spec *worldpkg.Spec, id string, stdout, stderr io.Writer) int {
	id = strings.TrimSpace(id)
	stateDir := spec.StateDir()
	pid := os.Getpid()
	agent, err := resolveBudgetedAgentID(spec, pid)
	if err != nil {
		fmt.Fprintf(stderr, "world: %v\n", err)
		return 1
	}

	lock, err := worldpkg.LockState(stateDir)
	if err != nil {
		fmt.Fprintf(stderr, "world: lock: %v\n", err)
		return 1
	}
	defer worldpkg.UnlockState(lock)

	st, err := worldpkg.LoadState(stateDir, spec.Config.Budget)
	if err != nil {
		fmt.Fprintf(stderr, "world: %v\n", err)
		return 1
	}

	if !st.ConsumeAgentGet(agent, spec.Config.AgentGetLimit) {
		fmt.Fprintf(stderr, "world: get limit exhausted (%d/%d this epoch). Use `world reset` to start a new epoch.\n", spec.Config.AgentGetLimit, spec.Config.AgentGetLimit)
		_ = worldpkg.AppendEvent(stateDir, worldpkg.Event{
			Time: worldpkg.NowString(), Action: "get", Cell: id,
			Agent: agent, PID: pid, BudgetAfter: st.BudgetRemaining, Result: "agent_get_limit_exhausted",
		})
		return 5
	}
	if err := st.Save(stateDir); err != nil {
		fmt.Fprintf(stderr, "world: save state: %v\n", err)
		return 1
	}

	// Check budget.
	if st.BudgetRemaining <= 0 {
		fmt.Fprintf(stderr, "world: budget exhausted (0/%d remaining). Use `world reset` to reset the environment and try again.\n", st.BudgetTotal)
		_ = worldpkg.AppendEvent(stateDir, worldpkg.Event{
			Time: worldpkg.NowString(), Action: "get", Cell: id,
			Agent: agent, PID: pid, BudgetAfter: 0, Result: "budget_exhausted",
		})
		return 3
	}

	// Look up cell.
	payload, lookupErr := spec.Payload(id)
	if lookupErr != nil {
		if errors.Is(lookupErr, worldpkg.ErrUnknownID) {
			fmt.Fprintf(stderr, "world: %v\n", lookupErr)
			return 2
		}
		fmt.Fprintf(stderr, "world: %v\n", lookupErr)
		return 1
	}

	// Deduct budget.
	st.BudgetRemaining--
	st.Collected[id] = worldpkg.CollectedCell{Value: payload, PID: pid}

	if err := st.Save(stateDir); err != nil {
		fmt.Fprintf(stderr, "world: save state: %v\n", err)
		return 1
	}

	_ = worldpkg.AppendEvent(stateDir, worldpkg.Event{
		Time: worldpkg.NowString(), Action: "get", Cell: id,
		Agent: agent, PID: pid, BudgetAfter: st.BudgetRemaining, Result: "ok",
	})

	result := worldpkg.FormatGetResult(id, payload, st.CurrentGeneration(), st.BudgetRemaining, st.BudgetTotal)
	fmt.Fprintln(stdout, result)
	return 0
}

func runBudgetedValidate(resultsPath string, stdout, stderr io.Writer) int {
	spec, _, err := worldpkg.LoadDefault()
	if err != nil {
		fmt.Fprintf(stderr, "world: %v\n", err)
		return 1
	}
	if !spec.Budgeted() {
		fmt.Fprintln(stderr, "world: validate is only available in budgeted mode")
		return 1
	}

	stateDir := spec.StateDir()
	pid := os.Getpid()
	agent, err := resolveBudgetedAgentID(spec, pid)
	if err != nil {
		fmt.Fprintf(stderr, "world: %v\n", err)
		return 1
	}

	lock, err := worldpkg.LockState(stateDir)
	if err != nil {
		fmt.Fprintf(stderr, "world: lock: %v\n", err)
		return 1
	}
	defer worldpkg.UnlockState(lock)

	// Reload the spec under the same lock used by reset so validate sees one
	// stable world epoch when checking the file.
	spec, _, err = worldpkg.LoadDefault()
	if err != nil {
		fmt.Fprintf(stderr, "world: %v\n", err)
		return 1
	}

	st, err := worldpkg.LoadState(stateDir, spec.Config.Budget)
	if err != nil {
		fmt.Fprintf(stderr, "world: %v\n", err)
		return 1
	}

	record, evalErr := worldpkg.EvaluateResultsFile(resultsPath, spec.Items)
	if evalErr != nil {
		fmt.Fprintf(stderr, "world: %v\n", evalErr)
		return 1
	}

	_ = worldpkg.AppendEvent(stateDir, worldpkg.Event{
		Time:        worldpkg.NowString(),
		Action:      "validate",
		Agent:       agent,
		PID:         pid,
		BudgetAfter: st.BudgetRemaining,
		Result:      record.EventResult(),
	})

	if record.Accepted {
		fmt.Fprintln(stdout, record.Message)
		return record.ExitCode()
	}
	fmt.Fprintln(stderr, record.Message)
	return record.ExitCode()
}

func runReset(stdout, stderr io.Writer) int {
	spec, _, err := worldpkg.LoadDefault()
	if err != nil {
		fmt.Fprintf(stderr, "world: %v\n", err)
		return 1
	}

	if !spec.Budgeted() {
		fmt.Fprintln(stderr, "world: reset is only available in budgeted mode")
		return 1
	}

	stateDir := spec.StateDir()
	pid := os.Getpid()
	agent, err := resolveBudgetedAgentID(spec, pid)
	if err != nil {
		fmt.Fprintf(stderr, "world: %v\n", err)
		return 1
	}

	lock, lockErr := worldpkg.LockState(stateDir)
	if lockErr != nil {
		fmt.Fprintf(stderr, "world: lock: %v\n", lockErr)
		return 1
	}
	defer worldpkg.UnlockState(lock)

	st, err := worldpkg.LoadState(stateDir, spec.Config.Budget)
	if err != nil {
		fmt.Fprintf(stderr, "world: %v\n", err)
		return 1
	}

	if spec.Config.ResetQuorum > 0 {
		votes, _ := st.RecordResetVote(agent)
		if err := st.Save(stateDir); err != nil {
			fmt.Fprintf(stderr, "world: save state: %v\n", err)
			return 1
		}
		if votes < spec.Config.ResetQuorum {
			_ = worldpkg.AppendEvent(stateDir, worldpkg.Event{
				Time: worldpkg.NowString(), Action: "reset",
				Agent: agent, PID: pid, BudgetAfter: st.BudgetRemaining, Result: "pending",
			})
			fmt.Fprintln(stdout, worldpkg.FormatResetPendingResult(st.CurrentGeneration(), votes, spec.Config.ResetQuorum))
			return 0
		}
	}

	// Regenerate all cell values.
	newItems, err := worldpkg.GenerateItems(spec.Config.Cells)
	if err != nil {
		fmt.Fprintf(stderr, "world: generate items: %v\n", err)
		return 1
	}
	spec.Items = newItems

	// Save the regenerated spec back to the spec file.
	specPath := worldpkg.DefaultSpecPath()
	if err := worldpkg.SaveSpec(spec, specPath); err != nil {
		fmt.Fprintf(stderr, "world: save spec: %v\n", err)
		return 1
	}

	// Reset state.
	st.BudgetRemaining = st.BudgetTotal
	st.Collected = make(map[string]worldpkg.CollectedCell)
	st.AgentGets = make(map[string]int)
	st.ResetVotes = make(map[string]bool)
	st.Resets++

	if err := st.Save(stateDir); err != nil {
		fmt.Fprintf(stderr, "world: save state: %v\n", err)
		return 1
	}

	_ = worldpkg.AppendEvent(stateDir, worldpkg.Event{
		Time: worldpkg.NowString(), Action: "reset", Cell: "",
		Agent: agent, PID: pid, BudgetAfter: st.BudgetTotal, Result: "ok",
	})

	result := worldpkg.FormatResetResult(st.BudgetTotal, st.Resets)
	fmt.Fprintln(stdout, result)
	return 0
}

func writeUsage(w io.Writer) {
	fmt.Fprintln(w, "Usage:")
	fmt.Fprintln(w, "  world --help")
	fmt.Fprintln(w, "  world get <id>")
	fmt.Fprintln(w, "  world validate <path>")
	fmt.Fprintln(w, "")
	fmt.Fprintln(w, "Description:")
	fmt.Fprintln(w, "  Return the raw payload for one item in the current world.")
	fmt.Fprintln(w, "")
	fmt.Fprintln(w, "World Source:")
	fmt.Fprintf(w, "  Reads its item mapping from embedded build data or from ./%s.\n", worldpkg.DefaultSpecFilename)
	fmt.Fprintln(w, "")
	fmt.Fprintln(w, "Output:")
	fmt.Fprintln(w, "  On success, writes the selected item's raw payload to stdout followed by a newline.")
	fmt.Fprintln(w, "")
	fmt.Fprintln(w, "Exit Status:")
	fmt.Fprintln(w, "  0  success")
	fmt.Fprintln(w, "  1  usage or runtime error")
	fmt.Fprintln(w, "  2  invalid id")
	fmt.Fprintln(w, "  5  get limit exhausted")
	fmt.Fprintln(w, "  4  validate rejected")
}
