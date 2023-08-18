package jwt

import (
	"context"
	"strconv"
	"time"

	"github.com/fzf-labs/fpkg/conv"
	jwts "github.com/golang-jwt/jwt"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
)

var (
	TokenMissing       = errors.New("token is missing")
	TokenExpired       = errors.New("token is expired")
	TokenInvalid       = errors.New("token is invalid")
	TokenCanNotRefresh = errors.New("token can not refresh")
	UIDNotRequest      = errors.New("uid not request")
	TokenStoreErr      = errors.New("token store fail")
	TokenGetErr        = errors.New("token get fail")
	TokenCheckErr      = errors.New("token check fail")
)

type Config struct {
	AccessSecret string // 秘钥
	AccessExpire int64  // 过期时间
	RefreshAfter int64  // 刷新时间 (小于过期时间,大于刷新时间,而小于过期时间,则刷新)
	Issuer       string // 签发人
}

type Token struct {
	Token     string `json:"token"`
	ExpiredAt int64  `json:"expired_at"`
	RefreshAt int64  `json:"refresh_at"`
}

type Jwt struct {
	Cfg   *Config
	redis *redis.Client
}

const (
	// jwt 官方定义
	JwtAudience  = "aud" // 受众
	JwtID        = "jti" // 编号
	JwtIssueAt   = "iat" // 签发时间
	JwtExpired   = "exp" // 过期时间
	JwtIssuer    = "iss" // 签发人
	JwtNotBefore = "nbf" // 生效时间，在此之前是无效的
	JwtSubject   = "sub" // 主题

	// 自定义
	JwtRefresh = "ref" // 刷新时间
	JwtUID     = "uid" // 用户标识
)

func NewJwt(redisClient *redis.Client, cfg *Config) *Jwt {
	return &Jwt{
		Cfg:   cfg,
		redis: redisClient,
	}
}

// GenerateToken token 生成
func (j *Jwt) GenerateToken(payloads map[string]interface{}) (*Token, jwts.MapClaims, error) {
	if payloads[JwtUID] == nil {
		return nil, nil, UIDNotRequest
	}
	now := time.Now()
	iat := now.Unix()
	expiredAt := iat + j.Cfg.AccessExpire
	refreshAt := iat + j.Cfg.RefreshAfter
	claims := make(jwts.MapClaims)
	claims[JwtID] = strconv.FormatInt(now.UnixNano(), 10)
	claims[JwtIssueAt] = iat
	claims[JwtIssuer] = j.Cfg.Issuer
	claims[JwtNotBefore] = iat - 1000
	claims[JwtExpired] = expiredAt
	claims[JwtRefresh] = refreshAt
	if len(payloads) > 0 {
		for k, v := range payloads {
			switch k {
			case JwtAudience, JwtExpired, JwtID, JwtIssueAt, JwtIssuer, JwtNotBefore, JwtSubject, JwtRefresh:
				// ignore the standard claims
			default:
				claims[k] = v
			}
		}
	}
	token := jwts.NewWithClaims(jwts.SigningMethodHS256, claims)
	signedString, err := token.SignedString([]byte(j.Cfg.AccessSecret))
	if err != nil {
		return nil, nil, err
	}
	return &Token{
		Token:     signedString,
		ExpiredAt: expiredAt,
		RefreshAt: refreshAt,
	}, claims, nil
}

// ParseToken token 解析
func (j *Jwt) ParseToken(tokenString string) (jwts.MapClaims, error) {
	token, err := jwts.Parse(tokenString, func(token *jwts.Token) (interface{}, error) {
		return []byte(j.Cfg.AccessSecret), nil
	})
	if err != nil {
		switch e := err.(type) {
		case *jwts.ValidationError:
			switch e.Errors {
			case jwts.ValidationErrorExpired: //过期
				return nil, TokenExpired
			default:
				return nil, errors.Wrap(TokenInvalid, "not jwt ValidationErrorExpired")
			}
		default:
			return nil, errors.Wrap(TokenInvalid, "not jwt ValidationError")
		}
	}
	if token == nil {
		return nil, errors.Wrap(TokenInvalid, "token is nil")
	}
	claims := token.Claims.(jwts.MapClaims)
	if _, ok := claims[JwtID]; !ok {
		return nil, errors.Wrap(TokenInvalid, "not set JwtID")
	}
	if _, ok := claims[JwtIssueAt]; !ok {
		return nil, errors.Wrap(TokenInvalid, "not set JwtIssueAt")
	}
	if _, ok := claims[JwtExpired]; !ok {
		return nil, errors.Wrap(TokenInvalid, "not set JwtExpired")
	}
	if _, ok := claims[JwtRefresh]; !ok {
		return nil, errors.Wrap(TokenInvalid, "not set JwtRefresh")
	}
	return claims, nil
}

func (j *Jwt) RefreshToken(oldClaims jwts.MapClaims) (*Token, jwts.MapClaims, error) {
	payloads := oldClaims
	if len(payloads) > 0 {
		for k := range payloads {
			switch k {
			case JwtAudience, JwtExpired, JwtID, JwtIssueAt, JwtIssuer, JwtNotBefore, JwtSubject, JwtRefresh:
				delete(payloads, k)
			default:
			}
		}
	}
	//生成新token
	newToken, newClaims, err := j.GenerateToken(payloads)
	if err != nil {
		return nil, nil, err
	}
	//token存入redis
	err = j.JwtTokenStore(newClaims)
	if err != nil {
		return nil, nil, err
	}
	//老的token写入黑名单中
	err = j.JwtBlackTokenStore(oldClaims, newToken.Token)
	if err != nil {
		return nil, nil, err
	}
	return newToken, newClaims, err
}

func (j *Jwt) GetPayloads(claims jwts.MapClaims) map[string]string {
	kv := make(map[string]string)
	if len(claims) > 0 {
		for k := range claims {
			switch k {
			case JwtAudience, JwtExpired, JwtID, JwtIssueAt, JwtIssuer, JwtNotBefore, JwtSubject, JwtRefresh:

			default:
				kv[k] = conv.String(claims[k])
			}
		}
	}
	return kv
}

func (j *Jwt) GetUID(claims jwts.MapClaims) int64 {
	return conv.Int64(claims[JwtUID])
}

type ContextWithValueKey string

func (j *Jwt) SetPayloadToContext(ctx context.Context, claims jwts.MapClaims) context.Context {
	for k, v := range claims {
		switch k {
		case JwtAudience, JwtExpired, JwtID, JwtIssueAt, JwtIssuer, JwtNotBefore, JwtSubject, JwtRefresh:
			// ignore the standard claims
		default:
			ctx = context.WithValue(ctx, ContextWithValueKey(k), v)
		}
	}
	return ctx
}

// 要做单一登录 即保存当前
func (j *Jwt) buildCacheKey(jwtUID string) string {
	return "jwt:" + j.Cfg.Issuer + ":" + jwtUID
}

// 黑名单的key
func (j *Jwt) buildBlackCacheKey(jwtUID, jwtID string) string {
	return "jwt_black:" + j.Cfg.Issuer + ":" + jwtUID + ":" + jwtID
}

// JwtBlackTokenStore 黑名单 防止刷新和请求时出现问题
func (j *Jwt) JwtBlackTokenStore(oldClaims jwts.MapClaims, newToken string) error {
	cacheKey := j.buildBlackCacheKey(conv.String(oldClaims[JwtUID]), conv.String(oldClaims[JwtID]))
	err := j.redis.Set(context.Background(), cacheKey, newToken, time.Second*time.Duration(10)).Err()
	if err != nil {
		return err
	}
	return nil
}

// JwtBlackTokenCheck Token黑名单检测  在黑名单中时 暂时允许通过
func (j *Jwt) JwtBlackTokenCheck(oldClaims jwts.MapClaims) (bool, *Token) {
	cacheKey := j.buildBlackCacheKey(conv.String(oldClaims[JwtUID]), conv.String(oldClaims[JwtID]))
	newToken, err := j.redis.Get(context.Background(), cacheKey).Result()
	if err != nil {
		return false, nil
	}
	//新的token不存在
	if newToken == "" {
		return false, nil
	}
	newClaims, err := j.ParseToken(newToken)
	if err != nil {
		return false, nil
	}
	return true, &Token{
		Token:     newToken,
		ExpiredAt: conv.Int64(newClaims[JwtExpired]),
		RefreshAt: conv.Int64(newClaims[JwtRefresh]),
	}
}

// JwtTokenStore 要做单一登录 即保存当前jwt的编号
func (j *Jwt) JwtTokenStore(claims jwts.MapClaims) error {
	cacheKey := j.buildCacheKey(conv.String(claims[JwtUID]))
	refreshTime := time.Unix(conv.Int64(claims[JwtRefresh]), 0)
	expiresAt := time.Until(refreshTime)
	err := j.redis.Set(context.Background(), cacheKey, conv.String(claims[JwtID]), expiresAt).Err()
	if err != nil {
		return err
	}
	return nil
}

// JwtTokenCheck Token检测
func (j *Jwt) JwtTokenCheck(claims jwts.MapClaims) error {
	cacheKey := j.buildCacheKey(conv.String(claims[JwtUID]))
	result, err := j.redis.Get(context.Background(), cacheKey).Result()
	if err != nil {
		if err != redis.Nil {
			return TokenGetErr
		}
	}
	jwtID := conv.Int64(claims[JwtID])
	if result != strconv.Itoa(int(jwtID)) {
		return TokenCheckErr
	}
	return nil
}

// JwtTokenClear Token清除
func (j *Jwt) JwtTokenClear(jwtUID string) error {
	cacheKey := j.buildCacheKey(jwtUID)
	err := j.redis.Del(context.Background(), cacheKey).Err()
	if err != nil {
		return err
	}
	return nil
}
