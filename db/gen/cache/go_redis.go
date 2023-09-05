package cache

import (
	"context"
	"strings"

	"github.com/fzf-labs/fpkg/conv"
	"github.com/redis/go-redis/v9"
)

type GoRedisCache struct {
	client redis.Client
}

func (r *GoRedisCache) Key(ctx context.Context, keys ...any) string {
	keyStr := make([]string, 0)
	for _, v := range keys {
		keyStr = append(keyStr, conv.String(v))
	}
	return strings.Join(keyStr, ":")
}

func (r *GoRedisCache) Take(ctx context.Context, keys []string, KvFn func() (map[string]string, error)) (map[string]string, error) {
	return KvFn()
}

func (r *GoRedisCache) Del(ctx context.Context, keys []string) error {
	//TODO implement me
	panic("implement me")
}
