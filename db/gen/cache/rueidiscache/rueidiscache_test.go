package rueidiscache

import (
	"context"
	"fmt"
	"testing"

	"github.com/fzf-labs/fpkg/cache/rueidiscache"
	"github.com/redis/rueidis"
	"github.com/stretchr/testify/assert"
)

func TestRueidisCache_Take(t *testing.T) {
	client, err := rueidiscache.NewRueidis(&rueidis.ClientOption{
		Username:    "",
		Password:    "123456",
		InitAddress: []string{"127.0.0.1:6379"},
		SelectDB:    0,
	})
	if err != nil {
		return
	}
	ctx := context.Background()
	rueidisCache := NewRueidisCache(client)
	take, err := rueidisCache.Fetch(ctx, "take_test", func() (string, error) {
		return "take", nil
	})
	fmt.Println(take)
	fmt.Println(err)
	assert.Equal(t, nil, err)
}

func TestRueidisCache_TakeBatch(t *testing.T) {
	client, err := rueidiscache.NewRueidis(&rueidis.ClientOption{
		Username:    "",
		Password:    "123456",
		InitAddress: []string{"127.0.0.1:6379"},
		SelectDB:    0,
	})
	if err != nil {
		return
	}
	ctx := context.Background()
	rueidisCache := NewRueidisCache(client)
	keys := []string{
		"a",
		"b",
		"c",
		"d",
	}
	take, err := rueidisCache.FetchBatch(ctx, keys, func(miss []string) (map[string]string, error) {
		return map[string]string{
			"a": "test1",
			"b": "test2",
			"c": "test3",
			"d": "test4",
		}, nil
	})
	fmt.Println(take)
	fmt.Println(err)
	assert.Equal(t, nil, err)
}

func TestRueidisCache_Del(t *testing.T) {
	client, err := rueidiscache.NewRueidis(&rueidis.ClientOption{
		Username:    "",
		Password:    "123456",
		InitAddress: []string{"127.0.0.1:6379"},
		SelectDB:    0,
	})
	if err != nil {
		return
	}
	ctx := context.Background()
	rueidisCache := NewRueidisCache(client)
	err = rueidisCache.Del(ctx, "a")
	if err != nil {
		return
	}
	assert.Equal(t, nil, err)
}

func TestRueidisCache_DelBatch(t *testing.T) {
	client, err := rueidiscache.NewRueidis(&rueidis.ClientOption{
		Username:    "",
		Password:    "123456",
		InitAddress: []string{"127.0.0.1:6379"},
		SelectDB:    0,
	})
	if err != nil {
		return
	}
	ctx := context.Background()
	rueidisCache := NewRueidisCache(client)
	err = rueidisCache.DelBatch(ctx, []string{"a", "b", "f"})
	if err != nil {
		return
	}
	assert.Equal(t, nil, err)
}
