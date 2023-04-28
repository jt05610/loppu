package sock_test

import (
	"context"
	"github.com/jt05610/loppu/comm"
	"github.com/jt05610/loppu/comm/sock"
	"io"
	"testing"
)

func TestCliSrv(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	srv := sock.NewServer("echo", func(w io.Writer, r []byte) error {
		_, err := w.Write(r)
		return err
	})
	err := srv.Open(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer srv.Close()
	go func() {
		err := srv.Run()
		if err != nil {
			t.Fail()
		}
	}()
	cli := sock.NewClient("echo")
	msg := comm.Packet{"msg": "ping"}
	rsp, err := cli.RoundTrip(ctx, msg)
	if err != nil {
		t.Fatal(err)
	}
	if len(rsp) != len(msg) {
		t.Fail()
	}
	for k, v := range rsp {
		if msg[k] != v {
			t.Fail()
		}
	}
}
