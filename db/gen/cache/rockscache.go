package cache

import (
	"context"
	"strings"

	"github.com/dtm-labs/rockscache"
	"github.com/fzf-labs/fpkg/conv"
)

type RocksCache struct {
	client rockscache.Client
}

func (r *RocksCache) Key(ctx context.Context, keys ...any) string {
	keyStr := make([]string, 0)
	for _, v := range keys {
		keyStr = append(keyStr, conv.String(v))
	}
	return strings.Join(keyStr, ":")
}

func (r *RocksCache) Take(ctx context.Context, keys []string, KvFn func() (map[string]string, error)) (map[string]string, error) {
	return KvFn()
}

func (r *RocksCache) Del(ctx context.Context, keys []string) error {
	//TODO implement me
	panic("implement me")
}
