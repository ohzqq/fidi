package fidi

import (
	"path/filepath"

	"github.com/ohzqq/fidi/tree"
	"github.com/spf13/afero"
)

var osFs = afero.Afero{Fs: afero.NewOsFs()}

type Filetree struct {
	tree.Trunk `yaml:",inline"`
	Fs         afero.Afero `yaml:"-"`
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
		Fs:    afero.Afero{fs},
	}, nil
}

func NewFromBasePath(rootDir string) (Filetree, error) {
	return NewFS(afero.NewBasePathFs(osFs, rootDir), "/")
}

func walkDirFs(fs afero.Afero, baseDir string, relativeDir string, parent *tree.Node) (int, error) {
	files, err := fs.ReadDir(baseDir)
	if err != nil {
		return 0, err
	}
	depth := 0
	parent.Children = make([]tree.Node, len(files))
	for i, f := range files {
		path := filepath.Join(baseDir, f.Name())
		parent.Children[i] = tree.NewNode(path, parent.Depth+1)
		if parent.Depth > 0 {
			parent.Children[i].Parents = append(parent.Children[i].Parents, parent.Parents...)
		}
		p := *parent
		p.Children = []tree.Node{}
		p.Parents = []tree.Node{}
		parent.Children[i].Parents = append(parent.Children[i].Parents, p)
		if !f.IsDir() {
			parent.Children[i].IsBranch = false
			//parent.Children[i].Path = relativeDir
		} else {
			depth++
			parent.Children[i].IsBranch = true
			//parent.Children[i].Path = filepath.Join(relativeDir, parent.Children[i].Basename)
			walkDirFs(fs, filepath.Join(baseDir, parent.Children[i].Basename), parent.Children[i].Path, &parent.Children[i])
			parent.Children[i].HasChildren = len(parent.Children[i].Children) > 0
		}
	}
	return depth, nil
}
