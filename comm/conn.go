package comm

import "io"

type Role string

const (
	Server     Role = "server"
	Client     Role = "client"
	Publisher  Role = "publisher"
	Subscriber Role = "subscriber"
	Streamer   Role = "streamer"
	Consumer   Role = "consumer"
)

type Conn interface {
	io.Reader
	io.Writer
	io.Closer
	Role() Role
}
