package redis

import (
	"context"
	"errors"
	"github.com/redis/go-redis/v9"
	"injector/bot"
	"strconv"
	"time"
)

type Stream struct {
	Device     string
	redis      *redis.Client
	injector   *bot.Injector
	consumeCh  chan *StreamItem
	SampleID   string
	Home       uint16
	RunConsume bool
	Interval   time.Duration `yaml:"interval_ms,omitempty"`
}

func (s *Stream) Open(ctx context.Context) error {
	s.redis = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
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

type StreamItem struct {
	Id    string
	Force float32
	Depth uint16
}

func (s *StreamItem) Dict() map[string]interface{} {
	return map[string]interface{}{
		"id":       s.Id,
		"force_mN": s.Force,
		"depth_um": s.Depth,
	}
}

func (s *Stream) doRequests(ctx context.Context) (*redis.XAddArgs, error) {
	d := &StreamItem{Id: s.SampleID}
	for {
		select {
		case <-ctx.Done():
			return nil, errors.New("timeout")
		default:
			var err error
			d.Depth, err = s.injector.Needle.Client.GetCurrentPos()
			if err != nil {
				panic(err)
			}
			d.Depth -= s.Home
			d.Force = s.injector.Pump.Read()
			return s.Format(d), nil
		}
	}
}

func (s *Stream) Format(i *StreamItem) *redis.XAddArgs {
	return &redis.XAddArgs{
		Stream: s.ID(),
		Values: i.Dict(),
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
	if !s.RunConsume {
		s.consumeCh = make(chan *StreamItem)
		s.RunConsume = true
		args := &redis.XReadArgs{
			Block:   time.Duration(1000) * time.Millisecond,
			Streams: []string{s.ID(), "$"},
		}
		go func() {
			defer func() {
				close(s.consumeCh)
				s.RunConsume = false
			}()
			for {
				select {
				case <-ctx.Done():
					return
				default:
					res, err := s.redis.XRead(ctx, args).
						Result()
					if err != nil {
						if err == context.Canceled || err == redis.ErrClosed {
							return
						}
						panic(err)
					}
					v := res[0].Messages[0].Values
					force := v["force_mN"].(string)
					depth := v["depth_um"].(string)
					forceFloat, err := strconv.ParseFloat(force, 32)
					depthInt, err := strconv.Atoi(depth)
					if err != nil {
						panic(err)
					}
					s.consumeCh <- &StreamItem{
						Id:    v["id"].(string),
						Force: float32(forceFloat),
						Depth: uint16(depthInt),
					}
				}
			}
		}()
	}
	return s.RunConsume
}

func (s *Stream) Recv(ctx context.Context) *StreamItem {
	for {
		select {
		case <-ctx.Done():
			return nil
		case item := <-s.consumeCh:
			return item
		}
	}
}

func NewRedisStream(device, id string, home uint16,
	interval time.Duration,
	injector *bot.Injector) *Stream {
	return &Stream{
		Device:     device,
		injector:   injector,
		SampleID:   id,
		Home:       home,
		Interval:   interval,
		RunConsume: false,
	}
}

type Handler interface {
	Handle(ctx context.Context, data *StreamItem)
}
