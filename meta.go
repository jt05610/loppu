package loppu

import (
	"context"
	"errors"
	"io"
	"strconv"
	"strings"
	"time"
)

// Addr is a generic address currently intended for IP or Modbus addresses.
type Addr []byte

// IP gives an IP representation of an Addr.
func (a Addr) IP() string {
	ret := ""
	for i, b := range a {
		if i > 0 {
			ret += "."
		}
		ret += strconv.Itoa(int(b))
	}
	return ret
}

// String gives a string representation of an Addr.
func (a Addr) String() string {
	return string(a)
}

// Bytes returns Addr as bytes.
func (a Addr) Bytes() []byte {
	return a
}

// Byte returns the first byte of an Addr.
func (a Addr) Byte() byte {
	return a[0]
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

type Packet interface {
	io.ReadWriter
	Addr() Addr
	Header() []byte
	Data() []byte
	CRC16() uint16
}

var ErrTimeout = errors.New("no response")

type ProtoRole string

const (
	Server     ProtoRole = "server"
	Client     ProtoRole = "client"
	Publisher  ProtoRole = "publisher"
	Subscriber ProtoRole = "subscriber"
	Streamer   ProtoRole = "streamer"
	Consumer   ProtoRole = "consumer"
)

type Proto interface {
	Role() ProtoRole
	Open() error
	Close()
}

type Requester interface {
	Request(ctx context.Context, addr Addr, data Packet) (Packet, error)
}

type Handler interface {
	Handle(ctx context.Context, w io.Writer, p Packet)
}

// MetaData represents common data needed for all Node instances.
type MetaData struct {
	// Node is the name of the node.
	Node string `yaml:"node"`
	// Desc is the description of the node.
	Desc string `yaml:"desc"`
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
	// Role is the primary public role played by the node in the network.
	Role ProtoRole `yaml:"role"`
}

// NewMetaData makes returns default MetaData for a new Node.
func NewMetaData(name string, addr string, port int) *MetaData {
	return &MetaData{
		Node:    name,
		Author:  Username(),
		Version: "0.1.0",
		Date:    time.Now(),
		Updated: time.Now(),
		Addr:    NewAddr(addr),
		Port:    50000 + port,
	}
}

func (m *MetaData) Patch() {
	m.Version = m.Version.Update(Patch)
	m.Updated = time.Now()
}

func (m *MetaData) Update() {
	m.Version = m.Version.Update(Minor)
	m.Updated = time.Now()

}

func (m *MetaData) Release() {
	m.Version = m.Version.Update(Major)
	m.Updated = time.Now()
}
