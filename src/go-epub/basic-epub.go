package main

import (
	"encoding/xml"
	"fmt"
	"gnuflag"
	"log"
	"os"
	"path/filepath"
	"strings"
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

var fileExtToInclude = []string{"xml", "textile", "markdown", "mk", "css", "html", "xhtml", ""}

func ismember(str string, strs *[]string) bool {
	for _, s := range *strs {
		if s == str {
			return true
		}
	}
	return false
}

// Processing a book
func (book *Book) process() {
	mimetype()
	fmt.Printf("%#v\n", book)

	toInclude := []string{"mimetype"} // mimetype MUST be the first file in the archive

	filepath.Walk(config.workDir, func(path string, info os.FileInfo, err error) error {
		if err == nil && ismember(filepath.Ext(path), &fileExtToInclude) {
			toInclude = append(toInclude, path)
		}
		return err
	})
	fmt.Printf("Files to include: %s\n", toInclude)

	

	// Process is as follows
	// 1. get book details: author, title + list of files to process
	// 2. Process files individually
	// 3. Generate OPF from spine and metadata
	// 4. zip generated files and auxiliary files in .epub (this should be
	//    done from the original "source" directory.

	// Step 4:
	if config.epubFile == "" {
		config.epubFile = strings.Replace(strings.Replace(book.title, " ", "-", -1), ",", "", -1)
	}

	if filepath.Ext(config.epubFile) != ".epub" {
		config.epubFile += ".epub"
	}

	fmt.Printf("Outputting to file: %s\n", config.epubFile)
}

// Initialises the working directory and returns the directory from which we
// came (so we can cd back, once the book has been processed.
// TODO Fix so that existing directories are not a problem...
func initWorkDir() (cwd string, err error) {
	cwd, err = os.Getwd()
	if err != nil {
		return
	}

	if err = os.MkdirAll(config.workDir, 0755); err != nil {
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

	flags.StringVar(&config.epubFile, "output", "", "Set the output .epub file name.")
	flags.StringVar(&config.epubFile, "o", "", "Set the output .epub file name.")

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

	book := &Book{"Thomas Jespersen", "Something, Something Dark Side", []string{"chap1", "chap2"}}

	cwd, err := initWorkDir()
	if err != nil {
		log.Fatal(err)
	}
	book.process()
	if err = os.Chdir(cwd); err != nil {
		log.Fatal(err)
	}

}
