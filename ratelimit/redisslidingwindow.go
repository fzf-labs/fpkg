package ratelimit

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

// redis 实现的滑动窗口
// to be compatible with aliyun redis, we cannot use `local key = KEYS[1]` to reuse the key
const periodScript = `local limit = tonumber(ARGV[1])
local window = tonumber(ARGV[2])
local current = redis.call("INCRBY", KEYS[1], 1)
if current == 1 then
    redis.call("expire", KEYS[1], window)
end
if current < limit then
    return 1
elseif current == limit then
    return 2
else
    return 0
end`

const (
	// Unknown means not initialized state.
	Unknown = iota
	// Allowed means allowed state.
	Allowed
	// HitQuota means this request exactly hit the quota.
	HitQuota
	// OverQuota means passed the quota.
	OverQuota

	internalOverQuota = 0
	internalAllowed   = 1
	internalHitQuota  = 2
)

// ErrUnknownCode is an error that represents unknown status code.
var ErrUnknownCode = errors.New("unknown status code")

type (
	// PeriodOption defines the method to customize a RedisSlidingWindow.
	PeriodOption func(l *RedisSlidingWindow)

	// A RedisSlidingWindow is used to limit requests during a period of time.
	RedisSlidingWindow struct {
		period     int
		quota      int
		limitStore *redis.Client
		keyPrefix  string
		align      bool
	}
)

// NewRedisSlidingWindow 返回具有给定参数的RedisSlidingWindow。
func NewRedisSlidingWindow(period, quota int, limitStore *redis.Client, keyPrefix string,
	opts ...PeriodOption) *RedisSlidingWindow {
	limiter := &RedisSlidingWindow{
		period:     period,
		quota:      quota,
		limitStore: limitStore,
		keyPrefix:  keyPrefix,
	}

	for _, opt := range opts {
		opt(limiter)
	}

	return limiter
}

// Take 请求许可证，它返回许可证状态。
func (h *RedisSlidingWindow) Take(key string) (int, error) {
	return h.TakeCtx(context.Background(), key)
}

// TakeCtx 请求带有上下文的许可，它将返回许可状态。
func (h *RedisSlidingWindow) TakeCtx(ctx context.Context, key string) (int, error) {
	resp, err := h.limitStore.Eval(ctx, periodScript, []string{h.keyPrefix + key}, []string{
		strconv.Itoa(h.quota),
		strconv.Itoa(h.calcExpireSeconds()),
	}).Result()
	if err != nil {
		return Unknown, err
	}

	code, ok := resp.(int64)
	if !ok {
		return Unknown, ErrUnknownCode
	}

	switch code {
	case internalOverQuota:
		return OverQuota, nil
	case internalAllowed:
		return Allowed, nil
	case internalHitQuota:
		return HitQuota, nil
	default:
		return Unknown, ErrUnknownCode
	}
}

func (h *RedisSlidingWindow) calcExpireSeconds() int {
	if h.align {
		now := time.Now()
		_, offset := now.Zone()
		unix := now.Unix() + int64(offset)
		return h.period - int(unix%int64(h.period))
	}

	return h.period
}

// Align returns a func to customize a RedisSlidingWindow with alignment.
// For example, if we want to limit end users with 5 sms verification messages every day,
// we need to align with the local timezone and the start of the day.
func Align() PeriodOption {
	return func(l *RedisSlidingWindow) {
		l.align = true
	}
}
