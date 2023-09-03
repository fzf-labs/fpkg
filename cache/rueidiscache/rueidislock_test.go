package rueidiscache

import (
	"context"
	"fmt"
	"testing"

	"github.com/redis/rueidis"
)

func TestLocker_AutoLock(t *testing.T) {
	client, err := NewRueidis(rueidis.ClientOption{
		Username:    "",
		Password:    "123456",
		InitAddress: []string{"127.0.0.1:6379"},
		SelectDB:    0,
	})
	if err != nil {
		return
	}
	locker := NewLocker(client)
	ctx := context.Background()
	err = locker.AutoLock(ctx, "test_lock", func() error {
		fmt.Println("test_lock do ")
		return nil
	})
	if err != nil {
		return
	}
}
