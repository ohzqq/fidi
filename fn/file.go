package fn

import (
	"mime"
	"path"
	"path/filepath"
	"strings"

	"github.com/danielgtaylor/casing"
)

type Filename struct {
	Name       string `yaml:"name,omitempty" json:"name,omitempty"`
	Basename   string `yaml:"basename,omitempty" json:"basename,omitempty"`
	Dir        string `yaml:"dir,omitempty" json:"dir,omitempty"`
	Path       string `yaml:"path,omitempty" json:"path,omitempty"`
	Ext        string `yaml:"ext,omitempty" json:"ext,omitempty"`
	Mimetype   string `yaml:"mimetype,omitempty" json:"mimetype,omitempty"`
	AbsPath    string `yaml:"absPath,omitempty" json:"absPath,omitempty"`
	RelPath    string `yaml:"relPath,omitempty" json:"relPath,omitempty"`
	CamelCase  string `yaml:"camelCase,omitempty" json:"camelCase,omitempty"`
	PascalCase string `yaml:"pascalCase,omitempty" json:"pascalCase,omitempty"`
	KebabCase  string `yaml:"kebabCase,omitempty" json:"kebabCase,omitempty"`
	SnakeCase  string `yaml:"snakeCase,omitempty" json:"snakeCase,omitempty"`
}

func New(name string) *Filename {
	n := &Filename{
		Path: name,
	}
	n.Ext = path.Ext(name)
	n.Dir, n.Basename = path.Split(name)
	if path.IsAbs(name) {
		n.AbsPath = name
	}
	n.RelPath = n.RelativizePath()
	n.Mimetype = mime.TypeByExtension(n.Ext)
	n.Name = strings.TrimSuffix(n.Basename, n.Ext)
	n.PascalCase = casing.Camel(n.Name)
	if n.PascalCase != "" {
		n.CamelCase = casing.LowerCamel(n.Name)
	}
	n.KebabCase = casing.Kebab(n.Name)
	n.SnakeCase = casing.Snake(n.Name)
	return n
}

func (n Filename) RelativizePath() string {
	if n.Path == "/" {
		return "./"
	}
	parts := strings.Split(strings.Trim(n.Dir, "/"), "/")
	dots := make([]string, len(parts))
	for i := range parts {
		dots[i] = ".."
	}
	dots = append(dots, n.Basename)
	return filepath.Join(dots...)
}

func (n *Filename) Matches(pat string) bool {
	m, err := path.Match(pat, n.Dir+n.Basename)
	if err != nil {
		return false
	}
	return m
}
