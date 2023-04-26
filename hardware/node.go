package hardware

import "github.com/jt05610/loppu"

type Node interface {
	loppu.Node
	Proto(p ...Proto) Proto
}
