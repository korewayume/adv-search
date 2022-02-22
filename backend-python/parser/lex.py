# -*- coding: utf-8 -*-
from io import BytesIO
import token as std_token
import tokenize as std_tokenize
from typing import Iterator

from ply.lex import LexToken

token_values = (
    std_tokenize.NAME,
    std_tokenize.STRING,
    std_tokenize.EQUAL,
    std_tokenize.NOTEQUAL,
    std_tokenize.AMPER,
    std_tokenize.VBAR,
    std_tokenize.LBRACE,
    std_tokenize.RBRACE,
    # std_tokenize.ENDMARKER,
)
tokens = tuple(std_token.tok_name[x] for x in token_values)


def tokenize(text: bytes):
    token_array = []
    std_token_array = list(std_tokenize.tokenize(BytesIO(text).readline))
    for index, token in enumerate(std_token_array):
        if token.type in (std_tokenize.ENCODING, std_tokenize.NEWLINE, std_tokenize.ENDMARKER):
            continue
        if token.type in (std_tokenize.ERRORTOKEN,) and token.string.strip() == '':
            continue
        token_type = token.type
        string = token.string
        start = token.start[1]
        end = token.end[1]
        line = token.line
        if string == '!':
            token_type = std_tokenize.NOTEQUAL
        elif string in ('&', '|') and index + 1 < len(std_token_array) and std_token_array[index + 1].string == string:
            string = string * 2
            end = end + 1
            if string == '&&':
                token_type = std_tokenize.AMPER
            else:
                token_type = std_tokenize.VBAR
        elif string in ('&', '|') and index - 1 > 0 and std_token_array[index - 1].string == string:
            continue
        elif string == '(':
            token_type = std_tokenize.LBRACE
        elif string == ')':
            token_type = std_tokenize.RBRACE
        elif string == '=':
            token_type = std_tokenize.EQUAL
        elif token_type != std_tokenize.STRING:
            token_type = std_tokenize.NAME

        token = LexToken()
        token.type = std_token.tok_name[token_type]
        token.value = string
        token.lineno = 0
        token.lexpos = end
        token.start = start
        token.end = end
        token.line = line
        token_array.append(token)
    return token_array


class Lexer(object):
    iter_token: Iterator[LexToken]

    def input(self, text):
        self.iter_token = iter(tokenize(text.encode()))

    def token(self):
        return next(self.iter_token, None)


lexer = Lexer()


def example():
    test_text = "".encode('utf8')
    for token in tokenize(test_text):
        print(token)


if __name__ == '__main__':
    example()
