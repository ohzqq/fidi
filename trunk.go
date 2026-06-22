package fidi

import (
	"fmt"
	"path/filepath"

	"github.com/ohzqq/fidi/tree"
	"github.com/spf13/afero"
)

var osFs = afero.Afero{Fs: afero.NewOsFs()}

type Filetree struct {
	tree.Trunk
}

func NewFS(fs afero.Fs, rootDir string) (Filetree, error) {
	node := tree.NewNode(rootDir, 0)
	node.Path = rootDir
	node.IsBranch = true
	m, err := walkDirFs(afero.Afero{fs}, rootDir, node.Path, &node)
	if err != nil {
		return Filetree{}, err
	}
	return Filetree{
		Trunk: tree.New(node, m),
	}, nil
}

func NewFromBasePath(rootDir string) (Filetree, error) {
	return NewFS(afero.NewBasePathFs(osFs, rootDir), "/")
}

func (t Filetree) GetNodesAtDepth(d int) ([]tree.Node, error) {
	if d > t.MaxDepth {
		return nil, fmt.Errorf("%d is greater than max depth", t.MaxDepth)
	}
	nodes := []tree.Node{}
	fn := func(node tree.Node) error {
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

func walkDirFs(fs afero.Afero, baseDir string, relativeDir string, parent *tree.Node) (int, error) {
	files, err := fs.ReadDir(baseDir)
	if err != nil {
		return 0, err
	}
	depth := 0
	parent.Children = make([]tree.Node, len(files))
	for i, f := range files {
		parent.Children[i] = tree.NewNode(f.Name(), parent.Depth+1)
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
		parent.Children[i].AbsPath = parent.Children[i].JoinPath()
		parent.Children[i].RelPath = parent.Children[i].RelativizePath()
	}
	return depth, nil
}
