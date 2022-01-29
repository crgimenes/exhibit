package compiler

func newLine(t *Tokenizer) (token Token, ok bool, err error) {
	if t.Content[t.CurrentOffset] == '\n' {
		token = Token{
			Literal: "\n",
			Type:    TokenType("NEWLINE"),
			Offset:  t.CurrentOffset,
			Line:    t.CurrentLine,
			Column:  t.CurrentColumn,
		}
		t.CurrentOffset++
		t.CurrentLine++
		t.CurrentColumn = 0
		ok = true
	}

	return
}

func notImplemented(t *Tokenizer) (token Token, ok bool, err error) {
	tokenStr := ""
	token = Token{
		Literal: "",
		Type:    TokenType("NOT_IMPLEMENTED"),
		Offset:  t.CurrentOffset,
		Line:    t.CurrentLine,
		Column:  t.CurrentColumn,
	}

	for t.CurrentOffset <= t.Length {
		char := t.Content[t.CurrentOffset]

		t.CurrentColumn++
		t.CurrentOffset++

		if isSeparator(char) {
			break
		}

		tokenStr += string(char)
	}

	token.Literal = tokenStr
	ok = true

	return
}

func titleH1(t *Tokenizer) (token Token, ok bool, err error) {
	const (
		beginPattern    = "# "
		lenBeginPattern = len(beginPattern)
	)

	tokenStr := beginPattern
	token = Token{
		Literal: "",
		Type:    TokenType("TITLE_H1"),
		Offset:  t.CurrentOffset,
		Line:    t.CurrentLine,
		Column:  t.CurrentColumn,
	}

	if isPattern(t, tokenStr) && t.CurrentColumn == 0 {
		t.CurrentOffset += lenBeginPattern
		t.CurrentColumn = lenBeginPattern

		for t.CurrentOffset <= t.Length {
			char := t.Content[t.CurrentOffset]

			t.CurrentColumn++
			t.CurrentOffset++

			if char == '\n' {
				break
			}

			tokenStr += string(char)
		}

		token.Literal = tokenStr
		ok = true
	}

	return token, ok, nil
}

func titleH2(t *Tokenizer) (token Token, ok bool, err error) {
	const (
		beginPattern    = "## "
		lenBeginPattern = len(beginPattern)
	)

	tokenStr := beginPattern
	token = Token{
		Literal: "",
		Type:    TokenType("TITLE_H2"),
		Offset:  t.CurrentOffset,
		Line:    t.CurrentLine,
		Column:  t.CurrentColumn,
	}

	if isPattern(t, beginPattern) && t.CurrentColumn == 0 {
		t.CurrentOffset += lenBeginPattern
		t.CurrentColumn = lenBeginPattern

		for t.CurrentOffset <= t.Length {
			char := t.Content[t.CurrentOffset]

			t.CurrentColumn++
			t.CurrentOffset++

			if char == '\n' {
				break
			}

			tokenStr += string(char)
		}

		token.Literal = tokenStr
		ok = true
	}

	return token, ok, nil
}

func code(t *Tokenizer) (token Token, ok bool, err error) {
	const (
		beginPattern    = "```"
		endPattern      = "```"
		lenEndPattern   = len(endPattern)
		lenBeginPattern = len(beginPattern)
	)

	tokenStr := beginPattern
	token = Token{
		Literal: "",
		Type:    TokenType("CODE"),
		Offset:  t.CurrentOffset,
		Line:    t.CurrentLine,
		Column:  t.CurrentColumn,
	}

	if isPattern(t, tokenStr) && t.CurrentColumn == 0 {
		t.CurrentOffset += lenBeginPattern
		t.CurrentColumn = lenBeginPattern

		for t.CurrentOffset < t.Length {
			char := t.Content[t.CurrentOffset]

			t.CurrentColumn++
			t.CurrentOffset++

			if char == '\n' {
				t.CurrentLine++
				t.CurrentColumn = 0
			}

			tokenStr += string(char)

			if isPattern(t, endPattern) {
				tokenStr += endPattern
				t.CurrentOffset += lenEndPattern
				break
			}
		}

		token.Literal = tokenStr
		ok = true
	}

	return token, ok, nil
}
