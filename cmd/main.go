package main

import (
	"log"

	"github.com/ohzqq/fidi/cmd/cmd"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	cmd.Execute()
}
