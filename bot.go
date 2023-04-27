package loppu

import (
	"errors"
	"github.com/jt05610/loppu/metadata"
	"net/http"
)

type Bot struct {
	MetaData *metadata.MetaData  `yaml:"metadata"`
	Nodes    map[string][]string `yaml:"nodes,omitempty"`
	nodes    []Node
}

func (b *Bot) Meta() *metadata.MetaData {
	//TODO implement me
	panic("implement me")
}

func (b *Bot) Register(srv *http.ServeMux) {
	//TODO implement me
	panic("implement me")
}

func (b *Bot) Run() error {
	//TODO implement me
	panic("implement me")
}

const DefaultName string = "newBot"

func NewProject(name string, host string) Node {
	return &Bot{
		MetaData: metadata.NewMetaData(name, host, 0),
		Nodes:    make(map[string][]string, 0),
	}
}

var ErrNodeExists = errors.New("node exists")

func (b *Bot) addNode(kind, name string) error {
	if b.Nodes[kind] == nil {
		b.Nodes[kind] = make([]string, 0)
	}
	if len(b.Nodes[kind]) > 0 {
		for _, n := range b.Nodes[kind] {
			if n == name {
				return ErrNodeExists
			}
		}
	}
	b.Nodes[kind] = append(b.Nodes[kind], name)
	return nil
}

func (b *Bot) NewSWNode(name string) error {
	return b.addNode("software", name)
}

func (b *Bot) NewHWNode(name string) error {
	return b.addNode("hardware", name)
}
