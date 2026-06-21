package fidi

import "testing"

func TestSite(t *testing.T) {
	//s := New(`tmp/video`)
	//s.Build()
	tree, err := New(`../tmp/video`)
	if err != nil {
		t.Fatal(err)
	}
	//list := NewList(tree)
	for _, c := range tree.nodes {
		t.Errorf("child %#v\n", c.Name)
		if len(c.nodes) > 0 {
			for _, ch := range c.nodes {
				t.Errorf("child %#v\n", ch.Name)
			}
		}
	}
}
