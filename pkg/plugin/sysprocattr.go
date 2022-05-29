//go:build !linux
// +build !linux

package plugin

import "syscall"

func GetSysProcAttr() *syscall.SysProcAttr {
	return nil
}
