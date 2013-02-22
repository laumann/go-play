package main

import "util/paths"
import "fmt"

var (
	p = paths.Get("~/Downloads")
	op = paths.Get("~")
	oop = paths.Get("~/Downloads")
	opf = paths.Get("~/Pictures")
)

func main() {
	p := paths.Get("~/Downloads")
	fmt.Println(p)
	fmt.Println(p.ToAbs())
	fmt.Println(p.NameCount())
	fmt.Println(p.Root())
	fmt.Println(p.IsAbs())

	fmt.Println("\n")
	fmt.Printf("p.StartsWith(%s): %t\n", op, p.StartsWith(op))
	fmt.Printf("p.StartsWith(%s): %t\n", oop, p.StartsWith(oop))
	fmt.Printf("p.StartsWith(%s): %t\n", opf, p.StartsWith(opf))

}
