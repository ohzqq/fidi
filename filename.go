package fidi

import (
	"fmt"
	"io/fs"
	"log"
	"path/filepath"
	"strings"
)

type Filename struct {
	Base      string
	Dir       string
	Extension string
	Name      string
	Root      string
	abs       string
	rel       string
	padding   string
	pad       bool
	prefix    string
	suffix    string
	FileInfo  fs.FileInfo
	num       int
	min       int
	max       int
}

func NewFilename(n string) *Filename {
	n = strings.TrimSuffix(n, "/")
	name := &Filename{
		padding: "%03d",
		Base:    filepath.Base(n),
		Dir:     filepath.Dir(n),
		Name:    n,
		num:     1,
	}
	return name
}

func (n Filename) Copy() *Filename {
	name := &Filename{
		Name:      n.Name,
		Dir:       n.Dir,
		Base:      n.Base,
		Extension: n.Extension,
		padding:   n.padding,
		pad:       n.pad,
	}
	return name
}

func (n Filename) Rename(root string) *Filename {
	name := n.Copy()
	name.Name = root
	return name
}

func (n Filename) Rel() string {
	s, err := filepath.Rel(n.Root, n.rel)
	if err != nil {
		log.Fatal(err)
	}
	return s
}

func (n Filename) Generate(bounds ...int) []*Filename {
	var min, max int
	switch len(bounds) {
	case 2:
		max = bounds[1]
		fallthrough
	case 1:
		min = bounds[0]
	default:
		min = 1
		max = 100
	}

	var names []*Filename
	for i := min; i <= max; i++ {
		n := n.Copy().Num(i).Ext(n.Extension)
		names = append(names, n)
	}

	return names
}

func (n *Filename) Ext(e string) *Filename {
	n.Extension = e
	return n
}

func (n *Filename) Suffix(suf string) *Filename {
	n.suffix = suf
	return n
}

func (n *Filename) Prefix(pre string) *Filename {
	n.prefix = pre
	return n
}

func (n *Filename) Pad() *Filename {
	n.pad = true
	return n
}

func (n *Filename) Fmt(p string) *Filename {
	n.padding = p
	return n
}

func (n *Filename) Zeros(z int) *Filename {
	n.padding = fmt.Sprintf("%%0%dd", z)
	n.Pad()
	return n
}

func (n *Filename) Num(i int) *Filename {
	n.Pad()
	n.num = i
	return n
}

func (n Filename) String() string {
	name := fmt.Sprintf("%s%s%s", n.prefix, n.Name, n.suffix)

	var padding string
	if n.pad {
		padding = fmt.Sprintf(n.padding, n.num)
	}

	name = fmt.Sprintf("%s%s%s", name, padding, n.Extension)

	return name
}
