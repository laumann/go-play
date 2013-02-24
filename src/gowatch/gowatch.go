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
	"util/directorywatcher"
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

func walkStat(now time.Time, events []directorywatcher.Event) {
	var msg, stat string
	bs, err := cmd.CombinedOutput()
	if err == nil {
		stat = "ok"
		msg = ""
	} else {
		stat = "failed"
		msg = string(bs)
	}

	fmt.Printf("[%s] regeneration: %d files changed (%s)\n", tstamp(now), len(events), stat)
	fmt.Printf(msg)

	// Renew
	cmd = renew_cmd(cmd)
}

// Configuration
var (
	help   bool
	config struct {
		cmd      string
		tick     uint64 // milliseconds
	}
)

func setup_flags() *gnuflag.FlagSet {
	flags := gnuflag.NewFlagSet(os.Args[0], gnuflag.ExitOnError)

	flags.BoolVar(&help, "help", false, "Print this help message.")
	flags.BoolVar(&help, "h", false, "Print this help message.")

	flags.StringVar(&config.cmd, "command", "", "Command to execute.")
	flags.StringVar(&config.cmd, "C", "", "Command to execute.")

	flags.Uint64Var(&config.tick, "tick", 2000, "How often to check (milliseconds).")
	flags.Uint64Var(&config.tick, "t", 2000, "How often to check (milliseconds).")

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

	dtw, _ = filepath.Abs(dir)
	fmt.Printf("Watching directory '%s'\n", dtw)

	dw, err := directorywatcher.New(dtw)
	if err != nil {
		fmt.Println(err)
		exit(1)
	}
	dw.Recursive = true
	dw.Pattern = "*.go"
	dw.Interval = config.tick
	ch := make(chan []directorywatcher.Event)
	dw.AddObserver(ch)

	dw.Start()
	for {
		select {
		case events := <-ch:
			walkStat(time.Now(), events)
		}
	}
}
