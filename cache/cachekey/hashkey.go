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

// CollectKey 收集的key
func (p *HashKey) CollectKey() string {
	return strings.Join([]string{p.keyPrefix.ServerName, p.keyPrefix.PrefixName, "collect", p.key}, ":")
}

// TTL 获取缓存key的过期时间time.Duration
func (p *HashKey) TTL() time.Duration {
	return p.keyPrefix.ExpirationTime
}

// TTLSecond 获取缓存key的过期时间 Second
func (p *HashKey) TTLSecond() int {
	return int(p.keyPrefix.ExpirationTime / time.Second)
}

// TTLUnix 获取缓存key的过期时间 unix
func (p *HashKey) TTLUnix() int64 {
	return time.Now().Add(p.keyPrefix.ExpirationTime).Unix()
}

// RocksCache  rocks缓存生成
func (p *HashKey) RocksCache(ctx context.Context, rd *redis.Client, rc *rockscache.Client, field string, fn func() (string, error)) (string, error) {
	cacheKey := p.Key(field)
	return rc.Fetch2(ctx, cacheKey, p.TTL(), func() (string, error) {
		result, err := fn()
		if err != nil {
			return "", err
		}
		_ = rd.HSet(ctx, p.CollectKey(), cacheKey, p.TTLUnix()).Err()
		_ = rd.Expire(ctx, p.CollectKey(), p.TTL()).Err()
		return result, nil
	})
}

// RocksCacheDel rocks缓存缓存删除
func (p *HashKey) RocksCacheDel(ctx context.Context, rd *redis.Client, rc *rockscache.Client) error {
	cacheKeys, err := rd.HGetAll(ctx, p.CollectKey()).Result()
	if err != nil {
		return err
	}
	delKeys := make([]string, 0)
	for k, v := range cacheKeys {
		tt, err := strconv.ParseUint(v, 10, 64)
		if err != nil {
			return err
		}
		if time.Now().Unix() < int64(tt) {
			delKeys = append(delKeys, k)
		}
	}
	if len(delKeys) > 0 {
		err = rc.TagAsDeletedBatch2(ctx, delKeys)
		if err != nil {
			return err
		}
	}
	err = rd.Del(ctx, p.CollectKey()).Err()
	if err != nil {
		return err
	}
	return nil
}
