package fidi

import (
	"bytes"
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
	Mime      string
	Depth     int
	Template  Template
	overwrite bool
	data      []byte
	file      *os.File
}

type Template interface {
	Execute(io.Writer, any) error
}

func NewFile(n string) File {
	n = filepath.Clean(n)
	file := File{
		Filename: NewFilename(n),
	}
	file.Root = filepath.Dir(n)

	stat, err := os.Stat(n)
	if err != nil {
		log.Fatal(err)
	}
	file.FileInfo = stat

	if !file.FileInfo.IsDir() {
		file.Extension = filepath.Ext(file.Name)
		file.Name = strings.TrimSuffix(file.Base, file.Extension)
		file.Mime = mime.TypeByExtension(file.Extension)
	} else {
		file.pad = true
		file.Name = file.Base
	}

	return file
}

func (f File) Path() string {
	if f.FileInfo.IsDir() {
		return filepath.Join(f.Dir, f.Name)
	}
	return filepath.Join(f.Dir, f.Filename.String())
}

func (f *File) Tmpl(tmpl Template) *File {
	f.Template = tmpl
	return f
}

func (f *File) Overwrite() *File {
	f.overwrite = true
	return f
}

func (f *File) Write(data []byte) *File {
	f.data = data
	return f
}

func (f *File) RenderTemplate(d any) (*File, error) {
	var buf bytes.Buffer
	if f.Template != nil {
		err := f.Template.Execute(&buf, d)
		if err != nil {
			return f, err
		}
		f.data = buf.Bytes()
	}
	return f, nil
}

func (f File) write(wr io.Writer) error {
	if len(f.data) == 0 {
		return fmt.Errorf("no data to write")
	}

	_, err := wr.Write(f.data)
	if err != nil {
		return err
	}

	return nil
}

func (f *File) Read() error {
	d, err := os.ReadFile(f.Path())
	if err != nil {
		return err
	}
	f.data = d

	return nil
}

func (f *File) Print() {
	f.write(os.Stdout)
}

func (f *File) Tmp() (*os.File, error) {
	file, err := os.CreateTemp("", f.Name)
	if err != nil {
		return nil, err
	}

	err = f.write(file)
	if err != nil {
		return nil, err
	}

	return file, nil
}

func (f File) String() string {
	return string(f.data)
}

func (f *File) Save(name string) error {
	if name == f.Path() {
		if !f.overwrite {
			return fmt.Errorf("can't save %s, because overwrite isn't set\n", f.Path())
		} else if f.FileInfo.IsDir() {
			return fmt.Errorf("can't save, because %s is a directory\n", f.Path())
		}
	}

	file, err := os.Create(name)
	if err != nil {
		return err
	}
	defer file.Close()

	err = f.write(file)
	if err != nil {
		return err
	}

	fmt.Printf("file saved to %s\n", name)

	return nil
}

func Exist(path string) bool {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return false
	}
	return true
}

func init() {
	mime.AddExtensionType(".ini", "text/plain")
	mime.AddExtensionType(".cue", "text/plain")
	mime.AddExtensionType(".m4b", "audio/mp4")
}
