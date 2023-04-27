package redis

import (
	"context"
	"fmt"
	"sync"
	"testing"
)

func TestSim_Run(t *testing.T) {
	s := NewSim("stream_test", "testData.csv")
	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup
	wg.Add(2)
	err := s.Stream.(*Node).Open(ctx)
	if err != nil {
		panic(err)
	}
	go func() {
		defer func() {
			wg.Done()
			cancel()
		}()
		s.Run(context.Background())
	}()

	go func() {
		defer wg.Done()
		for s.Stream.(*Node).Consume(ctx) {
			select {
			case <-ctx.Done():
				return
			default:
				item := s.Stream.(*Node).Recv(ctx)
				if item != nil {
					fmt.Printf("%v\n", item)
				}
			}
		}
	}()
	wg.Wait()
}
