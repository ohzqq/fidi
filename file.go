package dirfile

import (
	"errors"
	"fmt"
	"io"
	"log"
	"mime"
	"os"
	"path/filepath"
	"strings"
)

type File struct {
	*Name
	Abs  string
	Base string
	Mime string
	Stat os.FileInfo
	data []byte
	file *os.File
}

type Name struct {
	Dir       string
	Extension string
	Root      string
	padding   string
	pad       bool
	prefix    string
	suffix    string
	num       int
	min       int
	max       int
}

func New(n string) File {
	abs, err := filepath.Abs(n)
	if err != nil {
		log.Fatal(err)
	}

	stat, err := os.Stat(abs)
	if err != nil {
		log.Fatal(err)
	}

	file := File{
		Stat: stat,
		Abs:  abs,
		Base: filepath.Base(abs),
		Name: &Name{},
	}

	if !file.Stat.IsDir() {
		file.Extension = filepath.Ext(abs)
		file.Root = strings.TrimSuffix(file.Base, file.Extension)
		file.Mime = mime.TypeByExtension(file.Extension)
		file.Dir = filepath.Dir(abs)
	} else {
		file.Root = file.Base
		file.pad = true
		file.Dir = abs
	}

	return file
}

func NewName(n string) *Name {
	name := &Name{
		padding: "%03d",
		Root:    n,
		num:     1,
	}
	return name
}

func (n Name) Copy(names ...string) *Name {
	name := &Name{
		Root:    n.Root,
		Dir:     n.Dir,
		padding: n.padding,
		pad:     n.pad,
	}
	if len(names) > 0 {
		name.Root = names[0]
	}
	return name
}

func (n Name) Rename(root string) *Name {
	name := n.Copy(root)
	name.Extension = n.Extension
	return name
}

func (n Name) Generate(bounds ...int) []*Name {
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

	var names []*Name
	for i := min; i <= max; i++ {
		n := n.Copy().Num(i).Ext(n.Extension)
		names = append(names, n)
	}

	return names
}

func (n *Name) Ext(e string) *Name {
	n.Extension = e
	return n
}

func (n *Name) Suffix(suf string) *Name {
	n.suffix = suf
	return n
}

func (n *Name) Prefix(pre string) *Name {
	n.prefix = pre
	return n
}

func (n *Name) Pad(p string) *Name {
	n.padding = p
	n.pad = true
	return n
}

func (n *Name) Num(i int) *Name {
	n.pad = true
	n.num = i
	return n
}

func (n Name) Path() string {
	name := fmt.Sprintf("%s%s%s", n.prefix, n.Root, n.suffix)

	var padding string
	if n.pad {
		padding = fmt.Sprintf(n.padding, n.num)
	}

	name = fmt.Sprintf("%s%s%s", name, padding, n.Extension)

	name = filepath.Join(n.Dir, name)

	return name
}

func (n Name) String() string {
	return n.Path()
}

func (i File) Write(wr io.Writer) error {
	_, err := wr.Write(i.data)
	if err != nil {
		return err
	}
	return nil
}

func (i File) Run() error {
	if i.file != nil {
		defer i.file.Close()
	}

	err := i.Write(i.file)
	if err != nil {
		return err
	}

	return nil
}

func (i *File) Tmp(data []byte) {
	file, err := os.CreateTemp("", i.Root)
	if err != nil {
		log.Fatal(err)
	}
	i.file = file
	i.data = data
}

func (i *File) Save(data []byte) {
	file, err := os.Create(i.String())
	if err != nil {
		log.Fatal(err)
	}
	i.file = file
	i.data = data
}

func Exist(path string) bool {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return false
	}
	return true
}
