package jsonutil

import (
	"encoding/json"
	"fmt"

	"github.com/bytedance/sonic"
)

// Encode 编码
func Encode(v interface{}) ([]byte, error) {
	return sonic.Marshal(v)
}

// EncodeToString 编码到字符串
func EncodeToString(v interface{}) (string, error) {
	marshal, err := sonic.Marshal(v)
	if err != nil {
		return "", err
	}
	return string(marshal), nil
}

// Decode 解码
func Decode(bts []byte, ptr interface{}) error {
	return sonic.Unmarshal(bts, ptr)
}

// DecodeString 解码字符串
func DecodeString(str string, ptr interface{}) error {
	return sonic.Unmarshal([]byte(str), ptr)
}

// Dump 打印
func Dump(v interface{}) {
	marshal, err := json.MarshalIndent(v, "", " ")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(marshal))
}
