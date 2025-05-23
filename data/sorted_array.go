package data

import (
	"iter"
	"slices"
)

type SortedArray[T any] struct {
	equalityPredicate func(a, b T) bool
	orderingPredicate func(a, b T) bool
	storage           []T
}

func NewSortedArray[T any](
	orderingPredicate func(a, b T) bool,
	equalityPredicate func(a, b T) bool,
) *SortedArray[T] {
	return &SortedArray[T]{
		orderingPredicate: orderingPredicate,
		equalityPredicate: equalityPredicate,
		storage:           make([]T, 0),
	}
}

func (sa *SortedArray[T]) findIndexOf(value T) (int, bool) {

	for i, v := range sa.storage {
		if sa.equalityPredicate(v, value) {
			return i, true
		}
		if !sa.orderingPredicate(v, value) {
			return i, false
		}
	}

	return len(sa.storage), false
}

func (sa *SortedArray[T]) Insert(value T) {
	index, found := sa.findIndexOf(value)
	if !found {
		sa.storage = append(sa.storage, value)
		if index < len(sa.storage)-1 {
			copy(sa.storage[index+1:], sa.storage[index:len(sa.storage)-1])
		}
		sa.storage[index] = value
	}
}

func (sa *SortedArray[T]) Remove(value T) {
	index, found := sa.findIndexOf(value)
	if found {
		sa.storage = slices.Delete(sa.storage, index, index+1)
	}
}

func (sa *SortedArray[T]) Size() int {
	return len(sa.storage)
}

func (sa *SortedArray[T]) Contains(value T) bool {
	_, found := sa.findIndexOf(value)
	return found
}

func (sa *SortedArray[T]) Iterate() iter.Seq[T] {
	return NewSliceIterator(sa.storage)
}
