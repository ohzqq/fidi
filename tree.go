// Copyright Ivan Sukharev

// The MIT License (MIT)

// Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package fidi

import (
	"fmt"
	"log"
	"os"
)

const StartDepth = 1

type TreeI interface {
	//GetNode(int) (TreeI, error)
	HasParents() bool
	//Parents() []Dir
	HasChildren() bool
	Children() []Dir
	Leaves() []File
	Rel() string
	Filter(filter Filter) []File
}

type Tree struct {
	Dir
	//File
}

func NewTree(path string) TreeI {
	dir, err := NewDir(path)
	if err != nil {
		log.Fatal(err)
	}
	//tree := Tree{
	//  Dir: dir,
	//}

	//err = dir.Scan(path, StartDepth, false)
	//if err != nil {
	//  log.Fatal(err)
	//}

	//for _, n := range tree.nodes {
	//  if strings.Contains(n.Abs, tree.Name) && tree.Name != n.Name {
	//    tree.Children = append(tree.Children, n)
	//  }
	//}

	return dir
}

//func (list Tree) Children() Tree {
//  var tree Tree
//  if len(list.nodes) > 0 {
//    //first := list.nodes[1]
//    //name := filepath.Join(first.Root, first.Name)
//    //tree = NewTree(name)
//    tree.nodes = list.nodes
//    //tree.Nodes = list.Nodes[list.Depth:]
//    //tree = t
//  }
//  //tree.Nodes = list.Nodes
//  //tree.nodes = list.nodes

//  //if len(list.Nodes) > 0 {
//  //  tree.Nodes = list.Nodes[1:]
//  //}
//  return tree
//}

func (list *Dir) Add(node Dir) {
	list.Nodes = append(list.Nodes, node)

	if list.Reverse == nil {
		list.Reverse = make(map[string]int)
	}

	list.Reverse[node.rel] = len(list.Nodes) - 1
}

func (list *Dir) Scan(path string, depth int, ignoreErr bool) error {
	path += string(os.PathSeparator)

	node := Dir{
		File: NewFile(path),
	}
	node.rel = path

	fillErr := node.sort()
	if fillErr != nil && !ignoreErr {
		return fillErr
	}

	node.Depth = depth

	list.Add(node)

	depth++

	for _, f := range node.Sub() {
		n := path + f.Name
		scanErr := list.Scan(n, depth, ignoreErr)
		if scanErr != nil && !ignoreErr {
			return scanErr
		}
	}

	return nil
}

func (tree Dir) GetNodesAtDepth(d int) []Dir {
	var nodes []Dir
	for _, node := range tree.Nodes {
		if node.Depth == d {
			nodes = append(nodes, node)
		}
	}

	return nodes
}

func (tree Dir) HasParents() bool {
	return len(tree.parents) > 0
}

func (tree Dir) HasChildren() bool {
	return len(tree.Children()) > 0
}

type NodeIndexDontExistsError struct {
	Index int
}

func (e *NodeIndexDontExistsError) Error() string {
	return fmt.Sprintf("Node with index [%v] not exists", e.Index)
}
