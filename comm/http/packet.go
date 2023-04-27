package http

import (
	"bytes"
	"encoding/json"
	"github.com/jt05610/loppu/comm"
	"io"
	"net/http"
)

// Packet represents an HTTP request or response.
type Packet struct {
	Err     error
	Hdr     http.Header
	Content map[string]interface{}
}

// Error returns any errors that occurred with the packet
func (p *Packet) Error() error {
	return p.Err
}

// Header returns the packet's header
func (p *Packet) Header(h ...map[string][]string) io.Reader {
	if h != nil {
		p.Hdr = h[0]
	}
	var buf bytes.Buffer
	if err := p.Hdr.Write(&buf); err != nil {
		panic(err)
	}
	return &buf
}

// Body returns the packet's header
func (p *Packet) Body() map[string]interface{} {
	return p.Content
}

// JSON returns the packet's content as a JSON reader.
func (p *Packet) JSON() io.Reader {
	var buf bytes.Buffer
	d := json.NewDecoder(&buf)
	err := d.Decode(&p.Content)
	if err != nil && err != io.EOF {
		panic(err)
	}
	return &buf
}

// PacketService facilitates conversion of http packets to comm.Packet.
type PacketService struct {
}

// Load loads a packet from an http Body.
func (p *PacketService) Load(r io.Reader) (comm.Packet, error) {
	res := &Packet{}
	d := json.NewDecoder(r)
	err := d.Decode(&res.Content)
	if err != nil && err != io.EOF {
		res.Err = err
	}
	return res, err
}

// Flush writes the Packet to the given io.Writer.
func (p *PacketService) Flush(w io.Writer, pack comm.Packet) error {
	e := json.NewEncoder(w)
	return e.Encode(pack.Body)
}

// NewPacketService creates a new packet service.
func NewPacketService() comm.PacketService {
	return &PacketService{}
}
