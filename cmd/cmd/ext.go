package cmd

import (
	"fmt"
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
		for idx, file := range files {
			fn := file.Filename.Copy()
			fn.Fmt("%03d").Num(idx + 1)
			name := NewName(cmd, fn)
			fmt.Printf("new name %s\n", name)
		}
	},
}

func init() {
	rnCmd.AddCommand(extCmd)
}
