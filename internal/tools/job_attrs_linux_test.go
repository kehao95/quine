//go:build linux

package tools

import (
	"syscall"
	"testing"
)

func TestJobSysProcAttrWorkspaceNamespaceRequestedOnlyWhenNeeded(t *testing.T) {
	direct := jobSysProcAttr(false, false)
	if direct == nil {
		t.Fatal("jobSysProcAttr returned nil")
	}
	if direct.Cloneflags&syscall.CLONE_NEWNS != 0 {
		t.Fatalf("direct workspace jobs must not request CLONE_NEWNS, got %#x", direct.Cloneflags)
	}
	if !direct.Setpgid {
		t.Fatal("non-detached jobs should still setpgid")
	}

	overlay := jobSysProcAttr(false, true)
	if overlay == nil {
		t.Fatal("jobSysProcAttr returned nil for overlay workspace")
	}
	if overlay.Cloneflags&syscall.CLONE_NEWNS == 0 {
		t.Fatalf("overlay workspace jobs must request CLONE_NEWNS, got %#x", overlay.Cloneflags)
	}
}
