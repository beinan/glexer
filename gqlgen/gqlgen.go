package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
)

const usage = `gqlgen [-dir directory] <schema>
`

var (
	flagSrcDir = flag.String("dir", "", "package source directory, useful for vendored code")
)

func main() {
	flag.Parse()

	if len(flag.Args()) < 1 {
		fmt.Fprint(os.Stderr, usage)
		os.Exit(2)
	}

	schemaFile := flag.Arg(0)
	fmt.Println("Reading schema file from ", schemaFile)
	schemaBytes, err := ioutil.ReadFile(schemaFile)
	if err != nil {
		fmt.Print(err)
	}

	schemaString := string(schemaBytes)

	fmt.Println(schemaString)
}
