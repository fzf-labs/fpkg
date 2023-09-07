package rueidiscache

import (
	"context"
	"strings"
	"time"

	"github.com/fzf-labs/fpkg/conv"
	"github.com/redis/rueidis"
)

type RueidisCache struct {
	client rueidis.Client
	ttl    time.Duration
}

type RueidisCacheOption func(cache *RueidisCache)

func WithTTL(ttl time.Duration) RueidisCacheOption {
	return func(r *RueidisCache) {
		r.ttl = ttl
	}
}

func NewRueidisCache(client rueidis.Client, opts ...RueidisCacheOption) *RueidisCache {
	r := &RueidisCache{
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

func (r *RueidisCache) Key(ctx context.Context, keys ...any) string {
	keyStr := make([]string, 0)
	for _, v := range keys {
		keyStr = append(keyStr, conv.String(v))
	}
	return strings.Join(keyStr, ":")
}

func (r *RueidisCache) Fetch(ctx context.Context, key string, KvFn func() (string, error)) (string, error) {
	cacheValue := r.client.DoCache(ctx, r.client.B().Get().Key(key).Cache(), r.ttl)
	//数据存在
	if !rueidis.IsRedisNil(cacheValue.Error()) {
		resp, err := cacheValue.ToString()
		if err != nil {
			return "", err
		}
		return resp, nil
	}
	resp, err := KvFn()
	if err != nil {
		return "", err
	}
	err = r.client.Do(ctx, r.client.B().Set().Key(key).Value(resp).Ex(r.ttl).Build()).Error()
	if err != nil {
		return "", err
	}
	return resp, nil
}

func (r *RueidisCache) FetchBatch(ctx context.Context, keys []string, KvFn func(miss []string) (map[string]string, error)) (map[string]string, error) {
	resp := make(map[string]string)
	commands := make([]rueidis.CacheableTTL, 0)
	for _, v := range keys {
		commands = append(commands, rueidis.CT(r.client.B().Get().Key(v).Cache(), r.ttl))
	}
	cacheValue := r.client.DoMultiCache(ctx, commands...)
	miss := make([]string, 0)
	for k, v := range cacheValue {
		if rueidis.IsRedisNil(v.Error()) {
			miss = append(miss, keys[k])
		}
		toString, _ := v.ToString()
		resp[keys[k]] = toString
	}
	if len(miss) > 0 {
		dbValue, err := KvFn(miss)
		if err != nil {
			return nil, err
		}
		completes := make([]rueidis.Completed, 0)
		for k, v := range dbValue {
			completes = append(completes, r.client.B().Set().Key(k).Value(v).Ex(r.ttl).Build())
			resp[k] = v
		}
		r.client.DoMulti(ctx, completes...)
	}
	return resp, nil
}

func (r *RueidisCache) Del(ctx context.Context, key string) error {
	return r.client.Do(ctx, r.client.B().Del().Key(key).Build()).Error()
}

func (r *RueidisCache) DelBatch(ctx context.Context, keys []string) error {
	completes := make([]rueidis.Completed, 0)
	for _, v := range keys {
		completes = append(completes, r.client.B().Del().Key(v).Build())
	}
	r.client.DoMulti(ctx, completes...)
	return nil
}
