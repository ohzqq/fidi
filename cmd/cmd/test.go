package cmd

import (
	"fmt"

	"github.com/ohzqq/fidi"
	"github.com/spf13/cobra"
)

// testCmd represents the test command
var testCmd = &cobra.Command{
	Use:   "test",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {
		input := args[0]
		//f := fidi.SanitizeFilename(input)
		//println(f)
		//f := dirfile.NewFile(input)
		//tmpl := template.Must(template.New("toot").Parse("test"))
		//f := fidi.NewFile(input)
		//name := f.Copy()
		//name.Ext(".txt").Num(2)

		//f.Tmpl(tmpl)
		//_, err := f.RenderTemplate("")
		//if err != nil {
		//log.Fatal(err)
		//}
		//err = f.Save(f.Path())
		//if err != nil {
		//log.Fatal(err)
		//}
		//f.Save("toot.txt")

		//f := dirfile.ExtFileNames(input, ".html")
		f := fidi.NewTree(input)
		fmt.Printf("tree path %+V\n", f.Rel())
		//fmt.Printf("tree dir %+V\n", f.Dir)
		//fmt.Printf("tree base %+V\n", f.Base)
		//fmt.Printf("tree name %+V\n", f.Name)
		//fmt.Printf("tree root %+V\n", f.Root)
		//fmt.Printf("tree ext %+V\n", f.Extension)

		for _, node := range f.Nodes {
			fmt.Printf("node path %+V\n", node.Rel())
			//  //fi := node.FilterFilesByExt(".html")
			//  fi := dirfile.ExtFilter(".html")
			//  //fi := dirfile.MimeFilter("image")
			//  files := node.Filter(fi)
			for _, file := range node.Tree.Nodes {
				fmt.Printf("child path %+V\n", file.Rel())
				fmt.Printf("child name %+V\n", file.Depth)
			}

		}

	},
}

func init() {
	rootCmd.AddCommand(testCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// testCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// testCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
