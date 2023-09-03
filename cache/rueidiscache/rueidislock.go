package rueidiscache

import (
	"context"

	"github.com/redis/rueidis"
	"github.com/redis/rueidis/rueidislock"
)

type Locker struct {
	client rueidis.Client
}

func NewLocker(client rueidis.Client) *Locker {
	return &Locker{client: client}
}

func (l *Locker) AutoLock(ctx context.Context, key string, do func() error) error {
	locker, err := rueidislock.NewLocker(rueidislock.LockerOption{
		ClientBuilder: func(option rueidis.ClientOption) (rueidis.Client, error) {
			return l.client, nil
		},
		FallbackSETPX: false, // redis>6.2
	})
	if err != nil {
		return err
	}
	defer locker.Close()
	ctx, cancel, err := locker.WithContext(ctx, key)
	if err != nil {
		return err
	}
	defer cancel()
	return do()
}
