package tree

import (
	"path/filepath"
	"strings"

	"github.com/ohzqq/fidi/fn"
)

type Node struct {
	*fn.Name    `yaml:",inline"`
	Depth       int            `yaml:"depth" json:"depth,omitempty"`
	IsBranch    bool           `yaml:"isBranch,omitempty" json:"isBranch,omitempty"`
	Meta        map[string]any `yaml:"meta,omitempty" json:"meta,omitempty"`
	Parents     []Node         `yaml:"parents,omitempty" json:"parents,omitempty"`
	Children    []Node         `yaml:"children,omitempty" json:"children,omitempty"`
	HasChildren bool           `yaml:"hasChildren,omitempty" json:"hasChildren,omitempty"`
}

func NewNode(name string, depth int) Node {
	node := Node{
		Name:  fn.New(name),
		Depth: depth,
	}
	if depth == 0 {
		node.RelPath = "./"
	}
	node.Mimetype = strings.Split(node.Mimetype, ";")[0]
	return node
}

func (n Node) JoinPath() string {
	return filepath.Join(n.Path, n.Basename)
}

func (n Node) RelativizePath() string {
	if n.Path == "/" {
		return "./"
	}
	parts := strings.Split(strings.TrimPrefix(n.Path, "/"), "/")
	dots := make([]string, len(parts))
	for i := range parts {
		dots[i] = ".."
	}
	dots = append(dots, n.Basename)
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

func (n Node) FilterByExt(ext string, recurse bool) ([]Node, error) {
	filter := func(n Node) bool {
		if n.IsBranch {
			return false
		}
		return n.Ext == ext
	}
	return n.Filter(filter, recurse)
}

type WalkNodeFunc func(node Node) error

type NodeFilterFunc func(node Node) bool
