package redis

import (
	"context"
	"github.com/jt05610/loppu/comm"
	"github.com/redis/go-redis/v9"
	"time"
)

type Consumer struct {
	client  *redis.Client
	ctx     context.Context
	cancel  context.CancelFunc
	Streams []*Stream
	ch      chan comm.Packet
	doneCh  chan struct{}
}

func (c *Consumer) Open(ctx context.Context) error {
	c.client = redis.NewClient(DefaultRedis)
	c.ctx, c.cancel = context.WithCancel(ctx)
	return nil
}

func (c *Consumer) Close() {
	c.cancel()
	_ = c.client.Close()
}

func (c *Consumer) Done() <-chan struct{} {
	return c.doneCh
}

func (c *Consumer) Consume() {
	names := make([]string, 0)
	for _, s := range c.Streams {
		names = append(names, s.Name, "$")
	}
	c.doneCh = make(chan struct{})
	defer close(c.doneCh)
	go func() {

		go func() {
			<-c.ctx.Done()
			close(c.doneCh)
		}()
		c.ch = make(chan comm.Packet)
		args := &redis.XReadArgs{
			Block:   time.Duration(1000) * time.Millisecond,
			Streams: names,
		}
		go func() {
			defer func() {
				close(c.ch)
				close(c.doneCh)
			}()
			for {
				res, err := c.client.XRead(c.ctx, args).Result()
				if err != nil {
					if err == context.Canceled || err == redis.ErrClosed {
						return
					}
					panic(err)
				}
				c.ch <- res[0].Messages[0].Values
			}
		}()

	}()
}

func (c *Consumer) Recv() <-chan comm.Packet {
	return c.ch
}

func NewConsumer() comm.Consumer {
	return &Consumer{}
}
