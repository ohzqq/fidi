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
	Dir
	Children []Dir
	Parents  []Dir
	nodes    []Dir
	Nodes    []Tree
	Reverse  map[string]int
	reverse  map[string]int
}

func NewTree(path string) Tree {
	dir, err := NewDir(path)
	if err != nil {
		log.Fatal(err)
	}
	tree := Tree{
		Dir: dir,
	}

	err = tree.Scan(path, StartDepth, false)
	if err != nil {
		log.Fatal(err)
	}

	for _, node := range tree.nodes {
		t := Tree{
			Dir:     node,
			nodes:   tree.nodes[node.Depth:],
			Parents: tree.nodes[:node.Depth-1],
		}
		for _, n := range t.nodes {
			if strings.Contains(n.Abs, node.Name) && node.Name != n.Name {
				t.Children = append(t.Children, n)
			}
		}
		tree.Nodes = append(tree.Nodes, t)
		node.Root = path
		for _, file := range node.Files {
			file.Root = path
		}
	}

	return tree
}

func (list *Tree) Add(node Dir) {
	list.nodes = append(list.nodes, node)

	if list.Reverse == nil {
		list.Reverse = make(map[string]int)
	}

	list.Reverse[node.rel] = len(list.nodes) - 1
}

func (list *Tree) Scan(path string, depth int, ignoreErr bool) error {
	path += string(os.PathSeparator)

	node, fillErr := NewDir(path)
	if fillErr != nil && !ignoreErr {
		return fillErr
	}
	node.Depth = depth

	list.Add(node)

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
	if len(list.nodes) < index+1 {
		return nil, &NodeIndexDontExistsError{Index: index}
	}

	return &list.nodes[index], nil
}

func (tree Tree) GetNodesAtDepth(d int) []Dir {
	var nodes []Dir
	for _, node := range tree.nodes {
		if node.Depth == d {
			nodes = append(nodes, node)
		}
	}

	return nodes
}

func (tree Tree) HasParents() bool {
	return len(tree.Parents) > 0
}

func (tree Tree) HasChildren() bool {
	return len(tree.Children) > 0
}

type NodeIndexDontExistsError struct {
	Index int
}

func (e *NodeIndexDontExistsError) Error() string {
	return fmt.Sprintf("Node with index [%v] not exists", e.Index)
}
