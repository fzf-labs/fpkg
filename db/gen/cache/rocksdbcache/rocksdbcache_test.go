package rocksdbcache

import (
	"context"
	"fmt"
	"testing"

	"github.com/fzf-labs/fpkg/cache/gorediscache"
	"github.com/fzf-labs/fpkg/cache/rockscache"
	"github.com/stretchr/testify/assert"
)

func TestRocksCache_Fetch(t *testing.T) {
	goRedis, err := gorediscache.NewGoRedis(gorediscache.GoRedisConfig{
		Addr:     "127.0.0.1:6379",
		Password: "123456",
		DB:       0,
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	rocksCacheClient := rockscache.NewWeakRocksCacheClient(goRedis)
	cache := NewRocksCache(rocksCacheClient)

	ctx := context.Background()

	fetch, err := cache.Fetch(ctx, "RocksCache_Fetch", func() (string, error) {
		fmt.Println(1)
		return "为什么程序员总是带着眼镜？", nil
	})
	fmt.Println(fetch)
	fmt.Println(err)
	assert.Equal(t, nil, err)
}

func TestRocksCache_FetchBatch(t *testing.T) {
	goRedis, err := gorediscache.NewGoRedis(gorediscache.GoRedisConfig{
		Addr:     "127.0.0.1:6379",
		Password: "123456",
		DB:       0,
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	rocksCacheClient := rockscache.NewWeakRocksCacheClient(goRedis)
	cache := NewRocksCache(rocksCacheClient)
	ctx := context.Background()
	keys := []string{
		"RocksCache_FetchBatch_a",
		"RocksCache_FetchBatch_b",
		"RocksCache_FetchBatch_c",
	}
	take, err := cache.FetchBatch(ctx, keys, func(miss []string) (map[string]string, error) {
		return map[string]string{
			"RocksCache_FetchBatch_a": "test1",
			"RocksCache_FetchBatch_b": "test2",
			"RocksCache_FetchBatch_c": "test3",
			"RocksCache_FetchBatch_d": "test4",
		}, nil
	})
	fmt.Println(take)
	fmt.Println(err)
	assert.Equal(t, nil, err)
}
