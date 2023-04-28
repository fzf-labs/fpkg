package collectioncache

import (
	"time"

	"github.com/fzf-labs/fpkg/cache/collectioncache"
)

// NewDefaultCollectionCache 默认进程内缓存
func NewDefaultCollectionCache() *Cache {
	return NewCollectionCache("default", time.Minute, 10000)
}

// NewCollectionCache 进程内缓存
func NewCollectionCache(name string, duration time.Duration, limit int) *Cache {
	cache, err := NewCache(duration, WithLimit(limit), WithName(name))
	if err != nil {
		panic("collectionCache err:" + err.Error())
	}
	return cache
}
