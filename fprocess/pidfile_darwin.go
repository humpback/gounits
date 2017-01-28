package fprocess

import (
	"syscall"
)

func processExists(pid int) bool {

	err := syscall.Kill(pid, 0)
	if err != nil {
		return false
	}
	return true
}
