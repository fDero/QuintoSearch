/*=================================== LICENSE =======================================

                                   Apache License
                             Version 2.0, January 2004
                          http://www.apache.org/licenses/

============================== BRIEF FILE DESCRIPTION ===============================

A "ReverseIndex" is an interface that describes every entity that can be used to
store documents and initialize queries to retrieve them. It is a key component of the
search engine, as it allows for efficient storage and retrieval of documents based
on their content. The "ReverseIndex" is responsable for managing the mapping between
a single term and the documents that contain it, togheter with the positions of the
term in the document. Is essential for the well functioning of the search engine
that terms are iterated in ascending order of document-id and position within a given
document. The "IterateOverTerms" must herby work in this way in every implementation.
==================================================================================*/

package misc

import (
	"iter"
)

type Token struct {
	StemmedText  string
	OriginalText string
	Position     int
}

type TermTracker struct {
	DocumentId uint64
	Position   int
}

type ReverseIndex interface {
	IterateOverTerms(term string) iter.Seq[TermTracker]
	StoreNewDocument(toks iter.Seq[Token]) (uint64, error)
}
