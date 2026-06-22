package fidi

import (
	"testing"
)

func TestTree(t *testing.T) {
	//s := New(`tmp/video`)
	//s.Build()
	tree, err := NewFromBasePath(`tmp/video`)
	if err != nil {
		t.Fatal(err)
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
	b, err = tree.FilterExt(".html", true)
	if len(b) != 5 {
		t.Errorf("%#v, depth %#v\n", b, tree.MaxDepth)
	}

	path := `/depth2/depth3/index.html`
	n, err := tree.GetNodeByPath(path, false)
	if err != nil {
		t.Error(err)
	}
	t.Errorf("%#v\n", n.relPath())
}
