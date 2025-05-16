/*=================================== LICENSE =======================================

                                   Apache License
                             Version 2.0, January 2004
                          http://www.apache.org/licenses/

============================== BRIEF FILE DESCRIPTION ===============================

A "ResultSet" is the way multiple instances of "SearchResult" can be stored and
iterated over in Quinto. It is designed to be as simple of an interface as possible.
==================================================================================*/

package misc

import (
	"iter"
)

type SearchResult struct {
	DocId DocumentId
	Score float64
}

type ResultSet interface {
	StoreNewResult(result SearchResult)
	Iterate() iter.Seq[SearchResult]
	SortedSlice() []SearchResult
}
