// gowatch
//
// Compilation tool. Do 'gowatch <srcdir>' to watch all *.go files listed in
// that directory
package main

import (
	"fmt"
	"gnuflag"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

// We want to keep track of the folders, and look at their update times
var gofiles map[string]os.FileInfo

var (
	dtw     string // dir to watch
	timeout int
	cmd     *exec.Cmd
)

func exit(code int) {
	os.Exit(code)
}

func isdir(d string) bool {
	fp, _ := filepath.Abs(d)
	fi, err := os.Stat(fp)
	return err == nil && fi.IsDir()
}

// format: [yyyy-MM-dd HH:MM:SS]
func tstamp(t time.Time) string {
	year, month, day := t.Date()
	hour, min, sec := t.Clock()
	return fmt.Sprintf("%04d-%02d-%02d %02d:%02d:%02d", year, month, day, hour, min, sec)
}

func is_go_file(info os.FileInfo) bool {
	return !info.IsDir() && filepath.Ext(info.Name()) == ".go"
}

func walkStat(now time.Time) {
	var nchanged = 0
	filepath.Walk(dtw, func(path string, info os.FileInfo, err error) error {
		if !is_go_file(info) {
			return err
		}
		ninfo, exist := gofiles[path]
		if !exist || info.ModTime().After(ninfo.ModTime()) {
			gofiles[path] = info
			nchanged++
		}
		return err
	})
	if nchanged > 0 {
		var msg, stat string

		bs, err := cmd.CombinedOutput()
		if err == nil {
			stat = "ok"
			msg = ""
		} else {
			stat = "failed"
			msg = string(bs)
		}

		fmt.Printf("[%s] regeneration: %d files changed (%s)\n", tstamp(now), nchanged, stat)
		fmt.Printf(msg)

		// Renew
		cmd = renewcmd(cmd) 
	}
}

// Configuration
var (
	help   bool
	config struct {
		cmd  string
		tick int // milliseconds
	}
)

func setup() *gnuflag.FlagSet {
	flags := gnuflag.NewFlagSet(os.Args[0], gnuflag.ExitOnError)

	flags.BoolVar(&help, "help", false, "Print this help message.")
	flags.BoolVar(&help, "h", false, "Print this help message.")

	flags.StringVar(&config.cmd, "command", "", "Command to execute.")
	flags.StringVar(&config.cmd, "C", "", "Command to execute.")

	flags.IntVar(&config.tick, "tick", 2000, "How often to check (milliseconds).")
	flags.IntVar(&config.tick, "t", 2000, "How often to check (milliseconds).")

	flags.Parse(true, os.Args[1:])

	return flags
}

// TODO Set up compilation
// TODO 
func main() {
	flags := setup()

	if help {
		printUsage(flags)
		exit(0)
	}

	if flags.NArg() < 1 {
		fmt.Printf("No directory was set.")
		exit(1)
	}

	dir := flags.Args()[0]
	if !isdir(dir) {
		fmt.Printf("'%s' is not a directory\n", dir)
		exit(1)
	}

	cmd = preparecmd(config.cmd)

	dtw, _ = filepath.Abs(dir)
	fmt.Printf("Watching directory '%s'\n", dtw)

	gofiles = make(map[string]os.FileInfo)
	walkStat(time.Now())

	ticker := time.Tick(time.Duration(config.tick) * time.Millisecond)
	for now := range ticker {
		walkStat(now)
	}
}
