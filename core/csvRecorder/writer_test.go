package csvRecorder

import (
	"context"
	"injector/softNode/stream/redis"
	"sync"
	"testing"
)

func TestCSVWriter_Handle(t *testing.T) {
	pass := false
	s := redis.NewSim("writer_test", "testData.csv")
	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup
	wg.Add(2)
	err := s.Stream.Open(ctx)
	defer s.Stream.Close()
	if err != nil {
		panic(err)
	}
	go func() {
		defer func() {
			wg.Done()
			cancel()
		}()
		s.Run(ctx)
	}()
	go func() {
		w := NewCSVWriter("out.csv", []string{"time", "force"})
		defer func() {
			w.(*CSVWriter).Close()
			pass = true
			wg.Done()
		}()
		for s.Stream.Consume(ctx) {
			select {
			case <-ctx.Done():
				return
			default:
				item := s.Stream.Recv(ctx)
				w.Handle(ctx, item)
			}
		}
	}()
	wg.Wait()
	if !pass {
		t.Fail()
	}
}
