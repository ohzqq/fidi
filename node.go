package fidi

import (
	"fmt"
	"path/filepath"
	"strings"
)

type Node struct {
	Depth    int
	Name     string
	Dir      string
	Path     string
	RelPath  string
	IsDir    bool
	Parents  []string
	Children []Node
}

func NewNode(name string, depth int) Node {
	return Node{
		Name:  name,
		Depth: depth,
	}
}

func (n Node) path() string {
	return filepath.Join(n.Dir, n.Name)
}

func (n Node) relPath() string {
	parts := strings.Split(strings.TrimPrefix(n.Dir, "/"), "/")
	dots := make([]string, len(parts))
	for i := range parts {
		dots[i] = ".."
	}
	dots = append(dots, n.Name)
	return filepath.Join(dots...)
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

func (n Node) WalkNode(fn func(node Node) error) error {
	err := fn(n)
	if err != nil {
		return err
	}
	for _, c := range n.Children {
		err := c.WalkNode(fn)
		if err != nil {
			return err
		}
	}
	return nil
}

func (n Node) GetNodeByPath(path string, dir bool) (Node, error) {
	branch := n
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	fn := func(node Node) error {
		for _, c := range node.Children {
			if dir && !c.IsDir {
				continue
			}
			if c.Path == path {
				branch = c
				return nil
			}
		}
		return nil
	}
	err := n.WalkNode(fn)
	if err != nil {
		return n, err
	}
	return branch, nil
}

func (n Node) Filter(filter func(n Node) bool, recurse bool) ([]Node, error) {
	nodes := []Node{}
	if !recurse {
		for _, l := range n.Leaves() {
			if filter(l) {
				nodes = append(nodes, l)
			}
		}
		return nodes, nil
	}
	fn := func(node Node) error {
		if filter(node) {
			nodes = append(nodes, node)
		}
		return nil
	}
	err := n.WalkNode(fn)
	if err != nil {
		return nil, err
	}
	return nodes, nil
}

func (n Node) FilterExt(ext string, recurse bool) ([]Node, error) {
	filter := func(n Node) bool {
		if n.IsDir {
			return false
		}
		return filepath.Ext(n.Name) == ext
	}
	return n.Filter(filter, recurse)
}

func (n Node) GetBranchByPath(path string) (Node, error) {
	return n.GetNodeByPath(path, true)
}

//func walkDirFs(fs afero.Afero, baseDir string, relativeDir string, parent *Node) error {
//  files, err := fs.ReadDir(baseDir)
//  if err != nil {
//    return err
//  }
//  parent.Children = make([]Node, len(files))
//  for i, f := range files {
//    parent.Children[i] = NewNode(f.Name(), parent.Depth+1)
//    if !f.IsDir() {
//      parent.Children[i].IsDir = false
//      parent.Children[i].Dir = relativeDir
//      continue
//    }
//    parent.Children[i].IsDir = true
//    parent.Children[i].Dir = filepath.Join(relativeDir, parent.Children[i].Name)
//    walkDirFs(fs, filepath.Join(baseDir, parent.Children[i].Name), parent.Children[i].Dir, &parent.Children[i])
//  }
//  return nil
//}

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
