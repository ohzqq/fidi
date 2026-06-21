package fidi

import (
	"io/fs"
	"path/filepath"
)

type Trunk struct {
	Node
}

func New(fs fs.ReadDirFS, rootDir string) (Trunk, error) {
	node := Trunk{
		Node: NewNode(rootDir, 0),
	}
	node.RelPath = "/"
	node.IsDir = true
	err := node.Scan(fs, rootDir, node.RelPath, &node.Node)
	if err != nil {
		return node, err
	}
	return node, err
}

func (t *Trunk) Scan(fs fs.ReadDirFS, baseDir string, relativeDir string, parent *Node) error {
	files, err := fs.ReadDir(baseDir)
	if err != nil {
		return err
	}
	parent.Children = make([]Node, len(files))
	for i, f := range files {
		parent.Children[i] = NewNode(f.Name(), parent.Depth+1)
		if !f.IsDir() {
			parent.Children[i].IsDir = false
			parent.Children[i].RelPath = relativeDir
			continue
		}
		parent.Children[i].IsDir = true
		parent.Children[i].RelPath = filepath.Join(relativeDir, parent.Children[i].Name)
		walkDirFs(fs, filepath.Join(baseDir, parent.Children[i].Name), parent.Children[i].RelPath, &parent.Children[i])
	}
	return nil
}
