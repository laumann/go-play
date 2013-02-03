package main

import "fmt"

const VERSION = "0.0.1"

func putVersion() {
	fmt.Printf("%s version %s\n", PROGNAME, VERSION)
}
