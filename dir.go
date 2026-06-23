package fidi

import (
	"path/filepath"
	"strings"

	"github.com/ohzqq/fidi/fn"
	"github.com/ohzqq/fidi/tree"
)

type Dir struct {
	*fn.Filename `yaml:",inline"`
	tree.Node
	isDir bool
}

func NewDir(name string, depth int) *Dir {
	node := &Dir{
		Filename: fn.New(name),
		Node:     tree.NewNode(name, depth),
	}
	node.Set("name", node.Filename)
	if depth == 0 {
		node.RelPath = "./"
	}
	node.Mimetype = strings.Split(node.Mimetype, ";")[0]
	return node
}

func (d Dir) FilterByExt(ext string, recurse bool) ([]Dir, error) {
	filter := func(n tree.Node) bool {
		if !n.HasChildren() {
			return false
		}
		return d.Ext == ext
	}
	nodes, err := d.Filter(filter, recurse)
	if err != nil {
		return nil, err
	}
	dirs := make([]Dir, len(nodes))
	//for i, n := range nodes {
	//dirs[i]
	//}
	return dirs, nil
}

func (d Dir) GetNodeByPath(path string, dir bool) (Dir, error) {
	branch := d
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	fn := func(node tree.Node) error {
		for _, c := range node.Children() {
			if dir && c.HasChildren() {
				continue
			}
			//if c.AbsPath == path {
			//  branch = c
			//  return nil
			//}

		}
		return nil
	}
	err := d.Walk(fn)
	if err != nil {
		return d, err
	}
	return branch, nil
}

func (d Dir) Fn() *fn.Filename {
	m := d.Get("name")
	return m.(*fn.Filename)
}

func (d Dir) RelativizePath() string {
	if d.Path == "/" {
		return "./"
	}
	parts := strings.Split(strings.TrimPrefix(d.Path, "/"), "/")
	dots := make([]string, len(parts))
	for i := range parts {
		dots[i] = ".."
	}
	dots = append(dots, d.Basename)
	return filepath.Join(dots...)
}

//func (d Dir) Children() []tree.Node {
//return nodezToNodes(d.children)
//}

//func (d Dir) Parents() []tree.Node {
//return nodezToNodes(d.parents)
//}

func nodesToDirs(nodes []tree.Node) []*Dir {
	n := make([]*Dir, len(nodes))
	for i, node := range nodes {
		n[i] = node.(*Dir)
	}
	return n
}
