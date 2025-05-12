package search

import (
	"quinto/misc"
	"testing"
)

func runTestCollectMatchesHelper(t *testing.T, queryString string) []Match {

	index := NewNaiveReverseIndex()
	index.StoreNewDocument(misc.NewSliceIterator(helloWorldDocument))
	index.StoreNewDocument(misc.NewSliceIterator(guitarDocument))
	index.StoreNewDocument(misc.NewSliceIterator(hobbyDocument))
	index.StoreNewDocument(misc.NewSliceIterator(toolsDocument))

	queryFragments := SplitQuery(queryString)
	query, err := ParseQuery(queryFragments)
	if err != nil {
		t.Fatalf("Failed to parse query: %v", err)
	}

	defer query.Close()
	results := []Match{}

	query.Init(index)

	for !query.Ended() {
		queryResult := query.Run()
		if queryResult.success {
			results = append(results, queryResult)
		}
		query.Advance()
	}

	return results
}

func TestFirstComplexQueryOverMultipleDocuments(t *testing.T) {
	matches := runTestCollectMatchesHelper(t, "hello AND world")
	if len(matches) != 1 {
		t.Errorf("Expected 1 match, got %d", len(matches))
	}
}

func TestSecondComplexQueryOverMultipleDocuments(t *testing.T) {
	matches := runTestCollectMatchesHelper(t, "guitar AND music")
	if len(matches) != 1 {
		t.Errorf("Expected 1 match, got %d", len(matches))
	}
}
