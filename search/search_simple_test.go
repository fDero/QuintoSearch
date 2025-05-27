package search

import (
	"iter"
	"quinto/data"
	"testing"
)

func TestDocumentIteration(t *testing.T) {
	documentIterator := data.NewSliceIterator(helloWorldDocument)
	next, close := iter.Pull(documentIterator)
	_, exists := next()
	defer close()
	if !exists {
		t.Fatalf("Expected to get a value, but got none")
	}
}

func TestTermIteration(t *testing.T) {
	documentIterator := data.NewSliceIterator(helloWorldDocument)
	index := NewNaiveReverseIndex()
	index.StoreNewDocument(documentIterator)
	termIterator := index.IterateOverTerms("hello")
	next, close := iter.Pull(termIterator)
	_, exists := next()
	defer close()
	if !exists {
		t.Fatalf("Expected to get a value, but got none")
	}
}

func runTestQueryHelper(t *testing.T, queryString string, success bool) {
	index := NewNaiveReverseIndex()
	tokenIterator := data.NewSliceIterator(helloWorldDocument)
	index.StoreNewDocument(tokenIterator)

	queryFragments, err1 := SplitQuery(queryString)
	query, err2 := ParseQuery(queryFragments)
	if err1 != nil || err2 != nil {
		t.Fatalf("Failed to parse query: %v %v", err1, err2)
	}

	query.Init(index)
	queryResult := query.Run()
	if queryResult.Success != success {
		t.Errorf("Expected success to be true, got false")
	}
}

func TestSearchOverallForOneWord(t *testing.T) {
	runTestQueryHelper(t, "hello", true)
}

func TestSearchOverallForTwoWordSuccess(t *testing.T) {
	index := NewNaiveReverseIndex()
	tokenIterator := data.NewSliceIterator(helloWorldDocument)
	index.StoreNewDocument(tokenIterator)
	runTestQueryHelper(t, "hello AND world", true)
}

func TestTermIterationFailure(t *testing.T) {
	index := NewNaiveReverseIndex()
	tokenIterator := data.NewSliceIterator(helloWorldDocument)
	index.StoreNewDocument(tokenIterator)
	runTestQueryHelper(t, "guitar", false)
}

func TestSearchOverallForTwoWordFailure(t *testing.T) {
	index := NewNaiveReverseIndex()
	tokenIterator := data.NewSliceIterator(helloWorldDocument)
	index.StoreNewDocument(tokenIterator)
	runTestQueryHelper(t, "guitar AND world", false)
}
