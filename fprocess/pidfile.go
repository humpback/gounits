package fprocess

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// PIDFile is a file used to store the process ID of a running process.
type PIDFile struct {
	PID  int
	path string
}

func checkPIDFileAlreadyExists(path string) error {

	if buf, err := ioutil.ReadFile(path); err == nil {
		pidstr := strings.TrimSpace(string(buf))
		if pid, err := strconv.Atoi(pidstr); err == nil {
			if processExists(pid) {
				return fmt.Errorf("pid file found, ensure process is not running or delete %s", path)
			}
		}
	}
	return nil
}

func New(path string) (*PIDFile, error) {

	if err := checkPIDFileAlreadyExists(path); err != nil {
		return nil, err
	}

	if err := os.MkdirAll(filepath.Dir(path), os.FileMode(0755)); err != nil {
		return nil, err
	}

	pid := os.Getpid()
	if err := ioutil.WriteFile(path, []byte(fmt.Sprintf("%d", pid)), 0644); err != nil {
		return nil, err
	}

	return &PIDFile{
		PID:  pid,
		path: path,
	}, nil
}

func (pf *PIDFile) Remove() error {

	if err := os.Remove(pf.path); err != nil {
		return err
	}
	return nil
}
