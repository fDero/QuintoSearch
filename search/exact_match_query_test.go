package search

import (
	"quinto/core"
	"testing"
)

func TestExactMatchQuerySuccess(t *testing.T) {

	query := ExactQuery{
		"test",
		func() (core.TermTracker, bool) {
			return core.TermTracker{DocId: 1, Position: 0}, true
		},
		func() {},
		func() {},
	}

	match := query.Run()

	if !match.Success {
		t.Errorf("Expected success to be true, got false")
	}
	if match.DocId != 1 {
		t.Errorf("Expected DocumentId to be 1, got %d", match.DocId)
	}
	if match.StartPosition != 0 {
		t.Errorf("Expected StartPosition to be 0, got %d", match.StartPosition)
	}
	if match.EndPosition != 0 {
		t.Errorf("Expected EndPosition to be 0, got %d", match.EndPosition)
	}
}
