package redis

import (
	"github.com/redis/rueidis"
	"github.com/redis/rueidis/rueidiscompat"
	"github.com/redis/rueidis/rueidisotel"
)

type RueidisConfig struct {
	Addr     string `json:"addr"`
	Password string `json:"password"`
	DB       int    `json:"db"`
}

func NewRueidis(conf *RueidisConfig) (rueidis.Client, error) {
	//初始化rueidis
	client, err := rueidis.NewClient(rueidis.ClientOption{
		InitAddress: []string{conf.Addr},
		Password:    conf.Password,
		SelectDB:    conf.DB,
	})
	if err != nil {
		return nil, err
	}
	// 链路追踪
	client = rueidisotel.WithClient(client)
	return client, nil
}

func NewRueidisAdapter(client rueidis.Client) (rueidiscompat.Cmdable, error) {
	compat := rueidiscompat.NewAdapter(client)
	return compat, nil
}
