package adv_parse

import (
	"encoding/json"
	"strconv"
	"testing"
)

func TestAdvPArse(t *testing.T) {
	text := `!ip="127.0.0.1" || 银行 && local && ! port="80" && (server="apache")`
	lex := ParseStringLexer(text)
	t.Logf("Suggests: %s\n", lex.Suggests)
	dict, err := json.Marshal(lex.Ast.ToDSL())
	if err != nil {
		t.Errorf("Error: %s\n", err)
		t.FailNow()
	} else {
		t.Logf("DSL: %s\n", dict)
	}
}

func Unquote(s string, quote byte) (string, error) {
	if len(s) < 2 {
		return "", strconv.ErrSyntax
	}
	if s[0] != '"' && s[0] != '\'' {
		return "", strconv.ErrSyntax
	}
	if s[0] != s[len(s)-1] {
		return "", strconv.ErrSyntax
	}
	var runeArray []rune
	var chr rune
	var err error
	tail := s[1 : len(s)-1]
	for {
		chr, _, tail, err = strconv.UnquoteChar(tail, quote)
		if err == nil && chr > 0 {
			runeArray = append(runeArray, chr)
		}
		if err != nil || len(tail) <= 0 {
			break
		}
	}
	return string(runeArray), err
}

func TestUnquote(t *testing.T) {
	strings := [...]string{
		`'a="b"'`,
		`'a=\'b\''`,
		`"a='b'"`,
		`"a=\"b\""`,
	}
	for _, s := range strings {
		out, err := Unquote(s, s[0])
		if err != nil {
			t.Errorf("Error: %s\n", err)
			t.FailNow()
		} else {
			t.Logf("Unquote: %s\n", out)
		}
	}
}
