package comm

import (
	"context"
)

type Role string

const (
	ServerRole     Role = "server"
	PublisherRole  Role = "publisher"
	SubscriberRole Role = "subscriber"
	StreamerRole   Role = "streamer"
	ConsumerRole   Role = "consumer"
)

type Client interface {
	Read(addr string) (Packet, error)
	Write(addr string, p Packet) error
	RoundTrip(a Addr, p Packet) (Packet, error)
	Close()
}

type Server interface {
	Listen() error
}

type Consumer interface {
	Consume(ctx context.Context) bool
	Recv() <-chan map[string]interface{}
}

type Streamer interface {
	Stream(ctx context.Context)
}
