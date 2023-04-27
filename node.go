package loppu

import (
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// Addr is a generic address currently intended for IP or Modbus addresses.
type Addr []byte

// String gives a string representation of an Addr.
func (a Addr) String() string {
	ret := ""
	for i, b := range a {
		if i > 0 {
			ret += "."
		}
		ret += strconv.Itoa(int(b))
	}
	return ret
}

// Bytes returns Addr as bytes.
func (a Addr) Bytes() []byte {
	return a
}

// Byte returns the first byte of an Addr.
func (a Addr) Byte() byte {
	return a[0]
}

// MetaData represents common data needed for all Node instances.
type MetaData struct {
	// Node is the name of the node.
	Node string `yaml:"node"`
	// Author is the author of the node.
	Author string `yaml:"author"`
	// Version is the version of the node.
	Version Version `yaml:"version"`
	// Date is the date the node was created.
	Date time.Time `yaml:"date"`
	// Updated is the IP address the node can be reached at.
	Updated time.Time `yaml:"updated"`
	// Addr is the IP address the node can be reached at.
	Addr Addr `yaml:"addr"`
	// Port is the node's port.
	Port int `yaml:"port"`
}

// NewAddr returns an Addr from a string.
func NewAddr(s string) Addr {
	ret := make(Addr, 0)
	for _, t := range strings.Split(s, ".") {
		b, err := strconv.Atoi(t)
		if err != nil {
			panic(err)
		}
		if b > 0xFF {
			panic(errors.New("invalid address"))
		}
		ret = append(ret, byte(b))
	}

	return ret
}

// NewMetaData makes returns default MetaData for a new Node.
func NewMetaData(name string, addr string, port int) *MetaData {
	return &MetaData{
		Node:    name,
		Author:  Username(),
		Version: "0.1.0",
		Date:    time.Time{},
		Addr:    NewAddr(addr),
		Port:    50000 + port,
	}
}

func (m *MetaData) Patch() {

}

func (m *MetaData) Update() {

}

func (m *MetaData) Release() {

}

// Node is the base interface for everything in a Loppu robot.
type Node interface {
	// Meta returns the Node's MetaData.
	Meta() *MetaData
	// Register registers the node with an HTTP request multiplexer.
	Register(srv *http.ServeMux)
	// Run runs the Node in a single function.
	Run() error
}

// NodeService is an interface to manage Nodes
type NodeService interface {
	// New creates a new Node with a usable default configuration.
	New() (Node, error)
	// Delete deletes Node.
	Delete(Node) error
	// Load loads a Node from an io.Reader.
	Load(r io.Reader) (Node, error)
	// Flush writes a Node to an io.Writer.
	Flush(w io.Writer, node Node) error
}
