package loppu

import (
	"io"
	"net/http"
)

type Node interface {
	Register(srv *http.ServeMux)
	Endpoints(base string) []*Endpoint
}

type NodeService interface {
	Load(r io.Reader) (Node, error)
	Flush(w io.Writer, node Node) error
}

type EndpointParam struct {
	Name        string
	NameCap     string
	Type        string
	Description string
	Tag         string
}

type Endpoint struct {
	Func        string         `yaml:"func"`
	Route       string         `yaml:"route"`
	Method      string         `yaml:"method"`
	Description string         `yaml:"description"`
	Param       *EndpointParam `yaml:"params,omitempty"`
}
