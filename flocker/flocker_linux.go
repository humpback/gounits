package flocker

import (
	"os"
	"syscall"
	"time"
)

func LockFile(file *os.File, wait time.Duration) error {

	for {
		err := syscall.Flock(int(file.Fd()), syscall.LOCK_EX|syscall.LOCK_NB)
		if err == nil {
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
