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

func (s *signature) Verify(sign, timeStamp, path, method string, params json.RawMessage) (err error) {
	if path == "" {
		return errors.New("请求路径不存在")
	}
	if method == "" {
		return errors.New("请求方法不存在")
	}
	methodName := strings.ToUpper(method)
	if !methods[methodName] {
		return errors.New("请求方法错误")
	}
	t, err := strconv.ParseInt(timeStamp, 10, 63)
	if err != nil {
		return err
	}
	if time.Since(time.Unix(t, 0)) > s.ttl {
		return errors.Errorf("接口超时,限时:%v", s.ttl)
	}
	buffer := bytes.NewBuffer(nil)
	buffer.WriteString(path)
	buffer.WriteString(delimiter)
	buffer.WriteString(methodName)
	buffer.WriteString(delimiter)
	buffer.WriteString(string(params))
	buffer.WriteString(delimiter)
	buffer.WriteString(timeStamp)
	// 对数据进行 hmac 加密，并进行 base64 encode
	hash := hmac.New(sha256.New, []byte(s.secret))
	hash.Write(buffer.Bytes())
	digest := base64.StdEncoding.EncodeToString(hash.Sum(nil))
	if sign != digest {
		return errors.New("签名校验不通过")
	}
	return nil
}
