package jsonutil

import (
	JsonIter "github.com/json-iterator/go"
)

var jsonIter = JsonIter.ConfigCompatibleWithStandardLibrary

// JsonIterEncode 编码
func JsonIterEncode(v interface{}) ([]byte, error) {
	return jsonIter.Marshal(v)
}

// JsonIterEncodeToString 编码到字符串
func JsonIterEncodeToString(v interface{}) (string, error) {
	marshal, err := jsonIter.Marshal(v)
	if err != nil {
		return "", err
	}
	return string(marshal), nil
}

// JsonIterDecode 解码
func JsonIterDecode(bts []byte, ptr interface{}) error {
	return jsonIter.Unmarshal(bts, ptr)
}

// JsonIterDecodeString 解码字符串
func JsonIterDecodeString(str string, ptr interface{}) error {
	return jsonIter.Unmarshal([]byte(str), ptr)
}
