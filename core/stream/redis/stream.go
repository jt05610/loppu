package redis

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/jt05610/loppu"
	"github.com/jt05610/loppu/yaml"
	"github.com/redis/go-redis/v9"
	"io"
	"log"
	"net/http"
	"sync"
	"syscall"
	"time"
)

type Request struct {
	Name string `yaml:"name"`
	Uri  string `yaml:"uri"`
}

type Stream struct {
	Name     string `yaml:"name"`
	ID       string
	Requests []*Request    `yaml:"requests"`
	Interval time.Duration `yaml:"interval"`
}

func (s *Stream) doRequests(ctx context.Context) (*redis.XAddArgs, error) {
	data := map[string]interface{}{
		"id": s.ID,
	}
	for _, r := range s.Requests {
		select {
		case <-ctx.Done():
			return nil, errors.New("timeout")
		default:
			resp, err := http.Get(r.Uri)
			var res map[string]interface{}
			d := json.NewDecoder(resp.Body)
			err = d.Decode(&res)
			if err != nil {
				panic(err)
			}
			data[r.Name] = res["result"]
		}
	}
	return &redis.XAddArgs{
		Stream: s.Name,
		Values: data,
	}, nil
}

func (s *Stream) Stream(ctx context.Context, out chan *redis.XAddArgs) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-time.After(s.Interval):
			args, err := s.doRequests(ctx)
			if err != nil {
				return err
			}
			out <- args
		}
	}
}

type MetaData struct {
	Node   string `yaml:"node"`
	Author string `yaml:"author"`
	Date   string `yaml:"date"`
	Host   string `yaml:"host"`
	Port   int    `yaml:"port"`
}

type Node struct {
	MetaData  *MetaData `yaml:"meta"`
	RedisAddr string    `yaml:"redis"`
	Password  string    `yaml:"password,omitempty"`
	DB        int       `yaml:"db,omitempty"`
	Streams   []*Stream `yaml:"streams"`
	redis     *redis.Client
	consumeCh chan map[string]interface{}
	keepAlive bool
}

func (s *Node) Addr() string {
	return s.MetaData.Host
}

func (s *Node) Port() int {
	return s.MetaData.Port
}

func (s *Node) Start() error {
	ctx := context.Background()
	s.Stream(ctx)
	return nil
}

func (s *Node) Stop() error {
	return syscall.Exec("kill", []string{"stream.pid"}, nil)
}

func (s *Node) Load(r io.Reader) error {
	l := yaml.NodeService[Node]{}
	v, err := l.Load(r)
	s.MetaData = v.MetaData
	s.DB = v.DB
	s.RedisAddr = v.RedisAddr
	s.Password = v.Password
	if v.Streams == nil {
		v.Streams = make([]*Stream, 0)
	}
	s.Streams = v.Streams
	return err
}

func (s *Node) Flush(w io.Writer) error {
	l := yaml.NodeService[Node]{}
	return l.Flush(w, s)

}

func (s *Node) Meta() loppu.MetaData {
	return s.MetaData
}

func (s *Node) Open(ctx context.Context) error {
	s.redis = redis.NewClient(&redis.Options{
		Addr:     s.RedisAddr,
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

func (s *Node) Close() {
	_ = s.redis.Close()
}

func (s *Node) Stream(ctx context.Context) {
	var wg sync.WaitGroup
	ch := make(chan *redis.XAddArgs)
	defer close(ch)
	err := s.Open(ctx)
	if err != nil {
		panic(err)
	}
	defer s.Close()
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				return
			case msg := <-ch:
				s.redis.XAdd(ctx, msg)
			}
		}
	}()
	for _, s := range s.Streams {
		wg.Add(1)
		go func(s *Stream) {
			defer wg.Done()
			e := s.Stream(ctx, ch)
			if e != nil {
				log.Fatal(e)
			}
		}(s)
	}
	wg.Wait()
}

func (s *Node) Consume(ctx context.Context) bool {
	names := make([]string, 0)
	for _, s := range s.Streams {
		names = append(names, s.Name, "$")
	}

	if !s.keepAlive {
		s.consumeCh = make(chan map[string]interface{})
		s.keepAlive = true
		args := &redis.XReadArgs{
			Block:   time.Duration(1000) * time.Millisecond,
			Streams: names,
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

func (s *Node) Recv(ctx context.Context) map[string]interface{} {
	for {
		select {
		case <-ctx.Done():
			return nil
		case item := <-s.consumeCh:
			return item
		}
	}
}

func NewStream(name, id string, delay time.Duration, reqs []*Request) *Stream {
	return &Stream{
		Name:     name,
		ID:       id,
		Requests: reqs,
		Interval: delay,
	}
}

func (s *Node) AddStream(stream *Stream) {
	s.Streams = append(s.Streams, stream)
}

func (s *Node) Register(srv *http.ServeMux) {
	//TODO implement me
	panic("implement me")
}

func (s *Node) Endpoints(base string) []*loppu.Endpoint {
	//TODO implement me
	panic("implement me")
}

type Handler interface {
	Handle(ctx context.Context, data map[string]interface{})
}

func NewRedisNode() loppu.Node {
	return &Node{
		MetaData: &MetaData{
			Node:   "stream",
			Author: loppu.Username(),
			Date:   time.Now().String(),
			Host:   "127.0.0.1",
			Port:   50001,
		},
		RedisAddr: "localhost:6379",
		Password:  "",
		DB:        0,
		Streams:   nil,
	}
}
