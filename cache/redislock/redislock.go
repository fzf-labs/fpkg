package redislock

import (
	"context"
	"time"

	"github.com/bsm/redislock"
	"github.com/redis/go-redis/v9"
)

func NewRedisLock(rd *redis.Client) *RedisLock {
	return &RedisLock{
		rd: rd,
	}
}

// RedisLock 实际key参数
type RedisLock struct {
	rd *redis.Client
}

func (p *RedisLock) ttl() time.Duration {
	return time.Second * 5
}

func (p *RedisLock) AutoLock(ctx context.Context, key string, fn func() error) error {
	locker, err := redislock.Obtain(ctx, p.rd, key, p.ttl(), nil)
	if err != nil {
		return err
	}
	defer func(locker *redislock.Lock, ctx context.Context) {
		_ = locker.Release(ctx)
	}(locker, ctx)
	return fn()
}

func (p *RedisLock) AutoLockRetry(ctx context.Context, key string, fn func() error) error {
	locker, err := redislock.Obtain(ctx, p.rd, key, p.ttl(), &redislock.Options{
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
