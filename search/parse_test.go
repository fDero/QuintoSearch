package search

import (
	"testing"
)

func TestParseExactQuery(t *testing.T) {

	// a [exact match]
	fragments := []queryFragment{
		{"a", false, 0},
	}

	query, err := ParseQuery(fragments)
	if err != nil {
		t.Fatalf("ParseQuery failed: %v", err)
	}

	if query.(*ExactQuery).term != "a" {
		t.Errorf("Expected term 'a', got '%s'", query.(*ExactQuery).term)
	}
}

func TestParseAndQuery(t *testing.T) {

	// a AND b
	fragments := []queryFragment{
		{"a", false, 0},
		{"AND", false, 0},
		{"b", false, 0},
	}

	query, err := ParseQuery(fragments)
	if err != nil {
		t.Fatalf("ParseQuery failed: %v", err)
	}

	if query.(*ComplexQuery) == nil {
		t.Errorf("Expected complex query, got something else")
	}
}

func TestParseOrThenAndQuery(t *testing.T) {

	// a OR [b AND c]
	fragments := []queryFragment{
		{"a", false, 0},
		{"OR", false, 0},
		{"b", false, 0},
		{"AND", false, 0},
		{"c", false, 0},
	}

	query, err := ParseQuery(fragments)
	if err != nil {
		t.Fatalf("ParseQuery failed: %v", err)
	}

	if query.(*ComplexQuery).rx.(*ComplexQuery) == nil {
		t.Errorf("Expected complex query with complex query on the right, got something else")
	}

	if term := query.(*ComplexQuery).lx.(*ExactQuery).term; term != "a" {
		t.Errorf("Expected complex query with exact match 'a' on the right, got something else: %v", term)
	}
}

func TestParseAndThenOrQuery(t *testing.T) {

	// [a AND b] OR c
	fragments := []queryFragment{
		{"a", false, 0},
		{"AND", false, 0},
		{"b", false, 0},
		{"OR", false, 0},
		{"c", false, 0},
	}

	query, err := ParseQuery(fragments)
	if err != nil {
		t.Fatalf("ParseQuery failed: %v", err)
	}

	if term := query.(*ComplexQuery).rx.(*ExactQuery).term; term != "c" {
		t.Errorf("Expected complex query with exact match 'c' on the right, got something else: %v", term)
	}

	if query.(*ComplexQuery).lx.(*ComplexQuery) == nil {
		t.Errorf("Expected complex query with complex query on the right, got something else")
	}
}

func TestParseUselessParenAndThenOrQuery(t *testing.T) {

	// (a AND b) OR c
	fragments := []queryFragment{
		{"(", false, 0},
		{"a", false, 0},
		{"AND", false, 0},
		{"b", false, 0},
		{")", false, 0},
		{"OR", false, 0},
		{"c", false, 0},
	}

	query, err := ParseQuery(fragments)
	if err != nil {
		t.Fatalf("ParseQuery failed: %v", err)
	}

	if term := query.(*ComplexQuery).rx.(*ExactQuery).term; term != "c" {
		t.Errorf("Expected complex query with exact match 'c' on the right, got something else: %v", term)
	}

	if query.(*ComplexQuery).lx.(*ComplexQuery) == nil {
		t.Errorf("Expected complex query with complex query on the right, got something else")
	}
}

func TestParseUsefulParenAndThenOrQuery(t *testing.T) {

	// a AND (b OR c)
	fragments := []queryFragment{
		{"a", false, 0},
		{"AND", false, 0},
		{"(", false, 0},
		{"b", false, 0},
		{"OR", false, 0},
		{"c", false, 0},
		{")", false, 0},
	}

	query, err := ParseQuery(fragments)
	if err != nil {
		t.Fatalf("ParseQuery failed: %v", err)
	}

	if term := query.(*ComplexQuery).lx.(*ExactQuery).term; term != "a" {
		t.Errorf("Expected complex query with exact match 'a' on the right, got something else: %v", term)
	}

	if query.(*ComplexQuery).rx.(*ComplexQuery) == nil {
		t.Errorf("Expected complex query with complex query on the right, got something else")
	}
}

func TestParseUselessParenOrThenAndQuery(t *testing.T) {

	// (a OR b) AND c
	fragments := []queryFragment{
		{"(", false, 0},
		{"a", false, 0},
		{"AND", false, 0},
		{"b", false, 0},
		{")", false, 0},
		{"OR", false, 0},
		{"c", false, 0},
	}

	query, err := ParseQuery(fragments)
	if err != nil {
		t.Fatalf("ParseQuery failed: %v", err)
	}

	if term := query.(*ComplexQuery).rx.(*ExactQuery).term; term != "c" {
		t.Errorf("Expected complex query with exact match 'c' on the right, got something else: %v", term)
	}

	if query.(*ComplexQuery).lx.(*ComplexQuery) == nil {
		t.Errorf("Expected complex query with complex query on the right, got something else")
	}
}

func TestParseUsefulParenOrThenAndQuery(t *testing.T) {

	// a OR (b AND c)
	fragments := []queryFragment{
		{"a", false, 0},
		{"OR", false, 0},
		{"(", false, 0},
		{"b", false, 0},
		{"AND", false, 0},
		{"c", false, 0},
		{")", false, 0},
	}

	query, err := ParseQuery(fragments)
	if err != nil {
		t.Fatalf("ParseQuery failed: %v", err)
	}

	if term := query.(*ComplexQuery).lx.(*ExactQuery).term; term != "a" {
		t.Errorf("Expected complex query with exact match 'a' on the right, got something else: %v", term)
	}

	if query.(*ComplexQuery).rx.(*ComplexQuery) == nil {
		t.Errorf("Expected complex query with complex query on the right, got something else")
	}
}
