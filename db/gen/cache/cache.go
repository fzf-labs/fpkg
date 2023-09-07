package cache

import "context"

type IDBCache interface {
	Key(ctx context.Context, fields ...any) string
	Fetch(ctx context.Context, key string, KvFn func() (string, error)) (string, error)
	FetchBatch(ctx context.Context, keys []string, KvFn func(miss []string) (map[string]string, error)) (map[string]string, error)
	Del(ctx context.Context, key string) error
	DelBatch(ctx context.Context, keys []string) error
}
