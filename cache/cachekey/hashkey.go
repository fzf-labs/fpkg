package cachekey

import (
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/dtm-labs/rockscache"
	"github.com/go-redis/redis/v8"
)

// HashKey 实际key参数
type HashKey struct {
	keyPrefix *KeyPrefix
	key       string
}

// Key 获取构建好的key
func (p *HashKey) Key(field string) string {
	return strings.Join([]string{p.keyPrefix.ServerName, p.keyPrefix.PrefixName, p.key, field}, ":")
}

// TTL 获取缓存key的过期时间time.Duration
func (p *HashKey) TTL() time.Duration {
	return p.keyPrefix.ExpirationTime
}

// TTLSecond 获取缓存key的过去时间 Second
func (p *HashKey) TTLSecond() int {
	return int(p.keyPrefix.ExpirationTime / time.Second)
}

// RocksCache  rocks缓存生成
func (p *HashKey) RocksCache(ctx context.Context, rd *redis.Client, rc *rockscache.Client, field string, fn func() (string, error)) (string, error) {
	cacheKey := p.Key(field)
	err := rd.HSet(ctx, p.key, cacheKey, p.TTL()).Err()
	if err != nil {
		return "", err
	}
	err = rd.Expire(ctx, p.key, p.TTL()).Err()
	if err != nil {
		return "", err
	}
	return rc.Fetch2(ctx, cacheKey, p.TTL(), fn)
}

// RocksCacheDel rocks缓存缓存删除
func (p *HashKey) RocksCacheDel(ctx context.Context, rd *redis.Client, rc *rockscache.Client) error {
	cacheKeys, err := rd.HGetAll(ctx, p.key).Result()
	if err != nil {
		return err
	}
	delKeys := make([]string, 0)
	for k, v := range cacheKeys {
		tt, err := strconv.Atoi(v)
		if err != nil {
			return err
		}
		if time.Now().Unix() < int64(tt) {
			delKeys = append(delKeys, k)
		}
	}
	return rc.TagAsDeletedBatch2(ctx, delKeys)
}
