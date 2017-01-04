/*
* (C) 2001-2017 humpback Inc.
*
* gounits source code
* version: 1.0.0
* author: bobliu0909@gmail.com
* datetime: 2015-10-14
* 生成规则：uuid
* 生成一个KeyFile，如果文件不存在则自动创建，否则直接读取.
 */

package system

import "github.com/google/uuid"

import (
	"io/ioutil"
	"os"
)

func MakeKey(random bool) string {

	if !random {
		key, err := uuid.NewUUID() //按时间戳生成uuid
		if err != nil {
			return ""
		}
		return key.String()
	}
	return uuid.New().String() //随机生成uuid
}

func MakeKeyFile(fpath string) (string, error) {

	_, err := os.Stat(fpath)
	if err != nil && !os.IsExist(err) {
		key := uuid.New().String()
		if err := ioutil.WriteFile(fpath, []byte(key), 0777); err != nil {
			return "", err
		}
		return key, nil
	}

	fp, err := os.Open(fpath)
	if err != nil {
		return "", err
	}

	defer fp.Close()
	data, err := ioutil.ReadAll(fp)
	if err != nil {
		return "", err
	}

	key := string(data)
	if _, err := uuid.Parse(key); err != nil {
		return "", err
	}
	return key, nil
}
