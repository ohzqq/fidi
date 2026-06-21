package fidi

import (
	"testing"

	"github.com/spf13/afero"
)

func TestTree(t *testing.T) {
	//s := New(`tmp/video`)
	//s.Build()
	tree, err := New(afero.NewIOFS(osFs.Fs), `tmp/video`)
	if err != nil {
		t.Fatal(err)
	}
	//list := NewList(tree)
	if g := len(tree.Children); g != 3 {
		t.Errorf("got %d, wanted %d\n", g, 3)
	}
	for _, c := range tree.Children {
		//t.Errorf("child %#v, reverse %#v\n", c.Name, c.Reverse)
		if c.IsDir {
			t.Errorf("child %#v, leaves %#v\n", c.Name, len(c.Parents))
		}
		for _, ch := range c.Children {
			if ch.IsDir {
				t.Errorf("child %#v, parents %#v, \n", ch.Name, len(ch.Parents))
			}
		}
	}
}
