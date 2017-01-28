package fprocess

import (
	"os"
	"path/filepath"
	"strconv"
)

func processExists(pid int) bool {

	if _, err := os.Stat(filepath.Join("/proc", strconv.Itoa(pid))); err != nil {
		return false
	}
	return true
}
