package fidi

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/afero"
)

type Trunk struct {
	Node
	MaxDepth int
	fs       afero.Afero
}

func New(rootDir string) (Trunk, error) {
	return NewFS(afero.NewBasePathFs(osFs, rootDir), "/")
}

func NewFS(fs afero.Fs, rootDir string) (Trunk, error) {
	trunk := Trunk{
		Node:     NewNode(rootDir, 0),
		MaxDepth: 0,
		fs:       afero.Afero{Fs: fs},
	}
	trunk.Dir = rootDir
	trunk.IsDir = true
	err := trunk.walkDir(rootDir, trunk.Dir, &trunk.Node)
	if err != nil {
		return trunk, err
	}
	return trunk, err
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
	err := t.WalkNode(fn)
	if err != nil {
		return nil, err
	}
	return nodes, nil
}

func (t *Trunk) walkDir(baseDir string, relativeDir string, parent *Node) error {
	files, err := t.fs.ReadDir(baseDir)
	if err != nil {
		return err
	}
	parent.Children = make([]Node, len(files))
	for i, f := range files {
		parent.Children[i] = NewNode(f.Name(), parent.Depth+1)
		if parent.Depth > 0 {
			parent.Children[i].Parents = append(parent.Children[i].Parents, parent.Parents...)
		}
		parent.Children[i].Parents = append(parent.Children[i].Parents, parent.Dir)
		if !f.IsDir() {
			parent.Children[i].IsDir = false
			parent.Children[i].Dir = relativeDir
		} else {
			t.MaxDepth++
			parent.Children[i].IsDir = true
			parent.Children[i].Dir = filepath.Join(relativeDir, parent.Children[i].Name)
			t.walkDir(filepath.Join(baseDir, parent.Children[i].Name), parent.Children[i].Dir, &parent.Children[i])
		}
		parent.Children[i].Path = parent.Children[i].path()
		parent.Children[i].RelPath = parent.Children[i].relPath()
	}
	return nil
}
