package oss

import "io"

// Driver oss驱动接口定义
type Driver interface {
	// Put 上传
	Put(objectName, localFileName string) error
	// PutObj 上传
	PutObj(objectName string, file io.Reader) error
	// Get 下载
	Get(objectName, downloadedFileName string) error
	// Del 删除
	Del(objectName string) error
}
