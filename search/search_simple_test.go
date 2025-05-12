package search

import (
	"iter"
	"quinto/misc"
	"testing"
)

func TestDocumentIteration(t *testing.T) {
	documentIterator := misc.NewSliceIterator(helloWorldDocument)
	next, close := iter.Pull(documentIterator)
	_, exists := next()
	defer close()
	if !exists {
		t.Fatalf("Expected to get a value, but got none")
	}
}

func TestTermIteration(t *testing.T) {
	documentIterator := misc.NewSliceIterator(helloWorldDocument)
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
	tokenIterator := misc.NewSliceIterator(helloWorldDocument)
	index.StoreNewDocument(tokenIterator)

	queryFragments := SplitQuery(queryString)
	query, err := ParseQuery(queryFragments)
	if err != nil {
		t.Fatalf("Failed to parse query: %v", err)
	}

	query.Init(index)
	queryResult := query.Run()
	if queryResult.success != success {
		t.Errorf("Expected success to be true, got false")
	}
}

func TestSearchOverallForOneWord(t *testing.T) {
	runTestQueryHelper(t, "hello", true)
}

func TestSearchOverallForTwoWordSuccess(t *testing.T) {
	index := NewNaiveReverseIndex()
	tokenIterator := misc.NewSliceIterator(helloWorldDocument)
	index.StoreNewDocument(tokenIterator)
	runTestQueryHelper(t, "hello AND world", true)
}

func TestTermIterationFailure(t *testing.T) {
	index := NewNaiveReverseIndex()
	tokenIterator := misc.NewSliceIterator(helloWorldDocument)
	index.StoreNewDocument(tokenIterator)
	runTestQueryHelper(t, "guitar", false)
}

func TestSearchOverallForTwoWordFailure(t *testing.T) {
	index := NewNaiveReverseIndex()
	tokenIterator := misc.NewSliceIterator(helloWorldDocument)
	index.StoreNewDocument(tokenIterator)
	runTestQueryHelper(t, "guitar AND world", false)
}
