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
	"errors"
	"quinto/data"
	"regexp"
	"strconv"
)

type queryFragment struct {
	txt string
	ord bool
	opt int
}

func extractParenthesis(query string, index *int, fragments *[]queryFragment) error {
	paren := query[*index]
	if paren == '(' {
		*fragments = append(*fragments, queryFragment{"(", false, 0})
		*index++
		return nil
	}
	if paren == ')' {
		*fragments = append(*fragments, queryFragment{")", false, 0})
		*index++
		return nil
	}
	return errors.New("invalid parenthesis")
}

func extractSimpleQueryFragment(query string, index *int, fragments *[]queryFragment) error {
	fragmentRegex := regexp.MustCompile(`([a-z]+)`)
	matches := fragmentRegex.FindStringSubmatch(query[*index:])
	if matches == nil {
		return errors.New("impossible match of simple query fragment")
	}
	result := queryFragment{
		txt: matches[1],
		ord: false,
		opt: 0,
	}
	*fragments = append(*fragments, result)
	*index += len(matches[0])
	return nil
}

func extractComplexQueryFragment(query string, index *int, fragments *[]queryFragment) error {
	fragmentRegex := regexp.MustCompile(`([A-Z]+)(?::([A-Z]+))?(?::(\d+))?`)
	matches := fragmentRegex.FindStringSubmatch(query[*index:])
	opt := 0
	if matches[3] != "" {
		var err error
		opt, err = strconv.Atoi(matches[3])
		if err != nil {
			return err
		}
	}
	allowedOperators := []string{"AND", "OR", "XOR", "NEAR", "NOT"}
	if !data.SliceToSet(allowedOperators).Contains(matches[1]) {
		return errors.New("invalid operator in complex query fragment: " + matches[1])
	}
	if matches[2] != "" && matches[2] != "ORD" {
		return errors.New("invalid order specifier in complex query fragment: " + matches[2])
	}
	result := queryFragment{
		txt: matches[1],
		ord: matches[2] == "ORD",
		opt: opt,
	}
	*fragments = append(*fragments, result)
	*index += len(matches[0])
	return nil
}

func SplitQuery(query string) ([]queryFragment, error) {
	var fragments []queryFragment
	var err error = nil
	for index := 0; index < len(query) && err == nil; {
		char := query[index]
		if char == ' ' || char == '\t' || char == '\n' || char == '\r' {
			index++
			continue
		}
		if char >= 'a' && char <= 'z' {
			err = extractSimpleQueryFragment(query, &index, &fragments)
			continue
		}
		if char >= 'A' && char <= 'Z' {
			err = extractComplexQueryFragment(query, &index, &fragments)
			continue
		}
		if char == '(' || char == ')' {
			err = extractParenthesis(query, &index, &fragments)
			continue
		}
		err = errors.New("invalid character in query: " + string(char))
	}
	return fragments, err
}
