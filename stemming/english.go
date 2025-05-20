package stemming

import (
	"quinto/data"
	"strings"
)

func stopWordsEnglish() data.Set[string] {
	return data.ToSet([]string{
		"i", "me", "my", "myself", "we", "our", "ours", "ourselves", "you", "your",
		"yours", "yourself", "yourselves", "he", "him", "his", "himself", "she", "her",
		"hers", "herself", "it", "its", "itself", "they", "them", "their", "theirs",
		"themselves", "what", "which", "who", "whom", "this", "that", "these", "those",
		"am", "is", "are", "was", "were", "be", "been", "being", "have", "has",
		"had", "having", "do", "does", "did", "doing", "a", "an", "the", "and",
		"but", "if", "or", "because", "as", "until", "while", "of", "at", "by",
		"for", "with", "about", "against", "between", "into", "through", "during",
		"before", "after", "above", "below", "to", "from", "up", "down", "in",
		"out", "on", "off", "over", "under", "again", "further", "then", "once",
		"here", "there", "when", "where", "why", "how", "all", "any", "both",
		"each", "few", "more", "most", "other", "some", "such", "no", "nor",
		"not", "only", "own", "same", "so", "than", "too", "very", "s", "t",
		"can", "will", "just", "don", "should", "now",
	})
}

func removePluralEnglish(text string) string {
	if strings.HasSuffix(text, "s") {
		if hasVowelBeforeLastNChars(text, 2) {
			return strings.TrimSuffix(text, "s")
		} else {
			return text
		}
	}
	return text
}

func removeSimpleSuffixEnglish(text string) string {
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

func stemEnglish(text string) string {
	text = removeSimpleSuffixEnglish(text)
	text = removePluralEnglish(text)
	text = removeLastVowel(text)
	text = removeLastRepeatedLetter(text)
	return text
}
