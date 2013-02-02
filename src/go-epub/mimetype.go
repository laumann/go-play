package main

import "os"
import "log"

const epubMimetype = "application/epub+zip"

// Write the string application/epub+zip to mimetype
func mimetype() {
	mimetype, err := os.OpenFile("mimetype", os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer mimetype.Close()

	mimetype.Write([]byte(epubMimetype))
}
