package cachekey

import (
	"testing"
	"time"

	"github.com/fzf-labs/fpkg/cache/collectioncache"
)

func TestCollectionKey_CollectionCache(t *testing.T) {
	keyManage := NewKeyManage("collection_test")
	key := keyManage.AddKey("local", time.Hour, "本地")
	cc := collectioncache.NewCollectionCache("demo", time.Hour, 1000)
	collectionKey := key.NewCollectionKey(cc)
	_, err := collectionKey.CollectionCache(collectionKey.BuildKey("demo"), func() (string, error) {
		return "11231", nil
	})
	if err != nil {
		return
	}
}
