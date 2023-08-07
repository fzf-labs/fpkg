package cachekey

import (
	"context"
	"strings"
	"time"

	"github.com/bsm/redislock"
	"github.com/fzf-labs/fpkg/conv"
	"github.com/redis/go-redis/v9"
)

func (p *KeyPrefix) NewLockKey(rd *redis.Client) *LockKey {
	return &LockKey{
		keyPrefix: p,
		rd:        rd,
	}
}

// LockKey 实际key参数
type LockKey struct {
	keyPrefix *KeyPrefix
	rd        *redis.Client
}

// BuildKey  获取key
func (p *LockKey) BuildKey(keys ...any) string {
	keyStr := make([]string, 0)
	for _, v := range keys {
		keyStr = append(keyStr, conv.String(v))
	}
	return strings.Join(keyStr, ":")
}

// FinalKey 获取实际key
func (p *LockKey) FinalKey(key string) string {
	return p.keyPrefix.Key(key)
}

// TTL 获取缓存key的过期时间time.Duration
func (p *LockKey) TTL() time.Duration {
	return p.keyPrefix.ExpirationTime
}

// TTLSecond 获取缓存key的过期时间 Second
func (p *LockKey) TTLSecond() int {
	return int(p.keyPrefix.ExpirationTime / time.Second)
}

// TTLUnix 获取缓存key的过期时间 unix
func (p *LockKey) TTLUnix() int64 {
	return time.Now().Add(p.keyPrefix.ExpirationTime).Unix()
}

func (p *LockKey) AutoLock(ctx context.Context, key string, fn func() error) error {
	locker, err := redislock.Obtain(ctx, p.rd, p.FinalKey(key), p.TTL(), nil)
	if err != nil {
		return err
	}
	defer func(locker *redislock.Lock, ctx context.Context) {
		_ = locker.Release(ctx)
	}(locker, ctx)
	return fn()
}

func (p *LockKey) AutoLockRetry(ctx context.Context, key string, fn func() error) error {
	locker, err := redislock.Obtain(ctx, p.rd, p.FinalKey(key), p.TTL(), &redislock.Options{
		RetryStrategy: redislock.LimitRetry(redislock.LinearBackoff(200*time.Millisecond), 5),
	})
	if err != nil {
		return err
	}
	defer func(locker *redislock.Lock, ctx context.Context) {
		_ = locker.Release(ctx)
	}(locker, ctx)
	return fn()
}
