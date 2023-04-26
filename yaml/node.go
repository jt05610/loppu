package yaml

import (
	"gopkg.in/yaml.v3"
	"io"
	"os"
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

func LoadFile[T any](file string) (*T, error) {
	df, err := os.OpenFile(file, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return nil, err
	}
	l := &NodeService[T]{}
	return l.Load(df)
}

func FlushFile[T any](file string, create bool, overwrite bool, t *T) error {
	l := &NodeService[T]{}
	perm := 0
	if create {
		perm |= os.O_CREATE
	}
	if overwrite {
		perm |= os.O_WRONLY
	}
	df, err := os.OpenFile(file, perm, os.ModePerm)
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = df.Close()
	}()

	err = l.Flush(df, t)
	return err
}
