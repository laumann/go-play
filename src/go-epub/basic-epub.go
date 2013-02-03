package main

import (
	"encoding/xml"
	"fmt"
	"gnuflag"
	"log"
	"os"
)

// Global configuration variables
var (
	help    bool
	version bool
)

type Book struct {
	author string
	title  string
	spine  []string // Should be a "list of things to process"
}

// Processing a book
func (book *Book) process() {
	mimetype()
	fmt.Printf("%#v\n", book)
}

// Initialises the working directory and returns the directory from which we
// came (so we can cd back, once the book has been processed. 
func initWorkDir() (cwd string, err error) {
	cwd, err = os.Getwd()
	if err != nil {
		return
	}

	if err = os.Mkdir(config.workDir, 0755); err != nil {
		return
	}

	if err = os.Chdir(config.workDir); err != nil {
		return
	}
	return
}

func setup() *gnuflag.FlagSet {
	flags := gnuflag.NewFlagSet(os.Args[0], gnuflag.ExitOnError)

	flags.BoolVar(&help, "help", false, "Print this help message")
	flags.BoolVar(&help, "h", false, "Print this help message")

	flags.BoolVar(&version, "version", false, "Print this help message")
	flags.BoolVar(&version, "V", false, "Print this help message")

	flags.StringVar(&config.workDir, "workdir", ".book", "Set the working directory.")
	flags.StringVar(&config.workDir, "w", ".book", "Set the working directory.")

	flags.Parse(true, os.Args[1:])

	return flags
}

// Since everything will eventually be packed up into a work directory, we need
// to do some kind of working directory on all this...
func main() {
	flags := setup()

	if help {
		printUsage(flags)
		os.Exit(0)
	}

	if version {
		putVersion()
		os.Exit(0)
	}

	opf := exampleOpfPackage()
	xml.HTMLAutoClose = []string{"item", "itemref"} // Not sure this does anything

	opfXml, err := xml.MarshalIndent(opf, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s", xml.Header)
	fmt.Printf("%s\n", opfXml)

	book := &Book{"Thomas Jespersen", "something something dark side", []string{"chap1", "chap2"}}


	cwd, err := initWorkDir()
	if err != nil {
		log.Fatal(err)
	}
	book.process()
	if err = os.Chdir(cwd); err != nil {
		log.Fatal(err)
	}

}
