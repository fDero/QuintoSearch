package frontend

import (
	"strings"
)

type stemmingPattern struct {
	suffix string
	replacement string
	minLen int
	maxLen int
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
		valid := index < len(text) - lastChars
		isVowel := strings.ContainsRune(vowels, char)
		if valid && isVowel {
			return true
		}
	}
	return false
}

func removeDoubleConsonantAwareSuffix(text string, suffix string) string {
	if strings.HasSuffix(text, suffix) {
		textLen := len(text)
		suffLen := len(suffix)
		if textLen > suffLen + 2 && text[textLen - suffLen + 1] == text[textLen - suffLen] {
			prev := string(text[textLen - suffLen])
			trimmed := strings.TrimSuffix(text, prev + suffix)
			return trimmed
		}
		if textLen > suffLen {
			return strings.TrimSuffix(text, suffix)
		}
	}
	return text
}

func removeLastVowel(text string) string {
	vowels := "aeiouyAEIOUY"
	lastRune := rune(text[len(text) - 1])
	lastChar := string(text[len(text) - 1])
	if strings.ContainsRune(vowels, lastRune) {
		return strings.TrimSuffix(text, lastChar)
	}
	return text
}

func removePlural(text string) string {
	pluarPatterns := []stemmingPattern {
		{suffix: "sses", replacement: "ss", minLen: 4, maxLen: 0},
		{suffix: "ies",  replacement: "y",  minLen: 0, maxLen: 4},
		{suffix: "es",   replacement: "e",  minLen: 0, maxLen: 5},
		{suffix: "es",   replacement: "",   minLen: 5, maxLen: 0},
	}
	oldText := text
	text = matchAndReplace(text, pluarPatterns)
	if oldText == text && strings.HasSuffix(text, "s") {
		if hasVowelBeforeLastNChars(text, 2) {
			return strings.TrimSuffix(text, "s")
		} else {
			return text
		}
	}
	return text
}

func removeYbasedSuffixes(text string) string {
	return matchAndReplace(text, []stemmingPattern {
		{suffix: "eedly", replacement: "ee",  minLen: 7, maxLen: 0},
		{suffix: "edly",  replacement: "ed",  minLen: 6, maxLen: 0},
		{suffix: "ingly", replacement: "ing", minLen: 0, maxLen: 0},
	})
}

func removePast(text string) string {
	return matchAndReplace(text, []stemmingPattern {
		{suffix: "ied", replacement: "y",  minLen: 4, maxLen: 0},
		{suffix: "ies", replacement: "ie", minLen: 0, maxLen: 5},
		{suffix: "ed",  replacement: "",   minLen: 5, maxLen: 0},
	})
}

func removeIngSuffix(text string) string {
	return removeDoubleConsonantAwareSuffix(text, "ing")
}

func removeErSuffix(text string) string {
	return removeDoubleConsonantAwareSuffix(text, "er")
}

func removeComplexSuffix(text string) string {
	return matchAndReplace(text, []stemmingPattern {
		{suffix: "tional",   replacement: "tion", minLen: 0, maxLen: 0},
		{suffix: "iful",     replacement: "y",    minLen: 0, maxLen: 0},
		{suffix: "ened",     replacement: "",     minLen: 0, maxLen: 0},
		{suffix: "alism",    replacement: "",     minLen: 0, maxLen: 0},
		{suffix: "fulness",  replacement: "",     minLen: 0, maxLen: 0},
		{suffix: "ful",      replacement: "",     minLen: 0, maxLen: 0},
		{suffix: "ness",     replacement: "",     minLen: 0, maxLen: 0},
		{suffix: "or",       replacement: "",     minLen: 0, maxLen: 0},
		{suffix: "ative",    replacement: "",     minLen: 0, maxLen: 0},
		{suffix: "ery",      replacement: "",     minLen: 5, maxLen: 0},
		{suffix: "ion",      replacement: "",     minLen: 5, maxLen: 0},
	})
}

func normalizeSingleToken(token string) string {
	token = removeComplexSuffix(token)
	token = removePlural(token)
	token = removeComplexSuffix(token)
	token = removeYbasedSuffixes(token)
	token = removePast(token)
	token = removeIngSuffix(token)
	token = removeErSuffix(token)
	token = removeLastVowel(token)
	return token
}

func NormalizeTokens(tokens []string) []string {
	var normalizedTokens []string
	for _, txt := range tokens {
		normtxt := normalizeSingleToken(txt)
		normalizedTokens = append(normalizedTokens, normtxt)
	}
	return normalizedTokens
}