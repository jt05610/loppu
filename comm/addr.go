package comm

import (
	"errors"
	"strconv"
	"strings"
)

// Addr is a generic address currently intended for IP or Modbus addresses.
type Addr []byte

// IP gives an IP representation of an Addr.
func (a Addr) IP() string {
	ret := ""
	for i, b := range a {
		if i > 0 {
			ret += "."
		}
		ret += strconv.Itoa(int(b))
	}
	return ret
}

// String gives a string representation of an Addr.
func (a Addr) String() string {
	return string(a)
}

// Bytes returns Addr as bytes.
func (a Addr) Bytes() []byte {
	return a
}

// Byte returns the first byte of an Addr.
func (a Addr) Byte() byte {
	return a[0]
}

// NewAddr returns an Addr from a string.
func NewAddr(s string) Addr {
	ret := make(Addr, 0)
	for _, t := range strings.Split(s, ".") {
		b, err := strconv.Atoi(t)
		if err != nil {
			panic(err)
		}
		if b > 0xFF {
			panic(errors.New("invalid address"))
		}
		ret = append(ret, byte(b))
	}
	return ret
}
