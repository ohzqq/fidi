package fidi

import (
	"path/filepath"

	"github.com/ohzqq/fidi/tree"
	"github.com/spf13/afero"
)

var osFs = afero.Afero{Fs: afero.NewOsFs()}

type Filetree struct {
	*Dir
	MaxDepth int
	Fs       afero.Afero `yaml:"-"`
}

func NewFS(fs afero.Fs, rootDir string) (*Filetree, error) {
	node := NewDir(rootDir, 0)
	node.isDir = true
	ft := &Filetree{
		Dir: node,
		Fs:  afero.Afero{fs},
	}
	err := ft.walkDir(rootDir, node.Filename().Path, node)
	if err != nil {
		return &Filetree{}, err
	}
	err = ft.Walk(tree.SortByLeavesFirst)
	if err != nil {
		return nil, err
	}
	return ft, nil
}

func NewFromBasePath(rootDir string) (*Filetree, error) {
	return NewFS(afero.NewBasePathFs(osFs, rootDir), "/")
}

func (ft *Filetree) walkDir(baseDir string, relativeDir string, parent *Dir) error {
	files, err := ft.Fs.ReadDir(baseDir)
	if err != nil {
		return err
	}
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
			ft.MaxDepth++
			child.isDir = true
			ft.walkDir(filepath.Join(baseDir, child.Filename().Basename), child.Filename().Path, child)
		}
	}
	return nil
}
