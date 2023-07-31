package redislocalcache

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/klauspost/compress/zlib"
	"github.com/redis/go-redis/v9"
	"golang.org/x/exp/slog"
)

var defaultPubSubKey = "RedisLocalCachePubSubKeyDel"

type rediser interface {
	Set(ctx context.Context, key string, value interface{}, ttl time.Duration) *redis.StatusCmd
	Get(ctx context.Context, key string) *redis.StringCmd
	Del(ctx context.Context, keys ...string) *redis.IntCmd

	Publish(ctx context.Context, channel string, message interface{}) *redis.IntCmd
	Subscribe(ctx context.Context, channels ...string) *redis.PubSub
}
type CacheOption func(cache *Cache)

type Cache struct {
	Name       string
	Redis      rediser
	LocalCache LocalCache
}

func New(name string, redis rediser, localCache LocalCache, opts ...CacheOption) *Cache {
	cache := &Cache{
		Name:       name,
		Redis:      redis,
		LocalCache: localCache,
	}
	if len(opts) > 0 {
		for _, v := range opts {
			v(cache)
		}
	}
	go func() {
		err := cache.subscribe(cache.getChannel())
		if err != nil {
			return
		}
	}()
	return cache
}

func (cd *Cache) Set(ctx context.Context, key string, value []byte, ttl time.Duration) error {
	cd.LocalCache.Set(key, value)
	if ttl == 0 {
		ttl = time.Hour
	}
	err := cd.Redis.Set(ctx, key, value, ttl).Err()
	_ = cd.publish(ctx, cd.getChannel(), key)
	return err
}

// Exists reports whether value for the given key exists.
func (cd *Cache) Exists(ctx context.Context, key string) bool {
	_, err := cd.getBytes(ctx, key)
	return err == nil
}

// Get gets the value for the given key.
func (cd *Cache) Get(ctx context.Context, key string) ([]byte, error) {
	return cd.getBytes(ctx, key)
}

func (cd *Cache) getBytes(ctx context.Context, key string) ([]byte, error) {
	b, ok := cd.LocalCache.Get(key)
	if ok {
		return b, nil
	}
	b, err := cd.Redis.Get(ctx, key).Bytes()
	if err != nil && err != redis.Nil {
		return nil, err
	}
	cd.LocalCache.Set(key, b)
	return b, nil
}

func (cd *Cache) Delete(ctx context.Context, key string) error {
	cd.LocalCache.Del(key)
	_, err := cd.Redis.Del(ctx, key).Result()
	_ = cd.publish(ctx, cd.getChannel(), key)
	return err
}

func (cd *Cache) DeleteFromLocalCache(key string) {
	cd.LocalCache.Del(key)
}

func (cd *Cache) getChannel() string {
	return strings.Join([]string{defaultPubSubKey, cd.Name}, ":")
}

func (cd *Cache) subscribe(channel string) error {
	ctx := context.Background()
	pubSub := cd.Redis.Subscribe(ctx, channel)
	// 使用完毕，记得关闭
	defer func(pubSub *redis.PubSub) {
		err := pubSub.Close()
		if err != nil {
			fmt.Printf("pubSub close err: %s", err)
		}
		fmt.Println("pubSub close success")
	}(pubSub)
	for {
		msg, err := pubSub.ReceiveMessage(ctx)
		if err != nil {
			slog.Error("pubSub ReceiveMessage err:", err.Error())
		}
		if msg.String() != "" {
			slog.Info("pubSub ReceiveMessage :", msg.String())
			cd.DeleteFromLocalCache(msg.String())
		}
	}
}

func (cd *Cache) publish(ctx context.Context, channel string, key string) error {
	return cd.Redis.Publish(ctx, channel, key).Err()
}

func ZlibCompress(data []byte) ([]byte, error) {
	var b bytes.Buffer
	w := zlib.NewWriter(&b)
	defer func(w *zlib.Writer) {
		_ = w.Close()
	}(w)
	_, err := w.Write(data)
	if err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

func ZlibUnCompress(data []byte) ([]byte, error) {
	var b bytes.Buffer
	w := bytes.NewReader(data)
	r, err := zlib.NewReader(w)
	if err != nil {
		return nil, err
	}
	defer func(r io.ReadCloser) {
		_ = r.Close()
	}(r)
	_, err = io.Copy(&b, r)
	if err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}
