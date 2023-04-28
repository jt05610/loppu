package loppu

import (
	"io"
)

// Node is the base interface for everything in a Loppu robot.
type Node interface {
	// Meta returns the Node's MetaData.
	Meta() *MetaData
	// Run runs the Node in a single function.
	Run() error
}

// NodeService is an interface to manage Nodes
type NodeService interface {
	// New creates a new Node with a usable default configuration.
	New() (Node, error)
	// Delete deletes Node.
	Delete(Node) error
	// Load loads a Node from an io.Reader.
	Load(r io.Reader) (Node, error)
	// Flush writes a Node to an io.Writer.
	Flush(w io.Writer, node Node) error
}
