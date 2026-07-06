package qcli

import (
	"os/exec"
	"testing"
	"time"
)

// A zombie passes signal-0, so pidLive/verifyLivePID must consult the /proc
// state field. Regression for the dogfood incident where an exited agent
// (unreaped child of the exec'd qcli) read as ●live forever.
func TestPidLiveTreatsZombieAsDead(t *testing.T) {
	cmd := exec.Command("true")
	if err := cmd.Start(); err != nil {
		t.Fatalf("start child: %v", err)
	}
	pid := cmd.Process.Pid
	defer cmd.Wait() // reap at the end so the zombie does not outlive the test

	deadline := time.Now().Add(2 * time.Second)
	for !pidZombie(pid) {
		if time.Now().After(deadline) {
			t.Fatalf("pid %d never became a zombie", pid)
		}
		time.Sleep(10 * time.Millisecond)
	}

	if pidLive(pid) {
		t.Fatalf("pidLive(%d) = true for a defunct process", pid)
	}
	if err := verifyLivePID(pid); err == nil {
		t.Fatalf("verifyLivePID(%d) = nil for a defunct process", pid)
	}
}

func TestPidLiveOnLiveAndDeadPIDs(t *testing.T) {
	cmd := exec.Command("sleep", "5")
	if err := cmd.Start(); err != nil {
		t.Fatalf("start child: %v", err)
	}
	pid := cmd.Process.Pid
	if !pidLive(pid) {
		t.Fatalf("pidLive(%d) = false for a running process", pid)
	}
	if err := cmd.Process.Kill(); err != nil {
		t.Fatalf("kill child: %v", err)
	}
	if err := cmd.Wait(); err == nil {
		t.Fatalf("expected non-nil Wait error after kill")
	}
	// Reaped: the pid is fully gone (or recycled — signal-0 to a recycled pid
	// is out of scope for this test, so only assert the common path).
	if pidLive(pid) {
		t.Fatalf("pidLive(%d) = true after the process was reaped", pid)
	}
}
