// gowatch
//
// Compilation tool. Do 'gowatch <srcdir>' to watch all *.go files listed in
// that directory
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// We want to keep track of the folders, and look at their update times
var goFiles map[string]os.FileInfo

var (
	dtw     string // dir to watch
	timeout int
)

func exit(code int) {
	os.Exit(code)
}

func isDir(d string) bool {
	fp, _ := filepath.Abs(d)
	fi, err := os.Stat(fp)
	return err == nil && fi.IsDir()
}

// format: [yyyy-MM-dd HH:MM:SS]
func fmtTimestamp(t time.Time) string {
	year, month, day := t.Date()
	hour, min, sec := t.Clock()
	return fmt.Sprintf("%04d-%02d-%02d %02d:%02d:%02d", year, month, day, hour, min, sec)
}

func isGoFile(info os.FileInfo) bool {
	return !info.IsDir() && filepath.Ext(info.Name()) == ".go"
}

func watch() {
	c := time.Tick(2 * time.Second)
	var filesChanged int
	for now := range c {
		filesChanged = 0
		filepath.Walk(dtw, func(path string, info os.FileInfo, err error) error {
			if !isGoFile(info) {
				return nil
			}

			inf, exist := goFiles[path]
			if !exist || info.ModTime().After(inf.ModTime()) {
				goFiles[path] = info
				filesChanged++
			}
			return nil
		})
		if filesChanged > 0 {
			fmt.Printf("[%s] regeneration: %d files changed\n", fmtTimestamp(now), filesChanged)
		}
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Must have two arguments")
		exit(1)
	}

	dir := os.Args[1]
	if !isDir(dir) {
		fmt.Printf("'%s' is not a directory\n", dir)
		exit(1)
	}

	dtw, _ = filepath.Abs(dir)
	fmt.Printf("Watching directory '%s'\n", dtw)

	goFiles = make(map[string]os.FileInfo)
	filepath.Walk(dtw, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() && filepath.Ext(path) == ".go" {
			goFiles[path] = info
		}
		return nil
	})
	fmt.Printf("%#v\n", goFiles)

	watch()
}
