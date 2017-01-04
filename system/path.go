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
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
)

func GetExecAbsolutePath() (string, error) {

	lp, err := exec.LookPath(os.Args[0])
	if err != nil {
		return "", err
	}

	path, err := filepath.Abs(lp)
	if err != nil {
		return "", err
	}
	return path, nil
}

func GetExecDir() (string, error) {

	lp, err := exec.LookPath(os.Args[0])
	if err != nil {
		return "", err
	}

	path, err := filepath.Abs(lp)
	if err != nil {
		return "", err
	}
	return filepath.Dir(path), nil
}

func GetExecExt() (string, error) {

	lp, err := exec.LookPath(os.Args[0])
	if err != nil {
		return "", err
	}
	return filepath.Ext(lp), nil
}

func PathExists(path string) (bool, error) {

	p, err := filepath.Abs(path)
	if err != nil {
		return false, err
	}

	if _, err := os.Stat(p); err == nil {
		return true, nil
	}

	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func MakeDirectory(path string) error {

	ret, err := PathExists(path)
	if err != nil {
		return err
	}

	if !ret {
		if err := os.MkdirAll(path, 0777); err != nil {
			return err
		}
	}
	return nil
}

func EmptyDirectory(path string) (bool, error) {

	dirpath, err := filepath.Abs(path)
	if err != nil {
		return false, err
	}

	dir, err := ioutil.ReadDir(dirpath)
	if err != nil {
		return false, err
	}

	if len(dir) == 0 {
		return true, nil
	}
	return false, nil
}
