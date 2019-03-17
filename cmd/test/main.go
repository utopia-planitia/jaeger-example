package main

import (
	"log"
	"os"

	"github.com/utopia-planitia/exocomp"
)

func main() {

	root := "."
	if len(os.Args) > 1 {
		root = os.Args[1]
	}

	err := exocomp.ImageDigests(root)
	if err != nil {
		log.Printf("failed to modify code: %s", err)
		os.Exit(1)
	}
}
