package main

import (
	"io"
	"text/scanner"
	"unicode"
)

const NEW_LINE = '\n'

var symbolTables = map[string]int{
	"-":              UNORDERED_LIST_MARKER,
	string(NEW_LINE): CR,
	"[":              LBRACKET,
	"]":              RBRACKET,
	"<":              LT,
	">":              GT,
}

type Block interface{}

type Inline interface{}

type UnorderedList struct {
	Items []UnorderedListItem
}

type UnorderedListItem struct {
	Inlines []Inline
}

type Line struct {
	Inlines []Inline
}

type Quotation struct {
	Cite    string
	Content []Block
}

type InlineText struct {
	Literal string
}

type InlineHttp struct {
	Url string
}

type Token struct {
	token   int
	literal string
}

type Lexer struct {
	scanner.Scanner
	result []Block
	err    *ParseError
}

type ParseError struct {
	Message string
	Line    int
	Column  int
}

func (e *ParseError) Error() string {
	return e.Message
}

func isIdent(ch rune, size int) bool {
	return unicode.IsGraphic(ch) && !isReserved(ch) && !unicode.IsSpace(ch)
}

func isReserved(ch rune) bool {
	_, ok := symbolTables[string(ch)]
	return ok
}

func isWhitespace(ch rune) bool {
	return unicode.IsSpace(ch) && ch != rune(NEW_LINE)
}

func NewLexer(in io.Reader) *Lexer {
	l := new(Lexer)
	l.Init(in)
	l.Mode &^= scanner.ScanInts | scanner.ScanFloats | scanner.ScanStrings | scanner.ScanComments | scanner.SkipComments
	l.IsIdentRune = isIdent
	l.Whitespace = 1<<' ' | 1<<'\t'
	return l
}

func (l *Lexer) skipBlank() {
	for isWhitespace(l.Peek()) {
		l.Next()
	}
}

func (l *Lexer) Lex(lval *yySymType) int {
	l.skipBlank()
	ch := l.Peek()
	if isReserved(ch) {
		_ = l.Next()
		s := string(ch)
		token := symbolTables[s]
		lval.token = Token{token: token, literal: s}
		return token
	} else {
		token := int(l.Scan())
		s := l.TokenText()
		if token == scanner.String || token == scanner.Ident {
			token = TEXT
		}
		lval.token = Token{token: token, literal: s}
		return token
	}
}

func (l *Lexer) Error(e string) {
	l.err = &ParseError{
		Message: e,
		Line:    l.Line,
		Column:  l.Column,
	}
}

func Parse(src io.Reader) ([]Block, error) {
	lex := NewLexer(src)
	if ok := yyParse(lex); ok == 0 {
		return lex.result, nil
	} else {
		return nil, lex.err
	}
}
