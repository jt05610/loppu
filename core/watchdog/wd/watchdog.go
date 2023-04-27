package wd

import (
	"context"
	"github.com/jt05610/loppu"
	"github.com/jt05610/loppu/core/stream/redis"
	"github.com/jt05610/loppu/yaml"
	"io"
	"net/http"
	"time"
)

type MetaData struct {
	Node   string    `yaml:"node"`
	Author string    `yaml:"author"`
	Date   time.Time `yaml:"date"`
	Host   string    `yaml:"host"`
	Port   int       `yaml:"port"`
}

type Reaction struct {
	Endpoint string                 `yaml:"endpoint"`
	Data     map[string]interface{} `yaml:"data"`
}

type WD struct {
	MetaData     *MetaData `yaml:"meta"`
	Watch        string    `yaml:"watch"`
	React        *Reaction `yaml:"react"`
	Limit        float32   `yaml:"limit"`
	CountTrigger int       `yaml:"countTrigger"`
	stream       *redis.Node
	fired        bool
	over         int
}

func (w *WD) Load(r io.Reader) error {
	l := yaml.NodeService[WD]{}
	v, err := l.Load(r)
	w.MetaData = v.MetaData
	w.React = v.React
	w.Limit = v.Limit
	w.Watch = v.Watch
	w.CountTrigger = v.CountTrigger
	return err
}

func (w *WD) Flush(wr io.Writer) error {
	l := yaml.NodeService[WD]{}
	return l.Flush(wr, w)

}

func (w *WD) Addr() string {
	return w.MetaData.Host
}

func (w *WD) Port() int {
	return w.MetaData.Port
}

func (w *WD) Start() error {
	ctx := context.Background()
	w.stream = redis.NewRedisNode().(*redis.Node)
	err := w.stream.Open(ctx)
	if err != nil {
		return err
	}
	for w.stream.Consume(ctx) {
		select {
		case <-ctx.Done():
			return nil
		case msg := <-w.stream.Recv():
			w.Handle(ctx, msg)
			if w.fired {
				return nil
			}
		}
	}
	return nil
}

func (w *WD) Stop() error {
	//TODO implement me
	panic("implement me")
}

func (w *WD) Meta() loppu.MetaData {
	return w.MetaData
}

func (w *WD) Register(srv *http.ServeMux) {
	srv.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

	})
}

func (w *WD) Endpoints(base string) []*loppu.Endpoint {
	//TODO implement me
	panic("implement me")
}

func (w *WD) Handle(ctx context.Context, data map[string]interface{}) {
	select {
	case <-ctx.Done():
		return
	default:
		if !w.fired {
			if data[w.Watch].(float32) > w.Limit {
				w.over++
			} else {
				w.over = 0
			}
			if w.over > w.CountTrigger {
				w.fired = true
				w.over = 0
			}
		}
	}
}

func (w *WD) Fired() bool {
	return w.fired
}

func NewWatchDog(watch string, limit float32, maxOver int) loppu.Node {
	return &WD{
		MetaData: &MetaData{
			Node:   "watchdog",
			Author: loppu.Username(),
			Date:   time.Now(),
			Host:   "127.0.0.1",
			Port:   50002,
		},
		Watch:        watch,
		React:        nil,
		Limit:        limit,
		CountTrigger: maxOver,
		fired:        false,
		over:         0,
	}
}
