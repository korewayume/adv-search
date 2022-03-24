# -*- coding: utf-8 -*-
import sys
import logging
import tokenize
import ast
from functools import wraps
from utils.adv_search_parser.lex import lexer, tokens
import ply.yacc as yacc

logger = logging.getLogger(__name__)
logger.setLevel(logging.INFO)


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
        logger.debug('{:-^30}'.format('Step Into ' + func.__name__))
        verbose = ''
        for i, c in enumerate(p.slice):
            if i == 0:
                verbose += c.type + ' = '
            else:
                verbose += c.type + ' '
        logger.debug(verbose.strip())
        rv = func(p)
        logger.debug("Result: {}".format(p[0]))
        logger.debug('{:-^30}'.format('Step Out ' + func.__name__))
        return rv

    return inner


precedence = (
    ('left', 'VBAR'),
    ('left', 'AMPER'),
    ('right', 'NOTEQUAL'),
    ('right', 'EQUAL'),
)


@debug
def p_term(p):
    """
    term : expr
    """
    p[0] = p[1]


@debug
def p_factor(p):
    """
    factor : NAME
    factor : NAME EQUAL STRING
    """
    if len(p) == 2:
        p[0] = SuggestFactor(p[1], p.slice[1].start, p.slice[1].end)
    elif len(p) == 4:
        if p[3][0] == p[3][-1] in ('"', "'"):
            p[0] = MatchFactor(p[1], ast.literal_eval(p[3]))
        else:
            raise ParseError("Invalid string: {!r}".format(p[3]))


@debug
def p_expr_vbar(p):
    """
    expr : expr VBAR expr
    """
    p[0] = LogicalExpr(p[1], p[3], '|')


@debug
def p_expr_amper(p):
    """
    expr : expr AMPER expr
    """
    p[0] = LogicalExpr(p[1], p[3], '&')


@debug
def p_expr_not(p):
    """
    expr : NOTEQUAL expr
    """
    p[0] = NotExpr(p[2])


@debug
def p_expr_brace(p):
    """
    expr : LBRACE expr RBRACE
    """
    p[0] = p[2]


@debug
def p_expr_factor(p):
    """
    expr : factor
    """
    p[0] = p[1]


class ParseErrorWrapper(Exception):
    pass


class ParseError(Exception):
    def __init__(self, error):
        self.error = error

    def __str__(self):
        if isinstance(self.error, ParseErrorWrapper):
            try:
                exc = self.error.args[0]
                cursor = ' ' * exc.start + '^'
                if exc.end - exc.start > 2:
                    cursor = cursor + ' ' * (exc.end - exc.start - 2) + '^'
                return f"语法错误:\n{exc.line}\n{cursor}"
            except:  # pylint: disable=bare-except
                return f"语法错误: {self.error}"
        elif isinstance(self.error, (tokenize.TokenError, tokenize.StopTokenizing, str)):
            return f"语法错误: {self.error}"
        else:
            return "语法错误"


# Error rule for syntax errors
def p_error(error):
    raise ParseErrorWrapper(error)


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
parser = yacc.yacc(start='term')


def parse(text):
    try:
        return parser.parse(input=text, lexer=lexer)
    except Exception as exc:  # pylint: disable=broad-except
        raise ParseError(exc)


if __name__ == '__main__':
    logger.addHandler(logging.StreamHandler(sys.stdout))
    logger.info(f"""AST: {parse("!a='b' || c='d' && e='f'")}""")
