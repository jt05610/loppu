package modbus_test

import (
	mb "github.com/jt05610/loppu/hardware/modbus"
	"testing"
)

func TestLoadFlushMBusNode(t *testing.T) {
	testNode := mb.NewMBusNode("fakeNode", 0xFE)
	fName := "fake_node.yaml"
	err := mb.FlushMBusNode(fName, true, true, testNode)
	if err != nil {
		t.Error(err)
	}
	load := mb.LoadMBusNode(fName)
	if load.Meta() == nil {
		t.Fail()
	}
	if len(load.Endpoints("")) != len(testNode.Endpoints("")) {
		t.Fail()
	}
}
