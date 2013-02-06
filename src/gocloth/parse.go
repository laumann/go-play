package gocloth

import (
	"fmt"
	"io"
	"regexp"
)

// Use the lexer to output... whatever the user asks for

var iow io.Writer

// What most people are interested in
func ToHtml(w io.Writer, textile string) {
	iow = w

	lex := lex("textile", textile)
	go lex.run()

	var currentBlock itemType

	first := <-lex.Items
	switch first.typ {
	case itemPar:
		currentBlock = itemPar

		writeByts("<p ") // TODO Write styles, class, etc out
		writeAttrs(first.val, 1)
		writeByts(">")
	case itemText:
		writeByts(first.val)
	default:
		writeByts(fmt.Sprintf("%s\n", first))
	}

	for item := range lex.Items {
		switch item.typ {
		// Blocks
		case itemPar:
			blockEnd(currentBlock)
			currentBlock = itemPar

			writeByts("<p ") // TODO Write styles, class, etc out
			writeAttrs(item.val, 1)
			writeByts(">")

		// Non-blocks
		case itemText:
			writeByts(item.val)
		case itemBreak:
			writeByts("<br/>")
		default:
			writeByts(fmt.Sprintf("%s\n", item))
		}
	}

	blockEnd(currentBlock)

}

// Attributes
var (
	style map[string]string
	class []string
	id    string
	lang  string
)

// 
var classIdRx = regexp.MustCompile("\\((?P<cls>([^#]*))?(#(?P<ids>([^\\s]*)))?\\)")

// TODO Consider parsing these in lex.go and outputting them in tokens
// As soon as parsing fails, output _all_ of it to 
func writeAttrs(attrStr string, pfxLen int) {
	l := len(attrStr)
	if l <= pfxLen {
		return
	}

	toParse := attrStr[pfxLen : l-1]
	//writeByts("class=\"" + toParse + "\"")

	for i := 0; i < len(toParse); i++ {
		switch toParse[i] {
		case '(': // Syntax is (c1 c2 c3#id)
		case '{':
		case '<':
		case '>':
		}
	}
}

func blockEnd(blockType itemType) {
	switch blockType {
	case itemPar:
		writeByts("</p>")
	default:
		writeByts("foo")
	}
	writeByts("\n")
}

func writeByts(byts string) {
	iow.Write([]byte(byts))
}
