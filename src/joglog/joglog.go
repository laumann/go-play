package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
)

// Simple program for keeping track of jogging
// Accepts commands:
// jogo list
// jogo add [<date>] <dist><unit> [<comment>]
// jogo remove

var logFile *os.File

type command func([]string)

var cmds = map[string]command{
	"list": list,
	"help": usage,
	"add": add,
	"print": listraw,
}

// Print the raw file
func listraw(args []string) {
	if _, err := io.Copy(os.Stdout, logFile); err != nil {
		log.Fatal(err)
	}
}

func add(args []string) {
	if len(args) < 3 {
		fmt.Println("Usage: joglog add <date> <distance> <unit> [<comment>]")
		os.Exit(1)
	}
	date, dist, unit := args[0], args[1], args[2]

	comt := ""
	if len(args) > 3 {
		comt = args[3]
	}

	if _, err := logFile.WriteString(fmt.Sprintf("%s,%s,%s,\"%s\"\n", date, dist, unit, comt)); err != nil {
		log.Fatal(err)
	}
}

func list(args []string) {
	reader := csv.NewReader(logFile)
	records, err := reader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}
	records = records[1:]
	for i, rec := range records {
		fmt.Printf("%d: %v\n", i+1, rec)
	}
}

func usage(args []string) {
	fmt.Println("joglog <command> <args>")
	fmt.Println("")
	fmt.Println("  Logging jogs.")
	fmt.Println("")
	fmt.Println("Commands")
	fmt.Println("  list    List all the logged jogs.")
	fmt.Println("  add     Add a jog to the list.")
	fmt.Println("  remove  Remove a jog from the log.")
	fmt.Println("  help    Print this message and exit.")
	os.Exit(0)
}

func init() {
	path := os.ExpandEnv("$HOME/joglog")
	f, err := os.OpenFile(path, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		log.Fatal(err)
	}
	logFile = f
}

func main() {
	defer logFile.Close()

	// Some people want -h|--help
	if len(os.Args) < 2 || os.Args[1] == "-h" || os.Args[1] == "--help" {
		usage(os.Args)
	}

	arg1 := os.Args[1]
	if cmd, ok := cmds[arg1]; ok {
		cmd(os.Args[2:])
	} else {
		fmt.Printf("'%s' is not a know command. See 'joglog help'.\n", arg1)
		os.Exit(1)
	}
}
