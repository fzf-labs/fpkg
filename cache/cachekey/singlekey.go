package cachekey

import (
	"context"
	"strings"
	"time"

	"github.com/dtm-labs/rockscache"
	"github.com/fzf-labs/fpkg/conv"
)

func (p *KeyPrefix) NewSingleKey(rc *rockscache.Client) *SingleKey {
	return &SingleKey{
		keyPrefix: p,
		rc:        rc,
	}
}

// SingleKey 实际key参数
type SingleKey struct {
	keyPrefix *KeyPrefix
	rc        *rockscache.Client
}

// BuildKey  获取key
func (p *SingleKey) BuildKey(keys ...any) string {
	keyStr := make([]string, 0)
	for _, v := range keys {
		keyStr = append(keyStr, conv.String(v))
	}
	return strings.Join(keyStr, ":")
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
	return p.rc.Fetch2(ctx, p.FinalKey(key), p.TTL(), fn)
}

// SingleCacheDel 缓存删除
func (p *SingleKey) SingleCacheDel(ctx context.Context, key string) error {
	return p.rc.TagAsDeleted2(ctx, p.FinalKey(key))
}
