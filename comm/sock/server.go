package sock

import (
	"context"
	"fmt"
	"github.com/jt05610/loppu/comm"
	"io"
	"log"
	"net"
	"os"
	"path/filepath"
)

type Handler func(w io.Writer, r []byte) error

type Server struct {
	socket string
	h      Handler
	done   chan struct{}
	err    error
	dir    string
}

func (s *Server) Close() {
	if err := os.RemoveAll(s.dir); err != nil {
		log.Fatal(err)
	}
}

func (s *Server) Run() error {
	defer s.Close()
	<-s.done
	return s.err
}

func (s *Server) Listen(ctx context.Context) error {
	err := s.Open(ctx)
	if err != nil {
		return err
	}
	return s.Run()
}

func (s *Server) Open(ctx context.Context) error {
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
	return os.Chmod(s.socket, os.ModeSocket|0666)
}

func (s *Server) Role() comm.Role {
	return comm.ServerRole
}

func (s *Server) Serve(ctx context.Context) (<-chan struct{}, error) {
	srv, err := net.Listen("unix", s.socket)
	fmt.Println(srv.Addr().String())
	done := make(chan struct{})
	if err != nil {
		return nil, err
	}
	go func() {
		go func() {
			<-ctx.Done()
			close(done)
			_ = srv.Close()
		}()

		for {
			conn, err := srv.Accept()
			if err != nil {
				return
			}

			go func() {
				defer func() {
					_ = conn.Close()
				}()
				for {
					buf := make([]byte, 1024)
					n, err := conn.Read(buf)
					if err != nil {
						return
					}
					err = s.h(conn, buf[:n])
					if err != nil {
						return
					}
				}
			}()
		}
	}()

	return done, nil
}

func NewServer(addr string, h Handler) comm.Server {
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
		socket: filepath.Join(p, fmt.Sprintf("%s.sock", addr)),
	}
}
