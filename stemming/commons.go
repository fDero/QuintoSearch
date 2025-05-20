package stemming

import (
	"strings"
)

type stemmingPattern struct {
	suffix      string
	replacement string
	minLen      int
	maxLen      int
}

func matchAndReplace(text string, patterns []stemmingPattern) string {
	for _, pattern := range patterns {
		textLen := len(text)
		if textLen < pattern.minLen {
			continue
		}
		if pattern.maxLen > 0 && textLen > pattern.maxLen {
			continue
		}
		if strings.HasSuffix(text, pattern.suffix) {
			return strings.TrimSuffix(text, pattern.suffix) + pattern.replacement
		}
	}
	return text
}

func hasVowelBeforeLastNChars(text string, lastChars int) bool {
	vowels := "aeiouyAEIOUY"
	for index, char := range text {
		valid := index < len(text)-lastChars
		isVowel := strings.ContainsRune(vowels, char)
		if valid && isVowel {
			return true
		}
	}
	return false
}

func removeLastVowel(text string) string {
	vowels := "aeiouyAEIOUY"
	lastRune := rune(text[len(text)-1])
	lastChar := string(text[len(text)-1])
	if strings.ContainsRune(vowels, lastRune) {
		return strings.TrimSuffix(text, lastChar)
	}
	return text
}

func removeLastRepeatedLetter(text string) string {
	if len(text) > 2 {
		lastChar := string(text[len(text)-1])
		prevChar := string(text[len(text)-2])
		if lastChar == prevChar {
			return strings.TrimSuffix(text, lastChar)
		}
	}
	return text
}
