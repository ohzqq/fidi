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

//func (n Trunk) WalkMaxDepth(fn WalkNodeFunc, depth int) error {
//  if n.Depth > depth {
//    return nil
//  }
//  err := fn(n)
//  if err != nil {
//    return err
//  }
//  for _, c := range n.Children {
//    err := c.WalkMaxDepth(fn, c.Depth)
//    if err != nil {
//      return err
//    }
//  }
//  return nil
//}

//func (t Trunk) GetNodesAtDepth(d int) ([]Nodez, error) {
//  if d > t.MaxDepth {
//    return nil, fmt.Errorf("%d is greater than max depth", t.MaxDepth)
//  }
//  nodes := []Nodez{}
//  fn := func(node Nodez) error {
//    if node.Depth == d {
//      nodes = append(nodes, node)
//    }
//    return nil
//  }
//  err := t.Walk(fn)
//  if err != nil {
//    return nil, err
//  }
//  return nodes, nil
//}
