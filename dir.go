package fidi

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

type Filter func(File) bool

type Dir struct {
	File
	fsys     fs.FS
	entries  []os.DirEntry
	children []Tree
	parents  []Tree
	nodes    []Dir
	id       int
	Reverse  map[string]int
}

func NewDir(path string, root ...string) (Dir, error) {
	dir := Dir{
		File: NewFile(path),
		fsys: os.DirFS(path),
	}
	dir.rel = path

	if len(root) > 0 {
		dir.Root = root[0]
	}

	entries, err := os.ReadDir(dir.Path())
	if err != nil {
		return dir, err
	}
	dir.entries = entries

	return dir, err
}

func (node Dir) Info() File {
	return node.File
}

func (node Dir) Glob(pattern string) ([]string, error) {
	return fs.Glob(node.fsys, pattern)
}

func (node Dir) ReadFile(name string) ([]byte, error) {
	return fs.ReadFile(node.fsys, name)
}

func (node Dir) ReadDir(name string) ([]fs.DirEntry, error) {
	return fs.ReadDir(node.fsys, name)
}

func (node Dir) Open(name string) (fs.File, error) {
	return node.fsys.Open(name)
}

func (node Dir) Children() []Tree {
	var children []Tree
	if len(node.nodes) > 0 {
		nodes := node.nodes[node.id+1:]
		for _, sub := range nodes {
			if sub.Depth > node.Depth {
				children = append(children, sub)
			}
		}
	}
	return children
}

func (node Dir) Parents() []Tree {
	var parents []Tree
	if len(node.nodes) > 0 {
		nodes := node.nodes[:node.Depth-1]
		for _, parent := range nodes {
			parents = append(parents, parent)
		}
	}
	return parents
}

func (node Dir) Branches() []Tree {
	var dirs []Tree
	for _, f := range node.entries {
		if f.IsDir() {
			rel := filepath.Join(node.rel, f.Name())
			d := NewTree(rel)
			dirs = append(dirs, d)
		}
	}
	return dirs
}

func (node Dir) Leaves() []File {
	var files []File
	for _, e := range node.entries {
		if !e.IsDir() {
			rel := filepath.Join(node.rel, e.Name())
			file := NewFile(rel)
			file.Root = node.Root
			file.rel = rel
			files = append(files, file)
		}
	}
	return files
}

func (tree Dir) HasParents() bool {
	return len(tree.parents) > 0
}

func (tree Dir) HasChildren() bool {
	return len(tree.Children()) > 0
}

func (node Dir) Filter(filter Filter) []File {
	return FilterFiles(node.Leaves(), filter)
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
