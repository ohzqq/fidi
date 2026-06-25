package fn

import (
	"mime"
	"path"
	"path/filepath"
	"strings"

	"github.com/danielgtaylor/casing"
)

type Filename struct {
	Name       string `yaml:"name,omitempty" json:"name,omitempty" mapstructure:"name"`
	Basename   string `yaml:"basename,omitempty" json:"basename,omitempty" mapstructure:"basename"`
	Dir        string `yaml:"dir,omitempty" json:"dir,omitempty" mapstructure:"dir"`
	Path       string `yaml:"path,omitempty" json:"path,omitempty" mapstructure:"path"`
	Ext        string `yaml:"ext,omitempty" json:"ext,omitempty" mapstructure:"ext"`
	Mimetype   string `yaml:"mimetype,omitempty" json:"mimetype,omitempty" mapstructure:"mimetype"`
	AbsPath    string `yaml:"absPath,omitempty" json:"absPath,omitempty" mapstructure:"absPath"`
	RelPath    string `yaml:"relPath,omitempty" json:"relPath,omitempty" mapstructure:"relPath"`
	CamelCase  string `yaml:"camelCase,omitempty" json:"camelCase,omitempty" mapstructure:"camelCase"`
	PascalCase string `yaml:"pascalCase,omitempty" json:"pascalCase,omitempty" mapstructure:"pascalCase"`
	KebabCase  string `yaml:"kebabCase,omitempty" json:"kebabCase,omitempty" mapstructure:"kebabCase"`
	SnakeCase  string `yaml:"snakeCase,omitempty" json:"snakeCase,omitempty" mapstructure:"snakeCase"`
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
	if n.Dir == "/" {
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
