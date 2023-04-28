package redis

import (
	"context"
	"errors"
	"github.com/jt05610/loppu/comm"
	"github.com/redis/go-redis/v9"
	"log"
	"sync"
	"time"
)

type Request struct {
	Name string `yaml:"name"`
	Uri  string `yaml:"uri"`
}

type Stream struct {
	client   comm.Client
	Name     string        `yaml:"name"`
	SampleID string        `yaml:"sampleID"`
	Requests []*Request    `yaml:"requests"`
	Interval time.Duration `yaml:"interval"`
}

func (s *Stream) doRequests(ctx context.Context, c comm.Client) (*redis.XAddArgs, error) {
	data := map[string]interface{}{
		"id": s.SampleID,
	}
	for _, r := range s.Requests {
		select {
		case <-ctx.Done():
			return nil, errors.New("timeout")
		default:
			resp, err := c.Read(ctx)
			if err != nil {
				return nil, err
			}
			data[r.Name] = resp
		}
	}

	return &redis.XAddArgs{
		Stream: s.Name,
		Values: data,
	}, nil
}

type Streamer struct {
	reqClient comm.Client
	client    *redis.Client
	Streams   []*Stream `yaml:"streams"`
	ctx       context.Context
	cancel    context.CancelFunc
}

func (s *Streamer) Open(ctx context.Context) error {
	s.client = redis.NewClient(DefaultRedis)
	s.ctx, s.cancel = context.WithCancel(ctx)
	return nil
}

func (s *Streamer) Close() {
	s.cancel()
	_ = s.client.Close()
}

func (s *Stream) Stream(ctx context.Context, out chan *redis.XAddArgs) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-time.After(s.Interval):
			args, err := s.doRequests(ctx, s.client)
			if err != nil {
				return err
			}
			out <- args
		}
	}
}

func (s *Streamer) Stream(ctx context.Context) {
	var wg sync.WaitGroup
	ch := make(chan *redis.XAddArgs)
	defer close(ch)
	defer func() {
		_ = s.client.Close()
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				return
			case msg := <-ch:
				s.client.XAdd(ctx, msg)
			}
		}
	}()
	for _, stream := range s.Streams {
		wg.Add(1)
		go func(stream *Stream) {
			defer wg.Done()
			err := stream.Stream(ctx, ch)
			if err != nil {
				log.Fatal(err)
			}
		}(stream)
	}
	wg.Wait()
}

var DefaultRedis = &redis.Options{
	Addr:     "localhost:6379",
	Password: "",
	DB:       0,
}

func (s *Streamer) Add(new *Stream) error {
	if s.Streams == nil {
		s.Streams = make([]*Stream, 0)
	}
	s.Streams = append(s.Streams, new)
	return nil
}

func NewStreamer() comm.Streamer {
	return &Streamer{}
}
