package byteutil

import (
	"bytes"
	"io"

	"github.com/klauspost/compress/zlib"
)

// IsLetterUpper 检查给定字节 b 是否为大写。
func IsLetterUpper(b byte) bool {
	if b >= byte('A') && b <= byte('Z') {
		return true
	}
	return false
}

// IsLetterLower 检查给定字节 b 是否为小写。
func IsLetterLower(b byte) bool {
	if b >= byte('a') && b <= byte('z') {
		return true
	}
	return false
}

// IsLetter 检查给定的字节 b 是否是一个字母。
func IsLetter(b byte) bool {
	return IsLetterUpper(b) || IsLetterLower(b)
}

// ZlibCompress zlib压缩
func ZlibCompress(data []byte) ([]byte, error) {
	var b bytes.Buffer
	w := zlib.NewWriter(&b)
	defer func(w *zlib.Writer) {
		_ = w.Close()
	}(w)
	_, err := w.Write(data)
	if err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

// ZlibUnCompress zlib解压
func ZlibUnCompress(data []byte) ([]byte, error) {
	var b bytes.Buffer
	w := bytes.NewReader(data)
	r, err := zlib.NewReader(w)
	if err != nil {
		return nil, err
	}
	defer func(r io.ReadCloser) {
		_ = r.Close()
	}(r)
	_, err = io.Copy(&b, r)
	if err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}
