package main

import (
	"fmt"
	"io/fs"
	"log"
	"os"

	"github.com/ohzqq/fidi"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	//cmd.Execute()
	input := os.Args[1]
	println(input)
	//tree(input)
	//f := fidi.NewFile(input)
	//printFileInfo(f)
	//f := fidi.SanitizeFilename(input)
	//println(f)
	//f := dirfile.NewFile(input)
	//tmpl := template.Must(template.New("toot").Parse("test"))
	//f := fidi.NewFile(input)
	//name := f.Copy()
	//name.Ext(".txt").Num(2)
	f := fidi.NewTree(input)
	m, err := fs.ReadFile(f, "Nested 1/meta.toml")
	//m, err := fs.Glob(f, "*")
	if err != nil {
		log.Fatal(err)
	}
	//fmt.Printf("%+V\n", m)
	fmt.Printf("%+V\n", string(m))

	printFileInfo(f.Info())
	printFileInfo(f.Children()[0].Info())
	for _, file := range f.Children()[0].Leaves() {
		//fmt.Printf("%d: leaf path %+V\n", file.Depth, file.Rel())
		printFileInfo(file)
	}
}

func tree(input string) {
	f := fidi.NewTree(input)
	fmt.Printf("tree path %+V\n", f.Info().Rel())
	//printFileInfo(f.Info())
	for _, node := range f.Children() {
		d := node.(fidi.Dir)
		fmt.Printf("%d: node path %+V\n", d.Depth, node.Info().Rel())
		//printFileInfo(node.Info())

		//for _, file := range node.Filter(fidi.MimeFilter("image")) {
		for _, file := range node.Parents() {
			d := file.(fidi.Dir)
			fmt.Printf("%d: parent path %+V\n", d.Depth, file.Info().Rel())
			//printFileInfo(file.Info())
		}
		for _, file := range node.Leaves() {
			fmt.Printf("%d: leaf path %+V\n", file.Depth, file.Rel())
			//printFileInfo(file)
		}
		for _, file := range node.Children() {
			fmt.Printf("child path %+V\n", file.Info().Rel())
			//printFileInfo(file.Info())
			for _, sub := range file.Parents() {
				fmt.Printf("sub parent path %+V\n", sub.Info().Rel())
				//printFileInfo(sub.Info())
			}
			for _, sub := range file.Leaves() {
				fmt.Printf("%d: sub leaf path %+V\n", sub.Depth, sub.Rel())
				//printFileInfo(sub)
			}
			//for _, file := range file.Parents() {
			//fmt.Printf("parent path %+V\n", file.Rel())
			//}
			//for _, file := range file.Children() {
			//fmt.Printf("chu path %+V\n", file.Rel())
			//}
			//for _, f := range file.Sub() {
			//fmt.Printf("sub path %+V\n", f.Base)
			//}
			//}
		}
	}
}

func printFileInfo(f fidi.File) {
	//fmt.Printf("file abs %+V\n", f.Abs)
	fmt.Printf("file base %+V\n", f.Base)
	fmt.Printf("file dir %+V\n", f.Dir)
	fmt.Printf("file ext %+V\n", f.Extension)
	fmt.Printf("file name %+V\n", f.Name)
	fmt.Printf("file root %+V\n", f.Root)
	fmt.Printf("file rel %+V\n", f.Rel())
	fmt.Printf("file path %+V\n", f.Path())
}
