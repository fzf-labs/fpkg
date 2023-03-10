package cache

import (
	"fmt"
	"strings"
	"time"
)

var KeyPrefixes = map[string]*KeyPrefix{}

// KeyPrefix 缓存key前缀管理
type KeyPrefix struct {
	PrefixName     string
	Remark         string
	ExpirationTime time.Duration
}

func NewCacheKey(prefixName string, expirationTime time.Duration, remark string) *KeyPrefix {
	if _, ok := KeyPrefixes[prefixName]; ok {
		panic(fmt.Sprintf("cache key %s is exsit, please change one", prefixName))
	}
	key := &KeyPrefix{
		PrefixName:     prefixName,
		Remark:         remark,
		ExpirationTime: expirationTime,
	}
	KeyPrefixes[prefixName] = key
	return key
}

// BuildCacheKey 构建一个带有前缀的缓存key 使用 ":" 分隔
func (p *KeyPrefix) BuildCacheKey(keys ...string) *Key {
	cacheKey := Key{
		keyPrefix: p,
	}
	if len(keys) == 0 {
		cacheKey.buildKey = p.PrefixName
	} else {
		cacheKey.buildKey = strings.Join(append([]string{p.PrefixName}, keys...), ":")
	}
	return &cacheKey
}
