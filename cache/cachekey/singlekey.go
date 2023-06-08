package cachekey

import (
	"context"
	"time"

	"github.com/dtm-labs/rockscache"
	"github.com/fzf-labs/fpkg/cache/collectioncache"
)

// SingleKey 实际key参数
type SingleKey struct {
	keyPrefix *KeyPrefix
	buildKey  string
}

// Key 获取构建好的key
func (p *SingleKey) Key() string {
	return p.buildKey
}

// TTL 获取缓存key的过期时间time.Duration
func (p *SingleKey) TTL() time.Duration {
	return p.keyPrefix.ExpirationTime
}

// TTLSecond 获取缓存key的过期时间 Second
func (p *SingleKey) TTLSecond() int {
	return int(p.keyPrefix.ExpirationTime / time.Second)
}

// RocksCache  rocks缓存生成
func (p *SingleKey) RocksCache(ctx context.Context, rc *rockscache.Client, fn func() (string, error)) (string, error) {
	return rc.Fetch2(ctx, p.Key(), p.TTL(), fn)
}

// RocksCacheDel rocks缓存缓存删除
func (p *SingleKey) RocksCacheDel(ctx context.Context, rc *rockscache.Client) error {
	return rc.TagAsDeleted2(ctx, p.Key())
}

// CollectionCache 进程内缓存生成
func (p *SingleKey) CollectionCache(cc *collectioncache.Cache, fn func() (string, error)) (string, error) {
	take, err := cc.TakeWithExpire(p.Key(), p.TTL()/20, func() (interface{}, error) {
		return fn()
	})
	if err != nil {
		return "", err
	}
	return take.(string), nil
}

// CollectionCacheDel 进程内缓存
func (p *SingleKey) CollectionCacheDel(cc *collectioncache.Cache) error {
	cc.Del(p.Key())
	return nil
}

// CollectionRocksCache 进程内缓存生成(该方法设计不完善,仅用于不更新的数据)
// 1.查询进程内的缓存,有则返回,无则去获取RocksCache
// 2.进程内缓存的过期时间请务必设置远小于redis.例小20倍
// 3.Redis缓存在数据发生更新时,未做删除处理,所以请务必谨慎.(一般需要去做订阅redis的pub/sub)
func (p *SingleKey) CollectionRocksCache(cc *collectioncache.Cache, rc *rockscache.Client, fn func() (string, error)) (string, error) {
	ccRes, err := cc.TakeWithExpire(p.Key(), p.TTL()/20, func() (interface{}, error) {
		rcRes, err := rc.Fetch(p.Key(), p.TTL(), fn)
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
func (p *SingleKey) CollectionRocksCacheDel(cc *collectioncache.Cache, rc *rockscache.Client) error {
	cc.Del(p.Key())
	err := rc.TagAsDeleted(p.Key())
	if err != nil {
		return err
	}
	return nil
}
