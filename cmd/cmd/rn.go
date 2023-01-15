package cmd

import (
	"github.com/spf13/cobra"
)

var (
	namePrefix string
	nameSuffix string
	nameExt    string
	nameDig    int
	namePad    bool
	recurseDir bool
	nameNum    int
	nameMin    int
	nameMax    int
	minMax     []int
)

// rnCmd represents the rn command
var rnCmd = &cobra.Command{
	Use:   "rn",
	Short: "batch rename files",
	Args:  cobra.MaximumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		//var dirs []fidi.Dir
		//if len(args) > 0 {
		//if d := args[0]; d != "" {
		//dirs = GetDirs(cmd, d)
		//}
		//}
		//if len(args) > 1 {
		//if ext := args[1]; ext != "" {
		//  files = dir.Filter(fidi.ExtFilter(ext))
		//}
		//}
		//fmt.Printf("dir %s\n", dir.Path())

		//for _, dir := range dirs {
		//name := NewName(cmd, dir.Filename)
		//fmt.Printf("name %s\n", name)
		//}
	},
}

func init() {
	rootCmd.AddCommand(rnCmd)

	rnCmd.PersistentFlags().StringVarP(&namePrefix, "prefix", "p", "", "prefix for new name")
	rnCmd.PersistentFlags().StringVarP(&nameSuffix, "suffix", "s", "", "suffix for new name")
	rnCmd.PersistentFlags().StringVarP(&nameExt, "ext", "e", "", "ext for new name")
	rnCmd.PersistentFlags().IntSliceVar(&minMax, "range", []int{}, "range of nums")
	rnCmd.PersistentFlags().IntVarP(&nameDig, "num-digits", "N", 0, "number of digits for padding new name")
	rnCmd.PersistentFlags().IntVarP(&nameNum, "num", "n", 0, "number for new name")
	rnCmd.PersistentFlags().BoolVar(&namePad, "pad", false, "pad new name")
	rnCmd.PersistentFlags().BoolVarP(&recurseDir, "recurse", "r", false, "pad new name")
	rnCmd.Flags().IntVar(&nameMin, "min", 0, "start count at...")
	rnCmd.Flags().IntVar(&nameMax, "max", 0, "end count with...")
}
