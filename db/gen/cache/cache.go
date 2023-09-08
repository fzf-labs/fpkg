package cache

import "context"

type IDBCache interface {
	Key(fields ...any) string
	Fetch(ctx context.Context, key string, fn func() (string, error)) (string, error)
	FetchBatch(ctx context.Context, keys []string, fn func(miss []string) (map[string]string, error)) (map[string]string, error)
	Del(ctx context.Context, key string) error
	DelBatch(ctx context.Context, keys []string) error
}
