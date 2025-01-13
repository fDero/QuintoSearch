package frontend

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

func removeLastS(text string) string {
	if strings.HasSuffix(text, "s") {
		if hasVowelBeforeLastNChars(text, 2) {
			return strings.TrimSuffix(text, "s")
		} else {
			return text
		}
	}
	return text
}

func removeLastRepeatedConsonant(text string) string {
	if len(text) > 2 {
		lastChar := string(text[len(text)-1])
		prevChar := string(text[len(text)-2])
		if lastChar == prevChar {
			return strings.TrimSuffix(text, lastChar)
		}
	}
	return text
}

func removeSimpleSuffix(text string) string {
	return matchAndReplace(text, []stemmingPattern{
		{suffix: "fulness", replacement: "", minLen: 0, maxLen: 0},
		{suffix: "lyhood", replacement: "", minLen: 8, maxLen: 0},
		{suffix: "tional", replacement: "", minLen: 0, maxLen: 0},
		{suffix: "erable", replacement: "", minLen: 0, maxLen: 0},
		{suffix: "ingly", replacement: "", minLen: 0, maxLen: 0},
		{suffix: "alism", replacement: "", minLen: 0, maxLen: 0},
		{suffix: "cation", replacement: "", minLen: 8, maxLen: 0},
		{suffix: "cated", replacement: "", minLen: 8, maxLen: 0},
		{suffix: "ative", replacement: "", minLen: 0, maxLen: 0},
		{suffix: "eedly", replacement: "", minLen: 7, maxLen: 0},
		{suffix: "sses", replacement: "", minLen: 4, maxLen: 0},
		{suffix: "edly", replacement: "", minLen: 6, maxLen: 0},
		{suffix: "iful", replacement: "", minLen: 0, maxLen: 0},
		{suffix: "ened", replacement: "", minLen: 0, maxLen: 0},
		{suffix: "hood", replacement: "", minLen: 8, maxLen: 0},
		{suffix: "like", replacement: "", minLen: 8, maxLen: 0},
		{suffix: "ling", replacement: "", minLen: 8, maxLen: 0},
		{suffix: "ical", replacement: "", minLen: 9, maxLen: 0},
		{suffix: "ness", replacement: "", minLen: 0, maxLen: 0},
		{suffix: "ful", replacement: "", minLen: 0, maxLen: 0},
		{suffix: "ern", replacement: "", minLen: 4, maxLen: 0},
		{suffix: "ies", replacement: "", minLen: 0, maxLen: 4},
		{suffix: "oid", replacement: "", minLen: 6, maxLen: 4},
		{suffix: "ery", replacement: "", minLen: 5, maxLen: 0},
		{suffix: "ion", replacement: "", minLen: 5, maxLen: 0},
		{suffix: "ied", replacement: "", minLen: 4, maxLen: 0},
		{suffix: "ial", replacement: "", minLen: 5, maxLen: 0},
		{suffix: "est", replacement: "", minLen: 5, maxLen: 0},
		{suffix: "ing", replacement: "", minLen: 5, maxLen: 0},
		{suffix: "ies", replacement: "", minLen: 0, maxLen: 5},
		{suffix: "or", replacement: "", minLen: 0, maxLen: 0},
		{suffix: "ed", replacement: "", minLen: 5, maxLen: 0},
		{suffix: "en", replacement: "", minLen: 5, maxLen: 0},
		{suffix: "es", replacement: "", minLen: 0, maxLen: 5},
		{suffix: "es", replacement: "", minLen: 5, maxLen: 0},
		{suffix: "al", replacement: "", minLen: 5, maxLen: 0},
		{suffix: "ly", replacement: "", minLen: 5, maxLen: 0},
		{suffix: "er", replacement: "", minLen: 5, maxLen: 0},
		{suffix: "ed", replacement: "", minLen: 5, maxLen: 0},
		{suffix: "ee", replacement: "", minLen: 5, maxLen: 0},
	})
}

func normalizeSingleToken(text string) string {
	text = removeSimpleSuffix(text)
	text = removeLastS(text)
	text = removeLastVowel(text)
	text = removeLastRepeatedConsonant(text)
	return text
}

func NormalizeTokens(tokens []string) []string {
	var normalizedTokens []string
	for _, txt := range tokens {
		normtxt := normalizeSingleToken(txt)
		normalizedTokens = append(normalizedTokens, normtxt)
	}
	return normalizedTokens
}
