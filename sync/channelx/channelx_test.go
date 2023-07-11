package channelx

import (
	"context"
	"testing"
)

func TestChannel_Generate(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ch := NewChannel[int]()
	stream := ch.Generate(ctx, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10)
	for v := range stream {
		t.Log(v)
	}
}

func TestChannel_Repeat(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ch := NewChannel[int]()
	stream := ch.Repeat(ctx, 1, 2, 3)
	for v := range stream {
		t.Log(v)
	}
}

func TestChannel_RepeatFn(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ch := NewChannel[int]()
	stream := ch.RepeatFn(ctx, func() int {
		return 1
	})
	for v := range stream {
		t.Log(v)
	}
}
