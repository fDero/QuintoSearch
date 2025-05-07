package search

import (
	"quinto/misc"
	"testing"
)

func TestExactMatchQuerySuccess(t *testing.T) {

	query := Exact{
		func() (misc.TermTracker, bool) {
			return misc.TermTracker{DocumentId: 1, Position: 0}, true
		},
		func() {},
		func() {},
	}

	match := query.Run()

	if !match.success {
		t.Errorf("Expected success to be true, got false")
	}
	if match.DocumentId != 1 {
		t.Errorf("Expected DocumentId to be 1, got %d", match.DocumentId)
	}
	if match.StartPosition != 0 {
		t.Errorf("Expected StartPosition to be 0, got %d", match.StartPosition)
	}
	if match.EndPosition != 0 {
		t.Errorf("Expected EndPosition to be 0, got %d", match.EndPosition)
	}
}
