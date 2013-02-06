package gocloth

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

type stateFn func(*lexer) stateFn

type item struct {
	typ itemType
	val string
}

type itemType int

type token int

// There's a difference
const (
	itemText itemType = iota	// Blocky things (Text being special)
	itemPar
	itemHeader
	itemBlockQuote
	itemBlockCode
	//itemStyle			// Non-blocky things
	itemBreak

)

var itemName = map[itemType]string{
	itemText:       "TXT",
	itemPar:        "P",
	itemHeader:     "H",
	itemBlockQuote: "BQ",
	itemBlockCode:  "BC",
	itemBreak:	"BR",
}

func (it item) String() string {
	return fmt.Sprintf("%s [%s]", itemName[it.typ], it.val)
}

const eof = -1

type lexer struct {
	input    string    // the input
	state    stateFn   // The current state function
	pos      int       // current position
	start    int       // start position of this item
	width    int       // width of last rune read from input
	blockTyp itemType  // 
	Items    chan item // channel of tokens 
}

func (l *lexer) next() (r rune) {
	if l.pos >= len(l.input) { // the length of input is fixed?
		l.width = 0
		return eof
	}
	r, l.width = utf8.DecodeRuneInString(l.input[l.pos:])
	l.pos += l.width
	return r
}

// Backup to a given pos
func (l *lexer) backupTo(pos int) {
	l.pos = pos
}

func (l *lexer) backup() {
	l.pos -= l.width
}

func (l *lexer) peek() rune {
	r := l.next()
	l.backup()
	return r
}

func (l *lexer) emit(t itemType) {
	l.Items <- item{t, l.input[l.start:l.pos]}
	l.start = l.pos
}

func (l *lexer) ignore() {
	l.start = l.pos
}

func Lexer(name, input string) *lexer {
	return lex(name, input)
}

func lex(name, input string) *lexer {
	return &lexer{
		input: input,
		state: lexOutsideBlock,
		Items: make(chan item, 2),
	}
}

func (l *lexer) Run() {
	l.run()
}

func (l *lexer) run() {
	for state := l.state; state != nil; {
		state = state(l)
	}
	close(l.Items)
}

/* State functions */
func lexText(l *lexer) stateFn {
	l.pos = len(l.input)
	l.emit(itemText)
	return nil
}

func lexQuoted(l *lexer) stateFn {
	return nil
}

func lexOutsideBlock(l *lexer) stateFn {
	switch r := l.next(); r {
	case 'p', 'h', 'b':
		l.backup()
		return lexBlockHeader
	case '\n':
		l.ignore()
		return lexOutsideBlock
	case 'f':
		return lexFootnote
	case eof:
		return nil
	}
	l.backup()
	l.emit(itemPar)
	return lexInsideBlock
}

func lexFootnote(l *lexer) stateFn {
	return nil
}

// Idea: match on '{' which can contain anything, and lex in separate function
// turning the lexer, in effect, into a recursive descent (for a bit)
func lexBlockHeader(l *lexer) stateFn {
	p := l.pos
	var blockTyp itemType

	r := l.next()
	rr := l.next()
	if r == 'h' && strings.ContainsRune("123456", rr) {
		blockTyp = itemHeader
	} else if r == 'b' && rr == 'q' {
		blockTyp = itemBlockQuote
	} else if r == 'b' && rr == 'c' {
		blockTyp = itemBlockCode
	} else {
		blockTyp = itemPar
		l.backup()
	}

	for r := l.next(); ; r = l.next() {
		switch r {
		case '.':
			if l.peek() != ' ' {
				goto revert
			}
			goto emit
		case ' ', '\n', eof:
			goto revert
		case '{':
			if !lexStyle(l) {
				goto revert
			}
		case '(':
			if !lexParens(l) {
				goto revert
			}
		case '[':
			if !lexBrackets(l) {
				goto revert
			}
		case '<', '>':
			/* continue */
		default:
			goto revert
		}
	}
revert: // Header didn't work out - output as p.
	l.backupTo(p)
	blockTyp = itemPar
emit:
	l.emit(blockTyp)
	// TODO If the next character is a dot, then set "persistent" blockTyp
	if l.next() == ' ' {
		l.ignore()
	} else {
		l.backup()
	}
	return lexInsideBlock
}

func lexStyle(l *lexer) bool {
	for r := l.next(); r != '}'; r = l.next() {
		if r == '\n' || r == eof {
			return false
		}
	}
	return true
}

func lexParens(l *lexer) bool {
	for r := l.next(); r != ')'; r = l.next() {
		if r == '\n' || r == eof {
			return false
		}
	}
	return true
}

func lexBrackets(l *lexer) bool {
	r := l.next()
	if r == ']' {
		return false
	}
	for ; r == ']'; r = l.next() {
		if r == '\n' || r == eof || r == ' ' || !unicode.IsLetter(r) {
			return false
		}
	}
	return true
}

/*
func lexHeader(l *lexer) stateFn {
//	l.backup()
//	p := l.pos


	return nil
	//for r := l.next(); ; r = l.next() {
//
//	}
}
*/

// Inside a block
func lexInsideBlock(l *lexer) stateFn {
	switch r := l.next(); r {
	case eof:
		l.emit(itemText)
		return nil
	case '\n':
		l.backup()
		l.emit(itemText)
		l.next()
		l.ignore()	// TODO if fold_lines { emit BR } else { l.ignore() }
		return lexInsideJustSeenNL
	default:
		return lexInsideBlock
	}
	return lexInsideBlock
}

func lexInsideJustSeenNL(l *lexer) stateFn {
	r := l.next()
	if r == '\n' {
		l.ignore()
		return lexOutsideBlock
	}
	l.backup()
	l.emit(itemBreak)
	return lexInsideBlock
}
