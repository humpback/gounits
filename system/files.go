package system

import (
	"io"
	"os"
)

func FileExist(filename string) bool {

	fi, err := os.Stat(filename)
	return (err == nil || os.IsExist(err)) && !fi.IsDir()
}

func DirExist(dirname string) bool {

	fi, err := os.Stat(dirname)
	return (err == nil || os.IsExist(err)) && fi.IsDir()
}

func FileCopy(source string, dest string) (int64, error) {

	sourcefile, err := os.Open(source)
	if err != nil {
		return 0, err
	}

	defer sourcefile.Close()
	sourceinfo, err := os.Stat(source)
	if err != nil {
		return 0, err
	}

	destfile, err := os.Create(dest)
	if err != nil {
		return 0, err
	}

	w, err := io.Copy(destfile, sourcefile)
	if err != nil {
		destfile.Close()
		return 0, err
	}

	if err := os.Chmod(dest, sourceinfo.Mode()); err != nil {
		destfile.Close()
		return 0, err
	}
	destfile.Close()
	return w, nil
}

func DirectoryCopy(source string, dest string) error {

	sourceinfo, err := os.Stat(source)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(dest, sourceinfo.Mode()); err != nil {
		return err
	}

	directory, err := os.Open(source)
	if err != nil {
		return err
	}

	defer directory.Close()
	objects, err := directory.Readdir(-1)
	if err != nil {
		return err
	}

	for _, obj := range objects {
		sourcefilepointer := source + "/" + obj.Name()
		destinationfilepointer := dest + "/" + obj.Name()
		if obj.IsDir() {
			if err := DirectoryCopy(sourcefilepointer, destinationfilepointer); err != nil {
				return err
			}
		} else {
			if _, err := FileCopy(sourcefilepointer, destinationfilepointer); err != nil {
				return err
			}
		}
	}
	return nil
}
