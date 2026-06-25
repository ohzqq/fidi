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

func (d *Dir) Filter(filters ...tree.FilterNodeFunc) []*Dir {
	return FilterDir(d, false, filters...)
}

func (d *Dir) FilterWalk(filters ...tree.FilterNodeFunc) []*Dir {
	return FilterDir(d, true, filters...)
}

func (d *Dir) FilterByExt(ext string, recurse bool) []*Dir {
	return FilterDirByExt(d, ext, recurse)
}

func (d *Dir) FilterByMimetype(mt string, recurse bool) []*Dir {
	return FilterDirByMimetype(d, mt, recurse)
}

func (d *Dir) FindChild(filters ...tree.FilterNodeFunc) (*Dir, bool) {
	return FindChild(d, filters...)
}

func (d *Dir) FindChildByPath(path string) (*Dir, bool) {
	return FindChildByPath(d, path)
}

func (d *Dir) FindChildByBasename(path string) (*Dir, bool) {
	return FindChildByBasename(d, path)
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

func FindChild(d *Dir, filters ...tree.FilterNodeFunc) (*Dir, bool) {
	files := FilterDir(d, false, filters...)
	if len(files) > 0 {
		return files[0], true
	}
	return nil, false
}

func FindChildByBasename(d *Dir, name string) (*Dir, bool) {
	return FindChild(d, FilterBasename(name))
}

func FindChildByPath(d *Dir, path string) (*Dir, bool) {
	return FindChild(d, FilterPath(path))
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
