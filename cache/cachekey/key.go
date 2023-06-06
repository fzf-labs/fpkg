package cachekey

import (
	"time"
)

// Key 实际key参数
type Key struct {
	keyPrefix *KeyPrefix
	buildKey  string
}

// Key 获取构建好的key
func (p *Key) Key() string {
	return p.buildKey
}

// TTL 获取缓存key的过期时间time.Duration
func (p *Key) TTL() time.Duration {
	return p.keyPrefix.ExpirationTime
}

// TTLSecond 获取缓存key的过去时间 Second
func (p *Key) TTLSecond() int {
	return int(p.keyPrefix.ExpirationTime / time.Second)
}
