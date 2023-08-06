package cachekey

import (
	"context"
	"strings"
	"time"

	"github.com/dtm-labs/rockscache"
)

func (p *KeyPrefix) NewBatchKey(rc *rockscache.Client) *BatchKey {
	return &BatchKey{
		keyPrefix: p,
		rc:        rc,
	}
}

// BatchKey 实际key参数
type BatchKey struct {
	keyPrefix *KeyPrefix
	rc        *rockscache.Client
}

// FinalKey 获取实际key
func (p *BatchKey) FinalKey(keys []string) []string {
	result := make([]string, 0)
	if len(keys) > 0 {
		for _, key := range keys {
			result = append(result, strings.Join([]string{p.keyPrefix.ServerName, p.keyPrefix.PrefixName, key}, ":"))
		}
	}
	return result
}

// TTL 获取缓存key的过期时间time.Duration
func (p *BatchKey) TTL() time.Duration {
	return p.keyPrefix.ExpirationTime
}

// TTLSecond 获取缓存key的过期时间 Second
func (p *BatchKey) TTLSecond() int {
	return int(p.keyPrefix.ExpirationTime / time.Second)
}

// BatchKeyCache 缓存生成
func (p *BatchKey) BatchKeyCache(ctx context.Context, keys []string, fn func() (map[string]string, error)) (map[string]string, error) {
	resp := make(map[string]string)
	finalKeys := p.FinalKey(keys)
	fetchBatch, err := p.rc.FetchBatch2(ctx, finalKeys, p.TTL(), func(ids []int) (map[int]string, error) {
		values := make(map[int]string)
		m, err := fn()
		if err != nil {
			return values, err
		}
		for _, i := range ids {
			values[i] = m[keys[i]]
		}
		return values, nil
	})
	if err != nil {
		return resp, err
	}
	for k, v := range fetchBatch {
		resp[finalKeys[k]] = v
	}
	return resp, nil
}

// BatchKeyCacheDel 缓存删除
func (p *BatchKey) BatchKeyCacheDel(ctx context.Context, keys []string) error {
	return p.rc.TagAsDeletedBatch2(ctx, p.FinalKey(keys))
}
