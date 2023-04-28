package http_test

import (
	"context"
	"github.com/jt05610/loppu/comm"
	. "github.com/jt05610/loppu/comm/http"
	"io"
	"net/http"
	"testing"
)

func TestClientServer(t *testing.T) {
	addr := ":60000"
	m := http.NewServeMux()
	s := NewServer(addr, m)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	m.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, err := io.Copy(w, r.Body)
		if err != nil {
			panic(err)
		}
	})

	go func() {
		err := s.Listen(ctx)
		if err != nil {
			panic(err)
		}
	}()

	c := NewClient("http://127.0.0.1" + addr)
	p := comm.Packet{
		"test":    1,
		"tester":  2,
		"testest": 3,
	}
	r, err := c.RoundTrip(ctx, p)
	if err != nil {
		t.Fatal(err)
	}
	expect, err := io.ReadAll(p.JSON())
	if err != nil {
		t.Fatal(err)
	}
	actual, err := io.ReadAll(r.JSON())
	if err != nil {
		t.Fatal(err)
	}
	if len(expect) != len(actual) {
		t.Fail()
	}
	for i := 0; i < len(expect); i++ {
		if actual[i] != expect[i] {
			t.Fail()
		}
	}
}
