package comm

type Role string

const (
	ServerRole     Role = "server"
	PublisherRole  Role = "publisher"
	SubscriberRole Role = "subscriber"
	StreamerRole   Role = "streamer"
	ConsumerRole   Role = "consumer"
)

type Client interface {
	Read(a Addr) (Packet, error)
	Write(a Addr, p Packet) error
	RoundTrip(a Addr, p Packet) (Packet, error)
	Close()
}

type Server interface {
	Listen() error
}
