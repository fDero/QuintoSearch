package search

import (
	"testing"
)

func TestSplitQuery(t *testing.T) {
	queryString := "a OR (b NEAR:ORD:10 c)"
	fragments, err := SplitQuery(queryString)

	if err != nil {
		t.Fatalf("Failed to split query: %v", err)
	}

	if len(fragments) != 7 {
		t.Fatalf("Expected 7 fragments, got %d", len(fragments))
	}

	if fragments[0].txt != "a" {
		t.Errorf("Expected fragment 'a', got '%s'", fragments[0].txt)
	}

	if fragments[1].txt != "OR" {
		t.Errorf("Expected fragment 'OR', got '%s'", fragments[1].txt)
	}

	if fragments[2].txt != "(" {
		t.Errorf("Expected fragment '(', got '%s'", fragments[2].txt)
	}

	if fragments[3].txt != "b" {
		t.Errorf("Expected fragment 'b', got '%s'", fragments[3].txt)
	}

	if fragments[4].txt != "NEAR" {
		t.Errorf("Expected fragment 'NEAR', got '%s'", fragments[4].txt)
	}

	if fragments[5].txt != "c" {
		t.Errorf("Expected fragment 'c', got '%s'", fragments[5].txt)
	}

	if fragments[6].txt != ")" {
		t.Errorf("Expected fragment ')', got '%s'", fragments[6].txt)
	}
}

func TestComplexSplitting(t *testing.T) {
	queryString := "NEAR:ORD:10 AND:ORD NEAR:3"
	fragments, err := SplitQuery(queryString)

	if err != nil {
		t.Fatalf("Failed to split query: %v", err)
	}

	if len(fragments) != 3 {
		t.Fatalf("Expected 3 fragments, got %d", len(fragments))
	}

	if fragments[0].txt != "NEAR" {
		t.Errorf("Expected fragment 'NEAR', got '%s'", fragments[0].txt)
	}

	if !fragments[0].ord {
		t.Error("Expected fragment 'NEAR' to be ordered")
	}

	if fragments[0].opt != 10 {
		t.Errorf("Expected fragment 'NEAR' to have option = 10, instead got %v", fragments[0].opt)
	}

	if fragments[1].txt != "AND" {
		t.Errorf("Expected fragment 'AND', got '%s'", fragments[0].txt)
	}

	if !fragments[1].ord {
		t.Error("Expected fragment 'AND' to be ordered")
	}

	if fragments[1].opt != 0 {
		t.Errorf("Expected fragment 'NEAR' to have option = 0, instead got %v", fragments[0].opt)
	}

	if fragments[2].txt != "NEAR" {
		t.Errorf("Expected fragment 'NEAR', got '%s'", fragments[0].txt)
	}

	if fragments[2].ord {
		t.Error("Expected fragment 'NEAR2' to be unordered")
	}

	if fragments[2].opt != 3 {
		t.Errorf("Expected fragment 'NEAR' to have option = 3, instead got %v", fragments[2].opt)
	}
}

func TestBadSplittingErrorRecovery(t *testing.T) {
	queryString := "mike ZAPP robert"
	fragments, err := SplitQuery(queryString)

	if err == nil {
		t.Logf("Parsed fragments: %v", fragments)
		t.Fatalf("Expected an error due to bad query syntax, but got none")
	}

	queryString = "mike NEAR:GROBB robert"
	fragments, err = SplitQuery(queryString)

	if err == nil {
		t.Logf("Parsed fragments: %v", fragments)
		t.Fatalf("Expected an error due to bad query syntax, but got none")
	}

	queryString = "mike NEAR:Z:10 robert"
	fragments, err = SplitQuery(queryString)

	if err == nil {
		t.Logf("Parsed fragments: %v", fragments)
		t.Fatalf("Expected an error due to bad query syntax, but got none")
	}

	
	queryString = "mike NEAR:ORD:10:B robert"
	fragments, err = SplitQuery(queryString)

	if err == nil {
		t.Logf("Parsed fragments: %v", fragments)
		t.Fatalf("Expected an error due to bad query syntax, but got none")
	}
}
