package cache

import (
	"github.com/coocood/freecache"
	"github.com/dtm-labs/rockscache"
	"time"
)

// Key 实际key参数
type Key struct {
	keyPrefix *KeyPrefix
	buildKey  string
}

// Key 获取构建好的key
func (p *Key) Key() string {
	return p.buildKey
}

// TTL 获取缓存key的过期时间time.Duration
func (p *Key) TTL() time.Duration {
	return p.keyPrefix.ExpirationTime
}

// TTLSecond 获取缓存key的过去时间 Second
func (p *Key) TTLSecond() int {
	return int(p.keyPrefix.ExpirationTime / time.Second)
}

// RocksCache  rocks缓存生成
func (p *Key) RocksCache(rc *rockscache.Client, fn func() (string, error)) (string, error) {
	return rc.Fetch(p.Key(), p.TTL(), fn)
}

// RocksCacheDel rocks缓存缓存删除
func (p *Key) RocksCacheDel(rc *rockscache.Client) error {
	return rc.TagAsDeleted(p.Key())
}

// CollectionCache 进程内缓存生成
func (p *Key) CollectionCache(fc *freecache.Cache, fn func() (string, error)) (string, error) {
	fcRes, err := fc.Get([]byte(p.Key()))
	if err != nil && err != freecache.ErrNotFound {
		return "", err
	}
	if len(fcRes) > 0 {
		return string(fcRes), nil
	}
	res, err := fn()
	if err != nil {
		return "", err
	}
	err = fc.Set([]byte(p.Key()), []byte(res), p.TTLSecond()/20)
	if err != nil {
		return "", err
	}
	return res, nil
}

// CollectionCacheDel 进程内缓存
func (p *Key) CollectionCacheDel(fc *freecache.Cache) error {
	fc.Del([]byte(p.Key()))
	return nil
}

// CollectionRocksCache 进程内缓存生成(该方法设计不完善,仅用于不更新的数据)
// 1.查询进程内的缓存,有则返回,无则去获取rockscache.
// 2.进程内缓存的过期时间请务必设置远小于redis.例小20倍
// 3.进程内缓存在数据发生更新时,未做删除处理,所以请务必谨慎.(一般需要去做订阅redis的pub/sub)
func (p *Key) CollectionRocksCache(fc *freecache.Cache, rc *rockscache.Client, fn func() (string, error)) (string, error) {
	fcRes, err := fc.Get([]byte(p.Key()))
	if err != nil && err != freecache.ErrNotFound {
		return "", err
	}
	if len(fcRes) > 0 {
		return string(fcRes), nil
	}
	rcRes, err := rc.Fetch(p.Key(), p.TTL(), fn)
	if err != nil {
		return "", err
	}
	err = fc.Set([]byte(p.Key()), []byte(rcRes), p.TTLSecond()/20)
	if err != nil {
		return "", err
	}
	return rcRes, nil
}

// CollectionRocksCacheDel 进程内缓存
func (p *Key) CollectionRocksCacheDel(fc *freecache.Cache, rc *rockscache.Client) error {
	fc.Del([]byte(p.Key()))
	err := rc.TagAsDeleted(p.Key())
	if err != nil {
		return err
	}
	return nil
}
