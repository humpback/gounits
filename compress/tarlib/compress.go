/*
* (C) 2001-2017 humpback Inc.
*
* gounits source code
* version: 1.0.0
* author: bobliu0909@gmail.com
* datetime: 2016-06-14
* tar压缩与解压缩
 */

package tarlib

import (
	"archive/tar"
	"compress/gzip"
	"errors"
	"io"
	"io/ioutil"
	"os"
	"path"
)

func EnCompress(src string, dest string, format StreamFormat, failIfExist bool) (err error) {

	// 清理路径字符串
	src = path.Clean(src)
	// 判断要打包的文件或目录是否存在
	if !Exists(src) {
		return errors.New("source file not found：" + src)
	}

	// 判断目标文件是否存在
	if FileExists(dest) {
		if failIfExist { // 不覆盖已存在的文件
			return errors.New("dest file is exists：" + dest)
		} else { // 覆盖已存在的文件
			if err := os.Remove(dest); err != nil {
				return err
			}
		}
	}

	// 创建空的目标文件
	fw, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer fw.Close()

	var writer *tar.Writer
	if format == STREAM_TAR_GZIP {
		gwriter := gzip.NewWriter(fw)
		writer = tar.NewWriter(gwriter)
	} else if format == STREAM_TAR_SIMPLE {
		writer = tar.NewWriter(fw)
	} else {
		return errors.New("format invalid：" + format.String())
	}

	defer func() {
		// 这里要判断 writer 是否关闭成功，如果关闭失败，则 .tar 文件可能不完整
		if er := writer.Close(); er != nil {
			err = er
		}
	}()

	// 获取文件或目录信息
	fi, err := os.Stat(src)
	if err != nil {
		return err
	}

	// 获取要打包的文件或目录的所在位置和名称
	srcBase, srcRelative := path.Split(path.Clean(src))
	// 开始打包
	if fi.IsDir() {
		if err := compressDirectory(srcBase, srcRelative, writer, fi); err != nil {
			return err
		}
	} else {
		if err := compressFile(srcBase, srcRelative, writer, fi); err != nil {
			return err
		}
	}
	return nil
}

// 因为要执行遍历操作，所以要单独创建一个函数
func compressDirectory(srcBase, srcRelative string, w *tar.Writer, fi os.FileInfo) (err error) {

	// 获取完整路径
	srcFull := srcBase + srcRelative
	// 在结尾添加 "/"
	last := len(srcRelative) - 1
	if srcRelative[last] != os.PathSeparator {
		srcRelative += string(os.PathSeparator)
	}

	// 获取 srcFull 下的文件或子目录列表
	fis, err := ioutil.ReadDir(srcFull)
	if err != nil {
		return err
	}

	// 开始遍历
	for _, fi := range fis {
		if fi.IsDir() {
			if err := compressDirectory(srcBase, srcRelative+fi.Name(), w, fi); err != nil {
				return err
			}
		} else {
			if err := compressFile(srcBase, srcRelative+fi.Name(), w, fi); err != nil {
				return err
			}
		}
	}

	// 写入目录信息
	if len(srcRelative) > 0 {
		hdr, err := tar.FileInfoHeader(fi, "")
		if err != nil {
			return err
		}

		hdr.Name = srcRelative
		if err = w.WriteHeader(hdr); err != nil {
			return err
		}
	}
	return nil
}

// 因为要在 defer 中关闭文件，所以要单独创建一个函数
func compressFile(srcBase, srcRelative string, w *tar.Writer, fi os.FileInfo) (err error) {

	// 获取完整路径
	srcFull := srcBase + srcRelative
	// 写入文件信息
	hdr, err := tar.FileInfoHeader(fi, "")
	if err != nil {
		return err
	}

	hdr.Name = srcRelative
	if err = w.WriteHeader(hdr); err != nil {
		return err
	}

	// 打开要打包的文件，准备读取
	fr, err := os.Open(srcFull)
	if err != nil {
		return err
	}
	defer fr.Close()

	// 将文件数据写入 w 中
	if _, err = io.Copy(w, fr); err != nil {
		return err
	}
	return nil
}

func DeCompress(src string, dest string, format StreamFormat) (err error) {

	// 清理路径字符串
	dest = path.Clean(dest) + string(os.PathSeparator)
	if !DirExists(dest) {
		if err := os.MkdirAll(dest, 0777); err != nil {
			return err
		}
	}

	// 打开要解包的文件
	fr, err := os.Open(src)
	if err != nil {
		return err
	}
	defer fr.Close()

	var reader *tar.Reader
	if format == STREAM_TAR_GZIP {
		greader, err := gzip.NewReader(fr)
		if err != nil {
			return err
		}
		reader = tar.NewReader(greader)
	} else if format == STREAM_TAR_SIMPLE {
		reader = tar.NewReader(fr)
	} else {
		return errors.New("format invalid：" + format.String())
	}

	// 遍历包中的文件
	for hdr, err := reader.Next(); err != io.EOF; hdr, err = reader.Next() {
		if err != nil {
			return err
		}
		// 获取文件信息
		fi := hdr.FileInfo()
		// 获取绝对路径
		dstFullPath := dest + hdr.Name
		if hdr.Typeflag == tar.TypeDir {
			// 创建目录
			os.MkdirAll(dstFullPath, fi.Mode().Perm())
			// 设置目录权限
			os.Chmod(dstFullPath, fi.Mode().Perm())
		} else {
			// 创建文件所在的目录
			os.MkdirAll(path.Dir(dstFullPath), os.ModePerm)
			// 将 reader 中的数据写入文件中
			if er := uncompressFile(dstFullPath, reader); er != nil {
				return er
			}
			// 设置文件权限
			os.Chmod(dstFullPath, fi.Mode().Perm())
		}
	}
	return nil
}

// 因为要在 defer 中关闭文件，所以要单独创建一个函数
func uncompressFile(destFile string, r *tar.Reader) error {

	// 创建空文件，准备写入解包后的数据
	fw, err := os.Create(destFile)
	if err != nil {
		return err
	}
	defer fw.Close()

	// 写入解包后的数据
	_, err = io.Copy(fw, r)
	if err != nil {
		return err
	}
	return nil
}

// 判断档案是否存在
func Exists(name string) bool {
	_, err := os.Stat(name)
	return err == nil || os.IsExist(err)
}

// 判断文件是否存在
func FileExists(filename string) bool {
	fi, err := os.Stat(filename)
	return (err == nil || os.IsExist(err)) && !fi.IsDir()
}

// 判断目录是否存在
func DirExists(dirname string) bool {
	fi, err := os.Stat(dirname)
	return (err == nil || os.IsExist(err)) && fi.IsDir()
}
