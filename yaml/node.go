package yaml

import "io"

type Loader[T any] interface {
	Load(reader io.Reader) (T, error)
}

func (n *Node) Load(r io.Reader) (node.Node, error) {
	b, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	ret := &hardware.Node{}
	return ret, yaml.Unmarshal(b, ret)
}

func (n *Node) Flush(w io.Writer, node node.Node) error {
	bytes, err := yaml.Marshal(node)
	if err != nil {
		return err
	}
	_, err = w.Write(bytes)
	return err
}
