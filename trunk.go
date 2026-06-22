package fidi

import (
	"path/filepath"

	"github.com/spf13/afero"
)

type Trunk struct {
	Node
	depth int
	nodes map[string]int
	fs    afero.Afero
}

func New(rootDir string) (Trunk, error) {
	return NewFS(afero.NewBasePathFs(osFs, rootDir), rootDir)
}

func NewFS(fs afero.Fs, rootDir string) (Trunk, error) {
	trunk := Trunk{
		Node:  NewNode(rootDir, 0),
		depth: 0,
		nodes: make(map[string]int),
		fs:    afero.Afero{Fs: fs},
	}
	trunk.RelPath = rootDir
	trunk.IsDir = true
	err := trunk.walkDir(rootDir, trunk.RelPath, &trunk.Node)
	if err != nil {
		return trunk, err
	}
	return trunk, err
}

func (t *Trunk) walkDir(baseDir string, relativeDir string, parent *Node) error {
	files, err := t.fs.ReadDir(baseDir)
	if err != nil {
		return err
	}
	parent.Children = make([]Node, len(files))
	for i, f := range files {
		parent.Children[i] = NewNode(f.Name(), parent.Depth+1)
		if !f.IsDir() {
			t.nodes[relativeDir] = parent.Depth
			parent.Children[i].IsDir = false
			parent.Children[i].RelPath = relativeDir
			continue
		}
		t.depth++
		if parent.Depth > 0 {
			parent.Children[i].parents = append(parent.Children[i].parents, parent.parents...)
		}
		parent.Children[i].parents = append(parent.Children[i].parents, parent.RelPath)
		parent.Children[i].parent = parent.Name
		parent.Children[i].IsDir = true
		parent.Children[i].RelPath = filepath.Join(relativeDir, parent.Children[i].Name)
		t.walkDir(filepath.Join(baseDir, parent.Children[i].Name), parent.Children[i].RelPath, &parent.Children[i])
	}
	return nil
}
