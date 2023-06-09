package cachekey

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/fzf-labs/fpkg/cache/redis"
	"github.com/stretchr/testify/assert"
)

func TestBatchKey_RocksCacheBatch(t *testing.T) {
	keyManage := NewKeyManage("test")
	key := keyManage.AddKey("batch", time.Hour, "批量测试")
	kv := map[string]string{
		"a": "1",
		"b": "2",
		"c": "3",
		"d": "4",
		"e": "5",
		"f": "6",
	}
	newGoRedis, err := redis.NewGoRedis(redis.GoRedisConfig{
		Addr:     "127.0.0.1:6379",
		Password: "123456",
		DB:       0,
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	buildBatchKey := key.NewBatchKey(newGoRedis)
	batch, err := buildBatchKey.BatchKeyCache(context.Background(), []string{"a", "b", "c", "d", "e", "f"}, func() (map[string]string, error) {
		fmt.Println("do....")
		return kv, nil
	})
	if err != nil {
		return
	}
	fmt.Println(batch)
}

func TestBatchKey_BatchKeyCacheDel(t *testing.T) {
	keyManage := NewKeyManage("test")
	key := keyManage.AddKey("batch", time.Hour, "批量测试")
	newGoRedis, err := redis.NewGoRedis(redis.GoRedisConfig{
		Addr:     "127.0.0.1:6379",
		Password: "123456",
		DB:       0,
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	buildKey := key.NewBatchKey(newGoRedis)
	err = buildKey.BatchKeyCacheDel(context.Background(), []string{"a", "b", "c", "d", "e", "f"})
	if err != nil {
		return
	}
	assert.Equal(t, err, nil)
}
