package tree

type Trunk struct {
	Node     `yaml:",inline" json:"root,omitempty"`
	MaxDepth int `yaml:"-" json:"maxDepth,omitempty"`
}

func New(node Node, depth int) Trunk {
	return Trunk{
		Node:     node,
		MaxDepth: depth,
	}
}
