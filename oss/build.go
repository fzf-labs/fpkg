package oss

import (
	"strings"
	"time"

	"github.com/segmentio/ksuid"
)

func BuildNewNameAndPath(ext string, category ...string) (newFileName, filePath string) {
	// 日期
	date := time.Now().Format("20060102")
	join := make([]string, 0)
	join = append(join, category...)
	join = append(join, date)
	dir := strings.Trim(strings.Join(join, "/"), "/")
	newFileName = ksuid.New().String() + ext
	filePath = dir + "/" + newFileName
	return
}
