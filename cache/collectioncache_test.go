package cache

import (
	"fmt"
	"testing"
)

func TestNewDefaultCollectionCache(t *testing.T) {
	cache := NewDefaultCollectionCache()
	get, err := cache.Get([]byte("a"))
	if err != nil {
		return
	}
	fmt.Println(get)
}
