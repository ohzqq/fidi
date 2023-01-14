package cmd

import (
	"fmt"
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
		p := fmt.Sprintf("%0%vd", nameDig)
		name.Pad(p)
	}

	if cmd.Flags().Changed("num") {
		name.Num(nameNum)
	}

	if cmd.Flags().Changed("pad") {
	}

	return name
}

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
