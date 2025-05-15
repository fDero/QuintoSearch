package search

import (
	"quinto/misc"
	"testing"
)

func TestNearQuerySuccess(t *testing.T) {

	lxQuerySuccess := NewExactQueryFromSlice([]misc.TermTracker{{DocumentId: 1, Position: 0}})
	rxQuerySuccess := NewExactQueryFromSlice([]misc.TermTracker{{DocumentId: 1, Position: 1}})

	nearQuery := ComplexQuery{
		lx:     &lxQuerySuccess,
		rx:     &rxQuerySuccess,
		ord:    false,
		policy: NearQueryPolicy(10),
	}

	match := nearQuery.Run()
	if !match.Success {
		t.Errorf("Expected success to be true, got false")
	}
}

func TestNearQuerySuccessReverse(t *testing.T) {

	lxQuerySuccess := NewExactQueryFromSlice([]misc.TermTracker{{DocumentId: 1, Position: 1}})
	rxQuerySuccess := NewExactQueryFromSlice([]misc.TermTracker{{DocumentId: 1, Position: 0}})

	andQuery := ComplexQuery{
		lx:     &lxQuerySuccess,
		rx:     &rxQuerySuccess,
		ord:    false,
		policy: NearQueryPolicy(10),
	}

	match := andQuery.Run()
	if !match.Success {
		t.Errorf("Expected success to be true, got false")
	}
}

func TestOrderedNearQueryFailureByWrongOrder(t *testing.T) {

	lxQuerySuccess := NewExactQueryFromSlice([]misc.TermTracker{{DocumentId: 1, Position: 1}})
	rxQuerySuccess := NewExactQueryFromSlice([]misc.TermTracker{{DocumentId: 1, Position: 0}})

	andQuery := ComplexQuery{
		lx:     &lxQuerySuccess,
		rx:     &rxQuerySuccess,
		ord:    true,
		policy: NearQueryPolicy(10),
	}

	match := andQuery.Run()
	if match.Success {
		t.Errorf("Expected success to be false, got true")
	}
}

func TestNearQueryFailureByExcessiveDistance(t *testing.T) {

	lxQuerySuccess := NewExactQueryFromSlice([]misc.TermTracker{{DocumentId: 1, Position: 0}})
	rxQuerySuccess := NewExactQueryFromSlice([]misc.TermTracker{{DocumentId: 1, Position: 11}})

	andQuery := ComplexQuery{
		lx:     &lxQuerySuccess,
		rx:     &rxQuerySuccess,
		ord:    true,
		policy: NearQueryPolicy(10),
	}

	match := andQuery.Run()
	if match.Success {
		t.Errorf("Expected success to be false, got true")
	}
}

func TestNearQueryFailureByExcessiveDistanceReverse(t *testing.T) {

	lxQuerySuccess := NewExactQueryFromSlice([]misc.TermTracker{{DocumentId: 1, Position: 11}})
	rxQuerySuccess := NewExactQueryFromSlice([]misc.TermTracker{{DocumentId: 1, Position: 0}})

	andQuery := ComplexQuery{
		lx:     &lxQuerySuccess,
		rx:     &rxQuerySuccess,
		ord:    true,
		policy: NearQueryPolicy(10),
	}

	match := andQuery.Run()
	if match.Success {
		t.Errorf("Expected success to be false, got true")
	}
}
