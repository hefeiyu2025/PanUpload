package internal

import (
	"io"
	"os"
)

// IsEmptyDir 判断目录是否为空
func IsEmptyDir(dirPath string) (bool, error) {
	dir, err := os.Open(dirPath)
	if err != nil {
		return false, err
	}
	defer dir.Close()
	//如果目录不为空，Readdirnames 会返回至少一个文件名
	_, err = dir.Readdirnames(1)
	if err == io.EOF {
		return true, nil
	}
	return false, err
}
