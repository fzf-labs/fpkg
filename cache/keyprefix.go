package cache

import (
	"fmt"
	"strings"
	"time"
)

type KeyPrefixes struct {
	List map[string]KeyPrefix
}

// NewKeyPrefixes 实例化缓存key前缀集合
func NewKeyPrefixes() *KeyPrefixes {
	return &KeyPrefixes{List: make(map[string]KeyPrefix)}
}

// NewCacheKeyPrefix 添加一个缓存prefix
func (p *KeyPrefixes) NewCacheKeyPrefix(prefix string, expirationTime time.Duration, remark string) *KeyPrefix {
	if _, ok := p.List[prefix]; ok {
		panic(fmt.Sprintf("cache key %s is exsit, please change one", prefix))
	}
	key := KeyPrefix{
		PrefixName:     prefix,
		Remark:         remark,
		ExpirationTime: expirationTime,
	}
	p.List[prefix] = key
	return &key
}

// KeyPrefix 缓存key前缀
type KeyPrefix struct {
	PrefixName     string
	Remark         string
	ExpirationTime time.Duration
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
