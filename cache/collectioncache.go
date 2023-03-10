package cache

import "github.com/coocood/freecache"

const (
	defaultCacheSize = 1024 * 1024 * 100
)

func NewDefaultCollectionCache() *freecache.Cache {
	return freecache.NewCache(defaultCacheSize)
}

func NewCollectionCache(size int) *freecache.Cache {
	return freecache.NewCache(size)
}
