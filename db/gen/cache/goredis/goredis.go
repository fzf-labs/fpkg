package goredis

import (
	"context"
	"strings"
	"time"

	"github.com/fzf-labs/fpkg/conv"
	"github.com/redis/go-redis/v9"
)

type GoRedisCache struct {
	client redis.Client
	ttl    time.Duration
}

func (r *GoRedisCache) Key(ctx context.Context, keys ...any) string {
	keyStr := make([]string, 0)
	for _, v := range keys {
		keyStr = append(keyStr, conv.String(v))
	}
	return strings.Join(keyStr, ":")
}

func (r *GoRedisCache) Fetch(ctx context.Context, key string, KvFn func() (string, error)) (string, error) {
	result, err := r.client.Get(ctx, key).Result()
	if err != nil && err != redis.Nil {
		return "", err
	}
	if result == "" && err == redis.Nil {
		result, err = KvFn()
		if err != nil {
			return "", err
		}
		err = r.client.Set(ctx, key, result, r.ttl).Err()
		if err != nil {
			return "", err
		}
	}
	return result, nil
}

func (r *GoRedisCache) FetchBatch(ctx context.Context, keys []string, KvFn func(miss []string) (map[string]string, error)) (map[string]string, error) {
	resp := make(map[string]string)
	miss := make([]string, 0)
	_, err := r.client.Pipelined(ctx, func(p redis.Pipeliner) error {
		for _, v := range keys {
			result, err := p.Get(ctx, v).Result()
			if err != nil && err != redis.Nil {
				return err
			}
			if result == "" && err == redis.Nil {
				miss = append(miss, v)
			}
			resp[v] = result
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	if len(miss) > 0 {
		dbValue, err := KvFn(miss)
		if err != nil {
			return nil, err
		}
		_, err = r.client.Pipelined(ctx, func(p redis.Pipeliner) error {
			for k, v := range dbValue {
				resp[k] = v
				err := p.Set(ctx, k, v, r.ttl).Err()
				if err != nil {
					return err
				}
			}
			return nil
		})
		if err != nil {
			return nil, err
		}
	}
	return resp, nil
}

func (r *GoRedisCache) Del(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}

func (r *GoRedisCache) DelBatch(ctx context.Context, keys []string) error {
	return r.client.Del(ctx, keys...).Err()
}
