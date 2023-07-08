package jsonutil

import (
	"github.com/bytedance/sonic"
)

// SonicEncode 编码
func SonicEncode(v interface{}) ([]byte, error) {
	return sonic.Marshal(v)
}

// SonicEncodeToString 编码到字符串
func SonicEncodeToString(v interface{}) (string, error) {
	marshal, err := sonic.Marshal(v)
	if err != nil {
		return "", err
	}
	return string(marshal), nil
}

// SonicDecode 解码
func SonicDecode(bts []byte, ptr interface{}) error {
	return sonic.Unmarshal(bts, ptr)
}

// SonicDecodeString 解码字符串
func SonicDecodeString(str string, ptr interface{}) error {
	return sonic.Unmarshal([]byte(str), ptr)
}
