package search

import (
	"quinto/misc"
	"testing"
)

func TestAndQuerySuccess(t *testing.T) {

	lxQuerySuccess := NewExactQueryFromSlice([]misc.TermTracker{{DocumentId: 1, Position: 0}})
	rxQuerySuccess := NewExactQueryFromSlice([]misc.TermTracker{{DocumentId: 1, Position: 1}})

	andQuery := ComplexQuery{
		lx:     &lxQuerySuccess,
		rx:     &rxQuerySuccess,
		ord:    true,
		policy: AndQueryPolicy,
	}

	match := andQuery.Run()
	if !match.success {
		t.Errorf("Expected success to be true, got false")
	}
}

func TestAndQueryFailureByIdMismatch(t *testing.T) {

	lxQuerySuccess := NewExactQueryFromSlice([]misc.TermTracker{{DocumentId: 1, Position: 0}})
	rxQuerySuccess := NewExactQueryFromSlice([]misc.TermTracker{{DocumentId: 3, Position: 1}})

	andQuery := ComplexQuery{
		lx:     &lxQuerySuccess,
		rx:     &rxQuerySuccess,
		ord:    true,
		policy: AndQueryPolicy,
	}

	match := andQuery.Run()
	if match.success {
		t.Errorf("Expected success to be false, got true")
	}
}

func TestAndQueryFailureByPolicy(t *testing.T) {

	lxQuerySuccess := NewExactQueryFromSlice([]misc.TermTracker{{DocumentId: 1, Position: 3}})
	rxQuerySuccess := NewExactQueryFromSlice([]misc.TermTracker{{DocumentId: 1, Position: 2}})

	andQuery := ComplexQuery{
		lx:     &lxQuerySuccess,
		rx:     &rxQuerySuccess,
		ord:    true,
		policy: AndQueryPolicy,
	}

	match := andQuery.Run()
	if match.success {
		t.Errorf("Expected success to be false, got true")
	}
}
