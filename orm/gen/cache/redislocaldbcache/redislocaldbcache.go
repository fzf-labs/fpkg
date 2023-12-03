package redislocaldbcache

import (
	"context"
	"log"
	"math/rand"
	"runtime/debug"
	"strings"
	"time"

	"github.com/dgraph-io/ristretto"
	"github.com/dtm-labs/rockscache"
	"github.com/fzf-labs/fpkg/conv"
	"github.com/redis/go-redis/v9"
	"golang.org/x/sync/singleflight"
)

type Cache struct {
	name       string             // 缓存名称
	rocksCache *rockscache.Client // 弱一致性RocksCache缓存客户端
	redisCache *redis.Client      // redis客户端
	redisTTL   time.Duration      // redis缓存过期时间
	localCache *ristretto.Cache   // 本地缓存
	localTTL   time.Duration      // 本地缓存过期时间
	channel    string             // 订阅频道
	sf         singleflight.Group // 防止缓存击穿
}

// newGoRedisLocalDBCache redis+localdb缓存
func newGoRedisLocalDBCache(client *redis.Client, opts ...CacheOption) *Cache {
	r := &Cache{
		name:       "GormCache",
		rocksCache: nil,
		redisCache: client,
		redisTTL:   time.Hour * 24,
		localCache: nil,
		localTTL:   time.Second * 10,
		channel:    "GormCacheGoRedisLocalCacheDelChannel",
	}
	if len(opts) > 0 {
		for _, v := range opts {
			v(r)
		}
	}
	if r.localCache == nil {
		localCache, err := newRistretto()
		if err != nil {
			panic(err)
		}
		r.localCache = localCache
	}
	if r.rocksCache == nil {
		r.rocksCache = NewRocksCacheClient(client)
	}
	return r
}

// newRistretto 本地缓存
func newRistretto() (*ristretto.Cache, error) {
	return ristretto.NewCache(&ristretto.Config{
		NumCounters: 1e7,       // number of keys to track frequency of (10M).
		MaxCost:     100000000, // maximum cost of cache (100M).
		BufferItems: 64,        // number of keys per Get buffer.
	})
}

// NewRocksCacheClient 弱一致性RocksCache缓存客户端
func NewRocksCacheClient(rdb *redis.Client) *rockscache.Client {
	rc := rockscache.NewClient(rdb, rockscache.NewDefaultOptions())
	// 常用参数设置
	// 1、强一致性(默认关闭强一致性，如果开启的话会影响性能)
	rc.Options.StrongConsistency = false
	// 2、redis出现问题需要缓存降级时设置为true
	rc.Options.DisableCacheRead = false   // 关闭缓存读，默认false；如果打开，那么Fetch就不从缓存读取数据，而是直接调用fn获取数据
	rc.Options.DisableCacheDelete = false // 关闭缓存删除，默认false；如果打开，那么TagAsDeleted就什么操作都不做，直接返回
	// 3、其他设置
	// 标记删除的延迟时间，默认10秒，设置为100毫秒秒表示：被删除的key在100毫秒后才从redis中彻底清除
	rc.Options.Delay = time.Millisecond * time.Duration(100)
	// 防穿透：若fn返回空字符串，空结果在缓存中的缓存时间，默认60秒
	rc.Options.EmptyExpire = time.Second * time.Duration(120)
	// 防雪崩,默认0.1,当前设置为0.1的话，如果设定为600的过期时间，那么过期时间会被设定为540s - 600s中间的一个随机数，避免数据出现同时到期
	rc.Options.RandomExpireAdjustment = 0.1 // 设置为默认就行
	return rc
}

type CacheOption func(cache *Cache)

// WithName 设置缓存名称
func WithName(name string) CacheOption {
	return func(r *Cache) {
		r.name = name
	}
}

// WithRedisTTL 设置redis缓存过期时间
func WithRedisTTL(ttl time.Duration) CacheOption {
	return func(r *Cache) {
		r.redisTTL = ttl
	}
}

// WithLocalTTL 设置本地缓存过期时间
func WithLocalTTL(ttl time.Duration) CacheOption {
	return func(r *Cache) {
		r.localTTL = ttl
	}
}

// WithChannel 设置订阅频道
func WithChannel(channel string) CacheOption {
	return func(r *Cache) {
		r.channel = channel
	}
}

// WithLocalCache 设置本地缓存客户端
func WithLocalCache(localCache *ristretto.Cache) CacheOption {
	return func(r *Cache) {
		r.localCache = localCache
	}
}

// WithRocksCache 设置RocksCache客户端
func WithRocksCache(rocksCache *rockscache.Client) CacheOption {
	return func(r *Cache) {
		r.rocksCache = rocksCache
	}
}

func (r *Cache) Key(keys ...any) string {
	keyStr := make([]string, 0)
	keyStr = append(keyStr, r.name)
	for _, v := range keys {
		keyStr = append(keyStr, conv.String(v))
	}
	return strings.Join(keyStr, ":")
}

func (r *Cache) TTL(ttl time.Duration) time.Duration {
	return ttl - time.Duration(rand.Float64()*0.1*float64(ttl)) //nolint:gosec
}

func (r *Cache) Fetch(ctx context.Context, key string, fn func() (string, error)) (string, error) {
	return r.Fetch2(ctx, key, fn, r.redisTTL)
}

func (r *Cache) Fetch2(ctx context.Context, key string, fn func() (string, error), expire time.Duration) (string, error) {
	do, err, _ := r.sf.Do(key, func() (any, error) {
		// 查询本地缓存
		result, ok := r.localCache.Get(key)
		if ok {
			return result, nil
		}
		result, err := r.rocksCache.Fetch2(ctx, key, expire, fn)
		if err != nil {
			return nil, err
		}
		r.localCache.SetWithTTL(key, result, 1, r.localTTL)
		return result, nil
	})
	if err != nil {
		return "", err
	}
	return do.(string), nil
}

func (r *Cache) FetchBatch(ctx context.Context, keys []string, fn func(miss []string) (map[string]string, error)) (map[string]string, error) {
	return r.FetchBatch2(ctx, keys, fn, r.redisTTL)
}

func (r *Cache) FetchBatch2(ctx context.Context, keys []string, fn func(miss []string) (map[string]string, error), expire time.Duration) (map[string]string, error) {
	resp := make(map[string]string)
	localMissKeys := make([]string, 0)
	// 查询本地缓存
	for _, v := range keys {
		localCacheResult, ok := r.localCache.Get(v)
		if ok {
			resp[v] = localCacheResult.(string)
		} else {
			localMissKeys = append(localMissKeys, v)
		}
	}
	if len(localMissKeys) == 0 {
		return resp, nil
	}
	rocksCacheResult, err := r.rocksCache.FetchBatch2(ctx, localMissKeys, expire, func(idx []int) (map[int]string, error) {
		result := make(map[int]string)
		rocksCacheMiss := make([]string, 0)
		for _, v := range idx {
			rocksCacheMiss = append(rocksCacheMiss, localMissKeys[v])
		}
		dbValue, err := fn(rocksCacheMiss)
		if err != nil {
			return nil, err
		}
		keyToInt := make(map[string]int)
		for k, v := range localMissKeys {
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
	for k, v := range rocksCacheResult {
		resp[keys[k]] = v
	}
	for _, v := range localMissKeys {
		r.localCache.SetWithTTL(v, resp[v], 1, r.localTTL)
	}
	return resp, nil
}

func (r *Cache) Del(ctx context.Context, key string) error {
	err := r.redisCache.Del(ctx, key).Err()
	if err != nil {
		return err
	}
	err = r.pub(ctx, r.channel, key)
	if err != nil {
		return err
	}
	r.localCache.Del(key)
	return nil
}

func (r *Cache) DelBatch(ctx context.Context, keys []string) error {
	err := r.redisCache.Del(ctx, keys...).Err()
	if err != nil {
		return err
	}
	err = r.pub(ctx, r.channel, strings.Join(keys, ","))
	if err != nil {
		return err
	}
	for _, v := range keys {
		r.localCache.Del(v)
	}
	return nil
}

func (r *Cache) Init() {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				log.Println(err)
				debug.PrintStack()
			}
		}()
		for {
			for r == nil || !r.ok() {
				time.Sleep(10 * time.Millisecond)
			}
			r.sub(r.channel)
		}
	}()
}

// ok 检查redis是否可用
func (r *Cache) ok() bool {
	_, err := r.redisCache.Ping(context.Background()).Result()
	return err == nil
}

// sub 订阅消息
func (r *Cache) sub(channel string) {
	ctx := context.Background()
	sub := r.redisCache.Subscribe(ctx, channel)
	// 使用完毕，记得关闭
	defer func(pubSub *redis.PubSub) {
		err := pubSub.Close()
		if err != nil {
			log.Println("sub close err:", err)
		}
		log.Println("sub close success")
	}(sub)
	for {
		msg, err := sub.ReceiveMessage(ctx)
		if err != nil {
			log.Println("sub ReceiveMessage err:", err.Error())
			time.Sleep(time.Second) // 等待一段时间再重试
			continue
		}
		if msg.String() != "" {
			log.Println("sub ReceiveMessage:", msg.String())
			keys := strings.Split(msg.String(), ":")
			for _, key := range keys {
				r.localCache.Del(key)
			}
		}
	}
}

// pub 发布消息
func (r *Cache) pub(ctx context.Context, channel, key string) error {
	return r.redisCache.Publish(ctx, channel, key).Err()
}
