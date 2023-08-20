package rueidis

import (
	"github.com/redis/rueidis"
	"github.com/redis/rueidis/rueidisaside"
	"github.com/redis/rueidis/rueidiscompat"
	"github.com/redis/rueidis/rueidisotel"
)

type RueidisConfig struct {
	Addr     string `json:"addr"`
	Password string `json:"password"`
	DB       int    `json:"db"`
}

// NewRueidis  redis客户端rueidis
// redis > 6.0
func NewRueidis(conf *RueidisConfig) (rueidis.Client, error) {
	// 初始化rueidis
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

// NewRueidisAside 缓存和数据一起存储
// redis > 7.0
func NewRueidisAside(conf *RueidisConfig) (rueidisaside.CacheAsideClient, error) {
	return rueidisaside.NewClient(rueidisaside.ClientOption{
		ClientOption: rueidis.ClientOption{
			InitAddress: []string{conf.Addr},
			Password:    conf.Password,
			SelectDB:    conf.DB,
		},
	})
}
