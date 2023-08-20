package wechat

import (
	"context"
	"time"

	goRedis "github.com/redis/go-redis/v9"
)

type RedisCache struct {
	RedisClient *goRedis.Client
}

func NewRedisCache(redisClient *goRedis.Client) *RedisCache {
	return &RedisCache{RedisClient: redisClient}
}

func (r *RedisCache) Get(key string) any {
	result, err := r.RedisClient.Get(context.Background(), key).Result()
	if err != nil {
		return nil
	}
	return result
}

func (r *RedisCache) Set(key string, val any, timeout time.Duration) error {
	return r.RedisClient.Set(context.Background(), key, val, timeout).Err()
}

func (r *RedisCache) IsExist(key string) bool {
	result, err := r.RedisClient.Exists(context.Background(), key).Result()
	if err != nil {
		return false
	}
	if result != 1 {
		return false
	}
	return true
}

func (r *RedisCache) Delete(key string) error {
	return r.RedisClient.Del(context.Background(), key).Err()
}
