package rueidiscache

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/redis/rueidis"
	"github.com/stretchr/testify/assert"
)

func TestNewRueiis(t *testing.T) {
	client, err := NewRueidis(rueidis.ClientOption{
		Username:    "",
		Password:    "123456",
		InitAddress: []string{"127.0.0.1:6379"},
		SelectDB:    0,
	})
	if err != nil {
		return
	}
	client.DoMulti(
		context.Background(),
		client.B().Hmset().Key("myhash").FieldValue().FieldValue("1", "a").FieldValue("2", "b").Build(),
		client.B().Expire().Key("myhash").Seconds(1000).Build(),
	)

	ticker := time.NewTicker(time.Second * 5)
	for range ticker.C {
		fmt.Println("do")
		array, err2 := client.DoCache(context.Background(), client.B().Hmget().Key("myhash").Field("1", "2").Cache(), time.Minute).ToArray()
		if err2 != nil {
			return
		}
		fmt.Printf("%+v \n", array)
	}
	assert.Equal(t, nil, err)
}

func TestNewRueidisAside(t *testing.T) {
	client, err := NewRueidisAside(rueidis.ClientOption{
		Username:    "",
		Password:    "123456",
		InitAddress: []string{"127.0.0.1:6379"},
		SelectDB:    0,
	})
	if err != nil {
		return
	}
	val, err := client.Get(context.Background(), time.Minute, "mykey", func(ctx context.Context, key string) (val string, err error) {
		return "abcd", nil
	})
	fmt.Println(err)
	fmt.Println(val)
	assert.Equal(t, nil, err)
}
