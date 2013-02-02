package main

import (
	"encoding/xml"
	"fmt"
	"log"
)

type Book struct {
	author string
	title  string
	spine  []string // Should be a "list of things to process"
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

	mimetype()
}
