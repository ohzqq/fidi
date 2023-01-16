package fidi

import (
	"log"
	"os"
	"path/filepath"
	"strings"
)

type Filter func(File) bool

type Dir struct {
	File
	All []File
	//Files   []File
	SubDirs []File
	//Sub          []Dir
	FilesCount   int
	SubDirsCount int
	//Children     []Dir
	parents []Dir
	Nodes   []Dir
	id      int
	Reverse map[string]int
	//Tree
}

func NewDir(path string) (Dir, error) {
	//println(path)
	dir := Dir{
		File: NewFile(path),
	}
	dir.rel = path

	err := dir.sort()

	err = dir.Scan(path, StartDepth, false)
	if err != nil {
		log.Fatal(err)
	}

	for i, _ := range dir.Nodes {
		n, _ := dir.GetNode(i)
		node := n.(*Dir)
		node.Nodes = dir.Nodes
		node.id = i
		node.Root = path
		for _, file := range node.Leaves() {
			file.Root = path
		}
	}

	return dir, err
}

func (node *Dir) sort() error {
	entries, err := os.ReadDir(node.Path())
	if err != nil {
		return err
	}

	node.All = make([]File, 0, len(entries))

	for _, entry := range entries {
		e := filepath.Join(node.rel, entry.Name())
		n := NewFile(e)
		n.rel = e
		node.All = append(node.All, n)
	}

	//node.FilesCount = len(node.Files())
	//node.SubDirsCount = len(node.Sub())

	return nil
}

func (node Dir) Children() []Dir {
	var children []Dir
	nodes := node.Nodes[node.id+1:]
	for _, sub := range nodes {
		if sub.Depth > node.Depth {
			children = append(children, sub)
		}
	}
	return children
}

func (node Dir) Parents() []Dir {
	parents := node.Nodes[:node.Depth-1]
	return parents
}

func (node Dir) Sub() []Dir {
	var dirs []Dir
	for _, f := range node.All {
		if f.Stat.IsDir() {
			d, _ := NewDir(f.rel)
			dirs = append(dirs, d)
		}
	}
	return dirs
}

func (node Dir) Leaves() []File {
	var files []File
	for _, f := range node.All {
		if !f.Stat.IsDir() {
			files = append(files, f)
		}
	}
	return files
}

func (node Dir) Filter(filter Filter) []File {
	return FilterFiles(node.Leaves(), filter)
}

func (list Dir) GetNode(index int) (TreeI, error) {
	if len(list.Nodes) < index+1 {
		return &Dir{}, &NodeIndexDontExistsError{Index: index}
	}

	return &list.Nodes[index], nil
}

func FilterFiles(files []File, filter Filter) []File {
	var keep []File
	for _, fn := range files {
		if filter(fn) {
			keep = append(keep, fn)
		}
	}
	return keep
}

func ExtFilter(exts ...string) Filter {
	filter := func(file File) bool {
		for _, ex := range exts {
			if strings.EqualFold(file.Extension, ex) {
				return true
			}
		}
		return false
	}
	return filter
}

func MimeFilter(mimes ...string) Filter {
	filter := func(file File) bool {
		for _, mt := range mimes {
			if strings.Contains(file.Mime, mt) {
				return true
			}
		}
		return false
	}
	return filter
}
