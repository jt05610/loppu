package hardware

import "github.com/jt05610/loppu"

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

type Node interface {
	loppu.Node
	Proto(p ...Proto) Proto
	Endpoints(base string) []*Endpoint
}
