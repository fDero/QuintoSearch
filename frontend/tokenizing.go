package frontend

import (
	"strings"
	"unicode"
)

func Split(input string) []string {
	var tokens []string
	var currentToken strings.Builder
	for _, char := range input {
		if unicode.IsSpace(char) || unicode.IsPunct(char) {
			if currentToken.Len() > 0 {
				tokens = append(tokens, currentToken.String())
				currentToken.Reset()
			}
		} else {
			currentToken.WriteRune(char)
		}
	}
	if currentToken.Len() > 0 {
		tokens = append(tokens, currentToken.String())
	}
	return tokens
}

func MakeLowerCase(tokens []string) []string {
	var lowerCasedTokens []string
	for _, txt := range tokens {
		lowerCaseTxt := strings.ToLower(txt)
		lowerCasedTokens = append(lowerCasedTokens, lowerCaseTxt)
	} 
	return lowerCasedTokens
}

func StopWordsFilter(tokens []string) []string {	
	stopwords := []string {
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
	}
	var filtered_tokens []string
	for _, txt := range tokens {
		isStopWord := false
        for _, stopword := range stopwords {
			if strings.ToLower(txt) == strings.ToLower(stopword) {
				isStopWord = true
				break
			}
		}
		if (!isStopWord) {
			filtered_tokens = append(filtered_tokens, txt)
		}
	}
	return filtered_tokens
}