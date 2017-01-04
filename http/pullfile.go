package http

import (
	"os"
	"path/filepath"
)

func HttpPullFile(savefile string, remoteurl string, resumable bool) error {

	fpath, err := filepath.Abs(savefile)
	if err != nil {
		return ErrFilePathInvalid
	}

	fd, err := os.OpenFile(fpath, os.O_CREATE|os.O_RDWR, 0755)
	if err != nil {
		return ErrFileOpenException
	}

	httpclient := NewHttpClient(nil)
	if err := httpclient.GetFile(fd, remoteurl, resumable); err != nil {
		fd.Close()
		os.Remove(fpath)
		return err
	}
	fd.Close()
	return nil
}
