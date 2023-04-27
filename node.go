package loppu

import (
	"io"
	"net/http"
)

type Addr interface {
	String() string
	Bytes() []byte
	Byte() byte
}

type MetaData interface {
}

type Node interface {
	Load(r io.Reader) error
	Flush(w io.Writer) error
	Addr() string
	Port() int
	Start() error
	Stop() error
	Meta() MetaData
	Register(srv *http.ServeMux)
	Endpoints(base string) []*Endpoint
}

type Loader[T any] interface {
	Load(r io.Reader) (*T, error)
}

type Flusher[T any] interface {
	Flush(w io.Writer, t *T) error
}

type LoadFlusher[T any] interface {
	Loader[T]
	Flusher[T]
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
