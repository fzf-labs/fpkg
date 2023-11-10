package rocksdbcache

import (
	"context"
	"strings"
	"time"

	"github.com/dtm-labs/rockscache"
	"github.com/fzf-labs/fpkg/conv"
)

type Cache struct {
	name   string
	client *rockscache.Client
	ttl    time.Duration
}

type CacheOption func(cache *Cache)

func WithName(name string) CacheOption {
	return func(r *Cache) {
		r.name = name
	}
}
func WithTTL(ttl time.Duration) CacheOption {
	return func(r *Cache) {
		r.ttl = ttl
	}
}

func NewRocksDBCache(client *rockscache.Client, opts ...CacheOption) *Cache {
	r := &Cache{
		name:   "GormCache",
		client: client,
		ttl:    time.Hour * 24,
	}
	if len(opts) > 0 {
		for _, v := range opts {
			v(r)
		}
	}
	return r
}

func (r *Cache) Key(keys ...any) string {
	keyStr := make([]string, 0)
	keyStr = append(keyStr, r.name)
	for _, v := range keys {
		keyStr = append(keyStr, conv.String(v))
	}
	return strings.Join(keyStr, ":")
}
func (r *Cache) Fetch(ctx context.Context, key string, fn func() (string, error)) (string, error) {
	return r.Fetch2(ctx, key, fn, r.ttl)
}
func (r *Cache) Fetch2(ctx context.Context, key string, fn func() (string, error), expire time.Duration) (string, error) {
	return r.client.Fetch2(ctx, key, expire, fn)
}
func (r *Cache) FetchBatch(ctx context.Context, keys []string, fn func(miss []string) (map[string]string, error)) (map[string]string, error) {
	return r.FetchBatch2(ctx, keys, fn, r.ttl)
}
func (r *Cache) FetchBatch2(ctx context.Context, keys []string, fn func(miss []string) (map[string]string, error), expire time.Duration) (map[string]string, error) {
	resp, err := r.client.FetchBatch2(ctx, keys, expire, func(idx []int) (map[int]string, error) {
		result := make(map[int]string)
		miss := make([]string, 0)
		for _, v := range idx {
			miss = append(miss, keys[v])
		}
		dbValue, err := fn(miss)
		if err != nil {
			return nil, err
		}
		keyToInt := make(map[string]int)
		for k, v := range keys {
			keyToInt[v] = k
		}
		for k, v := range dbValue {
			result[keyToInt[k]] = v
		}
		return result, nil
	})
	if err != nil {
		return nil, err
	}
	result := make(map[string]string)
	for k, v := range resp {
		result[keys[k]] = v
	}
	return result, nil
}
func (r *Cache) Del(ctx context.Context, key string) error {
	return r.client.TagAsDeleted2(ctx, key)
}
func (r *Cache) DelBatch(ctx context.Context, keys []string) error {
	return r.client.TagAsDeletedBatch2(ctx, keys)
}
