package main

import (
	"bytes"
	"fmt"
	"os"
	"sort"
)

import "gnuflag"

/**
 * This is taken from gnuflag (because it doesn't export it for some reason...)
 * 
 * flagsByLength is a list of Flag, implementing the Sort interface
 */
type flagsByLength []*gnuflag.Flag

func (f flagsByLength) Less(i, j int) bool {
	s1, s2 := f[i].Name, f[j].Name
	if len(s1) != len(s2) {
		return len(s1) < len(s2)
	}
	return s1 < s2
}

func (f flagsByLength) Swap(i, j int) {
	f[i], f[j] = f[j], f[i]
}

func (f flagsByLength) Len() int {
	return len(f)
}

func flagWithMinus(str string) string {
	if len(str) > 1 {
		return "--" + str
	}
	return "-" + str
}

func printUsage(fs *gnuflag.FlagSet) {
	fmt.Printf("Usage: gowatch [OPTIONS] [PATH]\n\n")
	fmt.Println("OPTIONS")

	flagMap := make(map[interface{}]flagsByLength)
	fs.VisitAll(func(f *gnuflag.Flag) {
		flagMap[f.Value] = append(flagMap[f.Value], f)
	})

	var sortedByName [][]*gnuflag.Flag
	for _, f := range flagMap {
		sort.Sort(f)
		sortedByName = append(sortedByName, f)
	}

	var line bytes.Buffer
	format := "  %-13s %s\n"
	for _, fs := range sortedByName {
		line.Reset()
		if len(fs) > 1 {
			for i, f := range fs {
				if i > 0 {
					line.WriteString(", ")
				}
				line.WriteString(flagWithMinus(f.Name))
			}
		} else {
			name := fs[0].Name
			if len(name) > 1 {
				line.WriteString("    --" + name)
			} else {
				line.WriteString("-" + name)
			}
		}
		fmt.Fprintf(os.Stdout, format, line.Bytes(), fs[0].Usage)
	}
}
