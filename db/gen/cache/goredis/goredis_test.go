package goredis

import (
	"context"
	"fmt"
	"testing"

	"github.com/fzf-labs/fpkg/cache/gorediscache"
)

func TestGoRedisCache_Fetch(t *testing.T) {
	goRedis, err := gorediscache.NewGoRedis(gorediscache.GoRedisConfig{
		Addr:     "127.0.0.1:6379",
		Password: "123456",
		DB:       0,
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	cache := NewGoRedisCache(goRedis)
	ctx := context.Background()
	fetch, err := cache.Fetch(ctx, "GoRedisCache_Fetch", func() (string, error) {
		fmt.Println(1)
		return "为什么程序员总是带着眼镜？", nil
	})
	fmt.Println(fetch)
	fmt.Println(err)
}

func TestGoRedisCache_FetchBatch(t *testing.T) {
	goRedis, err := gorediscache.NewGoRedis(gorediscache.GoRedisConfig{
		Addr:     "127.0.0.1:6379",
		Password: "123456",
		DB:       0,
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	cache := NewGoRedisCache(goRedis)
	ctx := context.Background()
	keys := []string{
		"GoRedisCache_Fetch_a",
		"GoRedisCache_Fetch_b",
		"GoRedisCache_Fetch_c",
		"GoRedisCache_Fetch_d",
	}
	fetch, err := cache.FetchBatch(ctx, keys, func(miss []string) (map[string]string, error) {
		fmt.Println("do FetchBatch")
		return map[string]string{
			"GoRedisCache_Fetch_a": "1",
			"GoRedisCache_Fetch_b": "2",
			"GoRedisCache_Fetch_c": "3",
		}, nil
	})
	fmt.Println(fetch)
	fmt.Println(err)
}
