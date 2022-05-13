package advsearch

import (
	"fmt"
	"go/scanner"
	"go/token"
	"strings"
	"unicode/utf8"
)

// AdvLex AST
type AdvLex struct {
	s        scanner.Scanner
	fset     *token.FileSet
	text     []byte
	Ast      AdvTerm
	Suggests []AdvTermSuggest
	posStart int
	posEnd   int
}

// Position token位置
func (x *AdvLex) Position(pos token.Pos, literal string) (int, int) {
	startInBytes := x.fset.Position(pos).Offset
	endInBytes := startInBytes + len(literal)
	startInUnicode := utf8.RuneCount(x.text[:startInBytes])
	endInUnicode := utf8.RuneCount(x.text[:endInBytes])
	return startInUnicode, endInUnicode
}

// Lex lex
func (x *AdvLex) Lex(context *AdvSymType) int {
	pos, tok, literal := x.s.Scan()
	if (tok == token.SEMICOLON && literal == "\n") || tok == token.EOF {
		return 0
	}
	if literal == "" {
		literal = tok.String()
	}
	start, end := x.Position(pos, literal)
	x.posStart = start
	x.posEnd = end
	context.token = LexToken{
		Literal: literal,
		Start:   start,
		End:     end,
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
	case token.STRING, token.CHAR:
		return STRING
	case token.ASSIGN:
		return int('=')
	case token.EQL:
		return int('=')
	case token.IDENT:
		return IDENT
	default:
		return IDENT
		//panic(fmt.Sprintf(
		//	"SyntaxError: (position=%+v, token=%+v, literal=%+v)",
		//	x.fset.Position(pos), tok.String(), literal,
		//))
	}
}

// Error 错误提示
func (x *AdvLex) Error(s string) {
	errorText := make([]byte, len(x.text))
	for i, _ := range errorText {
		if i == x.posStart || i == x.posEnd-1 {
			errorText[i] = '^'
		} else {
			errorText[i] = ' '
		}
	}
	panic(fmt.Sprintf(
		"语法错误: %s, 位置: (%d, %d)\n%s\n%s",
		s, x.posStart, x.posEnd, x.text, string(errorText)),
	)
}

// NewLexer 构造方法
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

// ParseStringLexer 从string解析Lex
func ParseStringLexer(text string) *AdvLex {
	lex := NewLexer([]byte(text))
	if strings.TrimSpace(text) == "" {
		return lex
	}
	AdvParse(lex)
	return lex
}

// ParseAst 解析AST
func ParseAst(text string) AdvTerm {
	if strings.TrimSpace(text) == "" {
		return nil
	}
	lex := NewLexer([]byte(text))
	AdvParse(lex)
	return lex.Ast
}

// ParseSuggest 解析联想词
func ParseSuggest(text string) []AdvTermSuggest {
	if strings.TrimSpace(text) == "" {
		return []AdvTermSuggest{}
	}
	lex := NewLexer([]byte(text))
	AdvParse(lex)
	return lex.Suggests
}
