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
	lx     misc.Query
	rx     misc.Query
	ord    bool
	policy func(misc.Match, misc.Match) bool
}

var (
	OrQueryPolicy  = func(lx, rx misc.Match) bool { return lx.Success || rx.Success }
	XorQueryPolicy = func(lx, rx misc.Match) bool { return lx.Success != rx.Success }
	AndQueryPolicy = func(lx, rx misc.Match) bool { return lx.Success && rx.Success }
)

func withinBound(m1, m2 misc.Match, dist int) bool {
	canSubtract := m1.StartPosition > m2.EndPosition
	return !canSubtract || (m1.StartPosition-m2.EndPosition) <= misc.TermPosition(dist)
}

func NearQueryPolicy(dist int) func(lx, rx misc.Match) bool {
	return func(lx, rx misc.Match) bool {
		withinBoundsForwards := withinBound(lx, rx, dist)
		withinBoundsBackwards := withinBound(rx, lx, dist)
		withinBounds := withinBoundsForwards && withinBoundsBackwards
		return lx.Success && rx.Success && withinBounds
	}
}

func (q *ComplexQuery) Init(index misc.ReverseIndex) {
	q.lx.Init(index)
	q.rx.Init(index)
}

func (q *ComplexQuery) Run() misc.Match {
	lxMatch := q.lx.Run()
	rxMatch := q.rx.Run()
	if !q.policy(lxMatch, rxMatch) {
		return misc.Match{Success: false}
	}

	if lxMatch.Success && rxMatch.Success {

		if lxMatch.DocId != rxMatch.DocId {
			return misc.Match{Success: false}
		}

		success := true
		if lxMatch.StartPosition > rxMatch.StartPosition {
			lxMatch, rxMatch = rxMatch, lxMatch
			success = !q.ord
		}

		lxMatch.InvolvedTokens.InsertAll(&rxMatch.InvolvedTokens)
		return misc.Match{
			Success:        success,
			DocId:          lxMatch.DocId,
			StartPosition:  lxMatch.StartPosition,
			EndPosition:    rxMatch.EndPosition,
			InvolvedTokens: lxMatch.InvolvedTokens,
		}
	}

	if lxMatch.Success {
		return lxMatch
	}

	if rxMatch.Success {
		return rxMatch
	}

	return misc.Match{Success: true}
}

func (q *ComplexQuery) Advance() {
	lxDocumentId, lxPosition := q.lx.Coordinates()
	rxDocumentId, rxPosition := q.rx.Coordinates()
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

func (q *ComplexQuery) Coordinates() (misc.DocumentId, misc.TermPosition) {
	lxDocumentId, lxPosition := q.lx.Coordinates()
	rxDocumentId, rxPosition := q.rx.Coordinates()
	if lxDocumentId < rxDocumentId {
		return lxDocumentId, lxPosition
	}
	if lxDocumentId == rxDocumentId && lxPosition < rxPosition {
		return lxDocumentId, lxPosition
	}
	return rxDocumentId, rxPosition
}
