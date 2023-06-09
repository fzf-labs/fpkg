package cachekey

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/fzf-labs/fpkg/cache/redis"
)

func TestLockKey_AutoLock(t *testing.T) {
	keyManage := NewKeyManage("lock_test")
	key := keyManage.AddKey("autolock", time.Hour, "测试")
	newGoRedis, err := redis.NewGoRedis(redis.GoRedisConfig{
		Addr:     "127.0.0.1:6379",
		Password: "123456",
		DB:       0,
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	buildBatchKey := key.NewLockKey(newGoRedis)
	err = buildBatchKey.AutoLock(context.Background(), "AutoLock", func() error {
		fmt.Println(1111111)
		return nil
	})
	fmt.Println(err)
	if err != nil {
		return
	}
}

func TestLockKey_AutoLockRetry(t *testing.T) {
	keyManage := NewKeyManage("lock_test")
	key := keyManage.AddKey("autolock", time.Hour, "测试")
	newGoRedis, err := redis.NewGoRedis(redis.GoRedisConfig{
		Addr:     "127.0.0.1:6379",
		Password: "123456",
		DB:       0,
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	buildBatchKey := key.NewLockKey(newGoRedis)
	err = buildBatchKey.AutoLockRetry(context.Background(), "AutoLockRetry", func() error {
		fmt.Println(1111111)
		return nil
	})
	fmt.Println(err)
	if err != nil {
		return
	}
}
