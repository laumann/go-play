package gocloth

import (
	//"unicode"
	"unicode/utf8"
	"fmt"
)

type stateFn func(*lexer) stateFn

type item struct {
	typ itemType
	val string
}

type itemType int

type token int

// Block modifiers
const (
	PAR token = iota // Should be in <p>...</p>
	H1               // h1
	H2               // h2
	H3               // h3
	H4               // h4
	H5               // h5
	H6               // h6
	BQ               // wrap in <blockquote>...</blockquote>
	BC               // wrap in <pre><code>...</pre></code>
)

// There's a difference
const (
	itemText itemType = iota
	itemParBegin
)

var itemName = map[itemType]string{
	itemText: "TXT",
	itemParBegin: "P",
}

func (it item) String() string {
	return fmt.Sprintf("%s [%s]", itemName[it.typ], it.val)
}

const eof = -1

type lexer struct {
	input string    // the input
	state stateFn   // The current state function
	pos   int       // current position
	start int       // start position of this item
	width int       // width of last rune read from input
	Items chan item // channel of tokens 
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
	switch r := l.next(); {
	case r == 'p':
		return lexParBegin
	case r == '\n':
		l.ignore()
		return lexOutsideBlock
	case r == eof:
		return nil
	default:
		return lexParBegin
	}
	return lexParBegin
}

func lexParBegin(l *lexer) stateFn {
	l.backup()
	p := l.pos

	for r := l.next(); ; r = l.next() {
		if r == '.' {
			if l.peek() != ' ' {
				l.backupTo(p)
			}
			break
		}
		if r == ' ' || r == '\n' || r == eof {
			l.backupTo(p)
			break
		}
	}

	l.emit(itemParBegin)
	if l.peek() == ' ' {
		l.next()
		l.ignore()
	}
	return lexInsideBlock
}

// Inside a block
func lexInsideBlock(l *lexer) stateFn {
	switch r := l.next(); {
	case r == eof:
		l.emit(itemText)
		return nil
	case r == '\n':
		return lexOutsideBlock
	default:
		return lexInsideBlock
	}
	return lexInsideBlock
}
