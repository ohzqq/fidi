package tree

import (
	"mime"
	"path/filepath"
	"strings"
)

type Node struct {
	Depth    int            `yaml:"depth,omitempty"`
	Name     string         `yaml:"name,omitempty"`
	Path     string         `yaml:"path,omitempty"`
	Ext      string         `yaml:"ext,omitempty"`
	Mimetype string         `yaml:"mimetype,omitempty"`
	AbsPath  string         `yaml:"absPath,omitempty"`
	RelPath  string         `yaml:"relPath,omitempty"`
	IsBranch bool           `yaml:"isBranch,omitempty"`
	Meta     map[string]any `yaml:"meta,omitempty"`
	Parents  []string       `yaml:"parents,omitempty"`
	Children []Node         `yaml:"children,omitempty"`
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

func (n Node) JoinPath() string {
	return filepath.Join(n.Path, n.Name)
}

func (n Node) RelativizePath() string {
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

func (n Node) Walk(fn WalkNodeFunc) error {
	err := fn(n)
	if err != nil {
		return err
	}
	for _, c := range n.Children {
		err := c.Walk(fn)
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
	err := n.Walk(fn)
	if err != nil {
		return n, err
	}
	return branch, nil
}

func (n Node) Filter(filter NodeFilterFunc, recurse bool) ([]Node, error) {
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
	err := n.Walk(fn)
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

type WalkNodeFunc func(node Node) error

type NodeFilterFunc func(node Node) bool
