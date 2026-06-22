package fidi

import (
	"fmt"
	"mime"
	"path/filepath"
	"strings"
)

type Node struct {
	Depth    int      `yaml:"depth,omitempty"`
	Name     string   `yaml:"name,omitempty"`
	Path     string   `yaml:"path,omitempty"`
	Ext      string   `yaml:"ext,omitempty"`
	Mimetype string   `yaml:"mimetype,omitempty"`
	AbsPath  string   `yaml:"absPath,omitempty"`
	RelPath  string   `yaml:"relPath,omitempty"`
	IsBranch bool     `yaml:"isBranch,omitempty"`
	Parents  []string `yaml:"parents,omitempty"`
	Children []Node   `yaml:"children,omitempty"`
}

func NewNode(name string, depth int) Node {
	node := Node{
		Name:  name,
		Depth: depth,
		Ext:   filepath.Ext(name),
	}
	node.Mimetype = strings.Split(mime.TypeByExtension(node.Ext), ";")[0]
	return node
}

func (n Node) path() string {
	return filepath.Join(n.Path, n.Name)
}

func (n Node) relPath() string {
	parts := strings.Split(strings.TrimPrefix(n.Path, "/"), "/")
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
		if !n.IsBranch {
			nodes = append(nodes, n)
		}
	}
	return nodes
}

func (n Node) Branches() []Node {
	nodes := []Node{}
	for _, n := range n.Children {
		if n.IsBranch {
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
			if dir && !c.IsBranch {
				continue
			}
			if c.AbsPath == path {
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
		if n.IsBranch {
			return false
		}
		return filepath.Ext(n.Name) == ext
	}
	return n.Filter(filter, recurse)
}

func (n Node) GetBranchByPath(path string) (Node, error) {
	return n.GetNodeByPath(path, true)
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
