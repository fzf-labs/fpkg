package file

import (
	"fmt"
	"os"
)

// FileExists 判断文件是否存在
func FileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

// MkdirPath 生成文件夹
func MkdirPath(relativePath string) error {
	if err := os.MkdirAll(relativePath, os.ModePerm); err != nil {
		return fmt.Errorf("create model pkg path(%s) fail: %s", relativePath, err)
	}
	return nil
}
