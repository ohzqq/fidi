package fidi

import (
	"path/filepath"

	"github.com/spf13/afero"
)

type Trunk struct {
	Node
	depth int
	fs    afero.Afero
}

func New(rootDir string) (Trunk, error) {
	return NewFS(afero.NewBasePathFs(osFs, rootDir), "/")
}

func NewFS(fs afero.Fs, rootDir string) (Trunk, error) {
	trunk := Trunk{
		Node:  NewNode(rootDir, 0),
		depth: 0,
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
		if parent.Depth > 0 {
			parent.Children[i].Parents = append(parent.Children[i].Parents, parent.Parents...)
		}
		parent.Children[i].Parents = append(parent.Children[i].Parents, parent.RelPath)
		if !f.IsDir() {
			parent.Children[i].IsDir = false
			parent.Children[i].RelPath = relativeDir
			continue
		}
		t.depth++
		parent.Children[i].IsDir = true
		parent.Children[i].RelPath = filepath.Join(relativeDir, parent.Children[i].Name)
		t.walkDir(filepath.Join(baseDir, parent.Children[i].Name), parent.Children[i].RelPath, &parent.Children[i])
	}
	return nil
}
