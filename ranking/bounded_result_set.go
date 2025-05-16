package ranking

import (
	"iter"
	"quinto/data"
	"quinto/misc"
)

type BoundedResultSet struct {
	storage *data.Heap[misc.SearchResult]
	maxSize int
}

func compareResults(a, b misc.SearchResult) bool {
	return a.DocId <= b.DocId
}

func NewBoundedResultSet(maxSize int) *BoundedResultSet {
	return &BoundedResultSet{
		storage: data.NewHeap(compareResults),
		maxSize: maxSize,
	}
}

func (brs *BoundedResultSet) StoreNewResult(result misc.SearchResult) {
	brs.storage.Push(result)
	if brs.storage.Size() > brs.maxSize {
		brs.storage.Pop()
	}
}

func (brs *BoundedResultSet) ToSortedSlice() []misc.SearchResult {
	var result = make([]misc.SearchResult, brs.storage.Size())
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

func (brs *BoundedResultSet) Iterate() iter.Seq[misc.SearchResult] {
	return data.NewSliceIterator(brs.ToSortedSlice())
}
