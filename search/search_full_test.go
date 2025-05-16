package search

import (
	"quinto/misc"
	"testing"
)

func runTestCollectMatchesHelper(t *testing.T, queryString string) map[misc.DocumentId]misc.Match {

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
	results := map[misc.DocumentId]misc.Match{}

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
