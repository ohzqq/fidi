package fidi

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/afero"
)

var osFs = afero.Afero{Fs: afero.NewOsFs()}

type Trunk struct {
	Node     `yaml:"root"`
	MaxDepth int
}

func NewFS(fs afero.Fs, rootDir string) (Trunk, error) {
	tree := Trunk{
		Node:     NewNode(rootDir, 0),
		MaxDepth: 0,
	}
	tree.Path = rootDir
	tree.IsBranch = true
	m, err := walkDirFs(afero.Afero{fs}, rootDir, tree.Path, &tree.Node)
	if err != nil {
		return tree, err
	}
	tree.MaxDepth = m
	return tree, err
}

func NewFromBasePath(rootDir string) (Trunk, error) {
	return NewFS(afero.NewBasePathFs(osFs, rootDir), "/")
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

func walkDirFs(fs afero.Afero, baseDir string, relativeDir string, parent *Node) (int, error) {
	files, err := fs.ReadDir(baseDir)
	if err != nil {
		return 0, err
	}
	depth := 0
	parent.Children = make([]Node, len(files))
	for i, f := range files {
		parent.Children[i] = NewNode(f.Name(), parent.Depth+1)
		if parent.Depth > 0 {
			parent.Children[i].Parents = append(parent.Children[i].Parents, parent.Parents...)
		}
		parent.Children[i].Parents = append(parent.Children[i].Parents, parent.Path)
		if !f.IsDir() {
			parent.Children[i].IsBranch = false
			parent.Children[i].Path = relativeDir
		} else {
			depth++
			//t.MaxDepth++
			parent.Children[i].IsBranch = true
			parent.Children[i].Path = filepath.Join(relativeDir, parent.Children[i].Name)
			walkDirFs(fs, filepath.Join(baseDir, parent.Children[i].Name), parent.Children[i].Path, &parent.Children[i])
		}
		parent.Children[i].AbsPath = parent.Children[i].path()
		parent.Children[i].RelPath = parent.Children[i].relPath()
	}
	return depth, nil
}
