package tokenizer

import (
	"strings"
	"unicode"
)

func Split(input string) []string {
	var tokens []string
	var currentToken strings.Builder

	for _, r := range input {
		if unicode.IsSpace(r) || unicode.IsPunct(r) {
			if currentToken.Len() > 0 {
				tokens = append(tokens, currentToken.String())
				currentToken.Reset()
			}
		} else {
			currentToken.WriteRune(r)
		}
	}
	if currentToken.Len() > 0 {
		tokens = append(tokens, currentToken.String())
	}

	return tokens
} 