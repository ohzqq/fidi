package fidi

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/fs"
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
	fileInfo  fs.FileInfo
	Template  Template
	overwrite bool
	data      []byte
	file      *os.File
}

type Template interface {
	Execute(io.Writer, any) error
}

func NewFile(n string) File {
	name := NewFilename(n)
	name.Root = filepath.Dir(n)

	stat, err := os.Stat(name.Abs)
	if err != nil {
		log.Fatal(err)
	}

	file := File{
		fileInfo: stat,
		Filename: name,
	}

	if !file.fileInfo.IsDir() {
		file.Extension = filepath.Ext(file.Abs)
		file.Name = strings.TrimSuffix(file.Base, file.Extension)
		file.Mime = mime.TypeByExtension(file.Extension)
		file.Dir = filepath.Dir(file.Abs)
	} else {
		file.Name = file.Base
		file.pad = true
		file.Dir = file.Abs
	}

	return file
}

func (f File) Path() string {
	if f.fileInfo.IsDir() {
		return f.Abs
	}
	return f.String()
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

func (f File) Read(b []byte) (int, error) {
	data, err := os.ReadFile(f.Path())
	if err != nil {
		return 0, err
	}
	b = data
	f.data = data

	return len(data), nil
}

func (f *File) Stat() (fs.FileInfo, error) {
	stat, err := os.Stat(f.Abs)
	if err != nil {
		return nil, err
	}
	f.fileInfo = stat

	return stat, nil
}

func (f *File) Close() error {
	err := f.file.Close()
	return err
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

func (f *File) Save(name string) error {
	if name == f.Path() {
		if !f.overwrite {
			return fmt.Errorf("can't save %s, because overwrite isn't set\n", f.Path())
		} else if f.fileInfo.IsDir() {
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
