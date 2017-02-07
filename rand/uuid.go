package rand

import "github.com/google/uuid"

import (
	"io/ioutil"
	"os"
)

// NewUUID returns a uuid string
// format:xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
func UUID(random bool) string {

	if !random {
		key, err := uuid.NewUUID() //按时间戳生成uuid
		if err != nil {
			return ""
		}
		return key.String()
	}
	return uuid.New().String() //随机生成uuid
}

// WriteUUIDToFile return a uuid string and write to file
func UUIDFile(fpath string) (string, error) {

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
