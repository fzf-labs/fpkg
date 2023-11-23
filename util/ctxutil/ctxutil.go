package ctxutil

import (
	"context"
	"time"
)

func SleepContext(ctx context.Context, delay time.Duration) {
	timer := time.NewTimer(delay)
	select {
	case <-ctx.Done():
		if !timer.Stop() {
			<-timer.C
		}
	case <-timer.C:
	}
}
