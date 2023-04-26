package loppu

import (
	"errors"
	"github.com/jt05610/loppu/yaml"
	"os"
	"path"
	"time"
)

type BotMetaData struct {
	Name    string    `yaml:"name"`
	Author  string    `yaml:"author"`
	Version Version   `yaml:"version"`
	Date    time.Time `yaml:"date"`
}

func NewBotMeta(name string) *BotMetaData {
	return &BotMetaData{
		Name:    name,
		Author:  Username(),
		Version: "0.0.1",
		Date:    time.Now(),
	}
}

type Project struct {
	MetaData *BotMetaData        `yaml:"meta"`
	Nodes    map[string][]string `yaml:"nodes,omitempty"`
}

type Bot struct {
	Nodes []Node
}

func NewProject(name string) *Project {
	return &Project{
		MetaData: NewBotMeta(name),
		Nodes:    map[string][]string{},
	}
}

var ErrNodeExists = errors.New("node exists")

func (p *Project) addNode(kind, name string) error {
	if p.Nodes[kind] == nil {
		p.Nodes[kind] = make([]string, 0)
	}
	if len(p.Nodes[kind]) > 0 {
		for _, n := range p.Nodes[kind] {
			if n == name {
				return ErrNodeExists
			}
		}
	}
	p.Nodes[kind] = append(p.Nodes[kind], name)
	return nil
}

func (p *Project) NewSWNode(name string) error {
	return p.addNode("software", name)
}

func (p *Project) NewHWNode(name string) error {
	return p.addNode("hardware", name)
}

func InitProject(dest string, name string, overwrite bool) error {
	p := NewProject(name)
	parPath := path.Join(dest, name)
	cfgPath := path.Join(parPath, "bot.yaml")
	nodePath := path.Join(parPath, "nodes")
	err := os.Mkdir(parPath, 0777)
	if err != nil && (os.IsExist(err) && !overwrite) {
		return err
	}
	err = os.Mkdir(nodePath, 0777)
	if err != nil && (os.IsExist(err) && !overwrite) {
		return err
	}
	return yaml.FlushFile[Project](cfgPath, true, true, p)
}

func (p *Project) Flush(file string, create bool, overwrite bool) error {
	return yaml.FlushFile[Project](file, create, overwrite, p)
}

func (p *Project) Load(file string) error {
	r, err := yaml.LoadFile[Project](path.Join(file, "bot.yaml"))
	if err != nil {
		return err
	}
	p.MetaData = r.MetaData
	if r.Nodes == nil {
		r.Nodes = make(map[string][]string, 0)
	}
	p.Nodes = r.Nodes
	return nil
}
