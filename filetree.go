package fidi

import (
	"path/filepath"

	"github.com/spf13/afero"
)

var osFs = afero.Afero{Fs: afero.NewOsFs()}

type Filetree struct {
	//tree.Trunk `yaml:",inline"`
	*Dir
	MaxDepth int
	Fs       afero.Afero `yaml:"-"`
}

func NewFS(fs afero.Fs, rootDir string) (*Filetree, error) {
	node := NewDir(rootDir, 0)
	node.isDir = true
	m, err := walkDirFs(afero.Afero{fs}, rootDir, node.Filename().Path, node)
	if err != nil {
		return &Filetree{}, err
	}
	return &Filetree{
		Dir:      node,
		MaxDepth: m,
		Fs:       afero.Afero{fs},
	}, nil
}

func NewFromBasePath(rootDir string) (*Filetree, error) {
	return NewFS(afero.NewBasePathFs(osFs, rootDir), "/")
}

func walkDirFs(fs afero.Afero, baseDir string, relativeDir string, parent *Dir) (int, error) {
	files, err := fs.ReadDir(baseDir)
	if err != nil {
		return 0, err
	}
	depth := 0
	for _, f := range files {
		path := filepath.Join(baseDir, f.Name())
		child := NewDir(path, parent.Depth()+1)
		parent.AddChild(child)
		if parent.Depth() > 0 {
			child.AddParent(parent.Parents()...)
		}
		p := *parent
		child.AddParent(p)
		if !f.IsDir() {
			child.isDir = false
		} else {
			depth++
			child.isDir = true
			walkDirFs(fs, filepath.Join(baseDir, child.Filename().Basename), child.Filename().Path, child)
		}
	}
	return depth, nil
}
