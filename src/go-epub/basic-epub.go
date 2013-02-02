package main

import (
	"encoding/xml"
	"fmt"
	"log"
	"os"
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

type Config struct {
	workDir string // The directory in which all the work is done...
}

var config = Config{
	workDir: ".book", // Default "book" directory
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

// Create working directory, remember dir to return to cd into it
func setup() {

}

// Since everything will eventually be packed up into a work directory, we need
// to do some kind of working directory on all this...
func main() {
	opf := exampleOpfPackage()

	xml.HTMLAutoClose = []string{"item", "itemref"}

	opfXml, err := xml.MarshalIndent(opf, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s", xml.Header)
	fmt.Printf("%s\n", opfXml)

	book := &Book{"Thomas Jespersen", "something something dark side", []string{"chap1", "chap2"}}

	setup()

	cwd, err := initWorkDir()
	if err != nil {
		log.Fatal(err)
	}
	book.process()
	if err = os.Chdir(cwd); err != nil {
		log.Fatal(err)
	}

}
