package main

import (
	"flag"
	"log"
	"sbom/generator"
)

var (
	path string
)

func main() {
	flag.StringVar(&path, "path", "", "path to a binary file")
	flag.Parse()

	generator.GenerateSBOM(path)

	log.Printf("SBOMs have been generated\n")
}
