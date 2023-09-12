package localcache

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewRistretto(t *testing.T) {
	cache, err := NewRistretto()
	if err != nil {
		return
	}
	// set a value with a cost of 1
	cache.Set("key", "value", 1)

	// wait for value to pass through buffers
	cache.Wait()

	// get value from cache
	value, found := cache.Get("key")
	if !found {
		panic("missing value")
	}
	fmt.Println(value)

	// del value from cache
	cache.Del("key")
	assert.Equal(t, nil, err)
}
