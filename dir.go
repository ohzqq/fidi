package fidi

import (
	"io/fs"
	"os"
	"path/filepath"
)

type Node struct {
	Depth    int
	Path     string
	Root     string
	rel      string
	fsys     fs.FS
	entries  []os.DirEntry
	children []Tree
	parents  []Tree
	nodes    []Node
	id       int
	Reverse  map[string]int
}

func NewDir(path string, root ...string) (Node, error) {
	dir := Node{
		Path: path,
		fsys: os.DirFS(path),
	}
	dir.rel = path

	if len(root) > 0 {
		dir.Root = root[0]
	}

	//entries, err := os.ReadDir(dir.Path())
	entries, err := dir.ReadDir(".")
	if err != nil {
		return dir, err
	}
	dir.entries = entries

	return dir, err
}

func (n Node) Glob(pattern string) ([]string, error) {
	return fs.Glob(n.fsys, pattern)
}

func (n Node) ReadFile(name string) ([]byte, error) {
	return fs.ReadFile(n.fsys, name)
}

func (n Node) ReadDir(name string) ([]fs.DirEntry, error) {
	return fs.ReadDir(n.fsys, name)
}

func (n Node) Open(name string) (fs.File, error) {
	return n.fsys.Open(name)
}

func (n Node) Children() []Tree {
	var children []Tree
	if len(n.nodes) > 0 {
		nodes := n.nodes[n.id+1:]
		for _, sub := range nodes {
			if sub.Depth > n.Depth {
				children = append(children, sub)
			}
		}
	}
	return children
}

func (n Node) Parents() []Tree {
	var parents []Tree
	if len(n.nodes) > 0 && n.Depth > 0 {
		nodes := n.nodes[:n.Depth-1]
		for _, parent := range nodes {
			parents = append(parents, parent)
		}
	}
	return parents
}

func (n Node) Branches() []Tree {
	var dirs []Tree
	for _, f := range n.entries {
		if f.IsDir() {
			rel := filepath.Join(n.rel, f.Name())
			d := NewTree(rel)
			dirs = append(dirs, d)
		}
	}
	return dirs
}

func (n Node) Leaves() []string {
	var files []string
	for _, e := range n.entries {
		if !e.IsDir() {
			rel := filepath.Join(n.rel, e.Name())
			files = append(files, rel)
		}
	}
	return files
}

func (n Node) HasParents() bool {
	return len(n.Parents()) > 0
}

func (n Node) HasChildren() bool {
	return len(n.Children()) > 0
}

//func (n Node) Filter(filters ...Filter) []File {
//  return FilterFiles(n.Leaves(), filters...)
//}

//func FilterFiles(files []File, filters ...Filter) []File {
//  re := make(map[string]File)
//  for _, filter := range filters {
//    for _, fn := range files {
//      if filter(fn) {
//        re[fn.Name] = fn
//      }
//    }
//  }

//  var keep []File
//  for _, file := range re {
//    keep = append(keep, file)
//  }

//  sort.Slice(keep, func(i, j int) bool {
//    return keep[i].Name < keep[j].Name
//  })

//  return keep
//}

//func ExtFilter(exts ...string) Filter {
//  filter := func(file File) bool {
//    for _, ex := range exts {
//      if strings.EqualFold(file.Extension, ex) {
//        return true
//      }
//    }
//    return false
//  }
//  return filter
//}

//func MimeFilter(mimes ...string) Filter {
//  filter := func(file File) bool {
//    for _, mt := range mimes {
//      if strings.Contains(file.Mime, mt) {
//        return true
//      }
//    }
//    return false
//  }
//  return filter
//}
