package comm

import "io"

type Packet interface {
	io.ReadWriter
	Addr() Addr
	Header() []byte
	Data() []byte
	CRC16() uint16
}

type PacketService interface {
	Load(r io.Reader) (Packet, error)
	Flush(w io.Writer, p Packet) error
}
