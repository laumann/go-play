package main

/**
 * gop is a simple program to initialise the standard Go project "workspace"
 * directory structure.
 *
 * The layout is initialised with the following three folders:
 *
 * project_root
 *   src/
 *   pkg/
 *   bin/
 *
 * Depending on the version control system you (might) want to use, gop can
 * initialise proper structure for sharing. For example, with git you'd like to
 * keep the pkg and bin folders, but probably ignore all the files underneath
 * (to avoid committing built binary files), so these would be initialised with
 * a .gitkeep file and an initial .gitignore file would be added with the
 * entries "bin/*" and "pkg/*"
 */

import (
	"fmt"
	"log"
	"os"
)

func initDirs() {
	if err := os.Mkdir("src", 0775); err != nil {
		log.Fatal(err)
	}

	if err := os.Mkdir("pkg", 0775); err != nil {
		log.Fatal(err)
	}

	if err := os.Mkdir("bin", 0775); err != nil {
		log.Fatal(err)
	}
}

func main() {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Initialising Go workspace at %s\n", cwd)
	initDirs()
}
