package cachekey

import (
	"context"
	"fmt"
	"math/rand"
	"strconv"
	"testing"
	"time"

	"github.com/fzf-labs/fpkg/cache/redis"
	"github.com/fzf-labs/fpkg/cache/rockscache"
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
	newRocksCache := rockscache.NewRocksCache(newGoRedis)
	buildBatchKey := key.BuildHashKey("user123456")
	batch, err := buildBatchKey.RocksCache(context.Background(), newGoRedis, newRocksCache, "field", func() (string, error) {
		return strconv.Itoa(rand.Int()), nil
	})
	if err != nil {
		return
	}
	fmt.Println(batch)
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
	newRocksCache := rockscache.NewRocksCache(newGoRedis)
	buildBatchKey := key.BuildHashKey("user123456")
	err = buildBatchKey.RocksCacheDel(context.Background(), newGoRedis, newRocksCache)
	if err != nil {
		return
	}
}
