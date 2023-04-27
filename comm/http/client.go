package http

import (
	"errors"
	"fmt"
	"github.com/jt05610/loppu/comm"
	"net/http"
)

// Client represents the http client.
type Client struct {
	client *http.Client
	ps     comm.PacketService
}

// Read from the given comm.Addr.
func (c *Client) Read(a comm.Addr) (comm.Packet, error) {
	resp, err := c.client.Get(a.String())
	if err != nil {
		return nil, err
	}
	return c.ps.Load(resp.Body)
}

// Write the given comm.Packet to the given comm.Addr. Makes sure status code is not an error, otherwise ignores the response.
func (c *Client) Write(a comm.Addr, p comm.Packet) error {
	var err error
	resp, err := c.client.Post(a.String(), "application/json", p.JSON())
	if err == nil && resp.StatusCode > 400 {
		err = errors.New(fmt.Sprintf("write failed with status code %v", resp.StatusCode))
	}
	return err
}

// RoundTrip sends the given comm.Packet to the given comm.Address and returns resulting comm.Packet
func (c *Client) RoundTrip(a comm.Addr, p comm.Packet) (comm.Packet, error) {
	var err error
	resp, err := c.client.Post(a.String(), "application/json", p.JSON())
	if err == nil && resp.StatusCode > 400 {
		err = errors.New(fmt.Sprintf("write failed with status code %v", resp.StatusCode))
	}
	return c.ps.Load(resp.Body)
}

func (c *Client) Close() {
	c.client.CloseIdleConnections()
}

func NewClient() comm.Client {
	return &Client{client: http.DefaultClient, ps: NewPacketService()}
}
