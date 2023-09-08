package goredis

import (
	"context"
	"math/rand"
	"strings"
	"time"

	"github.com/fzf-labs/fpkg/conv"
	"github.com/redis/go-redis/v9"
	"golang.org/x/sync/singleflight"
)

type GoRedisCache struct {
	name   string
	client *redis.Client
	ttl    time.Duration
	sf     singleflight.Group
}

func NewGoRedisCache(client *redis.Client, opts ...CacheOption) *GoRedisCache {
	r := &GoRedisCache{
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

type CacheOption func(cache *GoRedisCache)

func WithName(name string) CacheOption {
	return func(r *GoRedisCache) {
		r.name = name
	}
}
func WithTTL(ttl time.Duration) CacheOption {
	return func(r *GoRedisCache) {
		r.ttl = ttl
	}
}
func (r *GoRedisCache) Key(ctx context.Context, keys ...any) string {
	keyStr := make([]string, 0)
	keyStr = append(keyStr, r.name)
	for _, v := range keys {
		keyStr = append(keyStr, conv.String(v))
	}
	return strings.Join(keyStr, ":")
}

func (r *GoRedisCache) TTL(ttl time.Duration) time.Duration {
	return ttl - time.Duration(rand.Float64()*0.1*float64(ttl))
}
func (r *GoRedisCache) Fetch(ctx context.Context, key string, KvFn func() (string, error)) (string, error) {
	do, err, _ := r.sf.Do(key, func() (interface{}, error) {
		result, err := r.client.Get(ctx, key).Result()
		if err != nil && err != redis.Nil {
			return "", err
		}
		if result == "" && err == redis.Nil {
			result, err = KvFn()
			if err != nil {
				return "", err
			}
			err = r.client.Set(ctx, key, result, r.TTL(r.ttl)).Err()
			if err != nil {
				return "", err
			}
		}
		return result, nil
	})
	if err != nil {
		return "", err
	}
	return do.(string), nil
}

func (r *GoRedisCache) FetchBatch(ctx context.Context, keys []string, KvFn func(miss []string) (map[string]string, error)) (map[string]string, error) {
	resp := make(map[string]string)
	miss := make([]string, 0)
	pipelined, err := r.client.Pipelined(ctx, func(p redis.Pipeliner) error {
		for _, v := range keys {
			_, err := p.Get(ctx, v).Result()
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil && err != redis.Nil {
		return nil, err
	}
	for k, cmder := range pipelined {
		if cmder.Err() == redis.Nil {
			miss = append(miss, keys[k])
		}
		resp[keys[k]] = cmder.(*redis.StringCmd).Val()
	}
	if len(miss) > 0 {
		dbValue, err := KvFn(miss)
		if err != nil {
			return nil, err
		}
		_, err = r.client.Pipelined(ctx, func(p redis.Pipeliner) error {
			for _, v := range miss {
				resp[v] = dbValue[v]
				err := p.Set(ctx, v, dbValue[v], r.TTL(r.ttl)).Err()
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
