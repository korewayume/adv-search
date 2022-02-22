package adv_parse

import (
	"fmt"
	"go/scanner"
	"go/token"
	"unicode/utf8"
)

type AdvLex struct {
	s        scanner.Scanner
	fset     *token.FileSet
	text     []byte
	Ast      AdvTerm
	Suggests []AdvTermSuggest
}

func (x *AdvLex) Position(pos token.Pos, literal string) (int, int) {
	startInBytes := x.fset.Position(pos).Offset
	endInBytes := startInBytes + len(literal)
	startInUnicode := utf8.RuneCount(x.text[:startInBytes])
	endInUnicode := utf8.RuneCount(x.text[:endInBytes])
	return startInUnicode, endInUnicode
}

func (x *AdvLex) Lex(context *AdvSymType) int {
	pos, tok, literal := x.s.Scan()
	if tok == token.SEMICOLON {
		return 0
	}
	start, end := x.Position(pos, literal)
	context.token = LexToken{
		literal: literal,
		start:   start,
		end:     end,
	}
	switch tok {
	case token.LAND:
		return AMPER
	case token.LOR:
		return VBAR
	case token.NOT:
		return int('!')
	case token.LPAREN:
		return int('(')
	case token.RPAREN:
		return int(')')
	case token.STRING:
		return STRING
	case token.IDENT:
		return IDENT
	case token.ASSIGN:
		return int('=')
	default:
		panic(fmt.Sprintf("SyntaxError: (position=%+v, token=%+v, literal=%+v)", x.fset.Position(pos), tok.String(), literal))
	}
}

func (x *AdvLex) Error(s string) {
	panic(fmt.Sprintf("ParseError: %s", s))
}

func NewLexer(text []byte) *AdvLex {
	var s scanner.Scanner
	fset := token.NewFileSet()                       // positions are relative to fset
	file := fset.AddFile("", fset.Base(), len(text)) // register input "file"
	s.Init(file, text, nil /* no error handler */, scanner.ScanComments)
	return &AdvLex{
		s:    s,
		fset: fset,
		text: text,
	}
}

func ParseStringLexer(text string) *AdvLex {
	lex := NewLexer([]byte(text))
	AdvParse(lex)
	return lex
}
