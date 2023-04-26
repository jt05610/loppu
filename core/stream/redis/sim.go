package redis

import (
	"context"
	"encoding/csv"
	"injector/bot"
	"os"
	"strconv"
	"time"
)

type Sim struct {
	Stream *Stream
	id     string
	data   string
}

func (s *Sim) Run(ctx context.Context) {
	err := s.Stream.Open(ctx)
	if err != nil {
		panic(err)
	}
	defer s.Stream.Close()
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
	data := &StreamItem{
		Id:    s.id,
		Force: 0,
		Depth: 1300,
	}
	for _, l := range ll[1:] {
		f, err := strconv.ParseFloat(l[1], 32)
		if err != nil {
			panic(err)
		}
		data.Force = float32(f)
		ok := s.Stream.redis.XAdd(ctx, s.Stream.Format(data))
		if ok.Err() != nil {
			panic(err)
		}
		time.Sleep(time.Duration(1) * time.Millisecond)
	}
}

func NewSim(id string, data string) *Sim {
	return &Sim{
		Stream: NewRedisStream("sim-injector", id, 375, time.Duration(100)*time.
			Millisecond, bot.NewInjector()),
		id:   id,
		data: data,
	}
}
