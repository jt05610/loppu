package sock

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/jt05610/loppu/comm"
	"net"
	"os"
	"path/filepath"
	"strings"
)

type Pub struct {
	client net.PacketConn
	socket string
	srv    net.Addr

	ctx context.Context
}

func (p *Pub) Open(ctx context.Context) (err error) {
	p.client, err = net.ListenPacket("unixgram", p.socket)
	if err != nil {
		return err
	}
	p.ctx = ctx
	return os.Chmod(p.socket, os.ModeSocket|0622)
}

func (p *Pub) Close() {
	_ = p.client.Close()
	_ = os.Remove(p.socket)
}

func (p *Pub) Read(ctx context.Context) (comm.Packet, error) {
	buf := make([]byte, 1024)
	n, _, err := p.client.ReadFrom(buf)
	if err != nil {
		return nil, err
	}
	var ret comm.Packet
	d := json.NewDecoder(bytes.NewBuffer(buf[:n]))
	return ret, d.Decode(&ret)
}

func (p *Pub) Write(ctx context.Context, pack comm.Packet) error {
	var buf bytes.Buffer
	e := json.NewEncoder(&buf)
	err := e.Encode(pack)
	if err != nil {
		panic(err)
	}
	_, err = p.client.WriteTo(buf.Bytes(), p.srv)
	return err
}

func (p *Pub) RoundTrip(ctx context.Context, pack comm.Packet) (comm.Packet,
	error) {
	err := p.Write(ctx, pack)
	if err != nil {
		return nil, err
	}
	return p.Read(ctx)

}

func NewPub(addr string) comm.Pub {
	dd, err := os.ReadDir("/tmp")
	if err != nil {
		panic(err)
	}
	var dir string
	for _, d := range dd {
		if strings.Contains(d.Name(), fmt.Sprintf("%s_unix", addr)) {
			dir = filepath.Join("/tmp", d.Name())
			break
		}
	}
	srv := filepath.Join(dir, fmt.Sprintf("%s-sub.sock", addr))
	ua, err := net.ResolveUnixAddr("unixgram", srv)
	if err != nil {
		panic(err)
	}

	return &Pub{
		socket: filepath.Join(dir, fmt.Sprintf("%s-pub.sock", addr)),
		srv:    ua,
	}
}
