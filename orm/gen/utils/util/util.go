package util

import (
	"golang.org/x/tools/go/packages"
)

// FillModelPkgPath 返回模型文件的包路径
func FillModelPkgPath(dir string) string {
	pkg, err := packages.Load(&packages.Config{
		Mode: packages.NeedName,
		Dir:  dir,
	})
	if err != nil {
		return ""
	}
	if len(pkg) > 0 {
		return pkg[0].PkgPath
	}
	return ""
}

// StrSliFind 判断字符串切片中是否存在某个元素
func StrSliFind(collection []string, element string) bool {
	for _, s := range collection {
		if s == element {
			return true
		}
	}
	return false
}

// SliRemove 删除字符串切片中的某个元素
func SliRemove(collection, element []string) []string {
	for _, s := range element {
		for i, v := range collection {
			if s == v {
				collection = append(collection[:i], collection[i+1:]...)
			}
		}
	}
	return collection
}
