/*
* (C) 2001-2017 humpback Inc.
*
* gounits source code
* version: 1.0.0
* author: bobliu0909@gmail.com
* datetime: 2016-05-13
*
 */

package flocker

import (
	"os"
	"syscall"
	"time"
)

func LockFile(file *os.File, wait time.Duration) error {

	h, err := syscall.LoadLibrary("kernel32.dll")
	if err != nil {
		return ERR_FileLockExecption
	}
	defer syscall.FreeLibrary(h)

	addr, err := syscall.GetProcAddress(h, "LockFile")
	if err != nil {
		return ERR_FileLockExecption
	}

	for {
		r0, _, _ := syscall.Syscall6(addr, 5, file.Fd(), 0, 0, 0, 1, 0)
		if 0 != int(r0) {
			break
		}

		if wait == 0 {
			return ERR_FileLocked
		} else {
			c := time.After(wait)
			select {
			case <-c:
				wait = time.Duration(0)
			}
		}
	}
	return nil
}
