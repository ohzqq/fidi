package fidi

import (
	"os"
	"testing"

	"go.yaml.in/yaml/v4"
)

func TestTree(t *testing.T) {
	//s := New(`tmp/video`)
	//s.Build()
	tree, err := NewFromBasePath(`testdata/video`)
	if err != nil {
		t.Fatal(err)
	}
	for _, c := range tree.Children {
		for _, ch := range c.Children {
			t.Errorf("abs %#v\n", ch.Name)
		}
	}
	//list := NewList(tree)
	if g := len(tree.Children); g != 3 {
		t.Errorf("got %d, wanted %d\n", g, 3)
	}
	b, err := tree.GetNodesAtDepth(1)
	if err != nil {
		t.Errorf("%#v, depth %#v\n", b, len(b))
	}
	if len(b) != 3 {
		t.Errorf("%#v, depth %#v\n", b, tree.MaxDepth)
	}
	b, err = tree.FilterByExt(".html", true)
	if len(b) != 5 {
		t.Errorf("%#v, depth %#v\n", b, tree.MaxDepth)
	}

	path := `/depth2/depth3/index.html`
	n, err := tree.GetNodeByPath(path, false)
	if err != nil {
		t.Error(err)
	}
	if n.RelPath != `../../index.html` {
		t.Errorf("%#v\n", n.RelativizePath())
	}
}

func TestTreeSerialize(t *testing.T) {
	t.SkipNow()
	tree, err := NewFromBasePath(`testdata/video`)
	if err != nil {
		t.Fatal(err)
	}
	dump, err := yaml.NewDumper(os.Stdout)
	if err != nil {
		t.Fatal(err)
	}
	dump.Dump(tree)
}
