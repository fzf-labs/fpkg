package cachekey

import (
	"context"
	"time"

	"github.com/fzf-labs/fpkg/cache/redislock"
	"github.com/go-redis/redis/v8"
)

// LockKey 实际key参数
type LockKey struct {
	keyPrefix *KeyPrefix
	buildKey  string
}

// Key 获取构建好的key
func (p *LockKey) Key() string {
	return p.buildKey
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

func (p *LockKey) AutoLock(ctx context.Context, rd *redis.Client, fn func() error) error {
	locker, err := redislock.Obtain(ctx, rd, p.Key(), p.TTL(), nil)
	if err != nil {
		return err
	}
	defer func(locker *redislock.Lock, ctx context.Context) {
		_ = locker.Release(ctx)
	}(locker, ctx)
	return fn()
}

func (p *LockKey) AutoLockRetry(ctx context.Context, rd *redis.Client, fn func() error) error {
	locker, err := redislock.Obtain(ctx, rd, p.Key(), p.TTL(), &redislock.Options{
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
