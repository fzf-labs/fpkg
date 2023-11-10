package cache

import (
	"context"
	"time"
)

type IDBCache interface {
	// Key returns a string key for the given fields
	Key(fields ...any) string
	// Fetch fetches the value from cache, if not found, call fn to fetch and set to cache
	Fetch(ctx context.Context, key string, fn func() (string, error)) (string, error)
	// Fetch2 fetches the value from cache,set ttl,if not found, call fn to fetch and set to cache
	Fetch2(ctx context.Context, key string, fn func() (string, error), expire time.Duration) (string, error)
	// FetchBatch fetches the values from cache, if not found, call fn to fetch and set to cache
	FetchBatch(ctx context.Context, keys []string, fn func(miss []string) (map[string]string, error)) (map[string]string, error)
	// FetchBatch2 fetches the values from cache,set ttl, if not found, call fn to fetch and set to cache
	FetchBatch2(ctx context.Context, keys []string, fn func(miss []string) (map[string]string, error), expire time.Duration) (map[string]string, error)
	// Del deletes the value from cache
	Del(ctx context.Context, key string) error
	// DelBatch deletes the values from cache
	DelBatch(ctx context.Context, keys []string) error
}
