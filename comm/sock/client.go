package sock

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jt05610/loppu/comm"
	"net"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Client struct {
	socket string
	conn   net.Conn
	buf    []byte
	to     time.Duration
	ctx    context.Context
}

func (c *Client) Open(ctx context.Context) (err error) {
	c.conn, err = net.Dial("unix", c.socket)
	c.ctx = ctx
	return err
}

func (c *Client) Do(ctx context.Context, f func() (comm.Packet,
	error)) (p comm.Packet, err error) {
	ctx, cancel := context.WithTimeout(ctx, c.to)
	defer cancel()
	done := make(chan struct{})
	go func() {
		select {
		case <-ctx.Done():
			err = errors.New("timeout")
		case <-done:
			return

		}
	}()
	go func() {
		p, err = f()
		close(done)
	}()
	<-done
	return p, err
}

func (c *Client) Read(ctx context.Context) (p comm.Packet, err error) {
	return c.Do(ctx, func() (comm.Packet, error) {
		var ret comm.Packet
		d := json.NewDecoder(c.conn)
		return ret, d.Decode(&ret)
	})
}

func (c *Client) Write(ctx context.Context, p comm.Packet) error {
	_, err := c.Do(ctx, func() (comm.Packet, error) {
		e := json.NewEncoder(c.conn)
		return nil, e.Encode(p)
	})
	return err
}

// RoundTrip sends the given comm.Packet to the given comm.Address and returns resulting comm.Packet
func (c *Client) RoundTrip(ctx context.Context, p comm.Packet) (comm.Packet,
	error) {
	err := c.Write(ctx, p)
	if err != nil {
		return nil, err
	}
	return c.Read(ctx)
}

func (c *Client) Close() {
	_ = c.conn.Close()
}

const DefaultTimeout = time.Duration(1000) * time.Millisecond

func NewClient(addr string, to ...time.Duration) comm.Client {
	t := DefaultTimeout
	if to != nil {
		t = to[0]
	}
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

	socket := filepath.Join(dir, fmt.Sprintf("%s.sock", addr))
	c := &Client{socket: socket, buf: make([]byte, 1024), to: t}
	c.conn, err = net.Dial("unix", socket)
	if err != nil {
		panic(err)
	}
	return c
}
