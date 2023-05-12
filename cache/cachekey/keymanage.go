package cachekey

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type KeyManage struct {
	ServerName string
	List       map[string]KeyPrefix
}

// NewKeyPrefixes 实例化缓存key前缀集合
func NewKeyPrefixes(serverName string) *KeyManage {
	return &KeyManage{
		ServerName: serverName,
		List:       make(map[string]KeyPrefix),
	}
}

// AddKey 添加一个缓存key prefix
func (p *KeyManage) AddKey(prefix string, expirationTime time.Duration, remark string) *KeyPrefix {
	if _, ok := p.List[prefix]; ok {
		panic(fmt.Sprintf("cache key %s is exsit, please change one", prefix))
	}
	key := KeyPrefix{
		ServerName:     p.ServerName,
		PrefixName:     prefix,
		Remark:         remark,
		ExpirationTime: expirationTime,
	}
	p.List[prefix] = key
	return &key
}

// Document 导出MD文档
func (p *KeyManage) Document() string {
	str := `|ServerName|PrefixName|TTL(s)|Remark` + "\n" + `|--|--|--|--|` + "\n"

	if len(p.List) > 0 {
		for _, m := range p.List {
			str += `|` + p.ServerName + `|` + m.PrefixName + `|` + strconv.FormatFloat(m.ExpirationTime.Seconds(), 'f', -1, 64) + `|` + m.Remark + `|` + "\n"
		}
	}
	return str
}

// KeyPrefix 缓存key前缀
type KeyPrefix struct {
	ServerName     string
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
		cacheKey.buildKey = strings.Join(append([]string{p.ServerName, p.PrefixName}), ":")
	} else {
		cacheKey.buildKey = strings.Join(append([]string{p.ServerName, p.PrefixName}, keys...), ":")
	}
	return &cacheKey
}
