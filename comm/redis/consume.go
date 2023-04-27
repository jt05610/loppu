package redis

import (
	"context"
	"github.com/jt05610/loppu/comm"
	"github.com/redis/go-redis/v9"
	"time"
)

type Consumer struct {
	client    *redis.Client
	Streams   []*Stream
	keepAlive bool
	ch        chan map[string]interface{}
}

func (c *Consumer) Consume(ctx context.Context) bool {
	names := make([]string, 0)
	for _, s := range c.Streams {
		names = append(names, s.Name, "$")
	}

	if !c.keepAlive {
		c.ch = make(chan map[string]interface{})
		c.keepAlive = true
		args := &redis.XReadArgs{
			Block:   time.Duration(1000) * time.Millisecond,
			Streams: names,
		}
		go func() {
			defer func() {
				close(c.ch)
				c.keepAlive = false
			}()
			for {
				select {
				case <-ctx.Done():
					return
				default:
					res, err := c.client.XRead(ctx, args).Result()
					if err != nil {
						if err == context.Canceled || err == redis.ErrClosed {
							return
						}
						panic(err)
					}
					c.ch <- res[0].Messages[0].Values
				}
			}
		}()
	}
	return c.keepAlive
}

func (c *Consumer) Recv() <-chan map[string]interface{} {
	return c.ch
}

func NewConsumer() comm.Consumer {
	return &Consumer{client: redis.NewClient(DefaultRedis)}
}
