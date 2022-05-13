%{
package advsearch
%}
%union {
    token LexToken
    term AdvTerm
}

%token <token> IDENT STRING '=' VBAR AMPER '!' '(' ')'
%type <term> term expr factor suggest

%left VBAR
%left AMPER
%right '!'
%right '='

%start term

%%
term : expr
{
$$ = $1
Advlex.(*AdvLex).Ast =  $$
}

suggest: IDENT
{
suggest := AdvTermSuggest{
    Literal: $1.Literal,
    Start: $1.Start,
    End: $1.End,
}
$$ = suggest
}
| suggest IDENT
{
suggest := AdvTermSuggest{
    Literal: $1.(AdvTermSuggest).Literal + $2.Literal,
    Start: $1.(AdvTermSuggest).Start,
    End: $2.End,
}
$$ = suggest
}

factor : suggest
{
$$ = $1
Advlex.(*AdvLex).Suggests =  append(Advlex.(*AdvLex).Suggests, $1.(AdvTermSuggest))
}
| IDENT '=' STRING
{
$$ = AdvTermSearch{
    Key: $1.Literal,
    Value: $3.Unquote(),
    Start: $1.Start,
    End: $3.End,
}
}

expr : expr VBAR expr
{
$$ = AdvTermVbar{
    Left: $1,
    Right: $3,
}
}
| expr AMPER expr
{
$$ = AdvTermAmper{
    Left: $1,
    Right: $3,
}
}
| '!' expr
{
$$ = AdvTermNot{
    Value: $2,
}
}
| '(' expr ')'
{
$$ = $2
}
| factor
{
$$ = $1
}
%%
