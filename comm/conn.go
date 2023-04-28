package comm

import (
	"context"
)

type Role string

const (
	ServerRole     Role = "server"
	ClientRole     Role = "client"
	PublisherRole  Role = "publisher"
	SubscriberRole Role = "subscriber"
	StreamerRole   Role = "streamer"
	ConsumerRole   Role = "consumer"
)

type Conn interface {
	Role() Role
}

type Client interface {
	Conn
	Read(ctx context.Context) (Packet, error)
	Write(ctx context.Context, p Packet) error
	RoundTrip(ctx context.Context, p Packet) (Packet, error)
	Close()
}

type Server interface {
	Conn
	Open(ctx context.Context) error
	Close()
	Serve(ctx context.Context) (<-chan struct{}, error)
	Listen(ctx context.Context) error
	Run() error
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
