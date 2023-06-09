package cachekey

import (
	"context"
	"strings"
	"time"

	"github.com/fzf-labs/fpkg/cache/rockscache"
	"github.com/go-redis/redis/v8"
)

func (p *KeyPrefix) NewSingleKey(rd *redis.Client) *SingleKey {
	return &SingleKey{
		keyPrefix: p,
		rd:        rd,
	}
}

// SingleKey 实际key参数
type SingleKey struct {
	keyPrefix *KeyPrefix
	rd        *redis.Client
}

// BuildKey  获取key
func (p *SingleKey) BuildKey(keys ...string) string {
	return strings.Join(keys, ":")
}

// FinalKey 获取实际key
func (p *SingleKey) FinalKey(key string) string {
	return strings.Join([]string{p.keyPrefix.ServerName, p.keyPrefix.PrefixName, key}, ":")
}

// TTL 获取缓存key的过期时间time.Duration
func (p *SingleKey) TTL() time.Duration {
	return p.keyPrefix.ExpirationTime
}

// TTLSecond 获取缓存key的过期时间 Second
func (p *SingleKey) TTLSecond() int {
	return int(p.keyPrefix.ExpirationTime / time.Second)
}

// SingleCache  缓存生成
func (p *SingleKey) SingleCache(ctx context.Context, key string, fn func() (string, error)) (string, error) {
	return rockscache.NewWeakRocksCacheClient(p.rd).Fetch2(ctx, p.FinalKey(key), p.TTL(), fn)
}

// SingleCacheDel 缓存删除
func (p *SingleKey) SingleCacheDel(ctx context.Context, key string) error {
	return rockscache.NewWeakRocksCacheClient(p.rd).TagAsDeleted2(ctx, p.FinalKey(key))
}
