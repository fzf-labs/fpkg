package oss

import (
	"context"
	"encoding/base64"
	"io"
	"mime/multipart"
	"net/http"
	"strings"
)

func (c *AliConfig) EnCrypt(file multipart.File) string {
	buffer := make([]byte, 500000)
	n, _ := file.Read(buffer)
	// base64压缩图片
	base64string := base64.StdEncoding.EncodeToString(buffer[:n])
	// 指定值 + 反转图片base64值
	base64string = c.Salt + Reverse(base64string)
	return base64string
}

func (c *AliConfig) EnCryptByBytes(file []byte) string {
	// base64压缩图片
	base64string := base64.StdEncoding.EncodeToString(file)
	// 指定值 + 反转图片base64值
	base64string = c.Salt + Reverse(base64string)
	return base64string
}

func (c *AliConfig) DeCrypt(url string) (string, error) {
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, url, http.NoBody)
	if err != nil {
		return "", err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	// 读取获取的[]byte数据
	data, _ := io.ReadAll(resp.Body)
	imageBase64 := base64.StdEncoding.EncodeToString(data)
	decodeString, err := base64.StdEncoding.DecodeString(imageBase64)
	if err != nil {
		return "", err
	}
	trim := strings.TrimLeft(string(decodeString), c.Salt)
	reverse := Reverse(trim)
	return "data:image/jpeg;base64," + reverse, nil
}

func (c *AliConfig) Base64(url string) (string, error) {
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, url, http.NoBody)
	if err != nil {
		return "", err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	// 读取获取的[]byte数据
	data, _ := io.ReadAll(resp.Body)
	imageBase64 := base64.StdEncoding.EncodeToString(data)
	return "data:image/jpeg;base64," + imageBase64, nil
}

func (c *AliConfig) DeCryptToByte(url string) ([]byte, error) {
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, url, http.NoBody)
	if err != nil {
		return nil, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	// 读取获取的[]byte数据
	data, _ := io.ReadAll(resp.Body)
	trim := strings.TrimLeft(string(data), c.Salt)
	reverse := Reverse(trim)
	decodeString, err := base64.StdEncoding.DecodeString(reverse)
	if err != nil {
		return nil, err
	}
	return decodeString, nil
}

// Reverse 字符串反转
func Reverse(str string) string {
	var result []byte
	tmp := []byte(str)
	length := len(str)
	for i := 0; i < length; i++ {
		result = append(result, tmp[length-i-1])
	}
	return string(result)
}
