/*=================================== LICENSE =======================================

                                   Apache License
                             Version 2.0, January 2004
                          http://www.apache.org/licenses/

============================== BRIEF FILE DESCRIPTION ===============================

The "Query" interface defines the structure and behavior of a full-text structured
query. It is supposed to be implemented as a tree-like structure, where each node
is either a leaf node (an "ExactQuery") or a complex query (like "AND","OR","NEAR",
"NOT", "XOR", "NEAR:ORD", ...). Every query is supposed to be initalized over some
reverse index. After initialization, the query will contain a bunch of iterators
that will be used to iterate over the terms in the indexed documents. The query
is supposed to be run, and then advanced, which will move the iterators to the next
position. Every possible configuration of the iterators the query internally holds
can result in a match. Therefore, a document can have multiple matches, and
multiple matches can be found in the same document. Multiple matches must be combined
into a single "SearchResult" before inserting it into the "ResultSet".

Calling "Init" and "Close" is mandatory for the well functioning of the query API.
==================================================================================*/

package misc

type Match struct {
	Success        bool
	DocId          DocumentId
	StartPosition  TermPosition
	EndPosition    TermPosition
	InvolvedTokens Set[Token]
}

type Query interface {
	Run() Match
	Advance()
	Ended() bool
	Close()
	Init(ReverseIndex)
	Coordinates() (DocumentId, TermPosition)
}
