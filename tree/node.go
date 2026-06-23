package tree

import (
	"slices"
	"strings"

	"github.com/ohzqq/fidi/fn"
)

type Node interface {
	ID() string
	Depth() int
	Parents() []Node
	Children() []Node
	AddNode(n Node)
	HasChildren() bool
	Walk(fn WalkNodeFunc) error
	Filter(filter FilterNodeFunc, recurse bool) ([]Node, error)
}

type Nodez struct {
	*fn.Filename `yaml:",inline"`
	depth        int            `yaml:"depth" json:"depth,omitempty"`
	isBranch     bool           `yaml:"isBranch" json:"isBranch,omitempty"`
	Meta         map[string]any `yaml:"meta,omitempty" json:"meta,omitempty"`
	parents      []Node         `yaml:"parents,omitempty" json:"parents,omitempty"`
	children     []Node         `yaml:"children,omitempty" json:"children,omitempty"`
	hasChildren  bool           `yaml:"hasChildren" json:"hasChildren,omitempty"`
	id           string
}

func (n Nodez) Depth() int {
	return n.depth
}

func (n Nodez) ID() string {
	return n.id
}

func (n Nodez) Parents() []Node {
	return n.parents
}

func (n Nodez) Children() []Node {
	return n.children
}

func (n Nodez) HasChildren() bool {
	return len(n.Children()) > 0
}

func (n Nodez) AddNode(node Node) {
	n.children = append(n.children, node)
}

func nodezToNodes(nodes []Nodez) []Node {
	n := make([]Node, len(nodes))
	for i, node := range nodes {
		n[i] = node
	}
	return n
}

func NewNode(name string, depth int) Nodez {
	node := Nodez{
		depth: depth,
		id:    name,
	}
	return node
}

func (n Nodez) Walk(fn WalkNodeFunc) error {
	err := fn(n)
	if err != nil {
		return err
	}
	for _, c := range n.children {
		err := c.Walk(fn)
		if err != nil {
			return err
		}
	}
	return nil
}

func (n Nodez) Filter(filter FilterNodeFunc, recurse bool) ([]Node, error) {
	nodes := []Node{}
	if !recurse {
		for _, l := range n.children {
			if l.HasChildren() {
				continue
			}
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

func (n Nodez) FilterByDepth(depth int) ([]Node, error) {
	return n.Filter(func(node Node) bool {
		return node.Depth() <= depth
	}, true)
}

func (n Nodez) GetNodeByPath(path string, dir bool) (Nodez, error) {
	branch := n
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	fn := func(node Node) error {
		for _, c := range node.Children() {
			if dir && c.HasChildren() {
				continue
			}
			//if c.AbsPath == path {
			//  branch = c
			//  return nil
			//}

		}
		return nil
	}
	err := n.Walk(fn)
	if err != nil {
		return n, err
	}
	return branch, nil
}

//func (n Nodez) FilterByExt(ext string, recurse bool) ([]Nodez, error) {
//  filter := func(n Nodez) bool {
//    if n.IsBranch {
//      return false
//    }
//    return n.Ext == ext
//  }
//  return n.Filter(filter, recurse)
//}

func SortLeavesFirst(n Node) error {
	slices.SortStableFunc(n.Children(), func(a, b Node) int {
		if a.HasChildren() && !b.HasChildren() {
			return -1
		}
		if !a.HasChildren() && b.HasChildren() {
			return 1
		}
		return 0
	})
	return nil
}

type WalkNodeFunc func(node Node) error

type FilterNodeFunc func(node Node) bool
