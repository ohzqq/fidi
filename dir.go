package fidi

import (
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
		n.RelPath = "/"
	}
	n.Mimetype = strings.Split(n.Mimetype, ";")[0]
	return n
}

func (d *Dir) Walk(fn tree.WalkNodeFunc) error {
	return tree.Walk(d, fn)
}

func (d *Dir) Filename() *fn.Filename {
	return d.Get("filename").(*fn.Filename)
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
	return d.Filter(filter, tree.FilterNodesByDepth(depth))
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

func FilterDir(dir *Dir, recurse bool, filters ...tree.FilterNodeFunc) []*Dir {
	files := []*Dir{}
	if !recurse {
		for _, d := range dir.children() {
			for _, f := range filters {
				if !f(d) {
					continue
				}
				files = append(files, d)
			}
		}
		return files
	}
	nodes, err := tree.Filter(dir, filters...)
	if err != nil {
		return files
	}
	return nodesToDirs(nodes)
}

func FindChildByBasename(d *Dir, name string) *Dir {
	files := FilterDir(d, false, FilterBasename(name))
	if len(files) > 0 {
		return files[0]
	}
	return d
}

func FindChildByPath(d *Dir, path string) *Dir {
	files := FilterDir(d, false, FilterBasename(path))
	if len(files) > 0 {
		return files[0]
	}
	return d
}

func FilterDirByExt(p *Dir, ext string, recurse bool) []*Dir {
	return FilterDir(p, recurse, FilterExt(ext))
}

func FilterExt(ext string) tree.FilterNodeFunc {
	return func(f tree.Node) bool {
		fn := f.Get("filename").(*fn.Filename)
		return fn.Ext == ext
	}
}

func FilterDirByMimetype(p *Dir, mt string, recurse bool) []*Dir {
	return FilterDir(p, recurse, FilterMimetype(mt))
}

func FilterMimetype(mt string) tree.FilterNodeFunc {
	return func(f tree.Node) bool {
		fn := f.Get("filename").(*fn.Filename)
		return strings.Contains(fn.Mimetype, mt)
	}
}

func FilterBasename(name string) tree.FilterNodeFunc {
	return func(f tree.Node) bool {
		fn := f.Get("filename").(*fn.Filename)
		return fn.Basename == name
	}
}

func FilterPath(name string) tree.FilterNodeFunc {
	return func(f tree.Node) bool {
		fn := f.Get("filename").(*fn.Filename)
		return fn.Path == name
	}
}

func nodesToDirs(nodes []tree.Node) []*Dir {
	n := make([]*Dir, len(nodes))
	for i, node := range nodes {
		n[i] = &Dir{Node: node}
	}
	return n
}
