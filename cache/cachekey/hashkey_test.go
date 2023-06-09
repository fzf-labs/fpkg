package cachekey

import (
	"context"
	"fmt"
	"math/rand"
	"strconv"
	"testing"
	"time"

	"github.com/fzf-labs/fpkg/cache/redis"
	"github.com/stretchr/testify/assert"
)

func TestHashKey_RocksCache(t *testing.T) {
	keyManage := NewKeyManage("hash_test")
	key := keyManage.AddKey("test", time.Hour, "测试")
	newGoRedis, err := redis.NewGoRedis(redis.GoRedisConfig{
		Addr:     "127.0.0.1:6379",
		Password: "123456",
		DB:       0,
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	hashKey := key.NewHashKey(newGoRedis)
	for i := 0; i < 100; i++ {
		cache, err := hashKey.HashCache(context.Background(), "k", strconv.Itoa(i), func() (string, error) {
			return strconv.Itoa(rand.Int()), nil
		})
		if err != nil {
			return
		}
		fmt.Println(cache)
	}
}

func TestHashKey_RocksCacheDel(t *testing.T) {
	keyManage := NewKeyManage("hash_test")
	key := keyManage.AddKey("test", time.Hour, "测试")
	newGoRedis, err := redis.NewGoRedis(redis.GoRedisConfig{
		Addr:     "127.0.0.1:6379",
		Password: "123456",
		DB:       0,
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	hashKey := key.NewHashKey(newGoRedis)
	err = hashKey.HashCacheDel(context.Background(), "k")
	if err != nil {
		return
	}
	assert.Equal(t, err, nil)
}
