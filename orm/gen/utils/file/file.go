package file

import (
	"os"
)

// Exists 判断文件是否存在
func Exists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

// MkdirPath 生成文件夹
func MkdirPath(relativePath string) error {
	return os.MkdirAll(relativePath, os.ModePerm)
}
