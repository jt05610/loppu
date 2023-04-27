package redis

import (
	"context"
	"encoding/csv"
	"github.com/jt05610/loppu"
	"github.com/redis/go-redis/v9"
	"os"
	"strconv"
	"time"
)

type Sim struct {
	Stream loppu.Node
	id     string
	data   string
}

func (s *Sim) Run(ctx context.Context) {
	err := s.Stream.(*Node).Open(ctx)
	if err != nil {
		panic(err)
	}
	defer s.Stream.(*Node).Close()
	file, err := os.Open(s.data)
	if err != nil {
		panic(err)
	}
	rdr := csv.NewReader(file)
	rdr.Comma = '\t'
	ll, err := rdr.ReadAll()
	if err != nil {
		panic(err)
	}
	data := map[string]interface{}{
		"id":    s.id,
		"force": float32(0),
		"depth": 1300,
	}
	for _, l := range ll[1:] {
		f, err := strconv.ParseFloat(l[1], 32)
		if err != nil {
			panic(err)
		}
		data["force"] = float32(f)

		args := &redis.XAddArgs{
			Stream: "stream-sim",
			Values: data,
		}

		ok := s.Stream.(*Node).redis.XAdd(ctx, args)
		if ok.Err() != nil {
			panic(err)
		}
		time.Sleep(time.Duration(1) * time.Millisecond)
	}
}

func NewSim(id string, data string) *Sim {
	return &Sim{
		Stream: NewRedisNode(),
		id:     id,
		data:   data,
	}
}
