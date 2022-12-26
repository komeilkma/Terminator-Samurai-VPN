// +build darwin,go1.11

package native-water

import "syscall"

func setNonBlock(fd int) error {
	return syscall.SetNonblock(fd, true)
}
