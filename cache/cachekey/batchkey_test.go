package cachekey

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/fzf-labs/fpkg/cache/redis"
	"github.com/fzf-labs/fpkg/cache/rockscache"
	"github.com/fzf-labs/fpkg/util/maputil"
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
	newRocksCache := rockscache.NewRocksCache(newGoRedis)
	buildBatchKey := key.BuildBatchKey(maputil.Keys(kv)...)
	batch, err := buildBatchKey.RocksCacheBatch(context.Background(), newRocksCache, func() (map[string]string, error) {
		return kv, nil
	})
	if err != nil {
		return
	}
	fmt.Println(batch)
}
