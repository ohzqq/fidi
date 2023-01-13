package fidi

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
)

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

var (
	separators  = regexp.MustCompile(`[ &_=+:]`)
	illegalName = regexp.MustCompile(`[^[:alnum:]-.]`)
	dashes      = regexp.MustCompile(`[\-]+`)
)

func SanitizeFilename(n string) string {
	// Remove any trailing space to avoid ending on -
	s = strings.Trim(s, " ")

	// Replace certain joining characters with a dash
	s = separators.ReplaceAllString(s, "-")

	// Remove all other unrecognised characters - NB we do allow any printable characters
	s = r.ReplaceAllString(s, "")

	// Remove any multiple dashes caused by replacements above
	s = dashes.ReplaceAllString(s, "-")

	return s
}
