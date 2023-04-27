package http

import (
	"github.com/jt05610/loppu/comm"
	"net/http"
)

// Server represents a http.Server with a multiplexer.
type Server struct {
	addr string
	mux  *http.ServeMux
}

// Listen is the main function run by a server.
func (s *Server) Listen() error {
	return http.ListenAndServe(s.addr, s.mux)
}

// NewServer creates a comm.Server at the given address and uses the given http.ServeMux to handle requests.
func NewServer(addr string, mux *http.ServeMux) comm.Server {
	return &Server{addr: addr, mux: mux}
}
