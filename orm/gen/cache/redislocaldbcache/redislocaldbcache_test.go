package redislocaldbcache

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/fzf-labs/fpkg/cache/gorediscache"
	"github.com/stretchr/testify/assert"
)

func TestCache_Fetch(t *testing.T) {
	goRedis, err := gorediscache.NewGoRedis(gorediscache.GoRedisConfig{
		Addr:     "127.0.0.1:6379",
		Password: "123456",
		DB:       0,
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	redisLocalDBCache := NewRedisLocalDBCache(goRedis, WithLocalTTL(time.Second*2))
	ctx := context.Background()
	var result string
	for i := 0; i < 10; i++ {
		result, err = redisLocalDBCache.Fetch(ctx, "test", func() (string, error) {
			return "test", nil
		})
		if err != nil {
			return
		}
		fmt.Println(result)
		time.Sleep(time.Second)
	}
	assert.True(t, true, err == nil)
}

func TestCache_FetchBatch(t *testing.T) {
	goRedis, err := gorediscache.NewGoRedis(gorediscache.GoRedisConfig{
		Addr:     "127.0.0.1:6379",
		Password: "123456",
		DB:       0,
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	redisLocalDBCache := NewRedisLocalDBCache(goRedis, WithLocalCacheDisable())
	ctx := context.Background()
	for i := 0; i < 10; i++ {
		result, err2 := redisLocalDBCache.FetchBatch(ctx, []string{"test1", "test1", "test2", "test3", "test4", "test5", "test6", "test7", "test8", "test9", "test10", "test11", "test12", "test13", "test14", "test15"}, func(miss []string) (map[string]string, error) {
			// 15个假数据
			return map[string]string{
				"test1":  "test1",
				"test2":  "test2",
				"test3":  "test",
				"test4":  "test",
				"test5":  "test",
				"test6":  "test",
				"test7":  "test",
				"test8":  "test",
				"test9":  "test",
				"test10": "test",
				"test11": "test",
				"test12": "test",
				"test13": "test",
				"test14": "test",
				"test15": "test",
			}, nil
		})
		if err2 != nil {
			return
		}
		fmt.Println(result)
		time.Sleep(time.Second)
	}
	assert.True(t, true, err == nil)
}

func TestCache_Del(t *testing.T) {
	goRedis, err := gorediscache.NewGoRedis(gorediscache.GoRedisConfig{
		Addr:     "127.0.0.1:6379",
		Password: "123456",
		DB:       0,
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	redisLocalDBCache := NewRedisLocalDBCache(goRedis)
	ctx := context.Background()
	for i := 0; i < 10; i++ {
		err = redisLocalDBCache.Del(ctx, "test")
		if err != nil {
			return
		}
		time.Sleep(time.Second)
	}
	time.Sleep(time.Minute)
	assert.True(t, true, err == nil)
}
