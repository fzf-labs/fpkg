package casbin

import (
	"sync"

	"github.com/casbin/casbin/v2"
)

var (
	c    Casbin
	Once sync.Once
)

type Casbin struct {
	e *casbin.Enforcer
}

func New(f func() *casbin.Enforcer) Casbin {
	Once.Do(func() {
		c = Casbin{e: f()}
	})
	return c
}

func (c *Casbin) Check(sub, obj, act string) (bool, error) {
	return c.e.Enforce(sub, obj, act)
}
