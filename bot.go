package loppu

import (
	"time"
)

type BotMetaData struct {
	Name    string    `yaml:"name"`
	Author  string    `yaml:"author"`
	Version Version   `yaml:"version"`
	Date    time.Time `yaml:"date"`
}

func NewBotMeta(name string) MetaData {
	return &BotMetaData{
		Name:    name,
		Author:  Username(),
		Version: "",
		Date:    time.Time{},
	}
}

type Project struct {
	MetaData
}

type Bot struct {
	Nodes []Node
}

func NewProject(name string) *Project {
	return &Project{MetaData: NewBotMeta(name)}
}

func (p *Project) Flush(file string, create bool, overwrite bool) error {
	panic("")
}

func (p *Project) Load(file string) error {

	panic("")
}
