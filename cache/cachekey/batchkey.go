package cachekey

import (
	"context"
	"strings"
	"time"

	"github.com/dtm-labs/rockscache"
)

// BatchKey 实际key参数
type BatchKey struct {
	keyPrefix *KeyPrefix
	keys      []string
}

// Keys 获取构建好的key
func (p *BatchKey) Keys() []string {
	result := make([]string, 0)
	if len(p.keys) > 0 {
		for _, key := range p.keys {
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

// RocksCacheBatch  rocks缓存生成
func (p *BatchKey) RocksCacheBatch(ctx context.Context, rc *rockscache.Client, fn func() (map[string]string, error)) (map[string]string, error) {
	resp := make(map[string]string)
	fetchBatch, err := rc.FetchBatch2(ctx, p.Keys(), p.TTL(), func(ids []int) (map[int]string, error) {
		values := make(map[int]string)
		m, err := fn()
		if err != nil {
			return values, err
		}
		for _, i := range ids {
			values[i] = m[p.keys[i]]
		}
		return values, nil
	})
	if err != nil {
		return resp, err
	}
	for k, v := range fetchBatch {
		resp[p.keys[k]] = v
	}
	return resp, nil
}

// RocksCacheBatchDel rocks缓存缓存删除
func (p *BatchKey) RocksCacheBatchDel(ctx context.Context, rc *rockscache.Client) error {
	return rc.TagAsDeletedBatch2(ctx, p.Keys())
}
