package main

import "gocloth"
import "strings"
import "fmt"
import "os"

var pars = []string{
	"p. foo",
	"p(c#id). foo",
	"p(c #id). foo",
	"pp. foo",
	"p . foo",
	"p.foo",
	"p{color: blue}. foo",
	"p{color: blue.foo",
}

var hdrs = []string{
	"h1. foo",
	"h1(c). foo",
	"h1.foo",
	"h5 .foo",
	"h6 . foo",
	"h9. foo",
}

var bqc = []string{
	"bc. block of code",
	"bc(clazz). Now with class",
	"bq. quote",
	"bqc. nonono",
}

func lexStrings(strs []string) {
	for _, s := range strs {
		lex := gocloth.Lexer("foo", s)
		go lex.Run()

		fmt.Println(s)
		fmt.Println(strings.Repeat("-", len(s))) 
		for item := range lex.Items {
			fmt.Println(item)
		}
		fmt.Printf("\n")
	}
}

func main() {
	lexStrings(pars)
	lexStrings(hdrs)
	lexStrings(bqc)


	lexComplex := []string{ "p(c#id). foo\nalso great\n\nbar" }
	lexStrings(lexComplex)

	fmt.Println("New interface!")
	gocloth.ToHtml(os.Stdout, lexComplex[0])
}
