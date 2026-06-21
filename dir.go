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

func New(rootDir string) (Trunk, error) {
	result := Node{
		Name:         rootDir,
		RelativePath: "/",
		IsDir:        true,
	}
	err := walkDir(rootDir, result.RelativePath, &result)
	return Trunk{result}, err
}

func walkDirFs(fs fs.ReadDirFS, baseDir string, relativeDirstring, parent *Node) error {
	files, err := os.ReadDir(baseDir)
	if err != nil {
		return err
	}
	parent.nodes = make([]Node, len(files))
	for i, f := range files {
		parent.nodes[i].Name = f.Name()
		if f.IsDir() {
			parent.nodes[i].IsDir = true
			parent.nodes[i].RelativePath = filepath.Join(relativeDir, parent.nodes[i].Name)
			walkDir(filepath.Join(baseDir, parent.nodes[i].Name), parent.nodes[i].RelativePath, &parent.nodes[i])
		} else {
			parent.nodes[i].IsDir = false
			parent.nodes[i].RelativePath = relativeDir
		}
	}
	return nil
}

func walkDir(baseDir string, relativeDir string, parent *Node) error {
	files, err := os.ReadDir(baseDir)
	if err != nil {
		return err
	}
	parent.nodes = make([]Node, len(files))
	for i, f := range files {
		parent.nodes[i].Name = f.Name()
		if f.IsDir() {
			parent.nodes[i].IsDir = true
			parent.nodes[i].RelativePath = filepath.Join(relativeDir, parent.nodes[i].Name)
			walkDir(filepath.Join(baseDir, parent.nodes[i].Name), parent.nodes[i].RelativePath, &parent.nodes[i])
		} else {
			parent.nodes[i].IsDir = false
			parent.nodes[i].RelativePath = relativeDir
		}
	}
	return nil
}

func NewNode(path string, root ...string) (Node, error) {
	dir := Node{
		Path: path,
		fsys: os.DirFS(path),
	}
	dir.rel = path

	if len(root) > 0 {
		dir.Root = root[0]
	}

	entries, err := os.ReadDir(dir.Path)
	//entries, err := dir.ReadDir(".")
	if err != nil {
		return dir, err
	}
	dir.entries = entries

	return dir, err
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
