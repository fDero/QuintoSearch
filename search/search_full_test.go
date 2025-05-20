package search

import (
	"quinto/core"
	"quinto/data"
	"testing"
)

func runTestCollectMatchesHelper(t *testing.T, queryString string) map[core.DocumentId]core.Match {

	index := NewNaiveReverseIndex()
	index.StoreNewDocument(data.NewSliceIterator(helloWorldDocument))
	index.StoreNewDocument(data.NewSliceIterator(guitarDocument))
	index.StoreNewDocument(data.NewSliceIterator(hobbyDocument))
	index.StoreNewDocument(data.NewSliceIterator(toolsDocument))

	queryFragments := SplitQuery(queryString)
	query, err := ParseQuery(queryFragments)
	if err != nil {
		t.Fatalf("Failed to parse query: %v", err)
	}

	defer query.Close()
	results := map[core.DocumentId]core.Match{}

	query.Init(index)

	counter := 0
	for !query.Ended() {
		if counter > 300 {
			t.Fatalf("Infinite loop detected in query execution")
		}
		counter++
		queryResult := query.Run()
		if queryResult.Success {
			if match, exists := results[queryResult.DocId]; exists {
				match.InvolvedTokens.InsertAll(&queryResult.InvolvedTokens)
			} else {
				results[queryResult.DocId] = queryResult
			}
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

func TestThirdComplexQueryOverMultipleDocuments(t *testing.T) {
	matches := runTestCollectMatchesHelper(t, "guitar OR music")
	if len(matches) != 2 {
		t.Errorf("Expected 2 match, got %d", len(matches))
	}
}

func TestFourthComplexQueryOverMultipleDocuments(t *testing.T) {
	matches := runTestCollectMatchesHelper(t, "music OR guitar")
	if len(matches) != 2 {
		t.Errorf("Expected 2 match, got %d", len(matches))
	}
}
