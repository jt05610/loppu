package comm

import "io"

type Packet interface {
	Header(h ...map[string][]string) io.Reader
	JSON() io.Reader
	Body() map[string]interface{}
	Error() error
}

type PacketService interface {
	Load(r io.Reader) (Packet, error)
	Flush(w io.Writer, p Packet) error
}
