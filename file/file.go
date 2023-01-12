package file

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

func FileExist(path string) bool {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return false
	}
	return true
}

type File struct {
	*FileName
	Abs      string
	Base     string
	File     string
	Mimetype string
}

type FileName struct {
	Path    string
	Ext     string
	Name    string
	Padding string
	data    []byte
	file    *os.File
}

func NewFile(n string) File {
	abs, err := filepath.Abs(n)
	if err != nil {
		log.Fatal(err)
	}

	f := File{
		Base:     filepath.Base(abs),
		Abs:      abs,
		FileName: NewFileName(),
	}

	f.Ext = filepath.Ext(abs)
	f.Mimetype = mime.TypeByExtension(f.Ext)
	f.Name = strings.TrimSuffix(f.Base, f.Ext)

	f.Path, f.File = filepath.Split(abs)

	return f
}

func NewFileName() *FileName {
	name := &FileName{
		Padding: "%03d",
	}
	return name
}

func (f File) NewName() *FileName {
	name := &FileName{
		Name:    f.Name,
		Path:    f.Path,
		Padding: f.Padding,
		//Ext:     f.Ext,
	}
	return name
}

func (f *FileName) WithExt(e string) *FileName {
	f.Ext = e
	return f
}

func (f *FileName) Suffix(suf string) *FileName {
	f.Name = f.Name + suf
	return f
}

func (f *FileName) Prefix(pre string) *FileName {
	f.Name = pre + f.Name
	return f
}

func (f *FileName) Pad(i int) *FileName {
	p := fmt.Sprintf(f.Padding, i)
	return f.Suffix(p)
}

func (f FileName) Join() string {
	return filepath.Join(f.Path, f.Name+f.Ext)
}

func (f FileName) Write(wr io.Writer) error {
	_, err := wr.Write(f.data)
	if err != nil {
		return err
	}
	return nil
}

func (f FileName) Run() error {
	if f.file != nil {
		defer f.file.Close()
	}

	err := f.Write(f.file)
	if err != nil {
		return err
	}

	return nil
}

func (f *FileName) Tmp(data []byte) {
	file, err := os.CreateTemp("", f.Name)
	if err != nil {
		log.Fatal(err)
	}
	f.file = file
	f.data = data
}

func (f *FileName) Save(data []byte) {
	file, err := os.Create(f.Join())
	if err != nil {
		log.Fatal(err)
	}
	f.file = file
	f.data = data
}
