// +build darwin,go1.11

package nativewater

import "syscall"

func setNonBlock(fd int) error {
	return syscall.SetNonblock(fd, true)
}
