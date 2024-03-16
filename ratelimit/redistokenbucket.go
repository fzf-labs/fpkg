package ratelimit

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/redis/go-redis/v9"
	xRate "golang.org/x/time/rate"
)

// redis 实现的令牌桶
const (
	// to be compatible with aliyun redis, we cannot use `local key = KEYS[1]` to reuse the key
	// KEYS[1] as tokens_key
	// KEYS[2] as timestamp_key
	script = `local rate = tonumber(ARGV[1])
local capacity = tonumber(ARGV[2])
local now = tonumber(ARGV[3])
local requested = tonumber(ARGV[4])
local fill_time = capacity/rate
local ttl = math.floor(fill_time*2)
local last_tokens = tonumber(redis.call("get", KEYS[1]))
if last_tokens == nil then
    last_tokens = capacity
end
local last_refreshed = tonumber(redis.call("get", KEYS[2]))
if last_refreshed == nil then
    last_refreshed = 0
end
local delta = math.max(0, now-last_refreshed)
local filled_tokens = math.min(capacity, last_tokens+(delta*rate))
local allowed = filled_tokens >= requested
local new_tokens = filled_tokens
if allowed then
    new_tokens = filled_tokens - requested
end
redis.call("setex", KEYS[1], ttl, new_tokens)
redis.call("setex", KEYS[2], ttl, now)
return allowed`
	pingInterval = time.Millisecond * 100
)

// RedisTokenBucket 控制事件在一秒钟内发生的频率。
type RedisTokenBucket struct {
	rate           int // 每秒生成几个令牌
	burst          int // 令牌桶最大值
	store          *redis.Client
	tokenKey       string
	timestampKey   string
	rescueLock     sync.Mutex
	redisAlive     uint32
	rescueLimiter  *xRate.Limiter
	monitorStarted bool
}

// NewTokenLimiter 返回一个新的令牌限制器，该令牌限制器允许事件达到速率，并允许最多突发令牌的突发。
func NewTokenLimiter(rate, burst int, store *redis.Client, key string) *RedisTokenBucket {
	tokenKey := fmt.Sprintf("%s:tokens", key)
	timestampKey := fmt.Sprintf("%s:ts", key)
	return &RedisTokenBucket{
		rate:          rate,
		burst:         burst,
		store:         store,
		tokenKey:      tokenKey,
		timestampKey:  timestampKey,
		redisAlive:    1,
		rescueLimiter: xRate.NewLimiter(xRate.Every(time.Second/time.Duration(rate)), burst),
	}
}

// Allow 是AllowN(time.Now()，1) 的简写。
func (lim *RedisTokenBucket) Allow() bool {
	return lim.AllowN(time.Now(), 1)
}

// AllowN 报告现在是否可能发生n个事件。如果您打算删除超出速率的跳过事件，请使用此方法。否则使用保留或等待。
func (lim *RedisTokenBucket) AllowN(now time.Time, n int) bool {
	return lim.reserveN(now, n)
}

func (lim *RedisTokenBucket) reserveN(now time.Time, n int) bool {
	if atomic.LoadUint32(&lim.redisAlive) == 0 {
		return lim.rescueLimiter.AllowN(now, n)
	}
	resp, err := lim.store.Eval(
		context.Background(),
		script,
		[]string{
			lim.tokenKey,
			lim.timestampKey,
		},
		[]string{
			strconv.Itoa(lim.rate),
			strconv.Itoa(lim.burst),
			strconv.FormatInt(now.Unix(), 10),
			strconv.Itoa(n),
		}).Result()
	// redis allowed == false
	// Lua boolean false -> r Nil bulk reply
	if errors.Is(err, redis.Nil) {
		return false
	} else if err != nil {
		fmt.Printf("fail to use rate limiter: %s, use in-process limiter for rescue", err)
		lim.startMonitor()
		return lim.rescueLimiter.AllowN(now, n)
	}
	code, ok := resp.(int64)
	if !ok {
		fmt.Printf("fail to eval redis script: %v, use in-process limiter for rescue", resp)
		lim.startMonitor()
		return lim.rescueLimiter.AllowN(now, n)
	}
	// redis allowed == true
	// Lua boolean true -> r integer reply with value of 1
	return code == 1
}

func (lim *RedisTokenBucket) startMonitor() {
	lim.rescueLock.Lock()
	defer lim.rescueLock.Unlock()
	if lim.monitorStarted {
		return
	}
	lim.monitorStarted = true
	atomic.StoreUint32(&lim.redisAlive, 0)
	go lim.waitForRedis()
}

func (lim *RedisTokenBucket) waitForRedis() {
	ticker := time.NewTicker(pingInterval)
	defer func() {
		ticker.Stop()
		lim.rescueLock.Lock()
		lim.monitorStarted = false
		lim.rescueLock.Unlock()
	}()
	for range ticker.C {
		val, _ := lim.store.Ping(context.Background()).Result()
		if val == "PONG" {
			atomic.StoreUint32(&lim.redisAlive, 1)
			return
		}
	}
}
