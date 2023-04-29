package jsonutil

import (
	"encoding/json"
)

// Encode 编码
func Encode(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

// EncodeToString 编码到字符串
func EncodeToString(v interface{}) (string, error) {
	marshal, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	return string(marshal), nil
}

// Decode 解码
func Decode(bts []byte, ptr interface{}) error {
	return json.Unmarshal(bts, ptr)
}

// DecodeString 解码字符串
func DecodeString(str string, ptr interface{}) error {
	return json.Unmarshal([]byte(str), ptr)
}
