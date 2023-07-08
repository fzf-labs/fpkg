package plugin

import (
	"github.com/go-gorm/caches"
)

func NewCaches() *caches.Caches {
	return &caches.Caches{Conf: &caches.Config{
		Easer: true,
	}}
}
