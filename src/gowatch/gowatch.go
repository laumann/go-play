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

func is_match_file(info os.FileInfo) bool {
	matched, _ := filepath.Match(config.mpattern, info.Name())
	return !info.IsDir() && matched
}

func walkStat(now time.Time) {
	nchanged := 0
	filepath.Walk(dtw, func(path string, info os.FileInfo, err error) error {
		if !is_match_file(info) {
			return err
		}
		ninfo, ok := gofiles[path]
		if !ok || info.ModTime().After(ninfo.ModTime()) {
			gofiles[path] = info
			nchanged++
		}
		return err
	})

	if nchanged == 0 {
		return
	}

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
	cmd = renew_cmd(cmd)
}

// Configuration
var (
	help   bool
	config struct {
		cmd      string
		tick     int // milliseconds
		mpattern string
	}
)

func setup_flags() *gnuflag.FlagSet {
	flags := gnuflag.NewFlagSet(os.Args[0], gnuflag.ExitOnError)

	flags.BoolVar(&help, "help", false, "Print this help message.")
	flags.BoolVar(&help, "h", false, "Print this help message.")

	flags.StringVar(&config.cmd, "command", "", "Command to execute.")
	flags.StringVar(&config.cmd, "C", "", "Command to execute.")

	flags.StringVar(&config.mpattern, "pattern", ".go", "The pattern to use for matching files.")
	flags.StringVar(&config.mpattern, "p", ".go", "The pattern to use for matching files.")

	flags.IntVar(&config.tick, "tick", 2000, "How often to check (milliseconds).")
	flags.IntVar(&config.tick, "t", 2000, "How often to check (milliseconds).")

	flags.Parse(true, os.Args[1:])

	return flags
}

func main() {
	flags := setup_flags()

	if help {
		printUsage(flags)
		exit(0)
	}

	if flags.NArg() < 1 {
		fmt.Printf("No directory was set.\n")
		printUsage(flags)
		exit(1)
	}

	dir := flags.Args()[0]
	if !isdir(dir) {
		fmt.Printf("'%s' is not a directory\n", dir)
		exit(1)
	}

	cmd = prepare_cmd(config.cmd)

	if _, err := filepath.Glob(config.mpattern); err != nil {
		fmt.Printf("Provided pattern doesn't work: %s - see 'gowatch --help'\n", config.mpattern)
		fmt.Printf("%s\n", err)
		exit(1)
	}

	dtw, _ = filepath.Abs(dir)
	fmt.Printf("Watching directory '%s'\n", dtw)

	gofiles = make(map[string]os.FileInfo)
	walkStat(time.Now())

	ticker := time.Tick(time.Duration(config.tick) * time.Millisecond)
	for now := range ticker {
		walkStat(now)
	}
}
