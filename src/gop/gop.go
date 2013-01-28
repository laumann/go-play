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
	"path"
	"path/filepath"
)

import /* GF */ "gnuflag"

var (
	/* Command-line arguments */
	help    bool
	verbose bool
	version bool
	vcs     string

	/* State */
	baseDir string // Could be something else
	folders map[string]string
)

/* Supported version control systems init */
var vcsFns = map[string]func(){
	"git": git,
}


func setpath(args []string) {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	if len(args) == 0 {
		baseDir = cwd
		return
	} else {
		baseDir = args[0]
	}

	if path.IsAbs(baseDir) {
		baseDir = path.Clean(args[0])
		return
	}

	/* Path is relative */
	baseDir = path.Clean(cwd + "/" + baseDir)
}

func initWorkspace() {
	if err := os.Mkdir(baseDir, 0775); err != nil {
		log.Fatal(err)
	}

	folders = map[string]string{
		"src": filepath.Join(baseDir, "src"),
		"pkg": filepath.Join(baseDir, "pkg"),
		"bin": filepath.Join(baseDir, "bin"),
	}

	for _, folder := range folders {
		if verbose {
			fmt.Printf("Creating folder: %s\n", folder)
		}
		if err := os.Mkdir(folder, 0775); err != nil {
			log.Fatal(err)
		}
	}
}

// TODO Dry-run cmd arg
func main() {
	flags := gnuflag.NewFlagSet("gop", gnuflag.ExitOnError)

	flags.BoolVar(&help, "help", false, "Print this help message")
	flags.BoolVar(&help, "h", false, "Print this help message")

	flags.BoolVar(&verbose, "verbose", false, "Be verbose")
	flags.BoolVar(&verbose, "v", false, "Be verbose")

	flags.BoolVar(&version, "version", false, "Print version and exit")
	flags.BoolVar(&version, "V", false, "Print version and exit")

	flags.StringVar(&vcs, "vcs", "", "Set the version control system")

	flags.Parse(true, os.Args[1:])

	if help {
		printUsage(flags)
		os.Exit(0)
	}

	if version {
		printVersion()
		os.Exit(0)
	}

	setpath(flags.Args())
	if verbose {
		fmt.Printf("Initialising Go workspace at %s\n", baseDir)
	}
	initWorkspace()

	if vcsFn, ok := vcsFns[vcs]; ok {
		vcsFn()
	}
}
