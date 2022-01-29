package compiler

import (
	"fmt"
)

type TokenType string

type Token struct {
	Literal string
	Type    TokenType
	Offset  int
	Line    int
	Column  int
}

type Tokenizer struct {
	Content       []rune
	CurrentOffset int
	CurrentLine   int
	CurrentColumn int
	Length        int
	Parse         []func(*Tokenizer) (token Token, ok bool, err error)
}

func (t Token) String() string {
	return fmt.Sprintf("%02d:%02d %s(%q)", t.Line, t.Column, t.Type, t.Literal)
}

func NewTokenizer(r []rune) *Tokenizer {
	return &Tokenizer{
		Content:       r,
		CurrentOffset: 0,
		CurrentLine:   0,
		CurrentColumn: 0,
		Length:        len(r),
		Parse: []func(*Tokenizer) (token Token, ok bool, err error){
			newLine,
			titleH1,
			titleH2,
			code,
			notImplemented,
		},
	}
}

func (t *Tokenizer) Tokenize() ([]Token, error) {
	tokens := []Token{}

	for {
		for _, p := range t.Parse {
			token, ok, err := p(t)
			if err != nil {
				return nil, err
			}

			if ok {
				tokens = append(tokens, token)
				break
			}
		}

		if t.CurrentOffset >= t.Length {
			break
		}
	}

	return tokens, nil
}

func isSeparator(c rune) bool {
	return c == ' ' || c == '\t' || c == '\n' || c == '\r'
}

func isPattern(t *Tokenizer, pattern string) (ok bool) {
	lenPattern := len(pattern)
	if t.CurrentOffset+lenPattern > t.Length {
		return false
	}

	return string(t.Content[t.CurrentOffset:t.CurrentOffset+lenPattern]) == pattern
}
