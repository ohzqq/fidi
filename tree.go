// Copyright Ivan Sukharev

// The MIT License (MIT)

// Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package dirfile

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

const StartDepth = 1

type Tree struct {
	File
	List    []TreeNode
	Reverse map[string]int
	exts    []string
}

type TreeNode struct {
	File
	Depth        int
	exts         []string
	Files        []File
	SubDirs      []File
	FilesCount   int
	SubDirsCount int
}

func NewTree(path string, exts ...string) Tree {
	tree := Tree{
		File: NewFile(path),
		exts: exts,
	}

	err := tree.Scan(path, StartDepth, true)
	if err != nil {
		log.Fatal(err)
	}

	for _, node := range tree.List {
		node.Root = path
		for _, file := range node.Files {
			file.Root = path
		}
	}

	return tree
}

func (list *Tree) Add(node TreeNode) {
	list.List = append(list.List, node)

	if list.Reverse == nil {
		list.Reverse = make(map[string]int)
	}

	list.Reverse[node.rel] = len(list.List) - 1
}

func (list *Tree) Scan(path string, depth int, ignoreErr bool) error {
	path += string(os.PathSeparator)

	node := TreeNode{
		exts: list.exts,
	}

	fillErr := node.Fill(path, depth)
	if fillErr != nil && !ignoreErr {
		return fillErr
	}

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

func (list *Tree) GetNode(index int) (*TreeNode, error) {
	if len(list.List) < index+1 {
		return nil, &NodeIndexDontExistsError{Index: index}
	}

	return &list.List[index], nil
}

func (node *TreeNode) Fill(path string, depth int) error {
	node.File = NewFile(path)

	f, err := os.Open(path)
	if err != nil {
		return err
	}
	entries, err := f.ReadDir(-1)
	if err != nil {
		println("entry error")
		return err
	}
	fmt.Printf("entries %+V\n", entries)
	filez, _ := f.Readdirnames(-1)
	fmt.Printf("filex %+V\n", filez)
	//if err != nil {
	//  println("file error")
	//  return err
	//}

	f.Close()

	files := make([]File, 0, len(entries))

	for _, entry := range entries {
		e := filepath.Join(path, entry.Name())
		n := NewFile(e)
		n.rel = e
		files = append(files, n)
	}

	node.rel = path
	node.Depth = depth

	for _, f := range files {
		if f.Stat.IsDir() {
			node.SubDirs = append(node.SubDirs, f)
		} else {
			node.Files = append(node.Files, f)
		}
	}

	node.FilesCount = len(node.Files)
	node.SubDirsCount = len(node.SubDirs)

	return nil
}

type NodeIndexDontExistsError struct {
	Index int
}

func (e *NodeIndexDontExistsError) Error() string {
	return fmt.Sprintf("Node with index [%v] not exists", e.Index)
}
