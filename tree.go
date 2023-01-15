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
	"strings"
)

const StartDepth = 1

type Tree struct {
	File
	//Dir
	Children []Dir
	Parents  []Dir
	Nodes    []Dir
	Nodez    []Tree
	Reverse  map[string]int
	reverse  map[string]int
}

func NewTree(path string) Tree {
	//dir, err := NewDir(path)
	//if err != nil {
	//  log.Fatal(err)
	//}
	tree := Tree{
		//Dir: dir,
		File: NewFile(path),
	}

	err := tree.Scan(path, StartDepth, false)
	if err != nil {
		log.Fatal(err)
	}

	for _, node := range tree.Nodes {
		t := Tree{
			File:    node.File,
			Nodes:   tree.Nodes[node.Depth:],
			Parents: tree.Nodes[:node.Depth-1],
		}
		for _, n := range t.Nodes {
			if strings.Contains(n.Abs, node.Name) && node.Name != n.Name {
				t.Children = append(t.Children, n)
			}
		}
		tree.Nodez = append(tree.Nodez, t)
		node.Root = path
		for _, file := range node.Files {
			file.Root = path
		}
	}

	return tree
}

func (list *Tree) Add(node Dir) {
	list.Nodes = append(list.Nodes, node)

	if list.Reverse == nil {
		list.Reverse = make(map[string]int)
	}

	list.Reverse[node.rel] = len(list.Nodes) - 1
}

func (list *Tree) AddNode(node Tree) {
	list.Nodez = append(list.Nodez, node)

	if list.reverse == nil {
		list.reverse = make(map[string]int)
	}

	list.reverse[node.rel] = len(list.Nodez) - 1
}

func (list *Tree) Scan(path string, depth int, ignoreErr bool) error {
	path += string(os.PathSeparator)

	node, fillErr := NewDir(path)
	if fillErr != nil && !ignoreErr {
		return fillErr
	}
	//list.Parents = list.Nodes[:depth-1]
	node.Depth = depth

	list.Add(node)

	//tree.Parents = list.Nodes[:depth-1]

	depth++

	for _, f := range node.SubDirs {
		scanErr := list.Scan(path+f.Name, depth, ignoreErr)
		if scanErr != nil && !ignoreErr {
			return scanErr
		}
	}

	return nil
}

func (list *Tree) GetNode(index int) (*Dir, error) {
	if len(list.Nodes) < index+1 {
		return nil, &NodeIndexDontExistsError{Index: index}
	}

	return &list.Nodes[index], nil
}

func (tree Tree) GetNodesAtDepth(d int) []Dir {
	var nodes []Dir
	for _, node := range tree.Nodes {
		if node.Depth == d {
			nodes = append(nodes, node)
		}
	}

	return nodes
}

func (tree Tree) GetNodezAtDepth(d int) []Tree {
	var nodes []Tree
	for _, node := range tree.Nodez {
		if node.Depth == d {
			nodes = append(nodes, node)
		}
	}

	return nodes
}

func (tree Tree) HasParents() bool {
	return len(tree.Parents) > 0
}

func (tree Tree) GetChildren() []Dir {
	//fmt.Printf("tree path %s\n", tree.Abs)
	var nodes []Dir
	for _, node := range tree.Nodes[tree.Depth:] {
		//fmt.Printf("current node %d\n", node.Name)
		sub := tree.GetNodesAtDepth(node.Depth + 1)
		for _, s := range sub {
			if strings.Contains(s.Abs, node.Name) {
				nodes = append(nodes, s)
			}
		}
	}
	return nodes
}

func (tree Tree) GetChildrenByDepth(d int) []Dir {
	if d == 0 {
		return tree.Nodes
	}

	//cur, err := tree.GetNode(d - 1)
	//if err != nil {
	//  log.Fatal(err)
	//}

	var nodes []Dir
	for i := d + 1; i < len(tree.Nodes); i++ {
		n := tree.GetNodesAtDepth(i)
		nodes = append(nodes, n...)
	}

	//return Tree{
	//  File:  cur.File,
	//  Nodes: nodes,
	//}
	return nodes
}

func (tree Tree) GetParentsByDepth(d int) Tree {
	if d == 0 {
		return tree
	}

	cur, err := tree.GetNode(d - 1)
	if err != nil {
		log.Fatal(err)
	}

	var nodes []Dir
	for i := d + 1; i < len(tree.Nodes); i-- {
		n := tree.GetNodesAtDepth(i)
		nodes = append(nodes, n...)
	}

	return Tree{
		//Dir:   *cur,
		File:  cur.File,
		Nodes: nodes,
	}
}

type NodeIndexDontExistsError struct {
	Index int
}

func (e *NodeIndexDontExistsError) Error() string {
	return fmt.Sprintf("Node with index [%v] not exists", e.Index)
}
