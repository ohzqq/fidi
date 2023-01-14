package cmd

import (
	"fmt"

	"github.com/ohzqq/fidi"
	"github.com/spf13/cobra"
)

var (
	namePrefix string
	nameSuffix string
	nameExt    string
	nameDig    int
	namePad    bool
	nameNum    int
	nameMin    int
	nameMax    int
	minMax     []int
)

// rnCmd represents the rn command
var rnCmd = &cobra.Command{
	Use:   "rn",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {
		dir := CWD()
		fmt.Printf("cwd %s\n", dir.Path())
		name := dir.Copy()
		fmt.Printf("cwd %s\n", name.String())

		names := []*fidi.Filename{dir.Filename}

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

		if min != 0 && max != 0 {
			names = dir.Generate(min, max)
		}

		for _, file := range names {
			if cmd.Flags().Changed("prefix") {
				file.Prefix(namePrefix)
			}

			if cmd.Flags().Changed("suffix") {
				file.Suffix(nameSuffix)
			}

			if cmd.Flags().Changed("ext") {
				file.Ext(nameExt)
			}

			if cmd.Flags().Changed("num-digits") {
				p := fmt.Sprintf("%0%vd", nameDig)
				file.Pad(p)
			}

			if cmd.Flags().Changed("num") {
				file.Num(nameNum)
			}

			if cmd.Flags().Changed("pad") {
			}

			fmt.Printf("name %s\n", file)
		}
	},
}

func init() {
	rootCmd.AddCommand(rnCmd)

	rnCmd.Flags().StringVarP(&namePrefix, "prefix", "p", "", "prefix for new name")
	rnCmd.Flags().StringVarP(&nameSuffix, "suffix", "s", "", "suffix for new name")
	rnCmd.Flags().StringVarP(&nameExt, "ext", "e", "", "ext for new name")
	rnCmd.Flags().IntSliceVarP(&minMax, "range", "r", []int{}, "range of nums")
	rnCmd.Flags().IntVarP(&nameDig, "num-digits", "N", 0, "number of digits for padding new name")
	rnCmd.Flags().IntVarP(&nameNum, "num", "n", 0, "number for new name")
	rnCmd.Flags().IntVar(&nameMin, "min", 0, "start count at...")
	rnCmd.Flags().IntVar(&nameMax, "max", 0, "end count with...")
	rnCmd.Flags().BoolVar(&namePad, "pad", false, "pad new name")
}
