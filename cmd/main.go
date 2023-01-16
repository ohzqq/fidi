package main

import (
	"fmt"
	"log"
	"os"

	"github.com/ohzqq/fidi"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	//cmd.Execute()
	input := os.Args[1]
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

	for _, node := range f.Children {
		fmt.Printf("node path %+V\n", node.Rel())
		//  //fi := node.FilterFilesByExt(".html")
		//  fi := dirfile.ExtFilter(".html")
		//  //fi := dirfile.MimeFilter("image")
		//  files := node.Filter(fi)
		for _, file := range node.Sub() {
			fmt.Printf("sub path %+V\n", file.Base)
			for _, f := range file.Sub() {
				fmt.Printf("sub path %+V\n", f.Base)
			}
		}
		//for _, file := range node.Children {
		//fmt.Printf("child path %+V\n", file.Rel())
		//}

	}

}
