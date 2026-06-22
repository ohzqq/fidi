package fidi

import (
	"testing"
)

func TestTree(t *testing.T) {
	//s := New(`tmp/video`)
	//s.Build()
	tree, err := New(`tmp/video`)
	if err != nil {
		t.Fatal(err)
	}
	//list := NewList(tree)
	if g := len(tree.Children); g != 3 {
		t.Errorf("got %d, wanted %d\n", g, 3)
	}
	fn := func(node Node) error {
		//if node.IsDir {
		t.Errorf("name %#v, parents %#v, depth %d\n", node.Path(), node.Parents, node.Depth)
		//}
		return nil
	}
	tree.Walk(fn)
	//path := `depth2/scene002-clip001.mp4`
	b, err := tree.FilterExt(".html", true)
	if err != nil {
		t.Errorf("%#v, depth %#v\n", b, len(b))
	}
	t.Errorf("%#v, depth %#v\n", b, len(b))
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
