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

type StreamFormat int

const (
	STREAM_TAR_SIMPLE StreamFormat = iota + 1 //tar格式
	STREAM_TAR_GZIP                           //tar.gz格式
)

func (fotmat StreamFormat) String() string {

	switch fotmat {
	case STREAM_TAR_SIMPLE:
		return "STREAM_TAR_SIMPLE"
	case STREAM_TAR_GZIP:
		return "STREAM_TAR_GZIP"
	}
	return ""
}

/*
 * TarGzipEnCompress
 * 将文件或目录打包成 .tar.gz 文件
 * src 是要打包的文件或目录的路径
 * dest 是要生成的 .tar.gz 文件的路径
 * failIfExist 标记如果 dest 文件存在，是否放弃打包，如果否，则会覆盖已存在的文件
 */
func TarGzipEnCompress(src string, dest string, failIfExist bool) (err error) {

	return EnCompress(src, dest, STREAM_TAR_GZIP, failIfExist)
}

/*
 * TarGzipDeCompress
 * 将.tar.gz格式文件解压
 * src 文件路径
 * dest 目标目录
 */
func TarGzipDeCompress(src string, dest string) (err error) {

	return DeCompress(src, dest, STREAM_TAR_GZIP)
}

/*
 * TarEnCompress
 * 将文件或目录打包成 .tar 文件
 * src 是要打包的文件或目录的路径
 * dest 是要生成的 .tar 文件的路径
 * failIfExist 标记如果 dest 文件存在，是否放弃打包，如果否，则会覆盖已存在的文件
 */
func TarEnCompress(src string, dest string, failIfExist bool) (err error) {

	return EnCompress(src, dest, STREAM_TAR_SIMPLE, failIfExist)
}

/*
 * TarDeCompress
 * 将.tar格式文件解压
 * src 文件路径
 * dest 目标目录
 */
func TarDeCompress(src string, dest string) (err error) {

	return DeCompress(src, dest, STREAM_TAR_SIMPLE)
}

/*
 * TarAutoDeCompress
 * 自动解压tar格式文件
 * src 是要打包的文件或目录的路径
 * dest 是要生成的 .tar 或 .tar.gz 文件的路径
 * 首先尝试以tar.gz解压格式解压，如果失败则尝试tar格式
 * 函数目的是一次调用动作支持两种格式
 */
func TarAutoDeCompress(src string, dest string) (err error) {

	if err := TarGzipDeCompress(src, dest); err != nil {
		if err := TarDeCompress(src, dest); err != nil {
			return err
		}
	}
	return nil
}
