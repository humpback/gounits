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
	"bytes"
	"errors"
	"os"
	"strconv"
	"time"
)

var (
	ERR_FILEINVALID       = errors.New("file invalid.")
	ERR_FileLocked        = errors.New("file locked.")
	ERR_FileLockExecption = errors.New("file lock execption.")
)

type FileLocker struct {
	FileName string
	Wait     time.Duration
	fd       *os.File
}

func NewFileLocker(fname string, wait time.Duration) *FileLocker {

	return &FileLocker{
		FileName: fname,
		Wait:     wait,
		fd:       nil,
	}
}

func (lock *FileLocker) Lock() error {

	fd, err := os.OpenFile(lock.FileName, os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0777)
	if err != nil {
		return ERR_FILEINVALID
	}

	if err := LockFile(fd, lock.Wait); err != nil {
		fd.Close()
		return err
	}

	pid := os.Getpid()
	buff := bytes.NewBuffer([]byte(strconv.Itoa(pid)))
	if _, err := fd.Write(buff.Bytes()); err == nil {
		fd.Sync()
	}
	return nil
}

func (lock *FileLocker) Unlock() {

	if lock.fd != nil {
		lock.fd.Close()
	}
}
