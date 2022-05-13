package advsearch

import "strconv"

const (
	Literal = iota
	Suggest
	Search
	Not
	Amper
	Vbar
)

func unquote(s string, quote byte) (string, error) {
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

// LexToken token
type LexToken struct {
	Literal string
	Start   int
	End     int
}

// GetType 类型
func (token LexToken) GetType() int {
	return Literal
}

// Unquote Unquote
func (token LexToken) Unquote() string {
	value, _ := unquote(token.Literal, token.Literal[0])
	return value
}

// AdvTerm AST抽象
type AdvTerm interface {
	GetType() int
}

// AdvTermSuggest 联想对象
type AdvTermSuggest struct {
	Literal string
	Start   int
	End     int
}

// GetType 类型
func (token AdvTermSuggest) GetType() int {
	return Suggest
}

// AdvTermSearch 查询类型
type AdvTermSearch struct {
	Key   string
	Value string
	Start int
	End   int
}

// GetType 类型
func (token AdvTermSearch) GetType() int {
	return Search
}

// AdvTermNot Not
type AdvTermNot struct {
	Value AdvTerm
}

// GetType 类型
func (token AdvTermNot) GetType() int {
	return Not
}

// AdvTermAmper Amper
type AdvTermAmper struct {
	Left  AdvTerm
	Right AdvTerm
}

// GetType 类型
func (token AdvTermAmper) GetType() int {
	return Amper
}

// AdvTermVbar Vbar
type AdvTermVbar struct {
	Left  AdvTerm
	Right AdvTerm
}

// GetType 类型
func (token AdvTermVbar) GetType() int {
	return Vbar
}
