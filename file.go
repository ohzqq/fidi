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
	*Filename
	Base string
	Mime string
	Stat os.FileInfo
	abs  string
	data []byte
	file *os.File
}

type Dir struct {
	*Filename
	Stat os.FileInfo
	abs  string
}

type Filename struct {
	Dir       string
	Extension string
	Name      string
	rel       string
	Root      string
	padding   string
	pad       bool
	prefix    string
	suffix    string
	num       int
	min       int
	max       int
}

func NewFile(n string) File {
	abs, err := filepath.Abs(n)
	if err != nil {
		log.Fatal(err)
	}

	stat, err := os.Stat(abs)
	if err != nil {
		log.Fatal(err)
	}

	file := File{
		Stat:     stat,
		abs:      abs,
		Base:     filepath.Base(abs),
		Filename: &Filename{},
	}

	if !file.Stat.IsDir() {
		file.Extension = filepath.Ext(abs)
		file.Name = strings.TrimSuffix(file.Base, file.Extension)
		file.Mime = mime.TypeByExtension(file.Extension)
		file.Dir = filepath.Dir(abs)
	} else {
		file.Name = file.Base
		file.pad = true
		file.Dir = abs
	}

	return file
}

func NewFilename(n string) *Filename {
	name := &Filename{
		padding: "%03d",
		Name:    n,
		num:     1,
	}
	return name
}

func (n Filename) Copy() *Filename {
	name := &Filename{
		Name:      n.Name,
		Dir:       n.Dir,
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
	return strings.ReplaceAll(n.rel, n.Root, ".")
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

func (n *Filename) Pad(p string) *Filename {
	n.padding = p
	n.pad = true
	return n
}

func (n *Filename) Num(i int) *Filename {
	n.pad = true
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

	name = filepath.Join(n.Dir, name)

	return name
}

func (f File) Path() string {
	if f.Stat.IsDir() {
		return f.abs
	}
	return f.String()
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
	file, err := os.CreateTemp("", i.Name)
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
