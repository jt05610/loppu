package wd

import (
	"context"
	"injector/softNode/stream/redis"
	"sync"
	"testing"
)

func TestWD_Handle(t *testing.T) {
	pass := false
	s := redis.NewSim("stream_test", "testData.csv")
	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup
	wg.Add(2)
	err := s.Stream.Open(ctx)
	if err != nil {
		panic(err)
	}
	go func() {
		defer func() {
			wg.Done()
		}()
		s.Run(ctx)
	}()
	go func() {
		wd := NewWatchDog(40, 5)
		defer func() {
			cancel()
			wg.Done()
		}()
		for s.Stream.Consume(ctx) {
			select {
			case <-ctx.Done():
				return
			default:
				item := s.Stream.Recv(ctx)
				wd.Handle(ctx, item)
				if wd.(*WD).Fired() {
					pass = true
				}
			}
		}
	}()
	wg.Wait()
	if !pass {
		t.Fail()
	}
}
