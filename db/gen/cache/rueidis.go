package cache

import (
	"context"
	"strings"

	"github.com/fzf-labs/fpkg/conv"
	"github.com/redis/rueidis"
)

type RueidisCache struct {
	client rueidis.Client
}

func NewRueidisCache(client rueidis.Client) *RueidisCache {
	return &RueidisCache{
		client: client,
	}
}

func (r *RueidisCache) Key(ctx context.Context, keys ...any) string {
	keyStr := make([]string, 0)
	for _, v := range keys {
		keyStr = append(keyStr, conv.String(v))
	}
	return strings.Join(keyStr, ":")
}

func (r *RueidisCache) Take(ctx context.Context, keys []string, KvFn func() (map[string]string, error)) (map[string]string, error) {
	return KvFn()
}

func (r *RueidisCache) Del(ctx context.Context, keys []string) error {
	//TODO implement me
	panic("implement me")
}
