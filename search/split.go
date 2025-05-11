/*=================================== LICENSE =======================================

                                   Apache License
                             Version 2.0, January 2004
                          http://www.apache.org/licenses/

============================== BRIEF FILE DESCRIPTION ===============================

This file contains the implementation of the SplitQuery function, which is responsible
for splitting a query string into its constituent fragments. The function uses regular
expressions to identify different types of fragments, including simple terms, complex
queries, and parentheses. The fragments are stored in a slice of queryFragment structs,
which contain the text of the fragment, a boolean indicating whether, in case of 
complex queries, the order is important, and an integer for additional options.
==================================================================================*/

package search

import (
	"regexp"
	"strconv"
)

type queryFragment struct {
	txt string
	ord bool
	opt int
}

func extractParenthesis(query string, index int, fragments []queryFragment) (int, []queryFragment) {
	if query[index] == '(' {
		fragments = append(fragments, queryFragment{"(", false, 0})
		index++
	} else {
		fragments = append(fragments, queryFragment{")", false, 0})
		index++
	}
	return index, fragments
}

func extractSimpleQueryFragment(query string, index int, fragments []queryFragment) (int, []queryFragment) {
	fragmentRegex := regexp.MustCompile(`([a-z]+)`)
	matches := fragmentRegex.FindStringSubmatch(query[index:])
	if matches == nil {
		return index, fragments
	}
	result := queryFragment{
		txt: matches[1],
		ord: false,
		opt: 0,
	}
	fragments = append(fragments, result)
	index += len(matches[0])
	return index, fragments
}
func extractComplexQueryFragment(query string, index int, fragments []queryFragment) (int, []queryFragment) {
	fragmentRegex := regexp.MustCompile(`([A-Z]+)(?::([A-Z]+))?(?::(\d+))?`)
	matches := fragmentRegex.FindStringSubmatch(query[index:])

	opt := 0
	if matches[3] != "" {
		var err error
		opt, err = strconv.Atoi(matches[3])
		if err != nil {
			return index, fragments
		}
	}

	result := queryFragment{
		txt: matches[1],
		ord: matches[2] == "ORD",
		opt: opt,
	}

	fragments = append(fragments, result)
	index += len(matches[0])
	return index, fragments
}

func SplitQuery(query string) []queryFragment {
	var fragments []queryFragment
	for index := 0; index < len(query); {
		switch query[index] {
		case '(', ')':
			index, fragments = extractParenthesis(query, index, fragments)
		case 'A', 'O', 'N', 'X':
			index, fragments = extractComplexQueryFragment(query, index, fragments)
		case ' ', '\t', '\n', '\r':
			index++
			continue
		default:
			index, fragments = extractSimpleQueryFragment(query, index, fragments)
		}
	}
	return fragments
}
