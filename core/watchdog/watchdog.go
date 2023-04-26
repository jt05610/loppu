package watchdog

import (
	"context"
	"injector/softNode/stream/redis"
	"time"
)

type MetaData struct {
	Name   string    `yaml:"name"`
	Author string    `yaml:"author"`
	Date   time.Time `yaml:"date"`
}

type WD struct {
	over    int
	limit   float32
	maxOver int
	fired   bool
}

func (w *WD) Handle(ctx context.Context, data *redis.StreamItem) {
	select {
	case <-ctx.Done():
		return
	default:
		if !w.fired {
			if data.Force > w.limit {
				w.over++
			} else {
				w.over = 0
			}
			if w.over > w.maxOver {
				w.fired = true
				w.over = 0
			}
		}
	}
}

func (w *WD) Fired() bool {
	return w.fired
}

func NewWatchDog(limit float32, maxOver int) redis.Handler {
	return &WD{over: 0, limit: limit, maxOver: maxOver}
}
