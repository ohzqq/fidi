package cmd

import (
	"log"
	"strings"

	"github.com/ohzqq/fidi"
	"github.com/spf13/cobra"
)

// extCmd represents the ext command
var extCmd = &cobra.Command{
	Use:   "ext",
	Short: "batch rename files by extension",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		ext := args[0]
		if !strings.HasPrefix(ext, ".") {
			log.Fatalf("prefixes should be prefixed with '.'")
		}

		dir := CWD()
		files := dir.Filter(fidi.ExtFilter(ext))
		for _, file := range files {
			println(file.Path())
		}
	},
}

func init() {
	rnCmd.AddCommand(extCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// extCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// extCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
