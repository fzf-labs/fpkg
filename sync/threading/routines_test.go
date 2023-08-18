package threading

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRoutineID(t *testing.T) {
	assert.True(t, RoutineID() > 0)
}

func TestRunSafe(t *testing.T) {
	i := 0
	defer func() {
		assert.Equal(t, 1, i)
	}()
	ch := make(chan struct{})
	go RunSafe(func() {
		defer func() {
			ch <- struct{}{}
		}()

		panic("panic")
	})
	<-ch
	i++
}

func TestRunSafeCtx(t *testing.T) {
	ctx := context.Background()
	ch := make(chan struct{})

	i := 0

	defer func() {
		assert.Equal(t, 1, i)
	}()

	go RunSafeCtx(ctx, func() {
		defer func() {
			ch <- struct{}{}
		}()

		panic("panic")
	})

	<-ch
	i++
}

func TestGoSafeCtx(t *testing.T) {
	ctx := context.Background()
	ch := make(chan struct{})

	i := 0

	defer func() {
		assert.Equal(t, 1, i)
	}()

	GoSafeCtx(ctx, func() {
		defer func() {
			ch <- struct{}{}
		}()

		panic("panic")
	})

	<-ch
	i++
}
