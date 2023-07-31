package redislocalcache

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
)

func TestNew(t *testing.T) {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "123456",
		DB:       0})
	localCache := NewTinyLFU(1000, time.Minute)

	cache := New("test", redisClient, localCache)
	ctx := context.TODO()
	key := "abc"
	value, err := cache.Get(ctx, key)
	if err != nil {
		return
	}
	fmt.Println(string(value))
	value2, err := cache.Get(ctx, key)
	if err != nil {
		return
	}
	fmt.Println(string(value2))
	time.Sleep(time.Hour)
}
