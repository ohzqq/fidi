package fidi

import (
	"path/filepath"
	"strings"

	"github.com/ohzqq/fidi/fn"
	"github.com/ohzqq/fidi/tree"
)

type Dir struct {
	tree.Node
	isDir bool
}

func NewDir(name string, depth int) *Dir {
	node := &Dir{
		Node: tree.NewNode(name, depth),
	}
	node.Set("filename", newDirName(name, depth))
	return node
}

func newDirName(name string, depth int) *fn.Filename {
	n := fn.New(name)
	if depth == 0 {
		n.RelPath = "./"
	}
	n.Mimetype = strings.Split(n.Mimetype, ";")[0]
	return n
}

func (d *Dir) Filter(filters ...tree.FilterNodeFunc) ([]*Dir, error) {
	nodes, err := tree.Filter(d, filters...)
	if err != nil {
		return nil, err
	}
	return nodesToDirs(nodes), nil
}

func (d *Dir) FilterByExt(ext string, depth int) ([]*Dir, error) {
	filter := func(n tree.Node) bool {
		name := fn.New(n.ID())
		if n.HasChildren() {
			return false
		}
		return name.Ext == ext
	}
	nodes, err := tree.Filter(d, filter, tree.FilterNodesByDepth(depth))
	if err != nil {
		return nil, err
	}
	return nodesToDirs(nodes), nil
}

func (d *Dir) FilterByMimetype(mt string, depth int) ([]*Dir, error) {
	filter := func(n tree.Node) bool {
		name := fn.New(n.ID())
		if n.HasChildren() {
			return false
		}
		return strings.Contains(name.Mimetype, mt)
	}
	return d.Filter(filter, tree.FilterNodesByDepth(depth))
}

func (d *Dir) Filename() *fn.Filename {
	return d.Get("filename").(*fn.Filename)
}

func (d Dir) RelativizePath() string {
	n := d.Filename()
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

func (d Dir) Files() []*Dir {
	nodes := []*Dir{}
	for _, n := range d.children() {
		if !n.HasChildren() {
			nodes = append(nodes, n)
		}
	}
	return nodes
}

func (d Dir) Dirs() []*Dir {
	nodes := []*Dir{}
	for _, n := range d.children() {
		if n.HasChildren() {
			nodes = append(nodes, n)
		}
	}
	return nodes
}

func (d *Dir) children() []*Dir {
	return nodesToDirs(d.Children())
}

func (d Dir) parents() []*Dir {
	return nodesToDirs(d.Parents())
}

func nodesToDirs(nodes []tree.Node) []*Dir {
	n := make([]*Dir, len(nodes))
	for i, node := range nodes {
		n[i] = &Dir{Node: node}
	}
	return n
}
