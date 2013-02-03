package main

import (
	"fmt"
	"os"
)

const VERSION = "0.0.1"

func putVersion() {
	fmt.Printf("%s version %s\n", os.Args[0], VERSION)
}
