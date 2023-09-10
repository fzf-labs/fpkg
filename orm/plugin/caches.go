package plugin

import (
	"github.com/go-gorm/caches/v2"
)

func NewCaches() *caches.Caches {
	return &caches.Caches{Conf: &caches.Config{
		Easer: true,
	}}
}
