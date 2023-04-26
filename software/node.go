package software

import (
	"github.com/jt05610/loppu"
	"net/http"
)

type Node struct {
}

func (n *Node) Register(srv *http.ServeMux) {
	//TODO implement me
	panic("implement me")
}

func (n *Node) Endpoints(base string) []*loppu.Endpoint {
	//TODO implement me
	panic("implement me")
}

func NewSoftwareNode() loppu.Node {
	return &Node{}
}
