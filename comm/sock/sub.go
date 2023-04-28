package sock

import (
	"bytes"
	"context"
	"fmt"
	"github.com/jt05610/loppu/comm"
	"net"
	"os"
	"path/filepath"
)

type Sub struct {
	socket string
	h      Handler
	done   chan struct{}
	err    error
	dir    string
}

func (s *Sub) Close() {
	_ = os.Remove(s.socket)
}

func (s *Sub) Run() error {
	defer s.Close()
	<-s.done
	return s.err
}

func (s *Sub) Listen(ctx context.Context) error {
	err := s.Open(ctx)
	if err != nil {
		return err
	}
	return s.Run()
}

func (s *Sub) Open(ctx context.Context) error {
	s.done = make(chan struct{})
	_, err := s.Serve(ctx)
	if err != nil {
		opErr, ok := err.(*net.OpError)
		if ok {
			s.err = opErr.Err
		} else {
			s.err = err
		}
		return s.err
	}
	return os.Chmod(s.socket, os.ModeSocket|0622)
}

func (s *Sub) Serve(ctx context.Context) (<-chan struct{}, error) {
	srv, err := net.ListenPacket("unixgram", s.socket)
	fmt.Println(srv.LocalAddr().String())
	done := make(chan struct{})
	if err != nil {
		return nil, err
	}
	go func() {
		go func() {
			<-ctx.Done()
			close(done)
			_ = srv.Close()
			_ = os.Remove(s.socket)
		}()
		buf := make([]byte, 1024)
		for {
			n, addr, err := srv.ReadFrom(buf)
			if err != nil {
				return
			}
			var rsp bytes.Buffer
			err = s.h(&rsp, buf[:n])
			if err != nil {
				return
			}
			_, err = srv.WriteTo(rsp.Bytes(), addr)
			if err != nil {
				return
			}
		}
	}()
	return done, nil
}

func NewSub(addr string, h Handler) comm.Sub {
	p := filepath.Join("/tmp", fmt.Sprintf("%s_unix", addr))
	err := os.Mkdir(p, 0777)
	if err != nil {
		if os.IsExist(err) {
			err = os.RemoveAll(p)
			if err != nil {
				panic(err)
			}
			err = os.Mkdir(p, 0777)
			if err != nil {
				panic(err)
			}
		} else {
			panic(err)
		}
	}
	return &Server{
		h:      h,
		dir:    p,
		socket: filepath.Join(p, fmt.Sprintf("%s-sub.sock", addr)),
	}
}
