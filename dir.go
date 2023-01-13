package dirfile

import (
	"os"
	"path/filepath"
	"strings"
)

type Filter func(File) bool

type Dir struct {
	File
	Depth        int
	Files        []File
	SubDirs      []File
	FilesCount   int
	SubDirsCount int
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

	allFiles := make([]File, 0, len(entries))

	for _, entry := range entries {
		e := filepath.Join(node.rel, entry.Name())
		n := NewFile(e)
		n.rel = e
		allFiles = append(allFiles, n)
	}

	for _, f := range allFiles {
		if f.Stat.IsDir() {
			node.SubDirs = append(node.SubDirs, f)
		} else {
			node.Files = append(node.Files, f)
		}
	}

	node.FilesCount = len(node.Files)
	node.SubDirsCount = len(node.SubDirs)

	return nil
}

func (node Dir) Filter(filter Filter) []File {
	return FilterFiles(node.Files, filter)
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
