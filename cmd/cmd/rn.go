package cmd

import (
	"fmt"

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
	Short: "batch rename files",
	Run: func(cmd *cobra.Command, args []string) {
		dir := CWD()
		fmt.Printf("cwd %s\n", dir.Path())
		name := dir.Copy()
		fmt.Printf("cwd %s\n", name.String())

		names := GenerateNames(cmd, dir.Filename)

		for _, file := range names {
			name := NewName(cmd, file)
			fmt.Printf("name %s\n", name)
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
