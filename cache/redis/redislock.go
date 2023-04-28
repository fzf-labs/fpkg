package redis

import (
	"context"
	"log"
	"math/rand"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/go-redis/redis/v8"
)

const (
	randomLen       = 16
	tolerance       = 500 // milliseconds
	millisPerSecond = 1000
	lockCommand     = `if redis.call("GET", KEYS[1]) == ARGV[1] then
    redis.call("SET", KEYS[1], ARGV[1], "PX", ARGV[2])
    return "OK"
else
    return redis.call("SET", KEYS[1], ARGV[1], "NX", "PX", ARGV[2])
end`
	delCommand = `if redis.call("GET", KEYS[1]) == ARGV[1] then
    return redis.call("DEL", KEYS[1])
else
    return 0
end`
)

// A RedisLock is a redis lock.
type RedisLock struct {
	store   *redis.Client
	seconds uint32
	key     string
	id      string
}

// NewRedisLock 返回一个 RedisLock。
func NewRedisLock(store *redis.Client, key string) *RedisLock {
	return &RedisLock{
		store: store,
		key:   key,
		id:    Random(randomLen),
	}
}

// Acquire 获取锁。
func (rl *RedisLock) Acquire() (bool, error) {
	return rl.AcquireCtx(context.Background())
}

// AcquireCtx 使用给定的 ctx 获取锁。
func (rl *RedisLock) AcquireCtx(ctx context.Context) (bool, error) {
	seconds := atomic.LoadUint32(&rl.seconds)
	resp, err := rl.store.Eval(ctx, lockCommand, []string{rl.key}, rl.id, strconv.Itoa(int(seconds)*millisPerSecond+tolerance)).Result()
	if err == redis.Nil {
		return false, nil
	} else if err != nil {
		log.Printf("Error on acquiring lock for %s, %s", rl.key, err.Error())
		return false, err
	} else if resp == nil {
		return false, nil
	}

	reply, ok := resp.(string)
	if ok && reply == "OK" {
		return true, nil
	}

	log.Printf("Unknown reply when acquiring lock for %s: %v", rl.key, resp)
	return false, nil
}

// Release 释放锁。
func (rl *RedisLock) Release() (bool, error) {
	return rl.ReleaseCtx(context.Background())
}

// ReleaseCtx 使用给定的 ctx 释放锁。
func (rl *RedisLock) ReleaseCtx(ctx context.Context) (bool, error) {
	resp, err := rl.store.Eval(ctx, delCommand, []string{rl.key}, rl.id).Result()
	if err != nil {
		return false, err
	}

	reply, ok := resp.(int64)
	if !ok {
		return false, nil
	}

	return reply == 1, nil
}

// SetExpire 设置过期时间。
func (rl *RedisLock) SetExpire(seconds int) {
	atomic.StoreUint32(&rl.seconds, uint32(seconds))
}

// Random 随机字符串
func Random(n int) string {
	cs := make([]byte, n)
	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	sl := len(str)
	for i := 0; i < n; i++ {
		// 1607400451937462000
		idx := rand.New(rand.NewSource(time.Now().UnixNano())).Intn(sl) // 0 - 25
		cs[i] = str[idx]
	}
	return string(cs)
}
