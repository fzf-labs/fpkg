package rueidiscache

import (
	"github.com/redis/rueidis"
	"github.com/redis/rueidis/rueidisaside"
	"github.com/redis/rueidis/rueidiscompat"
	"github.com/redis/rueidis/rueidisotel"
)

// NewRueidis  redis客户端rueidis
// redis > 6.0
func NewRueidis(clientOption *rueidis.ClientOption) (rueidis.Client, error) {
	// 初始化rueidis
	client, err := rueidis.NewClient(*clientOption)
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
func NewRueidisAside(clientOption *rueidis.ClientOption) (rueidisaside.CacheAsideClient, error) {
	return rueidisaside.NewClient(rueidisaside.ClientOption{
		ClientOption: *clientOption,
	})
}
