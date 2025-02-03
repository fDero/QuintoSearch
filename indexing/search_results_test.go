package indexing

import (
	"testing"
)

func TestIncrementScoreAndGetBestMatches(t *testing.T) {
	search_results := NewSearchResult()
	search_results.incrementScore(1, 40)
	search_results.incrementScore(2, 30)
	search_results.incrementScore(1, 10)
	bestMatches := search_results.GetBestMatches(10, 0)
	bestMatchesCount := len(bestMatches)
	if bestMatchesCount != 2 {
		t.Errorf("Expected 2 matches, got %d", bestMatchesCount)
	}
	rank1score := search_results.scoreByDocumentId[bestMatches[0]]
	rank2score := search_results.scoreByDocumentId[bestMatches[1]]
	if rank1score < rank2score {
		t.Error("Matches are not in descending order of score")
	}
}

func TestIncrementScoreAndGetSizeInPages(t *testing.T) {
	search_results := NewSearchResult()
	search_results.incrementScore(1, 40)
	search_results.incrementScore(2, 30)
	search_results.incrementScore(3, 10)
	search_results.incrementScore(4, 10)
	if size := search_results.GetSizeInPages(2); size != 2 {
		t.Errorf("Expected 4 docs in 2 pages of 2 element per page, got %d pages", size)
	}
	search_results.incrementScore(5, 30)
	if size := search_results.GetSizeInPages(2); size != 3 {
		t.Errorf("Expected 5 docs in 3 pages: 2-2-1, got %d pages", size)
	}
}
