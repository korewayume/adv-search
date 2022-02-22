# -*- coding: utf-8 -*-
from flask import Flask, jsonify, request
from parser.yacc import parse as ast_parse, extract_suggest, ParseError

app = Flask(__name__)


@app.route('/api/suggest', methods=['POST'])
def suggest_view():
    text = request.json.get('input') or ''
    cursor = request.json.get('cursor')
    try:
        result = ast_parse(text)
        suggests = [x for x in extract_suggest(result) if cursor is None or x.start <= cursor <= x.end]
        suggestions = {
            'host': ['localhost', '127.0.0.1', '192.168.1.1'],
            'port': ['80', '443', '8080', '8443', '9200'],
            'protocol': ['tcp', 'udp', 'http', 'https', 'rpc'],
            'server': ['nginx', 'apache', 'Apache-Tomcat'],
            'keyword': ['工商银行', '招商银行', '建设银行', '交通银行'],
        }
        rv = []
        for suggest in suggests:
            for key, key_suggestions in suggestions.items():
                for key_suggestion in key_suggestions:
                    if suggest.query in key_suggestion:
                        rv.append(dict(
                            suggest=f'{key}="{key_suggestion}"',
                            key=key,
                            query=suggest.query,
                            start=suggest.start,
                            end=suggest.end,
                        ))
        return jsonify(dict(data=rv, error=None))
    except Exception as exc:
        return jsonify(dict(data=None, error="Error: {!r}".format(exc)))


@app.route('/api/parse', methods=['POST'])
def parse_view():
    try:
        text = request.json.get('input') or ''
        result = ast_parse(text)
        return jsonify(dict(data=dict(result=repr(result)), error=None))
    except Exception as exc:
        return jsonify(dict(data=None, error="Error: {!r}".format(exc)))


if __name__ == '__main__':
    app.run(host='localhost', port=8765, debug=True)
