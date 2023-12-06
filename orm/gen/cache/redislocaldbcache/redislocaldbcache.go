package redislocaldbcache

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"runtime/debug"
	"strings"
	"time"

	"github.com/dgraph-io/ristretto"
	"github.com/dtm-labs/rockscache"
	"github.com/fzf-labs/fpkg/conv"
	"github.com/redis/go-redis/v9"
	"golang.org/x/sync/errgroup"
	"golang.org/x/sync/singleflight"
)

type Cache struct {
	name              string             // 缓存名称
	localCacheEnable  bool               // 是否启用本地缓存
	localCache        *ristretto.Cache   // 本地缓存 默认开启
	localCacheTTL     time.Duration      // 本地缓存过期时间
	redisCache        *redis.Client      // redis客户端
	rocksCache        *rockscache.Client // 弱一致性RocksCache缓存客户端
	redisTTL          time.Duration      // redis缓存过期时间
	redisLuaBatchSize int                // redis lua 批量查询数量  默认100 有些云厂商对lua的keys有限制
	redisDelChan      string             // 订阅频道
	sf                singleflight.Group // 防止缓存击穿
}

// NewRedisLocalDBCache redis+localdb缓存
func NewRedisLocalDBCache(client *redis.Client, opts ...CacheOption) *Cache {
	r := &Cache{
		name:              "GormCache",
		localCache:        nil,
		localCacheTTL:     time.Second * 10,
		localCacheEnable:  true,
		redisCache:        client,
		rocksCache:        nil,
		redisTTL:          time.Hour * 24,
		redisLuaBatchSize: 100,
		redisDelChan:      "",
	}
	if len(opts) > 0 {
		for _, v := range opts {
			v(r)
		}
	}
	if r.localCacheEnable && r.localCache == nil {
		localCache, err := newRistretto()
		if err != nil {
			panic(err)
		}
		r.localCache = localCache
	}
	if !r.localCacheEnable && r.localCache != nil {
		r.localCache.Close()
		r.localCache = nil
	}
	if r.rocksCache == nil {
		r.rocksCache = newRocksCacheClient(client)
	}
	if r.redisDelChan == "" {
		r.redisDelChan = "RedisLocalCacheDelChannel" + r.name
	}
	if r.localCacheEnable {
		r.init()
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

// newRocksCacheClient 弱一致性RocksCache缓存客户端
func newRocksCacheClient(rdb *redis.Client) *rockscache.Client {
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

// WithLocalCache 设置本地缓存客户端
func WithLocalCache(localCache *ristretto.Cache) CacheOption {
	return func(r *Cache) {
		r.localCache = localCache
	}
}

// WithLocalTTL 设置本地缓存过期时间
func WithLocalTTL(ttl time.Duration) CacheOption {
	return func(r *Cache) {
		r.localCacheTTL = ttl
	}
}

// WithLocalCacheDisable 禁用本地缓存
func WithLocalCacheDisable() CacheOption {
	return func(r *Cache) {
		r.localCacheEnable = false
	}
}

// WithRedisTTL 设置redis缓存过期时间
func WithRedisTTL(ttl time.Duration) CacheOption {
	return func(r *Cache) {
		r.redisTTL = ttl
	}
}

// WithRedisDelChan 设置订阅频道
func WithRedisDelChan(channel string) CacheOption {
	return func(r *Cache) {
		r.redisDelChan = channel
	}
}

// WithRocksCache 设置RocksCache客户端
func WithRocksCache(rocksCache *rockscache.Client) CacheOption {
	return func(r *Cache) {
		r.rocksCache = rocksCache
	}
}

// WithRedisLuaBatchSize 设置RocksCache批量查询数量
func WithRedisLuaBatchSize(batchSize int) CacheOption {
	return func(r *Cache) {
		r.redisLuaBatchSize = batchSize
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
	if r.localCacheEnable {
		// 查询本地缓存
		localCacheValue, ok := r.localCache.Get(key)
		if ok {
			return localCacheValue.(string), nil
		}
	}
	// 查询redis缓存
	rocksCacheValue, err := r.rocksCache.Fetch2(ctx, key, expire, fn)
	if err != nil {
		return "", err
	}
	if r.localCacheEnable {
		// 设置本地缓存
		r.localCache.SetWithTTL(key, rocksCacheValue, 1, r.localCacheTTL)
	}
	return rocksCacheValue, nil
}

func (r *Cache) FetchBatch(ctx context.Context, keys []string, fn func(miss []string) (map[string]string, error)) (map[string]string, error) {
	return r.FetchBatch2(ctx, keys, fn, r.redisTTL)
}

func (r *Cache) FetchBatch2(ctx context.Context, keys []string, fn func(miss []string) (map[string]string, error), expire time.Duration) (map[string]string, error) {
	resp := make(map[string]string)
	// keyList
	keyList := unique(keys)
	if r.localCacheEnable {
		localMissKeys := make([]string, 0)
		// 查询本地缓存
		for _, v := range keyList {
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
		keyList = localMissKeys
	}
	// 查询redis缓存
	batch := chunk(keyList, r.redisLuaBatchSize)
	// 使用`errgroup`并发查询
	g, ctx := errgroup.WithContext(ctx)
	g.SetLimit(100)
	// 创建一个channel用于接收每个goroutine的结果
	resultCh := make(chan map[string]string, len(batch))
	for k := range batch {
		i := k
		g.Go(func() error {
			rocksCacheResult, err := r.fetchBatchItem(ctx, batch[i], fn, expire)
			if err != nil {
				return err
			}
			// 将结果发送到channel
			resultCh <- rocksCacheResult
			return nil
		})
	}
	// 等待所有goroutine执行完毕
	err := g.Wait()
	if err != nil {
		return nil, err
	}
	// 关闭channel
	close(resultCh)
	// 从channel中读取结果
	for result := range resultCh {
		for k, v := range result {
			resp[k] = v
			if r.localCacheEnable {
				// 设置本地缓存
				r.localCache.SetWithTTL(k, v, 1, r.localCacheTTL)
			}
		}
	}
	return resp, nil
}

// fetchBatchItem 批量查询
func (r *Cache) fetchBatchItem(ctx context.Context, keys []string, fn func(miss []string) (map[string]string, error), expire time.Duration) (map[string]string, error) {
	resp := make(map[string]string)
	// 查询redis缓存
	rocksCacheResult, err := r.rocksCache.FetchBatch2(ctx, keys, expire, func(idx []int) (map[int]string, error) {
		result := make(map[int]string)
		miss := make([]string, 0)
		for _, v := range idx {
			result[v] = ""
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
	for k, v := range rocksCacheResult {
		resp[keys[k]] = v
	}
	return resp, nil
}

func (r *Cache) Del(ctx context.Context, key string) error {
	err := r.rocksCache.TagAsDeleted2(ctx, key)
	if err != nil {
		return err
	}
	if r.localCacheEnable {
		err = r.pub(ctx, r.redisDelChan, key)
		if err != nil {
			return err
		}
		r.localCache.Del(key)
	}
	return nil
}

func (r *Cache) DelBatch(ctx context.Context, keys []string) error {
	keys = unique(keys)
	batch := chunk(keys, r.redisLuaBatchSize)
	g, ctx := errgroup.WithContext(ctx)
	g.SetLimit(100)
	for k := range batch {
		i := k
		g.Go(func() error {
			err := r.delBatchItem(ctx, batch[i])
			if err != nil {
				return err
			}
			return nil
		})
	}
	err := g.Wait()
	if err != nil {
		return err
	}
	return nil
}

func (r *Cache) delBatchItem(ctx context.Context, keys []string) error {
	err := r.rocksCache.TagAsDeletedBatch2(ctx, keys)
	if err != nil {
		return err
	}
	if !r.localCacheEnable {
		err = r.pub(ctx, r.redisDelChan, strings.Join(keys, ","))
		if err != nil {
			return err
		}
		for _, v := range keys {
			r.localCache.Del(v)
		}
	}
	return nil
}

func (r *Cache) init() {
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
			fmt.Println("sub redisDelChan:", r.redisDelChan)
			r.sub(r.redisDelChan)
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
			if len(keys) > 1 {
				for i, key := range keys {
					if i == 0 {
						continue
					}
					r.localCache.Del(key)
				}
			}
		}
	}
}

// pub 发布消息
func (r *Cache) pub(ctx context.Context, channel, key string) error {
	log.Println("pub:", key)
	return r.redisCache.Publish(ctx, channel, key).Err()
}

// unique 去重
func unique(slice []string) []string {
	if len(slice) == 0 {
		return slice
	}
	// here no use map filter. if use it, the result slice element order is random, not same as origin slice
	var result []string
	for i := 0; i < len(slice); i++ {
		v := slice[i]
		skip := true
		for j := range result {
			if v == result[j] {
				skip = false
				break
			}
		}
		if skip {
			result = append(result, v)
		}
	}
	return result
}

// chunk 将一个数组分成多个数组，每个数组包含size个元素，最后一个数组可能包含少于size个元素。
func chunk(collection []string, size int) [][]string {
	if size <= 0 {
		panic("Second parameter must be greater than 0")
	}
	chunksNum := len(collection) / size
	if len(collection)%size != 0 {
		chunksNum += 1
	}
	result := make([][]string, 0, chunksNum)
	for i := 0; i < chunksNum; i++ {
		last := (i + 1) * size
		if last > len(collection) {
			last = len(collection)
		}
		result = append(result, collection[i*size:last])
	}
	return result
}
