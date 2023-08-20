package ziputil

import (
	azip "archive/zip"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/klauspost/compress/zip"
)

func ZipFiles(filename string, files []string, oldForm, newForm string) error {
	newZipFile, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer func() {
		_ = newZipFile.Close()
	}()

	zipWriter := azip.NewWriter(newZipFile)
	defer func() {
		_ = zipWriter.Close()
	}()
	// 把files添加到zip中
	for _, file := range files {
		err = func(file string) error {
			zipFile, err2 := os.Open(file)
			if err2 != nil {
				return err2
			}
			defer zipFile.Close()
			// 获取file的基础信息
			info, err2 := zipFile.Stat()
			if err2 != nil {
				return err2
			}

			header, err2 := azip.FileInfoHeader(info)
			if err2 != nil {
				return err2
			}

			// 使用上面的FileInfoHeader() 就可以把文件保存的路径替换成我们自己想要的了，如下面
			header.Name = strings.ReplaceAll(file, oldForm, newForm)

			// 优化压缩
			// 更多参考see http://golang.org/pkg/archive/zip/#pkg-constants
			header.Method = azip.Deflate

			writer, err2 := zipWriter.CreateHeader(header)
			if err2 != nil {
				return err2
			}
			if _, err2 = io.Copy(writer, zipFile); err2 != nil {
				return err2
			}
			return nil
		}(file)
		if err != nil {
			return err
		}
	}
	return nil
}

// ZipFolder 压缩
func ZipFolder(srcFile, destZip string) error {
	// 预防：旧文件无法覆盖
	err := os.RemoveAll(destZip)
	if err != nil {
		return err
	}
	zipFile, err := os.Create(destZip)
	if err != nil {
		return err
	}
	defer zipFile.Close()
	archive := azip.NewWriter(zipFile)
	defer archive.Close()
	err = filepath.Walk(srcFile, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// 获取: 文件头信息
		header, err := azip.FileInfoHeader(info)
		if err != nil {
			return err
		}
		header.Name = strings.TrimPrefix(path, filepath.Dir(srcFile)+"/")
		// 判断文件是否为文件夹
		if info.IsDir() {
			header.Name += "/"
		} else {
			// 设置: zip 的文件压缩算法
			header.Method = azip.Deflate
		}
		// 创建: 压缩包头部信息
		writer, err := archive.CreateHeader(header)
		if err != nil {
			return err
		}
		if !info.IsDir() {
			file, err2 := os.Open(path)
			if err2 != nil {
				return err2
			}
			defer file.Close()
			_, err2 = io.Copy(writer, file)
			if err2 != nil {
				return err2
			}
		}
		return err
	})
	if err != nil {
		return err
	}
	return err
}

// Zip create zip file, fpath could be a single file or a directory
func Zip(fpath, destPath string) error {
	zipFile, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer zipFile.Close()

	archive := zip.NewWriter(zipFile)
	defer archive.Close()

	err = filepath.Walk(fpath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		header.Name = strings.TrimPrefix(path, filepath.Dir(fpath)+"/")

		if info.IsDir() {
			header.Name += "/"
		} else {
			header.Method = zip.Deflate
		}

		writer, err := archive.CreateHeader(header)
		if err != nil {
			return err
		}

		if !info.IsDir() {
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()
			_, err = io.Copy(writer, file)
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

// UnZip unzip the file and save it to destPath
func UnZip(zipFile, destPath string) error {
	zipReader, err := zip.OpenReader(zipFile)
	if err != nil {
		return err
	}
	defer zipReader.Close()
	for _, f := range zipReader.File {
		path := filepath.Join(destPath, f.Name)
		if f.FileInfo().IsDir() {
			err = os.MkdirAll(path, os.ModePerm)
			if err != nil {
				return err
			}
			continue
		}
		if err2 := os.MkdirAll(filepath.Dir(path), os.ModePerm); err2 != nil {
			return err2
		}
		inFile, err2 := f.Open()
		if err2 != nil {
			return err2
		}
		outFile, err2 := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err2 != nil {
			inFile.Close() // 关闭输入文件
			return err2
		}
		_, err2 = io.Copy(outFile, inFile)
		outFile.Close() // 关闭输出文件
		inFile.Close()  // 关闭输入文件
		if err2 != nil {
			return err2
		}
	}
	return nil
}
