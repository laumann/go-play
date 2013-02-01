package main

import (
	"fmt"
	"encoding/xml"
	"log"
)

type Book struct {
	author string
	title  string
	spine  []string // Should be a "list of things to process"
}


func main() {
	opf := exampleOpfPackage()

	xml.HTMLAutoClose = []string{ "item", "itemref" } 

	opfXml, err := xml.MarshalIndent(opf, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\n", opfXml)
}
