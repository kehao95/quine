//go:build linux

package tools

import (
	"os"
	"syscall"
)

func jobSysProcAttr(detach, useWorkspace bool) *syscall.SysProcAttr {
	attr := &syscall.SysProcAttr{}
	if detach {
		attr.Setsid = true
	} else {
		attr.Setpgid = true
	}
	if useWorkspace {
		attr.Cloneflags = syscall.CLONE_NEWNS
		uid := os.Getuid()
		gid := os.Getgid()
		// Root can mount overlayfs directly inside a fresh mount namespace.
		// Non-root still needs a user namespace to gain mount capability.
		if uid != 0 {
			attr.Cloneflags |= syscall.CLONE_NEWUSER
			attr.UidMappings = []syscall.SysProcIDMap{
				{ContainerID: 0, HostID: uid, Size: 1},
			}
			attr.GidMappingsEnableSetgroups = false
			attr.GidMappings = []syscall.SysProcIDMap{
				{ContainerID: 0, HostID: gid, Size: 1},
			}
		}
	}
	return attr
}
