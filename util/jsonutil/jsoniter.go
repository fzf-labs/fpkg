package jsonutil

import (
	JSONIter "github.com/json-iterator/go"
)

var jsonIter = JSONIter.ConfigCompatibleWithStandardLibrary

// JSONIterEncode 编码
func JSONIterEncode(v any) ([]byte, error) {
	return jsonIter.Marshal(v)
}

// JSONIterEncodeToString 编码到字符串
func JSONIterEncodeToString(v any) (string, error) {
	marshal, err := jsonIter.Marshal(v)
	if err != nil {
		return "", err
	}
	return string(marshal), nil
}

// JSONIterDecode 解码
func JSONIterDecode(bts []byte, ptr any) error {
	return jsonIter.Unmarshal(bts, ptr)
}

// JSONIterDecodeString 解码字符串
func JSONIterDecodeString(str string, ptr any) error {
	return jsonIter.Unmarshal([]byte(str), ptr)
}
