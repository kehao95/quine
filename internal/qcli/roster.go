package qcli

import (
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

type RosterEntry struct {
	PID           int    `json:"pid"`
	Session       string `json:"session"`
	RunID         string `json:"run_id"`
	ParentSession string `json:"parent_session"`
	Depth         int    `json:"depth"`
	Model         string `json:"model"`
	Live          bool   `json:"live"`
	Pending       int    `json:"pending"`
	AgentRoot     string `json:"agent_root"`
	Attached      bool   `json:"attached"`
}

func ScanRoster(runtimeRoot string, attached *Agent) []RosterEntry {
	root := strings.TrimSpace(runtimeRoot)
	if root == "" {
		return nil
	}
	entries, err := os.ReadDir(filepath.Join(root, "pid"))
	if err != nil {
		return nil
	}
	seen := map[string]bool{}
	var roster []RosterEntry
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		pid, err := strconv.Atoi(e.Name())
		if err != nil || pid <= 0 {
			continue
		}
		resolved, err := filepath.EvalSymlinks(filepath.Join(root, "pid", e.Name()))
		if err != nil {
			continue
		}
		agent, err := loadAgentFromRoot(resolved)
		if err != nil {
			continue
		}
		key := agent.Session
		if key == "" {
			key = strconv.Itoa(agent.PID)
		}
		if seen[key] {
			continue
		}
		seen[key] = true
		pending := 0
		if st, err := ReadStatus(agent, false); err == nil {
			pending = st.Pending
		}
		roster = append(roster, RosterEntry{
			PID:           agent.PID,
			Session:       agent.Session,
			RunID:         agent.RunID,
			ParentSession: agent.ParentSession,
			Depth:         agent.Depth,
			Model:         agent.Model,
			Live:          pidLive(agent.PID),
			Pending:       pending,
			AgentRoot:     agent.AgentRoot,
			Attached:      attached != nil && attached.AgentRoot == agent.AgentRoot,
		})
	}
	sort.Slice(roster, func(i, j int) bool {
		if roster[i].PID != roster[j].PID {
			return roster[i].PID < roster[j].PID
		}
		return roster[i].Session < roster[j].Session
	})
	return roster
}

func DefaultAgent(runtimeRoot string) (Agent, error) {
	peers := ScanRoster(runtimeRoot, nil)
	if len(peers) != 1 {
		return Agent{}, ErrTargetNotFound
	}
	return Resolve(ResolveOptions{RuntimeRoot: runtimeRoot, PID: peers[0].PID})
}
