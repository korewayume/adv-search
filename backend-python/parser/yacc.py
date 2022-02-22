# -*- coding: utf-8 -*-
from parser.lex import lexer, tokens
import ply.yacc as yacc
from functools import wraps, partial


class Ast(object):
    tokens = tokens


class AstFactor(Ast):
    pass


class AstExpr(Ast):
    pass


class MatchFactor(AstFactor):
    def __init__(self, first, second):
        self.first = first
        self.second = second

    def __repr__(self):
        return f"{self.__class__.__name__}<{self.first}={self.second}>"


class SuggestFactor(AstFactor):
    def __init__(self, query, start, end):
        self.query = query
        self.start = start
        self.end = end

    def __repr__(self):
        return f"{self.__class__.__name__}<{self.query}>"


class NotExpr(AstExpr):
    def __init__(self, value):
        self.value = value

    def __repr__(self):
        return f"{self.__class__.__name__}<{self.value!r}>"


class LogicalExpr(AstExpr):
    def __init__(self, first, second, op):
        self.first = first
        self.second = second
        self.op = op

    def __repr__(self):
        return f"{self.__class__.__name__}<{self.first!r} {self.op} {self.second!r}>"


def debug(func):
    @wraps(func)
    def inner(p):
        # logging.info("Step Into " + func.__name__)
        verbose = ''
        for i, c in enumerate(p.slice):
            if i == 0:
                verbose += c.type + ' = '
            else:
                verbose += c.type + ' '
        # logging.info(verbose)
        rv = func(p)
        # logging.info("Result: {}".format(p[0]))
        # logging.info("Step Out " + func.__name__)
        return rv

    return inner


@debug
def p_expression_term(p):
    """
    expression : term
    """
    p[0] = p[1]


@debug
def p_term_vbar(p):
    """
    term : term VBAR factor
    """
    p[0] = LogicalExpr(p[1], p[3], '|')


@debug
def p_term_amper(p):
    """
    term : term AMPER factor
    """
    p[0] = LogicalExpr(p[1], p[3], '&')


@debug
def p_term_factor(p):
    """
    term : factor
    """
    p[0] = p[1]


@debug
def p_factor_num(p):
    """
    factor : NAME EQUAL STRING
    factor : NAME
    """
    if len(p) == 2:
        p[0] = SuggestFactor(p[1], p.slice[1].start, p.slice[1].end)
    elif len(p) == 4:
        if p[3][0] == p[3][-1] == '"':
            p[0] = MatchFactor(p[1], p[3][1:-1])
        else:
            raise ParseError("Invalid string: {!r}".format(p[3]))


@debug
def p_factor_not_expr(p):
    """
    factor : NOTEQUAL factor
    """
    p[0] = NotExpr(p[2])


@debug
def p_factor_expr(p):
    """
    factor : LBRACE expression RBRACE
    """
    p[0] = p[2]


class ParseError(Exception):
    pass


# Error rule for syntax errors
def p_error(error):
    raise ParseError(error)


def extract_suggest(result):
    if isinstance(result, NotExpr):
        return extract_suggest(result.value)
    elif isinstance(result, LogicalExpr):
        return extract_suggest(result.first) + extract_suggest(result.second)
    elif isinstance(result, MatchFactor):
        return []
    elif isinstance(result, SuggestFactor):
        return [result]


# Build the parser
parser = yacc.yacc(start='expression')
parse = partial(parser.parse, lexer=lexer)
