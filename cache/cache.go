package cache

import (
	"context"

	"github.com/dtm-labs/rockscache"
	"github.com/fzf-labs/fpkg/cache/cachekey"
	"github.com/fzf-labs/fpkg/cache/collectioncache"
	"github.com/go-redis/redis/v8"
)

type Cache struct {
	r  *redis.Client
	rc *rockscache.Client
	cc *collectioncache.Cache
}

func NewCache(r *redis.Client, rc *rockscache.Client, cc *collectioncache.Cache) *Cache {
	return &Cache{r: r, rc: rc, cc: cc}
}

// RocksCache  rocks缓存生成
func (c *Cache) RocksCache(ctx context.Context, key *cachekey.Key, fn func() (string, error)) (string, error) {
	return c.rc.Fetch2(ctx, key.Key(), key.TTL(), fn)
}

// RocksCacheDel rocks缓存缓存删除
func (c *Cache) RocksCacheDel(ctx context.Context, key *cachekey.Key) error {
	return c.rc.TagAsDeleted2(ctx, key.Key())
}

// CollectionCache 进程内缓存生成
func (c *Cache) CollectionCache(key *cachekey.Key, fn func() (string, error)) (string, error) {
	take, err := c.cc.TakeWithExpire(key.Key(), key.TTL()/20, func() (interface{}, error) {
		return fn()
	})
	if err != nil {
		return "", err
	}
	return take.(string), nil
}

// CollectionCacheDel 进程内缓存
func (c *Cache) CollectionCacheDel(key *cachekey.Key) error {
	c.cc.Del(key.Key())
	return nil
}

// CollectionRocksCache 进程内缓存生成(该方法设计不完善,仅用于不更新的数据)
// 1.查询进程内的缓存,有则返回,无则去获取RocksCache
// 2.进程内缓存的过期时间请务必设置远小于redis.例小20倍
// 3.Redis缓存在数据发生更新时,未做删除处理,所以请务必谨慎.(一般需要去做订阅redis的pub/sub) 或者使用redis 6.0的客户端缓存
func (c *Cache) CollectionRocksCache(ctx context.Context, key *cachekey.Key, fn func() (string, error)) (string, error) {
	ccRes, err := c.cc.TakeWithExpire(key.Key(), key.TTL()/20, func() (interface{}, error) {
		rcRes, err := c.rc.Fetch2(ctx, key.Key(), key.TTL(), fn)
		if err != nil {
			return nil, err
		}
		return rcRes, nil
	})
	if err != nil {
		return "", err
	}
	return ccRes.(string), nil
}

// CollectionRocksCacheDel 进程内缓存
func (c *Cache) CollectionRocksCacheDel(ctx context.Context, key *cachekey.Key) error {
	c.cc.Del(key.Key())
	err := c.rc.TagAsDeleted2(ctx, key.Key())
	if err != nil {
		return err
	}
	return nil
}
