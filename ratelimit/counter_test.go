package ratelimit

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_counter_Allow(t *testing.T) {
	counter := NewCounter(10, time.Minute)
	i := 1
	for {
		allow := counter.Allow()
		fmt.Println(i, allow)
		if !allow {
			break
		}
		i++
	}
	assert.Equal(t, 11, i)
}
