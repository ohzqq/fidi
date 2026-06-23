package fidi

import (
	"path/filepath"
	"strings"

	"github.com/ohzqq/fidi/fn"
	"github.com/ohzqq/fidi/tree"
	"github.com/spf13/afero"
)

var osFs = afero.Afero{Fs: afero.NewOsFs()}

type Filetree struct {
	tree.Trunk `yaml:",inline"`
	Fs         afero.Afero `yaml:"-"`
}

type Nodez struct {
	*fn.Filename `yaml:",inline"`
	tree.Node
	parents  []Nodez `yaml:"parents,omitempty" json:"parents,omitempty"`
	children []Nodez `yaml:"children,omitempty" json:"children,omitempty"`
	isDir    bool
}

func NewNode(name string, depth int) Nodez {
	node := Nodez{
		Filename: fn.New(name),
		Node:     tree.NewNode(name, depth),
	}
	if depth == 0 {
		node.RelPath = "./"
	}
	node.Mimetype = strings.Split(node.Mimetype, ";")[0]
	return node
}

func NewFS(fs afero.Fs, rootDir string) (Filetree, error) {
	node := NewNode(rootDir, 0)
	node.Path = rootDir
	node.isDir = true
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

func (n Nodez) RelativizePath() string {
	if n.Path == "/" {
		return "./"
	}
	parts := strings.Split(strings.TrimPrefix(n.Path, "/"), "/")
	dots := make([]string, len(parts))
	for i := range parts {
		dots[i] = ".."
	}
	dots = append(dots, n.Basename)
	return filepath.Join(dots...)
}

func (n Nodez) Children() []tree.Node {
	return nodezToNodes(n.children)
}

func (n Nodez) Parents() []tree.Node {
	return nodezToNodes(n.parents)
}

func nodezToNodes(nodes []Nodez) []tree.Node {
	n := make([]tree.Node, len(nodes))
	for i, node := range nodes {
		n[i] = node
	}
	return n
}

func walkDirFs(fs afero.Afero, baseDir string, relativeDir string, parent *Nodez) (int, error) {
	files, err := fs.ReadDir(baseDir)
	if err != nil {
		return 0, err
	}
	depth := 0
	parent.children = make([]Nodez, len(files))
	for i, f := range files {
		path := filepath.Join(baseDir, f.Name())
		parent.children[i] = NewNode(path, parent.Depth()+1)
		parent.AddNode(parent.children[i])
		if parent.Depth() > 0 {
			parent.children[i].parents = append(parent.children[i].parents, parent.parents...)
		}
		p := *parent
		p.children = []Nodez{}
		p.parents = []Nodez{}
		parent.children[i].parents = append(parent.children[i].parents, p)
		if !f.IsDir() {
			parent.children[i].isDir = false
		} else {
			depth++
			parent.children[i].isDir = true
			walkDirFs(fs, filepath.Join(baseDir, parent.children[i].Basename), parent.children[i].Path, &parent.children[i])
		}
	}
	return depth, nil
}
