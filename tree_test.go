package fidi

import (
	"fmt"
	"os"
	"testing"

	"github.com/ohzqq/fidi/fn"
	"github.com/ohzqq/fidi/tree"
	"go.yaml.in/yaml/v4"
)

func TestTree(t *testing.T) {
	//s := New(`tmp/video`)
	//s.Build()
	ft, err := NewFromBasePath(`testdata/video`)
	if err != nil {
		t.Fatal(err)
	}
	//list := NewList(tree)
	if g := len(ft.Children()); g != 3 {
		t.Fatalf("got %d, wanted %d\n", g, 3)
	}
	if dir := ft.Children()[0]; dir.Get("filename").(*fn.Filename).Basename != "index.html" {
		t.Errorf("got %#v\n", dir.Get("filename"))
	}
	b, err := ft.Filter(tree.FilterNodesByDepth(0))
	if err != nil {
		t.Fatal(err)
	}
	if g := len(b); g != 4 {
		for _, n := range b {
			t.Errorf("%#v, depth %#v\n", n.ID(), n.Depth())
		}
	}
	//if len(b) != 3 {
	//t.Errorf("%#v, depth %#v\n", b, tree.MaxDepth)
	//}
	filtered, err := ft.FilterByExt(".html", -1)
	if err != nil {
		t.Fatal(err)
	}
	if len(filtered) != 5 {
		t.Errorf("%#v, depth %#v\n", filtered, ft.MaxDepth)
	}

	mf, err := ft.FilterByMimetype("video", -1)
	if err != nil {
		t.Fatal(err)
	}
	if len(mf) != 3 {
		t.Errorf("%#v\n", mf)
	}

	walk := func(n tree.Node) error {
		//println(n.ID())
		return nil
	}
	err = ft.Walk(walk)
	if err != nil {
		t.Fatal(err)
	}

	//path := `/depth2/depth3/index.html`
	//n, err := tree.GetNodeByPath(path, false)
	//if err != nil {
	//t.Error(err)
	//}
	//if n.RelPath != `../../index.html` {
	//t.Errorf("%#v\n", n.RelativizePath())
	//}
}

func TestFilterDirFiles(t *testing.T) {
	ft, err := NewFromBasePath(`testdata/video`)
	if err != nil {
		t.Fatal(err)
	}
	filtered, err := ft.FilterByExt(".html", -1)
	if err != nil {
		t.Fatal(err)
	}
	if len(filtered) != 5 {
		t.Errorf("%#v, depth %#v\n", filtered[0].Filename(), ft.MaxDepth)
	}
	filtered = FilterDir(ft.Dir, false, FilterExt(".html"))
	if len(filtered) != 1 {
		t.Errorf("%#v, depth %#v\n", filtered[0].Filename(), ft.MaxDepth)
	}
}

func TestFilterDirFilesBasename(t *testing.T) {
	ft, err := NewFromBasePath(`testdata/video`)
	if err != nil {
		t.Fatal(err)
	}
	file := FindChildByBasename(ft.Dir, "index.html")
	if file.Filename().Basename != "index.html" {
		t.Errorf("%#v, %v\n", file.Filename(), file.ID())
	}
}

func TestTreeSerialize(t *testing.T) {
	t.SkipNow()
	ft, err := NewFromBasePath(`testdata/video`)
	if err != nil {
		t.Fatal(err)
	}
	err = tree.Walk(ft, tree.SortByLeavesFirst)
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
	err = tree.Walk(ft, tree.SortByLeavesFirst)
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
