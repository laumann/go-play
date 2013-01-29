package main

import "pathutil"
import "fmt"

func main() {
	p := pathutil.Get("~/Downloads")
	fmt.Println(p)
	fmt.Println(p.ToAbs())
	fmt.Printf("%#v\n", p)
}
