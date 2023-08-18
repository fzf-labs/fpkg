package channelx

import (
	"context"
	"sync"
)

// Channel 一些并发模型
type Channel[T any] struct {
}

// NewChannel 返回一个Channel实例
func NewChannel[T any]() *Channel[T] {
	return &Channel[T]{}
}

// Generate 创建通道，然后将值放入通道中。
func (c *Channel[T]) Generate(ctx context.Context, values ...T) <-chan T {
	dataStream := make(chan T)

	go func() {
		defer close(dataStream)

		for _, v := range values {
			select {
			case <-ctx.Done():
				return
			case dataStream <- v:
			}
		}
	}()

	return dataStream
}

// Repeat 创建通道，将值反复放入通道，直到取消上下文。
func (c *Channel[T]) Repeat(ctx context.Context, values ...T) <-chan T {
	dataStream := make(chan T)

	go func() {
		defer close(dataStream)
		for {
			for _, v := range values {
				select {
				case <-ctx.Done():
					return
				case dataStream <- v:
				}
			}
		}
	}()
	return dataStream
}

// RepeatFn 创建一个通道，重复执行fn，并将结果放入通道中
func (c *Channel[T]) RepeatFn(ctx context.Context, fn func() T) <-chan T {
	dataStream := make(chan T)

	go func() {
		defer close(dataStream)
		for {
			select {
			case <-ctx.Done():
				return
			case dataStream <- fn():
			}
		}
	}()
	return dataStream
}

// Take 创建一个通道，其值取自另一个有数量限制的通道。
func (c *Channel[T]) Take(ctx context.Context, valueStream <-chan T, number int) <-chan T {
	takeStream := make(chan T)

	go func() {
		defer close(takeStream)

		for i := 0; i < number; i++ {
			select {
			case <-ctx.Done():
				return
			case takeStream <- <-valueStream:
			}
		}
	}()

	return takeStream
}

// FanIn 扇入模型 将多个通道合并为一个通道。
func (c *Channel[T]) FanIn(ctx context.Context, channels ...<-chan T) <-chan T {
	out := make(chan T)

	go func() {
		var wg sync.WaitGroup
		wg.Add(len(channels))

		for _, ts := range channels {
			go func(c <-chan T) {
				defer wg.Done()
				for v := range c {
					select {
					case <-ctx.Done():
						return
					case out <- v:
					}
				}
			}(ts)
		}
		wg.Wait()
		close(out)
	}()

	return out
}

// Tee 将一个channel拆分为两个通道，直到取消上下文。
func (c *Channel[T]) Tee(ctx context.Context, in <-chan T) (ch1, ch2 <-chan T) {
	out1 := make(chan T)
	out2 := make(chan T)
	go func() {
		defer close(out1)
		defer close(out2)
		for val := range c.OrDone(ctx, in) {
			var out1, out2 = out1, out2
			for i := 0; i < 2; i++ {
				select {
				case <-ctx.Done():
				case out1 <- val:
					out1 = nil
				case out2 <- val:
					out2 = nil
				}
			}
		}
	}()
	return out1, out2
}

// Bridge 链接多个频道为一个频道。
func (c *Channel[T]) Bridge(ctx context.Context, chanStream <-chan <-chan T) <-chan T {
	valStream := make(chan T)

	go func() {
		defer close(valStream)

		for {
			var stream <-chan T
			select {
			case maybeStream, ok := <-chanStream:
				if !ok {
					return
				}
				stream = maybeStream
			case <-ctx.Done():
				return
			}

			for val := range c.OrDone(ctx, stream) {
				select {
				case valStream <- val:
				case <-ctx.Done():
				}
			}
		}
	}()

	return valStream
}

// Or 将一个或多个通道读入一个通道，当任何读入通道关闭时将关闭。
func (c *Channel[T]) Or(channels ...<-chan T) <-chan T {
	switch len(channels) {
	case 0:
		return nil
	case 1:
		return channels[0]
	}

	orDone := make(chan T)

	go func() {
		defer close(orDone)

		switch len(channels) {
		case 2:
			select {
			case <-channels[0]:
			case <-channels[1]:
			}
		default:
			select {
			case <-channels[0]:
			case <-channels[1]:
			case <-channels[2]:
			case <-c.Or(append(channels[3:], orDone)...):
			}
		}
	}()

	return orDone
}

// OrDone 读取一个通道到另一个通道，将关闭，直到取消上下文。
func (c *Channel[T]) OrDone(ctx context.Context, channel <-chan T) <-chan T {
	valStream := make(chan T)

	go func() {
		defer close(valStream)
		for {
			select {
			case <-ctx.Done():
				return
			case v, ok := <-channel:
				if !ok {
					return
				}
				select {
				case valStream <- v:
				case <-ctx.Done():
				}
			}
		}
	}()

	return valStream
}
