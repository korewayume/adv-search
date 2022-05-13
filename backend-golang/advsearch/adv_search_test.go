package advsearch

import (
	"testing"
)

func TestAdvSearch(t *testing.T) {
	var ast AdvTerm
	text := `!ip="127.0.0.1" || 银行 && local && ! port="80" && (server="apache")`
	lex := ParseStringLexer(text)
	t.Logf("Suggests:\n%#v\n", lex.Suggests)
	t.Logf("Parse Results:\n%#v\n", lex.Ast)
	ast = ParseAst("  abc  ")
	t.Logf("ParseAst Results:\n%#v\n", ast)
	ast = ParseAst("   ")
	t.Logf("ParseAst Results:\n%#v\n", ast)
	ast = ParseAst("123")
	t.Logf("ParseAst Results:\n%#v\n", ast)
	ast = ParseAst(`{}[]`)
	t.Logf("ParseAst Results:\n%#v\n", ast)
	ast = ParseAst(`ip='127.0.0.1'`)
	t.Logf("ParseAst Results:\n%#v\n", ast)
	ast = ParseAst(`ip="127.0.0.1"`)
}
