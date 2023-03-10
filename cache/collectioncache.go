package cache

import (
	"github.com/fzf-labs/fpkg/cache/collectioncache"
	"time"
)

// NewDefaultCollectionCache 默认进程内缓存
func NewDefaultCollectionCache() *collectioncache.Cache {
	return NewCollectionCache("default", time.Minute, 10000)
}

// NewCollectionCache 进程内缓存
func NewCollectionCache(name string, duration time.Duration, limit int) *collectioncache.Cache {
	cache, err := collectioncache.NewCache(duration, collectioncache.WithLimit(limit), collectioncache.WithName(name))
	if err != nil {
		panic("collectionCache err:" + err.Error())
	}
	return cache
}
