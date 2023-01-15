package cmd

import (
	"log"
	"os"

	"github.com/ohzqq/fidi"
	"github.com/spf13/cobra"
)

func GenerateNames(cmd *cobra.Command, name *fidi.Filename) []*fidi.Filename {
	var min, max int
	switch {
	case cmd.Flags().Changed("range"):
		switch len(minMax) {
		case 2:
			max = minMax[1]
			fallthrough
		case 1:
			min = minMax[0]
		}
	case cmd.Flags().Changed("min"):
		min = nameMin
	case cmd.Flags().Changed("max"):
		max = nameMax
	}

	names := []*fidi.Filename{name}
	if min != 0 && max != 0 {
		names = name.Generate(min, max)
	}

	return names
}

func NewName(cmd *cobra.Command, name *fidi.Filename) *fidi.Filename {
	if cmd.Flags().Changed("prefix") {
		name.Prefix(namePrefix)
	}

	if cmd.Flags().Changed("suffix") {
		name.Suffix(nameSuffix)
	}

	if cmd.Flags().Changed("ext") {
		name.Ext(nameExt)
	}

	if cmd.Flags().Changed("num-digits") {
		name.Zeros(nameDig)
	}

	if cmd.Flags().Changed("num") {
		name.Num(nameNum)
	}

	if cmd.Flags().Changed("pad") {
		name.Pad()
	}

	return name
}

//func GetDirs(cmd *cobra.Command, path ...string) []fidi.Dir {

//  var dir fidi.Dir
//  if len(path) > 0 {
//    var err error
//    dir, err = fidi.NewDir(path[0])
//    if err != nil {
//      log.Fatal(err)
//    }
//  } else {
//    dir = CWD()
//  }

//  if cmd.Flags().Changed("recurse") {
//    tree := fidi.NewTree(dir.Path())
//    return tree.Nodes
//  }

//  return []fidi.Dir{dir}
//}

func CWD() fidi.Dir {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	dir, err := fidi.NewDir(cwd)
	if err != nil {
		log.Fatal(err)
	}

	return dir
}
