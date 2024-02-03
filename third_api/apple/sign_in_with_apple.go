package apple

import (
	"context"
	"fmt"
	"strings"

	"github.com/Timothylock/go-signin-with-apple/apple"
	"github.com/pkg/errors"
	"github.com/pyihe/apple_validator"
)

//https://blog.csdn.net/qq_36770474/article/details/118340500

// SignInWithAppleConfig 苹果配置
type SignInWithAppleConfig struct {
	Secret      string
	KeyID       string
	TeamID      string
	ClientID    string
	RedirectURL string
}

type SignInWithApple struct {
	Cfg *SignInWithAppleConfig
}

// NewSignInWithApple 初始化
func NewSignInWithApple(cfg *SignInWithAppleConfig) *SignInWithApple {
	return &SignInWithApple{Cfg: cfg}
}

// CheckByAuthorizationCode 验证authorizationCode 直接请求校验接口
func (s *SignInWithApple) CheckByAuthorizationCode(code string) (uniqueID string, err error) {
	privateKey := s.FormatPrivateKey(s.Cfg.Secret)
	// 生成用于向 Apple 验证服务器进行身份验证的客户端密码
	clientSecret, err := apple.GenerateClientSecret(privateKey, s.Cfg.TeamID, s.Cfg.ClientID, s.Cfg.KeyID)
	if err != nil {
		return "", errors.New("error generating secret: " + err.Error())
	}
	// 生成新的验证客户端
	client := apple.New()
	vReq := apple.AppValidationTokenRequest{
		ClientID:     s.Cfg.ClientID,
		ClientSecret: clientSecret,
		Code:         code,
	}
	var resp apple.ValidationResponse
	// 进行验证
	err = client.VerifyAppToken(context.Background(), vReq, &resp)
	if err != nil {
		return "", errors.New("error verifying: " + err.Error())
	}
	if resp.Error != "" {
		return "", fmt.Errorf("apple returned an error: %s - %s", resp.Error, resp.ErrorDescription)
	}
	// Get the unique user ID
	unique, err := apple.GetUniqueID(resp.IDToken)
	if err != nil {
		return "", errors.New("failed to get unique ID: " + err.Error())
	}
	return unique, nil
}

// CheckIdentityToken 验证jwt 需要去服务端获取公钥,然后来验证格式是否正确
func (s *SignInWithApple) CheckIdentityToken(token string) (uniqueID string, err error) {
	validator := apple_validator.NewValidator()
	jwtToken, err := validator.CheckIdentityToken(token)
	if err != nil {
		return "", errors.New("CheckIdentityToken err: " + err.Error())
	}
	ok, err := jwtToken.IsValid()
	if !ok {
		return "", errors.New("CheckIdentityToken IsValid err: " + err.Error())
	}
	return jwtToken.Sub(), nil
}

func (s *SignInWithApple) FormatPrivateKey(privateKey string) (pKey string) {
	var buffer strings.Builder
	buffer.WriteString("-----BEGIN RSA PRIVATE KEY-----\n")
	rawLen := 64
	keyLen := len(privateKey)
	raws := keyLen / rawLen
	temp := keyLen % rawLen
	if temp > 0 {
		raws++
	}
	start := 0
	end := start + rawLen
	for i := 0; i < raws; i++ {
		if i == raws-1 {
			buffer.WriteString(privateKey[start:])
		} else {
			buffer.WriteString(privateKey[start:end])
		}
		buffer.WriteByte('\n')
		start += rawLen
		end = start + rawLen
	}
	buffer.WriteString("-----END RSA PRIVATE KEY-----\n")
	pKey = buffer.String()
	return
}
