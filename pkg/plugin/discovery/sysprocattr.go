// +build !linux

package discovery

import "syscall"

func GetSysProcAttr() *syscall.SysProcAttr {
	return nil
}
