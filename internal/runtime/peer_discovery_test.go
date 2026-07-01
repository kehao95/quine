package runtime

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestPeerDiscoveryHeartbeatPrunesStalePIDLock(t *testing.T) {
	cfg := testCfg(t)
	cfg.SessionID = "heartbeat-live"
	cfg.RunID = "heartbeat-run"
	cfg.PeerDiscoveryEnabled = true
	cfg.PeerDiscoveryHeartbeatMS = 20

	rt := NewWithProvider(cfg, &mockProvider{})
	silenceRuntime(rt)
	defer rt.stopPeerDiscoveryHeartbeat()
	if err := rt.agentRegistry.Register(); err != nil {
		t.Fatalf("Register failed: %v", err)
	}
	defer rt.agentRegistry.Deregister()
	if err := rt.agentRegistry.PublishSelfPID(); err != nil {
		t.Fatalf("PublishSelfPID failed: %v", err)
	}
	defer rt.agentRegistry.UnpublishSelfPID()
	rt.startPeerDiscoveryHeartbeat()

	root := cfg.RuntimeRoot()
	const stalePID = 99997
	stalePIDText := fmt.Sprint(stalePID)
	staleSession := "heartbeat-stale"
	staleRun := "heartbeat-stale-run"
	stalePublic := filepath.Join(root, "agent", staleSession, "public")
	staleStatus := filepath.Join(root, "agent", staleSession, "status", "session.json")
	if err := os.MkdirAll(filepath.Dir(staleStatus), 0o755); err != nil {
		t.Fatalf("mkdir stale status dir: %v", err)
	}
	if err := os.WriteFile(staleStatus, []byte(`{"session_id":"`+staleSession+`","run_id":"`+staleRun+`","pid":`+stalePIDText+`}`), 0o644); err != nil {
		t.Fatalf("write stale status: %v", err)
	}
	if err := os.MkdirAll(stalePublic, 0o755); err != nil {
		t.Fatalf("mkdir stale public: %v", err)
	}
	if err := replaceSymlink(filepath.Join(root, "pid", stalePIDText), stalePublic); err != nil {
		t.Fatalf("seed stale pid route: %v", err)
	}
	staleRegPath := filepath.Join(cfg.LockDir(), staleRun+".agent")
	staleReg, err := json.Marshal(agentRegistration{SessionID: staleSession, RunID: staleRun, PID: stalePID})
	if err != nil {
		t.Fatalf("marshal stale registration: %v", err)
	}
	if err := os.WriteFile(staleRegPath, append(staleReg, '\n'), 0o644); err != nil {
		t.Fatalf("seed stale registration: %v", err)
	}
	pidLockPath := filepath.Join(cfg.LockDir(), "agents", fmt.Sprintf("%d.agent.lock", stalePID))
	if err := os.MkdirAll(filepath.Dir(pidLockPath), 0o755); err != nil {
		t.Fatalf("mkdir stale lock dir: %v", err)
	}
	if err := os.WriteFile(pidLockPath, nil, 0o644); err != nil {
		t.Fatalf("seed stale pid lock: %v", err)
	}

	waitForCondition(t, 2*time.Second, func() bool {
		_, lockErr := os.Lstat(pidLockPath)
		_, pidErr := os.Lstat(filepath.Join(root, "pid", stalePIDText))
		_, regErr := os.Lstat(staleRegPath)
		return os.IsNotExist(lockErr) && os.IsNotExist(pidErr) && os.IsNotExist(regErr)
	}, "heartbeat stale pid cleanup")

	if _, err := os.Lstat(filepath.Join(root, "agent", staleSession)); !os.IsNotExist(err) {
		t.Fatalf("expected stale agent root removed, got err=%v", err)
	}
}

func TestPeerDiscoverySafePointReportsLocalTopologyDiff(t *testing.T) {
	cfg := testCfg(t)
	cfg.SessionID = "observer-live"
	cfg.RunID = "observer-run"
	cfg.PeerDiscoveryEnabled = true

	rt := NewWithProvider(cfg, &mockProvider{})
	silenceRuntime(rt)
	if err := rt.agentRegistry.Register(); err != nil {
		t.Fatalf("Register failed: %v", err)
	}
	defer rt.agentRegistry.Deregister()
	if err := rt.agentRegistry.PublishSelfPID(); err != nil {
		t.Fatalf("PublishSelfPID failed: %v", err)
	}
	defer rt.agentRegistry.UnpublishSelfPID()

	root := cfg.RuntimeRoot()
	const peerPID = 424245
	peerSession := "peer-session"
	peerRun := "peer-run"
	peerPublic := filepath.Join(root, "agent", peerSession, "public")
	peerStatus := filepath.Join(root, "agent", peerSession, "status", "session.json")
	if err := os.MkdirAll(peerPublic, 0o755); err != nil {
		t.Fatalf("mkdir peer public: %v", err)
	}
	if err := os.MkdirAll(filepath.Dir(peerStatus), 0o755); err != nil {
		t.Fatalf("mkdir peer status: %v", err)
	}
	if err := os.WriteFile(peerStatus, []byte(`{"session_id":"`+peerSession+`","run_id":"`+peerRun+`","pid":424245}`), 0o644); err != nil {
		t.Fatalf("write peer status: %v", err)
	}
	if err := replaceSymlink(filepath.Join(root, "pid", fmt.Sprint(peerPID)), peerPublic); err != nil {
		t.Fatalf("seed peer pid route: %v", err)
	}

	first := runtimeToolResultMessage("tool-1", "sh", "completed", map[string]any{"stdout": "ok"})
	if !rt.appendPeerDiscoveryStatus(&first) {
		t.Fatalf("expected initial peer discovery snapshot")
	}
	firstRuntime := toolMap(t, decodeToolContent(t, first.StructuredContent), "runtime")
	firstDiscovery := toolMap(t, firstRuntime, "peer_discovery")
	if firstDiscovery["type"] != peerDiscoveryMessageType {
		t.Fatalf("peer discovery type = %#v", firstDiscovery["type"])
	}
	joined, _ := firstDiscovery["joined"].([]any)
	if len(joined) != 1 {
		t.Fatalf("expected one joined peer, got %#v", firstDiscovery["joined"])
	}
	joinedPeer, _ := joined[0].(map[string]any)
	if toolInt(t, joinedPeer, "pid") != peerPID || joinedPeer["session_id"] != peerSession || joinedPeer["run_id"] != peerRun {
		t.Fatalf("joined peer = %#v", joinedPeer)
	}

	if err := os.Remove(filepath.Join(root, "pid", fmt.Sprint(peerPID))); err != nil {
		t.Fatalf("remove peer pid route: %v", err)
	}
	second := runtimeToolResultMessage("tool-2", "sh", "completed", map[string]any{"stdout": "ok"})
	if !rt.appendPeerDiscoveryStatus(&second) {
		t.Fatalf("expected peer discovery leave snapshot")
	}
	secondRuntime := toolMap(t, decodeToolContent(t, second.StructuredContent), "runtime")
	secondDiscovery := toolMap(t, secondRuntime, "peer_discovery")
	left, _ := secondDiscovery["left"].([]any)
	if len(left) != 1 {
		t.Fatalf("expected one left peer, got %#v", secondDiscovery["left"])
	}
	leftPeer, _ := left[0].(map[string]any)
	if toolInt(t, leftPeer, "pid") != peerPID || leftPeer["session_id"] != peerSession || leftPeer["run_id"] != peerRun {
		t.Fatalf("left peer = %#v", leftPeer)
	}
}
