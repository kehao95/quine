//go:build linux

package tools

import (
	"os"
	"syscall"
	"testing"
)

func TestJobSysProcAttrWorkspaceNamespaceRequestedOnlyWhenNeeded(t *testing.T) {
	host := jobSysProcAttr(false, false, false)
	if host == nil {
		t.Fatal("jobSysProcAttr returned nil")
	}
	if host.Cloneflags&syscall.CLONE_NEWNS != 0 {
		t.Fatalf("non-workspace jobs must not request CLONE_NEWNS, got %#x", host.Cloneflags)
	}
	if !host.Setpgid {
		t.Fatal("non-detached jobs should still setpgid")
	}

	overlay := jobSysProcAttr(false, true, false)
	if overlay == nil {
		t.Fatal("jobSysProcAttr returned nil for overlay workspace")
	}
	if overlay.Cloneflags&syscall.CLONE_NEWNS == 0 {
		t.Fatalf("overlay workspace jobs must request CLONE_NEWNS, got %#x", overlay.Cloneflags)
	}
}

func TestJobSysProcAttrNonRootWorkspaceUsesUserNamespace(t *testing.T) {
	if os.Getuid() == 0 {
		t.Skip("root workspace jobs can mount directly without user namespace mapping")
	}
	attr := jobSysProcAttr(false, true, false)
	if attr.Cloneflags&syscall.CLONE_NEWUSER == 0 {
		t.Fatalf("non-root workspace jobs must request CLONE_NEWUSER, got %#x", attr.Cloneflags)
	}
	if len(attr.UidMappings) != 1 || attr.UidMappings[0].HostID != os.Getuid() || attr.UidMappings[0].ContainerID != 0 {
		t.Fatalf("unexpected uid mapping: %#v", attr.UidMappings)
	}
	if len(attr.GidMappings) != 1 || attr.GidMappings[0].HostID != os.Getgid() || attr.GidMappings[0].ContainerID != 0 {
		t.Fatalf("unexpected gid mapping: %#v", attr.GidMappings)
	}
	if attr.GidMappingsEnableSetgroups {
		t.Fatal("non-root user namespace should disable setgroups before writing gid_map")
	}
}

func TestJobSysProcAttrDetachedJobsStillCreateSession(t *testing.T) {
	attr := jobSysProcAttr(true, true, false)
	if attr == nil {
		t.Fatal("jobSysProcAttr returned nil")
	}
	if !attr.Setsid {
		t.Fatal("detached jobs should request setsid")
	}
}

func TestJobSysProcAttrNetworkNamespaceRequestedOnlyWhenNeeded(t *testing.T) {
	host := jobSysProcAttr(false, false, false)
	if host.Cloneflags&syscall.CLONE_NEWNET != 0 {
		t.Fatalf("host-network jobs must not request CLONE_NEWNET, got %#x", host.Cloneflags)
	}

	isolated := jobSysProcAttr(false, false, true)
	if isolated.Cloneflags&syscall.CLONE_NEWNET == 0 {
		t.Fatalf("network-isolated jobs must request CLONE_NEWNET, got %#x", isolated.Cloneflags)
	}
}
