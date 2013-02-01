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


//var autoClosed = 

func main() {
	opf := exampleOpfPackage()
	fmt.Printf("%#v\n", opf)

	xml.HTMLAutoClose = []string{ "item", "itemref" } 
//	for _, tag := range autoClosed {
//		xml.HTMLAutoClose = append(xml.HTMLAutoClose, tag)
//	}

	fmt.Printf("HTMLAutoClose: ", xml.HTMLAutoClose)

	opfXml, err := xml.MarshalIndent(opf, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\n", opfXml)
}
