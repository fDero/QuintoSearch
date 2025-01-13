package frontend

import (
	"bufio"
    "log"
    "os"
)

func ProcessInputText(inputText string) []string {
	var splitted []string = Split(inputText)
	var lowerCased []string = MakeLowerCase(splitted)
	var filtered []string = StopWordsFilter(lowerCased)
	var normalized []string = NormalizeTokens(filtered)
	return normalized
}

func ProcessInputFile(inputFilePath string) []string {
	file, err := os.Open(inputFilePath)
    
	if err != nil {
        log.Fatal(err)
    }
    
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var tokens []string

	for scanner.Scan() {
		currentLineTokens := ProcessInputText(scanner.Text())
		tokens = append(tokens, currentLineTokens...)
    }

	if err := scanner.Err(); err != nil {
        log.Fatal(err)
    }	

	return tokens
}