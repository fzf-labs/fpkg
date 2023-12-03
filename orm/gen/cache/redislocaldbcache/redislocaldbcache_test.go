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
	redisLocalDBCache := NewRedisLocalDBCache(goRedis)
	ctx := context.Background()

	for i := 0; i < 10; i++ {
		result, err2 := redisLocalDBCache.FetchBatch(ctx, []string{"test1", "test2"}, func(miss []string) (map[string]string, error) {
			return map[string]string{
				"test1": "test1",
				"test2": "test2",
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
