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
	node.Set("name", newDirName(name, depth))
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

func (d *Dir) FilterByExt(ext string, depth int) ([]*Dir, error) {
	filter := func(n tree.Node) bool {
		name := fn.New(n.ID())
		if n.HasChildren() {
			return false
		}
		return name.Ext == ext
	}
	nodes, err := d.Filter(depth, filter)
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
	nodes, err := d.Filter(depth, filter)
	if err != nil {
		return nil, err
	}
	return nodesToDirs(nodes), nil
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

func (d *Dir) Filename() *fn.Filename {
	//fmt.Printf("%#v\n", d.Get("name"))
	return d.Get("name").(*fn.Filename)
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

//func (d Dir) Children() []tree.Node {
//return nodezToNodes(d.children)
//}

//func (d Dir) Parents() []tree.Node {
//return nodezToNodes(d.parents)
//}

func nodesToDirs(nodes []tree.Node) []*Dir {
	n := make([]*Dir, len(nodes))
	for i, node := range nodes {
		n[i] = &Dir{Node: node}
	}
	return n
}
