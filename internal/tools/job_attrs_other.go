//go:build !linux

package tools

import "syscall"

func jobSysProcAttr(detach, _ bool) *syscall.SysProcAttr {
	if detach {
		return &syscall.SysProcAttr{Setsid: true}
	}
	return &syscall.SysProcAttr{Setpgid: true}
}
