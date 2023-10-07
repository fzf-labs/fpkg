package jsonutil

import (
	"encoding/json"
	"fmt"

	"github.com/bytedance/sonic"
)

// Marshal 编码
func Marshal(v any) ([]byte, error) {
	return sonic.Marshal(v)
}

// MarshalToString 编码到字符串
func MarshalToString(v any) (string, error) {
	return sonic.MarshalString(v)
}

// Unmarshal 解码
func Unmarshal(bts []byte, ptr any) error {
	return sonic.Unmarshal(bts, ptr)
}

// UnmarshalString 解码字符串
func UnmarshalString(str string, ptr any) error {
	return sonic.UnmarshalString(str, ptr)
}

// Dump 打印
func Dump(v any) {
	marshal, err := json.MarshalIndent(v, "", " ")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(marshal))
}
