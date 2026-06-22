package fidi

import (
	"testing"
)

func TestTree(t *testing.T) {
	//s := New(`tmp/video`)
	//s.Build()
	tree, err := NewFS(osFs.Fs, `tmp/video`)
	if err != nil {
		t.Fatal(err)
	}
	//list := NewList(tree)
	if g := len(tree.Children); g != 3 {
		t.Errorf("got %d, wanted %d\n", g, 3)
	}
	fn := func(node Node) error {
		if node.IsDir {
			t.Errorf("child %#v, parents %#v\n", node.Name, node.parents)
		}
		return nil
	}
	tree.Walk(fn)
	t.Errorf("%#v, depth %#v\n", tree.nodes, len(tree.Children))
	//for _, c := range tree.Children {
	//  //t.Errorf("child %#v, reverse %#v\n", c.Name, c.Reverse)
	//  if c.IsDir {
	//    t.Errorf("child %#v, leaves %#v\n", c.Name, len(c.Children))
	//  }
	//  for _, ch := range c.Children {
	//    if ch.IsDir {
	//      t.Errorf("child %#v, parents %#v, \n", ch.Name, len(ch.Children))
	//    }
	//  }
	//}

}
