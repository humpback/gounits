package system

import (
	"crypto/md5"
	"fmt"
	"io"
	"os"
)

//ReadFileMD5Code is exported
//read file real md5 code string.
func ReadFileMD5Code(filepath string) (string, error) {

	f, err := os.Open(filepath)
	if err != nil {
		return "", err
	}

	defer f.Close()
	md5hash := md5.New()
	if _, err := io.Copy(md5hash, f); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", md5hash.Sum(nil)), nil
}
