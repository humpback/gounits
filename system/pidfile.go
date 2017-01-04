/*
* (C) 2001-2017 humpback Inc.
*
* gounits source code
* version: 1.0.0
* author: bobliu0909@gmail.com
* datetime: 2015-10-14
*
 */

package system

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"syscall"
)

var (
	ErrProcessRunning = errors.New("process is running")
	ErrFileStale      = errors.New("pidfile exists but process is not running")
	ErrFileInvalid    = errors.New("pidfile has invalid contents")
)

func Remove(fname string) error {
	return os.RemoveAll(fname)
}

func Write(fname string) error {
	return WriteControl(fname, os.Getpid(), false)
}

func WriteControl(fname string, pid int, overwrite bool) error {

	oldpid, err := pidfileContents(fname)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	if err == nil {
		if pidIsRunning(oldpid) {
			return ErrProcessRunning
		}
		if !overwrite {
			return ErrFileStale
		}
	}
	return ioutil.WriteFile(fname, []byte(fmt.Sprintf("%d\n", pid)), 0644)
}

func pidfileContents(fname string) (int, error) {
	contents, err := ioutil.ReadFile(fname)
	if err != nil {
		return 0, err
	}

	pid, err := strconv.Atoi(strings.TrimSpace(string(contents)))
	if err != nil {
		return 0, ErrFileInvalid
	}

	return pid, nil
}

func pidIsRunning(pid int) bool {
	process, err := os.FindProcess(pid)
	if err != nil {
		return false
	}

	err = process.Signal(syscall.Signal(0))

	if err != nil && err.Error() == "no such process" {
		return false
	}

	if err != nil && err.Error() == "os: process already finished" {
		return false
	}

	return true
}
