package comm

import (
	"context"
)

type Conn interface {
	Open(ctx context.Context) error
	Close()
}

type Client interface {
	Conn
	Read(ctx context.Context) (Packet, error)
	Write(ctx context.Context, p Packet) error
	RoundTrip(ctx context.Context, p Packet) (Packet, error)
}

type Server interface {
	Conn
	Serve(ctx context.Context) (<-chan struct{}, error)
	Listen(ctx context.Context) error
	Run() error
}

type Sub interface {
	Conn
	Serve(ctx context.Context) (<-chan struct{}, error)
	Listen(ctx context.Context) error
	Run() error
}

type Pub interface {
	Conn
	Read(ctx context.Context) (Packet, error)
	Write(ctx context.Context, p Packet) error
	RoundTrip(ctx context.Context, p Packet) (Packet, error)
}

type Consumer interface {
	Conn
	Consume(ctx context.Context) bool
	Recv() <-chan map[string]interface{}
}

type Streamer interface {
	Conn
	Stream(ctx context.Context)
}
