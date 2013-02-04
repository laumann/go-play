package main

import "gocloth"
import "fmt"

var pars = []string{
	"p(c#id). foo",
	"p(c #id). foo",
	"pp. foo",
	"p . foo",
	"p.foo",
	"p{color: blue}. foo",
}

func main() {
	for _, p := range pars {
		lex := gocloth.Lexer("foo", p)
		go lex.Run()

		fmt.Println(p,":")
		for item := range lex.Items {
			fmt.Println(item)
		}
	}
}
