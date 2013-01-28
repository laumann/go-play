package main

import "fmt"
import "os"
import "path/filepath"
import "log"

func git() {
	if verbose {
		fmt.Printf("Initialising git structure\n")
	}

	touch := func(file, data string) error {
		f, err := os.OpenFile(file, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
		if err != nil {
			return err
		}
		defer f.Close()

		if _, err := f.Write([]byte(data)); err != nil {
			return err
		}
		return nil
	}

	// pkg/.gitkeep
	pkg := filepath.Join(folders["pkg"], ".gitkeep")
	if err := touch(pkg, ""); err != nil {
		log.Fatal(err)
	}

	// bin/.gitkeep
	bin := filepath.Join(folders["bin"], ".gitkeep")
	if err := touch(bin, ""); err != nil {
		log.Fatal(err)
	}

	// .gitignore
	ignore := filepath.Join(baseDir, ".gitignore")
	if err := touch(ignore, "bin/*\npkg/*\n"); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Now do 'git init %s' and 'git commit -a'\n", baseDir)
}
