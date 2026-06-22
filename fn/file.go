package fn

import (
	"mime"
	"path"
	"strings"

	"github.com/danielgtaylor/casing"
)

type Name struct {
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

func New(name string) *Name {
	n := &Name{
		Path: name,
	}
	if path.IsAbs(name) {
		n.AbsPath = name
	} else {
		n.RelPath = name
	}
	n.Dir, n.Basename = path.Split(name)
	n.Ext = path.Ext(name)
	n.Mimetype = mime.TypeByExtension(n.Ext)
	n.Name = strings.TrimSuffix(name, n.Ext)
	n.PascalCase = casing.Camel(n.Name)
	if n.PascalCase != "" {
		n.CamelCase = casing.LowerCamel(n.Name)
	}
	n.KebabCase = casing.Kebab(n.Name)
	n.SnakeCase = casing.Snake(n.Name)
	return n
}

func (n *Name) Matches(pat string) bool {
	m, err := path.Match(pat, n.Dir+n.Basename)
	if err != nil {
		return false
	}
	return m
}
