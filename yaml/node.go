package yaml

import (
	"github.com/jt05610/loppu"
	"gopkg.in/yaml.v3"
	"io"
)

type NodeService[T any] struct {
}

func (s *NodeService[T]) Load(r io.Reader) (*T, error) {
	var ret T
	b, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	return &ret, yaml.Unmarshal(b, &ret)
}

func (s *NodeService[T]) Flush(w io.Writer, t *T) error {
	bytes, err := yaml.Marshal(t)
	if err != nil {
		return err
	}
	_, err = w.Write(bytes)
	return err
}

func NewYAMLLoader[T any]() loppu.LoadFlusher[T] {
	return &NodeService[T]{}
}
