package search

type ComplexQuery struct {
	lx     Query
	rx     Query
	ord    bool
	policy func(Match, Match) bool
}

var (
	OrQueryPolicy  = func(lx, rx Match) bool { return lx.success || rx.success }
	XorQueryPolicy = func(lx, rx Match) bool { return lx.success != rx.success }
	AndQueryPolicy = func(lx, rx Match) bool { return lx.success && rx.success }
)

func NearQueryPolicy(dist int) func(lx, rx Match) bool {
	return func(lx, rx Match) bool {
		withinBoundsForwards := (rx.StartPosition - lx.EndPosition) <= dist
		withinBoundsBackwards := (lx.StartPosition - rx.EndPosition) <= dist
		withinBounds := withinBoundsForwards && withinBoundsBackwards
		return lx.success && rx.success && withinBounds
	}
}

func (q *ComplexQuery) Run() Match {
	lxMatch := q.lx.Run()
	rxMatch := q.rx.Run()
	if !q.policy(lxMatch, rxMatch) {
		return Match{success: false}
	}

	if lxMatch.success && rxMatch.success {

		if lxMatch.DocumentId != rxMatch.DocumentId {
			return Match{success: false}
		}

		success := true
		if lxMatch.StartPosition > rxMatch.StartPosition {
			lxMatch, rxMatch = rxMatch, lxMatch
			success = !q.ord
		}

		return Match{
			success:       success,
			DocumentId:    lxMatch.DocumentId,
			StartPosition: lxMatch.StartPosition,
			EndPosition:   rxMatch.EndPosition,
		}
	}

	if lxMatch.success {
		return lxMatch
	}

	if rxMatch.success {
		return rxMatch
	}

	return Match{success: true}
}

func (q *ComplexQuery) Advance() {
	if q.lx.lowestDocumentId() < q.rx.lowestDocumentId() {
		q.lx.Advance()
	} else {
		q.rx.Advance()
	}
}

func (q *ComplexQuery) Close() {
	q.lx.Close()
	q.rx.Close()
}

func (q *ComplexQuery) lowestDocumentId() uint64 {
	return min(
		q.lx.lowestDocumentId(),
		q.rx.lowestDocumentId(),
	)
}
