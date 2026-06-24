package tree

import (
	"fmt"
	"slices"
)

type Node interface {
	ID() string
	Depth() int
	Parents() []Node
	Children() []Node
	AddChild(n ...Node)
	AddParent(n ...Node)
	HasChildren() bool
	Walk(fn WalkNodeFunc, filters ...FilterNodeFunc) error
	Filter(filters ...FilterNodeFunc) bool
	Meta() map[string]any
	Get(k string) any
	Set(k string, v any)
	Has(k string) bool
}

type node struct {
	depth       int            `yaml:"depth" json:"depth,omitempty"`
	isBranch    bool           `yaml:"isBranch" json:"isBranch,omitempty"`
	meta        map[string]any `yaml:"meta,omitempty" json:"meta,omitempty"`
	parents     []Node         `yaml:"parents,omitempty" json:"parents,omitempty"`
	children    []Node         `yaml:"children,omitempty" json:"children,omitempty"`
	hasChildren bool           `yaml:"hasChildren" json:"hasChildren,omitempty"`
	id          string
}

func NewNode(name string, depth int) *node {
	node := &node{
		depth: depth,
		id:    name,
		meta:  make(map[string]any),
	}
	return node
}

func (n *node) Depth() int {
	return n.depth
}

func (n *node) ID() string {
	return n.id
}

func (n *node) Parents() []Node {
	return n.parents
}

func (n *node) Children() []Node {
	return n.children
}

func (n *node) HasChildren() bool {
	return len(n.Children()) > 0
}

func (n *node) AddChild(node ...Node) {
	n.children = append(n.children, node...)
}

func (n *node) AddParent(node ...Node) {
	n.parents = append(n.parents, node...)
}

func (n *node) Meta() map[string]any {
	return n.meta
}

func (n *node) Set(k string, v any) {
	n.meta[k] = v
}

func (n *node) Get(k string) any {
	if n.Has(k) {
		return n.meta[k]
	}
	return nil
}

func (n *node) Has(k string) bool {
	_, ok := n.meta[k]
	return ok
}

func (n *node) Walk(fn WalkNodeFunc, filters ...FilterNodeFunc) error {
	walk := n.Filter(filters...)
	err := fn(n)
	if err != nil {
		return err
	}
	if walk {
		for _, c := range n.children {
			err := c.Walk(fn)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// Filter returns false if any of the filters return false.
func (n *node) Filter(filters ...FilterNodeFunc) bool {
	f := true
	for _, filter := range filters {
		if !filter(n) {
			fmt.Printf("%#v\n", n.Get("name"))
			f = false
			//return false
		}
	}
	return f
}

func NodesAtDepth(n Node, depth int) ([]Node, error) {
	filter := func(node Node) bool {
		return node.Depth() == depth
	}
	nodes := []Node{}
	walk := func(node Node) error {
		nodes = append(nodes, node)
		return nil
	}
	err := n.Walk(walk, filter)
	if err != nil {
		return nil, err
	}
	return nodes, nil
}

func FilterByDepth(n Node, depth int) ([]Node, error) {
	filter := func(node Node) bool {
		return node.Depth() <= depth
	}
	nodes := []Node{}
	walk := func(node Node) error {
		nodes = append(nodes, node)
		return nil
	}
	err := n.Walk(walk, filter)
	if err != nil {
		return nil, err
	}

	return nodes, nil
}

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
