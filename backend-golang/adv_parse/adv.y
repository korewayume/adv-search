%{
package adv_parse
%}
%union {
    token LexToken
    term AdvTerm
}

%token <token> IDENT STRING '=' VBAR AMPER '!' '(' ')'
%type <term> term expr factor

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

factor : IDENT
{
suggest := AdvTermSuggest{
    literal: $1.literal,
    start: $1.start,
    end: $1.end,
}
$$ = suggest
Advlex.(*AdvLex).Suggests =  append(Advlex.(*AdvLex).Suggests, suggest)
}
| IDENT '=' STRING
{
$$ = AdvTermSearch{
    key: $1.literal,
    value: $3.literal,
    start: $1.start,
    end: $3.end,
}
}

expr : expr VBAR expr
{
$$ = AdvTermVbar{
    left: $1,
    right: $3,
}
}
| expr AMPER expr
{
$$ = AdvTermAmper{
    left: $1,
    right: $3,
}
}
| '!' expr
{
$$ = AdvTermNot{
    value: $2,
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
