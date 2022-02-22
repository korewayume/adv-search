package adv_parse

import (
	"fmt"
	"strings"
)

const (
	Literal = iota
	Suggest
	Search
	Not
	Amper
	Vbar
)

type LexToken struct {
	literal string
	start   int
	end     int
}

func (token LexToken) Value() string {
	starts := strings.HasPrefix(token.literal, `"`) // true
	ends := strings.HasSuffix(token.literal, `"`)
	if starts && ends {
		return token.literal[1 : len(token.literal)-1]
	} else {
		return token.literal
	}
}

func (token LexToken) GetType() int {
	return Literal
}

func (token LexToken) ToDSL() map[string]interface{} {
	return map[string]interface{}{
		"query_string": map[string]interface{}{
			"query": token.Value(),
		},
	}
}

type AdvTerm interface {
	ToDSL() map[string]interface{}
	GetType() int
}

type AdvTermSuggest struct {
	literal string
	start   int
	end     int
}

func (token AdvTermSuggest) Value() string {
	starts := strings.HasPrefix(token.literal, `"`) // true
	ends := strings.HasSuffix(token.literal, `"`)
	if starts && ends {
		return token.literal[1 : len(token.literal)-1]
	} else {
		return token.literal
	}
}

func (token AdvTermSuggest) Start() int {
	return token.start
}

func (token AdvTermSuggest) End() int {
	return token.end
}

func (token AdvTermSuggest) String() string {
	return fmt.Sprintf("Suggest<%s %d,%d>", token.Value(), token.start, token.end)
}

func (token AdvTermSuggest) GetType() int {
	return Suggest
}

func (token AdvTermSuggest) ToDSL() map[string]interface{} {
	return map[string]interface{}{
		"query_string": map[string]interface{}{
			"query": token.Value(),
		},
	}
}

type AdvTermSearch struct {
	key   string
	value string
	start int
	end   int
}

func (token AdvTermSearch) Value() string {
	starts := strings.HasPrefix(token.value, `"`) // true
	ends := strings.HasSuffix(token.value, `"`)
	if starts && ends {
		return token.value[1 : len(token.value)-1]
	} else {
		return token.value
	}
}

func (token AdvTermSearch) GetType() int {
	return Search
}

func (token AdvTermSearch) ToDSL() map[string]interface{} {
	return map[string]interface{}{
		"term": map[string]interface{}{
			token.key: token.Value(),
		},
	}
}

type AdvTermNot struct {
	value AdvTerm
}

func (token AdvTermNot) GetType() int {
	return Not
}

func (token AdvTermNot) ToDSL() map[string]interface{} {
	return map[string]interface{}{
		"bool": map[string]interface{}{
			"must":                 [0]map[string]interface{}{},
			"must_not":             [1]map[string]interface{}{token.value.ToDSL()},
			"should":               [0]map[string]interface{}{},
			"minimum_should_match": 1,
		},
	}
}

type AdvTermAmper struct {
	left  AdvTerm
	right AdvTerm
}

func (token AdvTermAmper) GetType() int {
	return Amper
}

func (token AdvTermAmper) ToDSL() map[string]interface{} {
	return map[string]interface{}{
		"bool": map[string]interface{}{
			"must": [2]map[string]interface{}{
				token.left.ToDSL(), token.right.ToDSL(),
			},
			"must_not":             [0]map[string]interface{}{},
			"should":               [0]map[string]interface{}{},
			"minimum_should_match": 1,
		},
	}
}

type AdvTermVbar struct {
	left  AdvTerm
	right AdvTerm
}

func (token AdvTermVbar) GetType() int {
	return Vbar
}

func (token AdvTermVbar) ToDSL() map[string]interface{} {
	return map[string]interface{}{
		"bool": map[string]interface{}{
			"should": [2]map[string]interface{}{
				token.left.ToDSL(), token.right.ToDSL(),
			},
			"must":                 [0]map[string]interface{}{},
			"must_not":             [0]map[string]interface{}{},
			"minimum_should_match": 1,
		},
	}
}
