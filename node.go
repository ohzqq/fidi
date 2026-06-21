package fidi

import (
	"fmt"
	"io/fs"
	"path/filepath"

	"github.com/spf13/afero"
)

var osFs = afero.Afero{Fs: afero.NewOsFs()}

type Node struct {
	Depth        int
	Path         string
	Root         string
	Name         string
	RelativePath string
	IsDir        bool
	Children     []Node
}

func New(fs fs.ReadDirFS, rootDir string) (Node, error) {
	node := NewNode(rootDir, 0)
	node.RelativePath = "/"
	node.IsDir = true
	err := walkDirFs(fs, rootDir, node.RelativePath, &node)
	if err != nil {
		return node, err
	}
	return node, err
}

func NewNode(name string, depth int) Node {
	return Node{
		Name:  name,
		Depth: depth,
	}
}

func (n Node) Leaves() []Node {
	nodes := []Node{}
	for _, n := range n.Children {
		if !n.IsDir {
			nodes = append(nodes, n)
		}
	}
	return nodes
}

func (n Node) Branches() []Node {
	nodes := []Node{}
	for _, n := range n.Children {
		if n.IsDir {
			nodes = append(nodes, n)
		}
	}
	return nodes
}

func walkDirFs(fs fs.ReadDirFS, baseDir string, relativeDir string, parent *Node) error {
	files, err := fs.ReadDir(baseDir)
	if err != nil {
		return err
	}
	parent.Children = make([]Node, len(files))
	for i, f := range files {
		parent.Children[i] = NewNode(f.Name(), parent.Depth+1)
		if !f.IsDir() {
			parent.Children[i].IsDir = false
			parent.Children[i].RelativePath = relativeDir
			continue
		}
		parent.Children[i].IsDir = true
		parent.Children[i].RelativePath = filepath.Join(relativeDir, parent.Children[i].Name)
		walkDirFs(fs, filepath.Join(baseDir, parent.Children[i].Name), parent.Children[i].RelativePath, &parent.Children[i])
	}
	return nil
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

type NodeIndexDontExistsError struct {
	Index int
}

func (e *NodeIndexDontExistsError) Error() string {
	return fmt.Sprintf("Node with index [%v] not exists", e.Index)
}
