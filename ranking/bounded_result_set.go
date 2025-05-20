/*=================================== LICENSE =======================================

                                   Apache License
                             Version 2.0, January 2004
                          http://www.apache.org/licenses/

============================== BRIEF FILE DESCRIPTION ===============================

This file contains the implementation of "BoundedResultSet", which is a concrete
implementation of the "ResultSet" interface, based on a heap data structure.

It is designed to store a limited number of search results, and when the limit is
reached, it will remove the least relevant result. The results are stored in a
heap, which allows for efficient insertion and removal of elements.
==================================================================================*/

package ranking

import (
	"iter"
	"quinto/core"
	"quinto/data"
)

type BoundedResultSet struct {
	storage *data.Heap[core.SearchResult]
	maxSize int
}

func compareResults(a, b core.SearchResult) bool {
	return a.DocId <= b.DocId
}

func NewBoundedResultSet(maxSize int) *BoundedResultSet {
	return &BoundedResultSet{
		storage: data.NewHeap(compareResults),
		maxSize: maxSize,
	}
}

func (brs *BoundedResultSet) StoreNewResult(result core.SearchResult) {
	brs.storage.Push(result)
	if brs.storage.Size() > brs.maxSize {
		brs.storage.Pop()
	}
}

func (brs *BoundedResultSet) ToSortedSlice() []core.SearchResult {
	var result = make([]core.SearchResult, brs.storage.Size())
	originalSize := brs.storage.Size()
	newStorage := data.NewHeap(compareResults)
	for i := originalSize - 1; i >= 0; i-- {
		popped, _ := brs.storage.Pop()
		result[i] = popped
		newStorage.Push(popped)
	}
	brs.storage = newStorage
	return result
}

func (brs *BoundedResultSet) Iterate() iter.Seq[core.SearchResult] {
	return data.NewSliceIterator(brs.ToSortedSlice())
}
