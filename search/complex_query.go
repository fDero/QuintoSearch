/*=================================== LICENSE =======================================

                                   Apache License
                             Version 2.0, January 2004
                          http://www.apache.org/licenses/

============================== BRIEF FILE DESCRIPTION ===============================

An "ComplexQuery" query is a type of search query that looks for some specific
condition encompassing two other queries. Specifically, it relies on a policy function
that defines wether the two match results of the two queries at hand should result
in a match or not. If you want to enforce the scenario in which matches from
the two queries are required to be ordered such that the first match must refer
to a term that appears before the second match in the document, you can set the
"ord" field to true. "ComplexQuery" implements the Query interface, which defines the
"Run", "Advance", and "Close" methods. Please refer to the documentation of the
"Query" interface for more details about its methods and their intended usage.
==================================================================================*/

package search

import (
	"quinto/misc"
)

type ComplexQuery struct {
	lx     Query
	rx     Query
	ord    bool
	policy func(Match, Match) bool
}

var (
	OrQueryPolicy  = func(lx, rx Match) bool { return lx.Success || rx.Success }
	XorQueryPolicy = func(lx, rx Match) bool { return lx.Success != rx.Success }
	AndQueryPolicy = func(lx, rx Match) bool { return lx.Success && rx.Success }
)

func NearQueryPolicy(dist int) func(lx, rx Match) bool {
	return func(lx, rx Match) bool {
		withinBoundsForwards := (rx.StartPosition - lx.EndPosition) <= dist
		withinBoundsBackwards := (lx.StartPosition - rx.EndPosition) <= dist
		withinBounds := withinBoundsForwards && withinBoundsBackwards
		return lx.Success && rx.Success && withinBounds
	}
}

func (q *ComplexQuery) Init(index misc.ReverseIndex) {
	q.lx.Init(index)
	q.rx.Init(index)
}

func (q *ComplexQuery) Run() Match {
	lxMatch := q.lx.Run()
	rxMatch := q.rx.Run()
	if !q.policy(lxMatch, rxMatch) {
		return Match{Success: false}
	}

	if lxMatch.Success && rxMatch.Success {

		if lxMatch.DocumentId != rxMatch.DocumentId {
			return Match{Success: false}
		}

		success := true
		if lxMatch.StartPosition > rxMatch.StartPosition {
			lxMatch, rxMatch = rxMatch, lxMatch
			success = !q.ord
		}

		return Match{
			Success:       success,
			DocumentId:    lxMatch.DocumentId,
			StartPosition: lxMatch.StartPosition,
			EndPosition:   rxMatch.EndPosition,
		}
	}

	if lxMatch.Success {
		return lxMatch
	}

	if rxMatch.Success {
		return rxMatch
	}

	return Match{Success: true}
}

func (q *ComplexQuery) Advance() {
	lxDocumentId, lxPosition := q.lx.coordinates()
	rxDocumentId, rxPosition := q.rx.coordinates()
	shouldGoLxByDocumentId := lxDocumentId < rxDocumentId
	shouldGoLxByPosition := lxDocumentId == rxDocumentId && lxPosition < rxPosition
	shouldGoLx := shouldGoLxByDocumentId || shouldGoLxByPosition
	if shouldGoLx && !q.lx.Ended() {
		q.lx.Advance()
	}
	if !shouldGoLx || q.lx.Ended() {
		q.rx.Advance()
	}
	if !shouldGoLx && !q.lx.Ended() && q.rx.Ended() {
		q.lx.Advance()
	}
}

func (q *ComplexQuery) Ended() bool {
	return q.lx.Ended() && q.rx.Ended()
}

func (q *ComplexQuery) Close() {
	q.lx.Close()
	q.rx.Close()
}

func (q *ComplexQuery) coordinates() (uint64, int) {
	lxDocumentId, lxPosition := q.lx.coordinates()
	rxDocumentId, rxPosition := q.rx.coordinates()
	if lxDocumentId < rxDocumentId {
		return lxDocumentId, lxPosition
	}
	if lxDocumentId == rxDocumentId && lxPosition < rxPosition {
		return lxDocumentId, lxPosition
	}
	return rxDocumentId, rxPosition
}
