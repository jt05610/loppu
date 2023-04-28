package http

import (
	"context"
	"github.com/jt05610/loppu/comm"
	"net/http"
	"time"
)

// Server represents a http.Server with a multiplexer.
type Server struct {
	addr string
	srv  http.Server
	done chan struct{}
	mux  *http.ServeMux
}

func (s *Server) Open(_ context.Context) error {
	s.done = make(chan struct{})
	s.srv = http.Server{
		Addr:        s.addr,
		Handler:     s.mux,
		ReadTimeout: time.Duration(10000) * time.Millisecond,
	}
	return nil
}

func (s *Server) Close() {
	_ = s.srv.Close()
}

func (s *Server) Serve(_ context.Context) (<-chan struct{}, error) {
	return nil, s.srv.ListenAndServe()
}

func (s *Server) Listen(_ context.Context) error {
	return s.srv.ListenAndServe()
}

func (s *Server) Run() error {
	//TODO implement me
	panic("implement me")
}

// NewServer creates a comm.Server at the given address and uses the given http.ServeMux to handle requests.
func NewServer(addr string, mux *http.ServeMux) comm.Server {
	return &Server{addr: addr, mux: mux}
}
