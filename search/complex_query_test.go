package search

import (
	"quinto/misc"
	"testing"
)

func TestAndQuerySuccess(t *testing.T) {

	lxQuerySuccess := Exact{
		func() (misc.TermTracker, bool) {
			return misc.TermTracker{DocumentId: 1, Position: 0}, true
		},
		func() {},
		func() {},
	}

	rxQuerySuccess := Exact{
		func() (misc.TermTracker, bool) {
			return misc.TermTracker{DocumentId: 1, Position: 1}, true
		},
		func() {},
		func() {},
	}

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

	lxQuerySuccess := Exact{
		func() (misc.TermTracker, bool) {
			return misc.TermTracker{DocumentId: 1, Position: 0}, true
		},
		func() {},
		func() {},
	}

	rxQuerySuccess := Exact{
		func() (misc.TermTracker, bool) {
			return misc.TermTracker{DocumentId: 3, Position: 1}, true
		},
		func() {},
		func() {},
	}

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

	lxQuerySuccess := Exact{
		func() (misc.TermTracker, bool) {
			return misc.TermTracker{DocumentId: 1, Position: 0}, true
		},
		func() {},
		func() {},
	}

	rxQuerySuccess := Exact{
		func() (misc.TermTracker, bool) {
			return misc.TermTracker{DocumentId: 0, Position: -1}, false
		},
		func() {},
		func() {},
	}

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
