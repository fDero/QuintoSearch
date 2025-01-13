package indexing

import (
	"testing"
)

func TestGenerateCursorsMapOnEmptyIndex(t *testing.T) {
	emptyIndex := InvertedIndex{invertedLists: make(map[string]*node)}
	cursorsMap := emptyIndex.generateCursorsMap([]string{"A", "B", "C"})
	if cursorsMapLen := len(cursorsMap); cursorsMapLen != 3 {
		t.Errorf("Expected 3 cursors, got %d", cursorsMapLen)
	}
	for term, cursor := range cursorsMap {
		if term != "A" && term != "B" && term != "C" {
			t.Errorf("Found a term that should not be in the index: %s", term)
		}
		if cursor != nil {
			t.Error("Expected nil cursor, found non-nil")
		}
	}
}

func TestGenerateCursorsMapOnNonEmptyIndex(t *testing.T) {
	index := InvertedIndex{invertedLists: make(map[string]*node)}
	index.invertedLists["A"] = &node{documentId: 88, position: 0, nextNode: nil}
	cursorsMap := index.generateCursorsMap([]string{"A", "B", "C"})
	if cursorsMapLen := len(cursorsMap); cursorsMapLen != 3 {
		t.Errorf("Expected 3 cursors, got %d", cursorsMapLen)
	}
	for term, cursor := range cursorsMap {
		if term != "A" && term != "B" && term != "C" {
			t.Errorf("Found a term that should not be in the index: %s", term)
		}
		if cursor != nil && term != "A" {
			t.Error("Expected nil cursor for non-A terms, found non-nil")
		}
		if cursor == nil && term == "A" {
			t.Error("Expected non-nil cursor for A terms, found nil")
		}
	}
}

func TestSearchOnEmptyIndex(t *testing.T) {
	emptyIndex := InvertedIndex{invertedLists: make(map[string]*node)}
	query := []string{"A", "B"}
	weights := map[string]uint64{"A": 4, "B": 2}
	res := emptyIndex.Search(query, weights)
	matches := res.GetBestMatches(10, 0)
	if matchesCount := len(matches); matchesCount != 0 {
		t.Errorf("Expected 0 matches, found: %d", matchesCount)
	}
}

func TestSearchOnNonEmptyIndex(t *testing.T) {
	index := InvertedIndex{invertedLists: make(map[string]*node)}
	index.invertedLists["A"] = &node{documentId: 88, position: 0, nextNode: &node{
		documentId: 99, position: 1, nextNode: nil,
	}}
	index.invertedLists["B"] = &node{documentId: 88, position: 0, nextNode: nil}
	query := []string{"A", "B"}
	weights := map[string]uint64{"A": 4, "B": 2}
	res := index.Search(query, weights)
	matches := res.GetBestMatches(10, 0)
	if matchesCount := len(matches); matchesCount != 2 {
		t.Errorf("Expected 2 matches, found: %d", matchesCount)
	}
}

func TestGenerateInsertionMapOnEmptyIndex(t *testing.T) {
	emptyIndex := InvertedIndex{invertedLists: make(map[string]*node)}
	insertionMap := emptyIndex.generateInsertionMap([]string{"A", "B", "C"}, 22)
	if insertionMapLen := len(insertionMap); insertionMapLen != 3 {
		t.Errorf("Expected 3 cursors, got %d", insertionMapLen)
	}
	for term, cursor := range insertionMap {
		if term != "A" && term != "B" && term != "C" {
			t.Errorf("Found a term that should not be in the index: %s", term)
		}
		if cursor != nil {
			t.Error("Expected nil cursor, found non-nil")
		}
	}
}

func TestGenerateInsertionMapOnNonEmptyIndex(t *testing.T) {
	index := InvertedIndex{invertedLists: make(map[string]*node)}
	index.invertedLists["A"] = &node{documentId: 33, position: 0, nextNode: &node{
		documentId: 55, position: 1, nextNode: nil,
	}}
	index.invertedLists["B"] = &node{documentId: 55, position: 0, nextNode: nil}
	index.invertedLists["C"] = &node{documentId: 11, position: 0, nextNode: nil}
	insertionMap := index.generateInsertionMap([]string{"A", "B", "C"}, 44)
	if cursor := insertionMap["A"]; cursor == nil {
		t.Error("Expected A-cursor to be non-nil, it was nil")
	}
	if cursor := insertionMap["B"]; cursor != nil {
		t.Error("Expected B-cursor to be nil, it was non-nil")
	}
	if cursor := insertionMap["C"]; cursor == nil {
		t.Error("Expected C-cursor to be non-nil, it was nil")
	}
}

func TestStoreThenSearch(t *testing.T) {
	index := InvertedIndex{invertedLists: make(map[string]*node)}
	index.Store([]string{"A", "B", "C"}, 22)
	index.Store([]string{"B", "C"}, 33)
	index.Store([]string{"A", "C"}, 11)
	index.Store([]string{"A", "B", "C"}, 44)
	index.Store([]string{"X", "Y"}, 55)
	index.Store([]string{"X"}, 66)
	query := []string{"A"}
	weights := map[string]uint64{"A": 4}
	res := index.Search(query, weights)
	resLen := len(res.GetBestMatches(10, 0))
	if resLen != 3 {
		t.Errorf("Expected 3 docs with term 'A', got: %d", resLen)
	}
	query = []string{"X"}
	weights = map[string]uint64{"X": 4}
	res = index.Search(query, weights)
	resLen = len(res.GetBestMatches(10, 0))
	if resLen != 2 {
		t.Errorf("Expected 3 docs with term 'X', got: %d", resLen)
	}
}
