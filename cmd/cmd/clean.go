package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// cleanCmd represents the clean command
var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {
		dir := CWD()
		fmt.Printf("cwd %s\n", dir.Path())

		//for _, file := range dir.All {
		//  if !strings.HasPrefix(file.Base, ".") {
		//    n := fidi.SanitizeFilename(file.Base)
		//    fmt.Printf("old: %s, new: %s\n", file.Base, n)
		//  }
		//}

	},
}

func init() {
	rootCmd.AddCommand(cleanCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// cleanCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// cleanCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
