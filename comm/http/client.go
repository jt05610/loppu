package http

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jt05610/loppu/comm"
	"net/http"
)

// Client represents the http client.
type Client struct {
	client *http.Client
	url    string
}

func (c *Client) Open(ctx context.Context) error {
	c.client = &http.Client{}
	return nil
}

func (c *Client) Read(ctx context.Context) (comm.Packet, error) {
	resp, err := c.client.Get(c.url)
	if err != nil {
		return nil, err
	}
	var res comm.Packet
	d := json.NewDecoder(resp.Body)
	return res, d.Decode(&res)
}

func (c *Client) Write(ctx context.Context, p comm.Packet) error {
	var err error
	resp, err := c.client.Post(c.url, "application/json", p.JSON())
	if err == nil && resp.StatusCode > 400 {
		err = errors.New(fmt.Sprintf("write failed with status code %v", resp.StatusCode))
	}
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
	c.client.CloseIdleConnections()
}

func NewClient(url string) comm.Client {
	return &Client{client: http.DefaultClient, url: url}
}
