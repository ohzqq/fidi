// Copyright Ivan Sukharev

// The MIT License (MIT)

// Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package fidi

import (
	"fmt"
	"io/fs"
	"log"
	"os"
)

const StartDepth = 1

type Tree interface {
	GetNode(int) (Tree, error)
	HasParents() bool
	Parents() []Tree
	HasChildren() bool
	Children() []Tree
	Leaves() []File
	Branches() []Tree
	Filter(filter Filter) []File
	Info() File
	fs.FS
}

func NewTree(path string) Tree {
	dir, err := NewDir(path)
	if err != nil {
		log.Fatal(err)
	}
	dir.Root = path

	err = dir.Scan(path, StartDepth, false)
	if err != nil {
		log.Fatal(err)
	}

	for i, _ := range dir.nodes {
		n, _ := dir.GetNode(i)
		node := n.(*Dir)
		node.nodes = dir.nodes
		node.id = i
		//node.Root = strings.TrimSuffix(path, "/")
	}

	return dir
}

func (list *Dir) Add(node Dir) {
	list.nodes = append(list.nodes, node)

	if list.Reverse == nil {
		list.Reverse = make(map[string]int)
	}

	list.Reverse[node.rel] = len(list.nodes) - 1
}

func (list *Dir) Scan(path string, depth int, ignoreErr bool) error {
	path += string(os.PathSeparator)

	node, fillErr := NewDir(path, list.Root)
	if fillErr != nil && !ignoreErr {
		return fillErr
	}

	node.Depth = depth

	list.Add(node)

	depth++

	for _, f := range node.entries {
		if f.IsDir() {
			n := path + f.Name()
			scanErr := list.Scan(n, depth, ignoreErr)
			if scanErr != nil && !ignoreErr {
				return scanErr
			}
		}
	}

	return nil
}

func (list Dir) GetNode(index int) (Tree, error) {
	if len(list.nodes) < index+1 {
		return &Dir{}, &NodeIndexDontExistsError{Index: index}
	}

	return &list.nodes[index], nil
}

type NodeIndexDontExistsError struct {
	Index int
}

func (e *NodeIndexDontExistsError) Error() string {
	return fmt.Sprintf("Node with index [%v] not exists", e.Index)
}
