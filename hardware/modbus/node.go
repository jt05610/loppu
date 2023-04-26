package modbus

import (
	"github.com/jt05610/loppu"
)

type ParamType string

type Param struct {
	Type        ParamType
	Description string `yaml:"desc"`
}

type Handler struct {
	Name        string              `yaml:"name"`
	Description string              `yaml:"desc"`
	Params      []map[string]*Param `yaml:"params,omitempty"`
}

type MBusNode struct {
	loppu.MetaData `yaml:"meta"`
	Tables         map[string][]*Handler `yaml:"tables"`
	Diag           []*Handler            `yaml:"diag"`
	Client         *Client
	rfLookup       map[string]map[string]func(uint16, uint16) *MBusPDU
	addrLookup     map[string]uint16
	paramLookup    map[string]string
}
