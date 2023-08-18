package signature

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
)

// Generate
// path 请求的路径 (不附带 querystring)
func (s *signature) Generate(path, method string, params json.RawMessage) (sign, timeStamp string, err error) {
	if path == "" {
		return "", "", errors.New("path required")
	}
	if method == "" {
		return "", "", errors.New("method required")
	}
	methodName := strings.ToUpper(method)
	if !methods[methodName] {
		return "", "", errors.New("method param error")
	}
	// Date
	timeStamp = strconv.FormatInt(time.Now().Unix(), 10)
	// 加密字符串规则
	buffer := bytes.NewBuffer(nil)
	buffer.WriteString(path)
	buffer.WriteString(delimiter)
	buffer.WriteString(methodName)
	buffer.WriteString(delimiter)
	buffer.WriteString(string(params))
	buffer.WriteString(delimiter)
	buffer.WriteString(timeStamp)
	// 对数据进行 sha256 加密，并进行 base64 encode
	hash := hmac.New(sha256.New, []byte(s.secret))
	hash.Write(buffer.Bytes())
	sign = base64.StdEncoding.EncodeToString(hash.Sum(nil))
	return sign, timeStamp, nil
}
