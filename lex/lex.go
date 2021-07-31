package lex

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
)

// TokenFunction type for parser functions
type tokenFunction func(l *lexer) (ok bool, err error)

// Token data
type Token struct {
	Type    string
	Literal string
	Line    int
	Column  int
	Offset  int
}

// lexer data
type lexer struct {
	tokRune       rune
	tokens        []Token
	currentLine   int
	currentColumn int
	offset        int
	tokenParsers  []tokenFunction
	reader        io.RuneScanner
}

// Parse is a simple tokenizer engine
func Parse(reader io.Reader) (tokens []Token, err error) {
	tokens, err = parseSimpleTokens(reader)
	if err != nil && err != io.EOF {
		return
	}
	tokens, err = parseIdentifer(tokens)
	return
}

func parseIdentifer(tokens []Token) (ret []Token, err error) {
	for i := range tokens {
		if tokens[i].Type == "IDENTIFER" {
			if isNumeric(tokens[i].Literal) {
				tokens[i].Type = "NUMBER"
			}
		}
	}
	ret = tokens
	return
}

/*
parseSimpleTokens is the first stage of converting files
and/or strings to a list of tokens, at this stage we only
deal with the most simple tokens, comments, comment
with multiple lines, begin and end of lists, strings,
separators and identifiers.
*/
func parseSimpleTokens(reader io.Reader) (tokens []Token, err error) {
	var l = lexer{
		tokenParsers: []tokenFunction{
			separatorTok,
			singleLineCommentTok,
			multipleLineCommentTok,
			listBeginTok,
			listEndTok,
			stringTok,
			identiferTok,
			notImplemented,
		},
		reader: bufio.NewReader(reader),
	}

	var currentLine, currentColumn int
	for {
		var size int
		l.tokRune, size, err = l.reader.ReadRune()
		if err == io.EOF {
			break
		}
		l.offset += size
		if l.tokRune == '\n' {
			currentLine++
			currentColumn = 0
		}
		l.currentLine = currentLine
		l.currentColumn = currentColumn
		currentColumn++
		for _, f := range l.tokenParsers {
			var ok bool
			ok, err = f(&l)
			if err != nil && err != io.EOF {
				return
			}
			if ok {
				break
			}
		}
	}
	tokens = l.tokens
	return
}

// appendTok add a new token to current token list
func (l *lexer) appendTok(literalValue, tokType string) {
	t := Token{Type: tokType,
		Literal: literalValue,
		Line:    l.currentLine,
		Offset:  l.offset,
		Column:  l.currentColumn,
	}
	l.tokens = append(l.tokens, t)
}

func separatorTok(l *lexer) (ok bool, err error) {
	switch l.tokRune {
	case '\n', '\r', ' ', '\t':
		l.appendTok(string(l.tokRune), "SEPARATOR")
		ok = true
	}
	return
}

func isIdentiferSeparator(v rune) (ok bool) {
	switch v {
	case '\n', '\t', '\v', '\f', '\r', ' ', 0x85, 0xA0:
		ok = true
	}
	return
}

func singleLineCommentTok(l *lexer) (ok bool, err error) {
	ok, err = l.ifRunePattern([]rune{'/', '/'})
	if !ok || err != nil {
		return
	}
	if l.tokRune == '/' {
		var value = "/"
		for l.tokRune != '\n' {
			var size int
			l.tokRune, size, err = l.reader.ReadRune()
			if err != nil {
				if err == io.EOF {
					break
				}
				return
			}
			value += string(l.tokRune)
			l.offset += size
		}
		ok = true
		l.appendTok(value, "SINGLE-LINE-COMMENT")
	}
	return
}

func notImplemented(l *lexer) (ok bool, err error) {
	panic(fmt.Sprintf("%q NOT IMPLEMENTED\n", l.tokRune))
}

func listBeginTok(l *lexer) (ok bool, err error) {
	if ok = l.tokRune == '('; ok {
		l.appendTok("(", "LIST-BEGIN")
	}
	return
}

func listEndTok(l *lexer) (ok bool, err error) {
	if ok = l.tokRune == ')'; ok {
		l.appendTok(")", "LIST-END")
	}
	return
}

func isNumeric(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}

func stringTok(l *lexer) (ok bool, err error) {
	if ok = l.tokRune == '"'; ok {
		var (
			value    string
			lastRune rune
		)
		for {
			var size int
			l.tokRune, size, err = l.reader.ReadRune()
			if err == io.EOF {
				break
			}
			l.offset += size
			if lastRune != '\\' && l.tokRune == '"' {
				break
			}
			value += string(l.tokRune)
			lastRune = l.tokRune
		}
		l.appendTok(value, "STRING")
	}
	return
}

func identiferTok(l *lexer) (ok bool, err error) {
	var identifer string
	for !isIdentiferSeparator(l.tokRune) {
		var size int
		identifer += string(l.tokRune)
		ok = true
		l.tokRune, size, err = l.reader.ReadRune()
		if err == io.EOF {
			break
		}
		if isIdentiferSeparator(l.tokRune) {
			err = l.reader.UnreadRune()
			break
		}
		l.offset += size
	}
	if ok {
		l.appendTok(identifer, "IDENTIFER")
	}
	return
}

func (l *lexer) peekRune() (v rune, err error) {
	v, _, err = l.reader.ReadRune()
	if err != nil {
		return
	}
	err = l.reader.UnreadRune()
	return
}

func (l *lexer) ifRunePattern(pattern []rune) (ok bool, err error) {
	if l.tokRune != pattern[0] {
		return
	}
	p, err := l.peekRune()
	ok = p == pattern[1]
	return
}

func multipleLineCommentTok(l *lexer) (ok bool, err error) {
	ok, err = l.ifRunePattern([]rune{'/', '*'})
	if !ok || err != nil {
		return
	}
	var value = "/"
	for {
		var size int
		l.tokRune, size, err = l.reader.ReadRune()
		if err != nil {
			return
		}
		value += string(l.tokRune)
		l.offset += size
		ok, err = l.ifRunePattern([]rune{'*', '/'})
		if err != nil {
			return
		}
		if ok {
			value += "/"
			_, _, err = l.reader.ReadRune()
			break
		}
	}
	l.appendTok(value, "MULTIPLE-LINE-COMMENT")
	return
}
