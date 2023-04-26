package loppu

import (
	"os"
	"testing"
)

func makeProj(name string, t *testing.T) {
	err := InitProject("testing", name, true)
	if err != nil {
		t.Fail()
	}
	err = InitProject("testing", name, false)
	if err == nil {
		t.Fail()
	}
	if !os.IsExist(err) {
		t.Fail()
	}
}

func teardownProj(name string, t *testing.T) {
	err := os.RemoveAll("testing/" + name)
	if err != nil {
		t.Error(err)
	}
}

func TestInitProject(t *testing.T) {
	makeProj("initTest", t)
	teardownProj("initTest", t)
}

func TestProject_AddNode(t *testing.T) {
	makeProj("addNodeTest", t)
	p := &Project{}
	err := p.Load("testing/addNodeTest")
	if err != nil {
		t.Error(err)
	}
	err = p.NewHWNode("fakeMB")
	if err != nil {
		t.Error(err)
	}
	err = p.Flush("testing/addNodeTest/bot.yaml", false, true)
	if err != nil {
		t.Error(err)
	}
	// make sure we can't add again
	err = p.NewHWNode("fakeMB")
	if err != ErrNodeExists {
		t.Error(err)
	}
}
