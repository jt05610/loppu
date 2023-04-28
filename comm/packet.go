package comm

import (
	"bytes"
	"encoding/json"
	"io"
)

type Packet map[string]interface{}

func (p *Packet) JSON() io.Reader {
	var buf bytes.Buffer
	d := json.NewDecoder(&buf)
	err := d.Decode(p)
	if err != nil {
		buf.WriteString("")
	}
	return &buf
}
