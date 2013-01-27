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

import /* GF */ "gnuflag"

var (
	/* Command-line arguments */
	help    bool
	verbose bool
	vcs     string

	/* State */
	path string // Could be something else
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

// TODO Dry-run cmd arg
// TODO First un-matched cmd arg is directory
// TODO Usage String: "Usage: gop [OPTIONS] [path]"
// TODO Save path as some sart of path representation... (File?)
func main() {
	flags := gnuflag.NewFlagSet("gop", gnuflag.ExitOnError)

	flags.BoolVar(&help, "help", false, "Print this help message")
	flags.BoolVar(&help, "h", false, "Print this help message")

	flags.BoolVar(&verbose, "verbose", false, "Be verbose")
	flags.BoolVar(&verbose, "v", false, "Be verbose")

	flags.StringVar(&vcs, "vcs", "", "Set the version control system")

	flags.Parse(true, os.Args[1:])

	if help {
		printUsage(flags)
	}

	if verbose {
		fmt.Printf("%#v\n", os.Args)
	}

	for i, arg := range flags.Args() {
		fmt.Printf("[%d] %s\n", i, arg)
	}

	os.Exit(0)

	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Initialising Go workspace at %s\n", cwd)
	initDirs()
}
