package cachekey

import (
	"context"
	"strings"
	"time"

	"github.com/fzf-labs/fpkg/cache/redislock"
	"github.com/go-redis/redis/v8"
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
func (p *LockKey) BuildKey(keys ...string) string {
	return strings.Join(keys, ":")
}

// FinalKey 获取实际key
func (p *LockKey) FinalKey(key string) string {
	return strings.Join([]string{p.keyPrefix.ServerName, p.keyPrefix.PrefixName, "LOCK", key}, ":")
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
