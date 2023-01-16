package fidi

import (
	"os"
	"path/filepath"
	"strings"
)

type Filter func(File) bool

type Dir struct {
	File
	All []File
	//Files   []File
	SubDirs []File
	//Sub          []Dir
	FilesCount   int
	SubDirsCount int
	//Tree
}

func NewDir(path string) (Dir, error) {
	dir := Dir{
		File: NewFile(path),
	}
	dir.rel = path

	err := dir.sort()

	return dir, err
}

func (node *Dir) sort() error {
	entries, err := os.ReadDir(node.Path())
	if err != nil {
		return err
	}

	node.All = make([]File, 0, len(entries))

	for _, entry := range entries {
		e := filepath.Join(node.rel, entry.Name())
		n := NewFile(e)
		n.rel = e
		node.All = append(node.All, n)
	}

	//node.FilesCount = len(node.Files())
	//node.SubDirsCount = len(node.Sub())

	return nil
}

func (node Dir) Sub() []Dir {
	var dirs []Dir
	for _, f := range node.All {
		if f.Stat.IsDir() {
			d, _ := NewDir(f.rel)
			dirs = append(dirs, d)
		}
	}
	return dirs
}

func (node Dir) Files() []File {
	var files []File
	for _, f := range node.All {
		if !f.Stat.IsDir() {
			files = append(files, f)
		}
	}
	return files
}

func (node Dir) Filter(filter Filter) []File {
	return FilterFiles(node.Files(), filter)
}

func FilterFiles(files []File, filter Filter) []File {
	var keep []File
	for _, fn := range files {
		if filter(fn) {
			keep = append(keep, fn)
		}
	}
	return keep
}

func ExtFilter(exts ...string) Filter {
	filter := func(file File) bool {
		for _, ex := range exts {
			if strings.EqualFold(file.Extension, ex) {
				return true
			}
		}
		return false
	}
	return filter
}

func MimeFilter(mimes ...string) Filter {
	filter := func(file File) bool {
		for _, mt := range mimes {
			if strings.Contains(file.Mime, mt) {
				return true
			}
		}
		return false
	}
	return filter
}
