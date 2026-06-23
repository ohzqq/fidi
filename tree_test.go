package fidi

import (
	"fmt"
	"os"
	"testing"

	"github.com/ohzqq/fidi/tree"
	"go.yaml.in/yaml/v4"
)

func TestTree(t *testing.T) {
	//s := New(`tmp/video`)
	//s.Build()
	tree, err := NewFromBasePath(`testdata/video`)
	if err != nil {
		t.Fatal(err)
	}
	//list := NewList(tree)
	if g := len(tree.Children()); g != 3 {
		t.Fatalf("got %d, wanted %d\n", g, 3)
	}
	dir := tree.Children()[0]
	t.Errorf("got %#v\n", dir)
	//b, err := tree.GetNodesAtDepth(1)
	//if err != nil {
	//t.Errorf("%#v, depth %#v\n", b, len(b))
	//}
	//if len(b) != 3 {
	//t.Errorf("%#v, depth %#v\n", b, tree.MaxDepth)
	//}
	//b, err = tree.FilterByExt(".html", true)
	//if len(b) != 5 {
	//t.Errorf("%#v, depth %#v\n", b, tree.MaxDepth)
	//}

	//path := `/depth2/depth3/index.html`
	//n, err := tree.GetNodeByPath(path, false)
	//if err != nil {
	//t.Error(err)
	//}
	//if n.RelPath != `../../index.html` {
	//t.Errorf("%#v\n", n.RelativizePath())
	//}
}

func TestTreeSerialize(t *testing.T) {
	t.SkipNow()
	ft, err := NewFromBasePath(`testdata/video`)
	if err != nil {
		t.Fatal(err)
	}
	err = ft.Walk(tree.SortLeavesFirst)
	if err != nil {
		t.Fatal(err)
	}
	dump, err := yaml.NewDumper(os.Stdout)
	if err != nil {
		t.Fatal(err)
	}
	dump.Dump(ft)
}

func TestSortNodes(t *testing.T) {
	ft, err := NewFromBasePath(`testdata/video`)
	if err != nil {
		t.Fatal(err)
	}
	err = ft.Walk(tree.SortLeavesFirst)
	if err != nil {
		t.Fatal(err)
	}
}

//func TestFilterNodeDepth(t *testing.T) {
//  ft, err := NewFromBasePath(`testdata/video`)
//  if err != nil {
//    t.Fatal(err)
//  }
//  nodes, err := ft.FilterByDepth(1)
//  if err != nil {
//    t.Fatal(err)
//  }

//  for _, node := range nodes {
//    printNode(node)
//  }
//  //err = ft.Walk(printNode)
//  //if err != nil {
//  //t.Fatal(err)
//  //}
//}

func printNode(n tree.Node) error {
	fmt.Printf("%#v\n", n.Depth())
	return nil
}
