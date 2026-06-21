package fidi

import (
	"fmt"
	"io/fs"
	"path/filepath"

	"github.com/spf13/afero"
)

var osFs = afero.Afero{Fs: afero.NewOsFs()}

type Node struct {
	Depth        int
	Path         string
	Root         string
	Name         string
	rel          string
	RelativePath string
	IsDir        bool
	fsys         fs.FS
	//children     []Tree
	//parents      []Tree
	nodes   []Node
	id      int
	Reverse map[string]int
}

func New(rootDir string) (Node, error) {
	result := Node{
		Name:         rootDir,
		RelativePath: "/",
		IsDir:        true,
		Depth:        1,
	}
	err := walkDirFs(afero.NewIOFS(osFs.Fs), rootDir, result.RelativePath, &result)
	return result, err
}

func walkDirFs(fs fs.ReadDirFS, baseDir string, relativeDir string, parent *Node) error {
	files, err := fs.ReadDir(baseDir)
	if err != nil {
		return err
	}
	//parent.nodes = make([]Node, len(files))
	for _, f := range files {
		child := Node{
			Depth: parent.Depth + 1,
			Name:  f.Name(),
		}
		if f.IsDir() {
			child.IsDir = true
			child.RelativePath = filepath.Join(relativeDir, child.Name)
			walkDirFs(fs, filepath.Join(baseDir, child.Name), child.RelativePath, &child)
		} else {
			child.IsDir = false
			child.RelativePath = relativeDir
		}
		parent.Add(child)
	}

	return nil
}

func (n *Node) Add(node Node) {
	n.nodes = append(n.nodes, node)
	if n.Reverse == nil {
		n.Reverse = make(map[string]int)
	}
	n.Reverse[node.RelativePath] = len(n.nodes) - 1
}

func (n Node) Leaves() []Node {
	nodes := []Node{}
	for _, n := range n.nodes {
		if !n.IsDir {
			nodes = append(nodes, n)
		}
	}
	return nodes
}

func (n Node) Branches() []Node {
	nodes := []Node{}
	for _, n := range n.nodes {
		if n.IsDir {
			nodes = append(nodes, n)
		}
	}
	return nodes
}

func (n Node) Children() []Node {
	var children []Node
	if len(n.nodes) > 0 {
		nodes := n.nodes[n.id+1:]
		for _, sub := range nodes {
			if sub.Depth > n.Depth {
				children = append(children, sub)
			}
		}
	}
	return children
}

func (n Node) Parents() []Node {
	var parents []Node
	if len(n.nodes) > 0 && n.Depth > 0 {
		nodes := n.nodes[:n.Depth-1]
		for _, parent := range nodes {
			parents = append(parents, parent)
		}
	}
	return parents
}

func (n Node) HasParents() bool {
	return len(n.Parents()) > 0
}

func (n Node) HasChildren() bool {
	return len(n.Children()) > 0
}

//func (n Node) Filter(filters ...Filter) []File {
//  return FilterFiles(n.Leaves(), filters...)
//}

//func FilterFiles(files []File, filters ...Filter) []File {
//  re := make(map[string]File)
//  for _, filter := range filters {
//    for _, fn := range files {
//      if filter(fn) {
//        re[fn.Name] = fn
//      }
//    }
//  }

//  var keep []File
//  for _, file := range re {
//    keep = append(keep, file)
//  }

//  sort.Slice(keep, func(i, j int) bool {
//    return keep[i].Name < keep[j].Name
//  })

//  return keep
//}

//func ExtFilter(exts ...string) Filter {
//  filter := func(file File) bool {
//    for _, ex := range exts {
//      if strings.EqualFold(file.Extension, ex) {
//        return true
//      }
//    }
//    return false
//  }
//  return filter
//}

//func MimeFilter(mimes ...string) Filter {
//  filter := func(file File) bool {
//    for _, mt := range mimes {
//      if strings.Contains(file.Mime, mt) {
//        return true
//      }
//    }
//    return false
//  }
//  return filter
//}

type NodeIndexDontExistsError struct {
	Index int
}

func (e *NodeIndexDontExistsError) Error() string {
	return fmt.Sprintf("Node with index [%v] not exists", e.Index)
}
