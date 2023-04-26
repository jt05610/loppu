package redis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"net/http"
	"time"
)

type Stream struct {
	SampleID  string        `yaml:"sample_id,omitempty"`
	Interval  time.Duration `yaml:"interval_ms"`
	Port      string        `yaml:"port"`
	Addr      string        `yaml:"addr"`
	Password  string        `yaml:"password,omitempty"`
	DB        int           `yaml:"db,omitempty"`
	Requests  []string      `yaml:"requests"`
	Device    string
	redis     *redis.Client
	consumeCh chan map[string]interface{}
	keepAlive bool
}

func (s *Stream) Open(ctx context.Context) error {
	s.redis = redis.NewClient(&redis.Options{
		Addr:     s.Addr,
		Password: s.Password,
		DB:       s.DB,
	})
	res, err := s.redis.Ping(ctx).Result()
	if err != nil {
		return err
	}
	if res != "PONG" {
		return errors.New("did not pong the ping")
	}
	return nil
}

func (s *Stream) Close() {
	_ = s.redis.Close()
}

func (s *Stream) doRequests(ctx context.Context) (*redis.XAddArgs, error) {
	data := map[string]interface{}{
		"id": s.SampleID,
	}
	for _, r := range s.Requests {
		select {
		case <-ctx.Done():
			return nil, errors.New("timeout")
		default:
			url := fmt.Sprintf("http://localhost:%v/%v", s.Port, r)
			resp, err := http.Get(url)
			var res map[string]interface{}
			d := json.NewDecoder(resp.Body)
			err = d.Decode(&res)
			if err != nil {
				panic(err)
			}
			data[r] = res["result"]
		}
	}
	return s.Format(data), nil
}

func (s *Stream) Format(d map[string]interface{}) *redis.XAddArgs {
	return &redis.XAddArgs{
		Stream: s.ID(),
		Values: d,
	}
}

func (s *Stream) Stream(ctx context.Context) error {
	err := s.Open(ctx)
	if err != nil {
		return err
	}
	defer s.Close()
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-time.After(s.Interval):
			args, err := s.doRequests(ctx)
			if err != nil {
				return err
			}
			s.redis.XAdd(ctx, args)
		}
	}
}

func (s *Stream) ID() string {
	return s.Device + ":" + s.SampleID
}

func (s *Stream) Consume(ctx context.Context) bool {
	if !s.keepAlive {
		s.consumeCh = make(chan map[string]interface{})
		s.keepAlive = true
		args := &redis.XReadArgs{
			Block:   time.Duration(1000) * time.Millisecond,
			Streams: []string{s.ID(), "$"},
		}
		go func() {
			defer func() {
				close(s.consumeCh)
				s.keepAlive = false
			}()
			for {
				select {
				case <-ctx.Done():
					return
				default:
					res, err := s.redis.XRead(ctx, args).Result()
					if err != nil {
						if err == context.Canceled || err == redis.ErrClosed {
							return
						}
						panic(err)
					}
					s.consumeCh <- res[0].Messages[0].Values
				}
			}
		}()
	}
	return s.keepAlive
}

func (s *Stream) Recv(ctx context.Context) map[string]interface{} {
	for {
		select {
		case <-ctx.Done():
			return nil
		case item := <-s.consumeCh:
			return item
		}
	}
}

func NewRedisStream(device, id string, interval time.Duration) *Stream {
	return &Stream{
		Device:    device,
		SampleID:  id,
		Interval:  interval,
		keepAlive: false,
	}
}

type Handler interface {
	Handle(ctx context.Context, data map[string]interface{})
}
