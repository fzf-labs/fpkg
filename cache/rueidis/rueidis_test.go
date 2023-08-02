package rueidis

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestNewRueiis(t *testing.T) {
	client, err := NewRueidis(&RueidisConfig{
		Addr:     "127.0.0.1:6379",
		Password: "123456",
		DB:       0,
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
		array, err := client.DoCache(context.Background(), client.B().Hmget().Key("myhash").Field("1", "2").Cache(), time.Minute).ToArray()
		if err != nil {
			return
		}
		fmt.Printf("%+v \n", array)
	}
}

func TestNewRueidisAside(t *testing.T) {
	client, err := NewRueidisAside(&RueidisConfig{
		Addr:     "127.0.0.1:6379",
		Password: "123456",
		DB:       0,
	})
	if err != nil {
		return
	}
	val, err := client.Get(context.Background(), time.Minute, "mykey", func(ctx context.Context, key string) (val string, err error) {
		return "abcd", nil
	})
	fmt.Println(err)
	fmt.Println(val)
}
