package redis

import (
	"context"
	"fmt"
	"testing"
)

func TestNewGoRedis(t *testing.T) {
	newGoRedis, err := NewGoRedis(GoRedisConfig{
		Addr:     "127.0.0.1:6379",
		Password: "123456",
		DB:       0,
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	err = newGoRedis.Get(context.Background(), "test").Err()
	if err != nil {
		return
	}
}
