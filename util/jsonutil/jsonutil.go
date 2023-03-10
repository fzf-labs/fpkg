package jsonutil

import (
	"github.com/bytedance/sonic"
)

// Encode 编码
func Encode(v interface{}) ([]byte, error) {
	return sonic.Marshal(v)
}

// EncodeToString 编码到字符串
func EncodeToString(v interface{}) (string, error) {
	return sonic.MarshalString(v)
}

// Decode 解码
func Decode(bts []byte, ptr interface{}) error {
	return sonic.Unmarshal(bts, ptr)
}

// DecodeString 解码字符串
func DecodeString(str string, ptr interface{}) error {
	return sonic.UnmarshalString(str, ptr)
}
