package runtime

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/kehao95/quine/internal/tape"
)

const peerDiscoveryMessageType = "quine.peer_discovery"
const defaultPeerDiscoveryHeartbeatInterval = 5 * time.Second

type peerDiscoveryPeer struct {
	PID        int    `json:"pid"`
	SessionID  string `json:"session_id,omitempty"`
	RunID      string `json:"run_id,omitempty"`
	PublicRoot string `json:"public_root,omitempty"`
	Role       string `json:"role,omitempty"`
}

type peerDiscoverySnapshot struct {
	Type      string              `json:"type"`
	Timestamp int64               `json:"timestamp"`
	Self      peerDiscoveryPeer   `json:"self"`
	Online    []peerDiscoveryPeer `json:"online,omitempty"`
	Joined    []peerDiscoveryPeer `json:"joined,omitempty"`
	Left      []peerDiscoveryPeer `json:"left,omitempty"`
	Summary   []string            `json:"summary,omitempty"`
}

func (r *Runtime) startPeerDiscoveryHeartbeat() {
	if r == nil || r.cfg == nil || r.agentRegistry == nil || !r.cfg.PeerDiscoveryEnabled {
		return
	}
	interval := defaultPeerDiscoveryHeartbeatInterval
	if r.cfg.PeerDiscoveryHeartbeatMS > 0 {
		interval = time.Duration(r.cfg.PeerDiscoveryHeartbeatMS) * time.Millisecond
	}

	r.peerDiscoveryHeartbeatMu.Lock()
	if r.peerDiscoveryHeartbeatStop != nil {
		r.peerDiscoveryHeartbeatMu.Unlock()
		return
	}
	stop := make(chan struct{})
	done := make(chan struct{})
	r.peerDiscoveryHeartbeatStop = stop
	r.peerDiscoveryHeartbeatDone = done
	r.peerDiscoveryHeartbeatMu.Unlock()

	go func() {
		defer close(done)
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for {
			select {
			case <-stop:
				return
			case <-ticker.C:
				if err := r.agentRegistry.PruneStale(); err != nil && r.log != nil {
					r.log("peer discovery heartbeat prune error: %v", err)
				}
			}
		}
	}()
}

func (r *Runtime) stopPeerDiscoveryHeartbeat() {
	if r == nil {
		return
	}
	r.peerDiscoveryHeartbeatMu.Lock()
	stop := r.peerDiscoveryHeartbeatStop
	done := r.peerDiscoveryHeartbeatDone
	r.peerDiscoveryHeartbeatStop = nil
	r.peerDiscoveryHeartbeatDone = nil
	r.peerDiscoveryHeartbeatMu.Unlock()
	if stop == nil {
		return
	}
	close(stop)
	if done == nil {
		return
	}
	select {
	case <-done:
	case <-time.After(time.Second):
		if r.log != nil {
			r.log("peer discovery heartbeat stop timed out")
		}
	}
}

func (r *Runtime) appendPeerDiscoveryStatus(msg *tape.Message) bool {
	if msg == nil || msg.Role != tape.RoleToolResult {
		return false
	}
	snapshot := r.observePeerDiscovery()
	if snapshot == nil {
		return false
	}
	updated, err := setRuntimeField(msg.StructuredContent, "peer_discovery", snapshot)
	if err != nil {
		if r.log != nil {
			r.log("tool result peer-discovery update error: %v", err)
		}
		return false
	}
	msg.StructuredContent = updated
	syncToolResultMessageContent(msg)
	return true
}

func (r *Runtime) observePeerDiscovery() *peerDiscoverySnapshot {
	if r == nil || r.cfg == nil || r.agentRegistry == nil || !r.cfg.PeerDiscoveryEnabled {
		return nil
	}

	current := r.scanPeerDiscoveryPeers()
	selfPID := r.agentRegistry.pid
	self := current[selfPID]
	if self.PID <= 0 {
		self = peerDiscoveryPeer{
			PID:       selfPID,
			SessionID: r.agentRegistry.sessionID,
			RunID:     r.agentRegistry.runID,
			Role:      "myself",
		}
	} else {
		self.Role = "myself"
		current[selfPID] = self
	}

	r.peerDiscoveryMu.Lock()
	defer r.peerDiscoveryMu.Unlock()

	joined := make([]peerDiscoveryPeer, 0)
	left := make([]peerDiscoveryPeer, 0)
	if !r.peerDiscoveryInitialized {
		for pid, peer := range current {
			if pid != selfPID {
				joined = append(joined, peer)
			}
		}
		r.peerDiscoveryInitialized = true
	} else {
		for pid, peer := range current {
			if pid == selfPID {
				continue
			}
			prev, ok := r.peerDiscoveryKnown[pid]
			if !ok || !samePeerDiscoveryPeer(prev, peer) {
				if ok {
					left = append(left, prev)
				}
				joined = append(joined, peer)
			}
		}
		for pid, peer := range r.peerDiscoveryKnown {
			if pid == selfPID {
				continue
			}
			if _, ok := current[pid]; !ok {
				left = append(left, peer)
			}
		}
	}
	r.peerDiscoveryKnown = current

	online := make([]peerDiscoveryPeer, 0, len(current))
	for pid, peer := range current {
		if pid == selfPID {
			continue
		}
		online = append(online, peer)
	}
	sortPeerDiscoveryPeers(online)
	sortPeerDiscoveryPeers(joined)
	sortPeerDiscoveryPeers(left)
	if len(online) == 0 && len(joined) == 0 && len(left) == 0 {
		return nil
	}

	return &peerDiscoverySnapshot{
		Type:      peerDiscoveryMessageType,
		Timestamp: time.Now().UnixMilli(),
		Self:      self,
		Online:    online,
		Joined:    joined,
		Left:      left,
		Summary:   peerDiscoverySummary(self, online, joined, left),
	}
}

func (r *Runtime) scanPeerDiscoveryPeers() map[int]peerDiscoveryPeer {
	peers := map[int]peerDiscoveryPeer{}
	if r == nil || r.agentRegistry == nil || strings.TrimSpace(r.agentRegistry.runtimeRoot) == "" {
		return peers
	}
	pidRoot := filepath.Join(r.agentRegistry.runtimeRoot, "pid")
	entries, err := os.ReadDir(pidRoot)
	if err != nil {
		return peers
	}
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		pid, err := strconv.Atoi(entry.Name())
		if err != nil || pid <= 0 {
			continue
		}
		linkPath := filepath.Join(pidRoot, entry.Name())
		target, err := os.Readlink(linkPath)
		if err != nil {
			continue
		}
		peer := peerDiscoveryPeerFromPIDTarget(pid, linkPath, target)
		if peer.PID > 0 {
			peers[pid] = peer
		}
	}
	return peers
}

func peerDiscoveryPeerFromPIDTarget(pid int, pidLinkPath string, target string) peerDiscoveryPeer {
	peer := peerDiscoveryPeer{PID: pid}
	if strings.TrimSpace(target) == "" {
		return peer
	}
	if !filepath.IsAbs(target) {
		target = filepath.Clean(filepath.Join(filepath.Dir(pidLinkPath), target))
	}
	publicRoot := filepath.Clean(target)
	if filepath.Base(publicRoot) != "public" {
		return peer
	}
	peer.PublicRoot = publicRoot
	agentRoot := filepath.Dir(publicRoot)
	data, err := os.ReadFile(filepath.Join(agentRoot, "status", "session.json"))
	if err != nil {
		return peer
	}
	var status struct {
		SessionID string `json:"session_id"`
		RunID     string `json:"run_id"`
	}
	if err := json.Unmarshal(data, &status); err != nil {
		return peer
	}
	peer.SessionID = status.SessionID
	peer.RunID = status.RunID
	return peer
}

func samePeerDiscoveryPeer(a, b peerDiscoveryPeer) bool {
	return a.PID == b.PID &&
		a.SessionID == b.SessionID &&
		a.RunID == b.RunID &&
		a.PublicRoot == b.PublicRoot
}

func sortPeerDiscoveryPeers(peers []peerDiscoveryPeer) {
	sort.Slice(peers, func(i, j int) bool {
		if peers[i].PID != peers[j].PID {
			return peers[i].PID < peers[j].PID
		}
		if peers[i].SessionID != peers[j].SessionID {
			return peers[i].SessionID < peers[j].SessionID
		}
		return peers[i].RunID < peers[j].RunID
	})
}

func peerDiscoverySummary(self peerDiscoveryPeer, online, joined, left []peerDiscoveryPeer) []string {
	lines := []string{fmt.Sprintf("Pid %d - myself", self.PID)}
	for _, peer := range online {
		lines = append(lines, fmt.Sprintf("Pid %d - online", peer.PID))
	}
	for _, peer := range joined {
		lines = append(lines, fmt.Sprintf("+ Pid %d - joined", peer.PID))
	}
	for _, peer := range left {
		lines = append(lines, fmt.Sprintf("- Pid %d - left", peer.PID))
	}
	return lines
}

func (r *AgentRegistry) tryAcquirePIDLockFile(lockPath string) (*os.File, bool, error) {
	f, err := os.OpenFile(lockPath, os.O_RDWR, agentRegistrySharedFileMode)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, false, nil
		}
		return nil, false, fmt.Errorf("open pid lock file %s: %w", lockPath, err)
	}
	err = syscall.Flock(int(f.Fd()), syscall.LOCK_EX|syscall.LOCK_NB)
	if err == nil {
		return f, true, nil
	}
	_ = f.Close()
	if errors.Is(err, syscall.EWOULDBLOCK) || errors.Is(err, syscall.EAGAIN) {
		return nil, false, nil
	}
	return nil, false, err
}
