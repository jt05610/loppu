package loppu

import (
	"github.com/jt05610/loppu/comm"
	"time"
)

// MetaData represents common data needed for all Node instances.
type MetaData struct {
	// Node is the name of the node.
	Node string `yaml:"node"`
	// Desc is the description of the node.
	Desc string `yaml:"desc"`
	// Author is the author of the node.
	Author string `yaml:"author"`
	// Version is the version of the node.
	Version Version `yaml:"version"`
	// Date is the date the node was created.
	Date time.Time `yaml:"date"`
	// Updated is the IP address the node can be reached at.
	Updated time.Time `yaml:"updated"`
	// Addr is the IP address the node can be reached at.
	Addr comm.Addr `yaml:"addr"`
	// Port is the node's port.
	Port int `yaml:"port"`
}

// Patch creates a patch and updates the metadata.
func (m *MetaData) Patch() {
	m.Version = m.Version.Update(Patch)
	m.Updated = time.Now()
}

// Update creates a minor update and updates the metadata.
func (m *MetaData) Update() {
	m.Version = m.Version.Update(Minor)
	m.Updated = time.Now()
}

// Release creates a major update and updates the metadata.
func (m *MetaData) Release() {
	m.Version = m.Version.Update(Major)
	m.Updated = time.Now()
}

// NewMetaData makes returns default MetaData for a new Node.
func NewMetaData(name string, addr string, port int) *MetaData {
	return &MetaData{
		Node:    name,
		Author:  Username(),
		Version: "0.1.0",
		Date:    time.Now(),
		Updated: time.Now(),
		Addr:    comm.NewAddr(addr),
		Port:    50000 + port,
	}
}
