package hardware

import (
	"context"
	"errors"
	"github.com/jt05610/loppu"
	"io"
)

type Packet interface {
	io.ReadWriter
	Addr() loppu.Addr
	Header() []byte
	Data() []byte
	CRC16() uint16
}

var ErrTimeout = errors.New("no response")

type Proto interface {
	Open() error
	Close()
	Request(ctx context.Context, addr loppu.Addr, data Packet) (Packet, error)
}
