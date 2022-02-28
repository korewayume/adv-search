package main

import (
	"encoding/json"
	"main/adv_parse"
	"testing"
)

func TestAdvPArse(t *testing.T) {
	text := `!ip="127.0.0.1" || 银行 && local && ! port="80" && (server="apache")`
	lex := adv_parse.ParseStringLexer(text)
	t.Logf("Suggests: %s\n", lex.Suggests)
	dict, err := json.Marshal(lex.Ast.ToDSL())
	if err != nil {
		t.Errorf("Error: %s\n", err)
		t.FailNow()
	} else {
		t.Logf("DSL: %s\n", dict)
	}
}
