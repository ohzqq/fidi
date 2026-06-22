package fidi

import (
	"path/filepath"
	"strings"
)

type Filename struct {
	Dir     string
	Ext     string
	Name    string
	RelPath string
	Root    string
}

func NewFilename(n string) *Filename {
	n = strings.TrimSuffix(n, "/")
	name := &Filename{
		Dir:  filepath.Dir(n),
		Name: n,
	}
	return name
}
