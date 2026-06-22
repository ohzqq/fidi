package tree

import (
	"fmt"
)

type Trunk struct {
	Node     `yaml:"root" json:"root,omitempty"`
	MaxDepth int `yaml:"-" json:"maxDepth,omitempty"`
}

func New(node Node, depth int) Trunk {
	return Trunk{
		Node:     node,
		MaxDepth: depth,
	}
}

func (t Trunk) GetNodesAtDepth(d int) ([]Node, error) {
	if d > t.MaxDepth {
		return nil, fmt.Errorf("%d is greater than max depth", t.MaxDepth)
	}
	nodes := []Node{}
	fn := func(node Node) error {
		if node.Depth == d {
			nodes = append(nodes, node)
		}
		return nil
	}
	err := t.Walk(fn)
	if err != nil {
		return nil, err
	}
	return nodes, nil
}
