package cachekey

import (
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/dtm-labs/rockscache"
	"github.com/fzf-labs/fpkg/conv"
	"github.com/redis/go-redis/v9"
)

func (p *KeyPrefix) NewHashKey(rd *redis.Client, rc *rockscache.Client) *HashKey {
	return &HashKey{
		keyPrefix: p,
		rd:        rd,
		rc:        rc,
	}
}

// HashKey 实际key参数
type HashKey struct {
	keyPrefix *KeyPrefix
	rd        *redis.Client
	rc        *rockscache.Client
}

// BuildKey  获取key
func (p *HashKey) BuildKey(keys ...any) string {
	keyStr := make([]string, 0)
	for _, v := range keys {
		keyStr = append(keyStr, conv.String(v))
	}
	return strings.Join(keyStr, ":")
}

// FinalKey 获取实际key
func (p *HashKey) FinalKey(key string, field string) string {
	return strings.Join([]string{p.keyPrefix.ServerName, p.keyPrefix.PrefixName, key, field}, ":")
}

// DelKey 收集的key
func (p *HashKey) DelKey(key string) string {
	return strings.Join([]string{p.keyPrefix.ServerName, p.keyPrefix.PrefixName, "DEL", key}, ":")
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

// HashCache  缓存生成
func (p *HashKey) HashCache(ctx context.Context, key string, field string, fn func() (string, error)) (string, error) {
	return p.rc.Fetch2(ctx, p.FinalKey(key, field), p.TTL(), func() (string, error) {
		result, err := fn()
		if err != nil {
			return "", err
		}
		_ = p.rd.HSet(ctx, p.DelKey(key), p.FinalKey(key, field), p.TTLUnix()).Err()
		_ = p.rd.Expire(ctx, p.DelKey(key), p.TTL()).Err()
		return result, nil
	})
}

// HashCacheDel 缓存删除
func (p *HashKey) HashCacheDel(ctx context.Context, key string) error {
	cacheKeys, err := p.rd.HGetAll(ctx, p.DelKey(key)).Result()
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
		err = p.rc.TagAsDeletedBatch2(ctx, delKeys)
		if err != nil {
			return err
		}
	}
	err = p.rd.Del(ctx, p.DelKey(key)).Err()
	if err != nil {
		return err
	}
	return nil
}
