/*=================================== LICENSE =======================================

                                   Apache License
                             Version 2.0, January 2004
                          http://www.apache.org/licenses/

============================== BRIEF FILE DESCRIPTION ===============================

This file contains a simple implementation of a sorted array data structure. It's
basically a wrapper around a slice that maintains the order of the elements according
to a user-defined ordering predicate. The sorted array supports insertion, removal,
and searching for elements. It also provides a way to iterate over the elements in
the array (in ascending order) using an iterator.
==================================================================================*/

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

func (sa *SortedArray[T]) Lowest() (T, bool) {
	if len(sa.storage) == 0 {
		var zero T
		return zero, false
	}
	return sa.storage[0], true
}

func (sa *SortedArray[T]) Highest() (T, bool) {
	if len(sa.storage) == 0 {
		var zero T
		return zero, false
	}
	return sa.storage[len(sa.storage)-1], true
}

func (sa *SortedArray[T]) findIndexOf(value T) (int, bool) {
	low, high := 0, len(sa.storage)-1
	for low <= high {
		pivotIndex := (low + high) / 2
		pivotValue := sa.storage[pivotIndex]
		switch {
		case sa.equalityPredicate(pivotValue, value):
			return pivotIndex, true
		case sa.orderingPredicate(pivotValue, value):
			low = pivotIndex + 1
		default:
			high = pivotIndex - 1
		}
	}
	return low, false
}

func (sa *SortedArray[T]) Insert(value T) bool {
	highestValue, exists := sa.Highest()
	if exists && sa.equalityPredicate(highestValue, value) {
		return false
	}
	if exists && sa.orderingPredicate(highestValue, value) {
		sa.storage = append(sa.storage, value)
		return true
	}
	index, found := sa.findIndexOf(value)
	if !found {
		sa.storage = append(sa.storage, value)
		if index < len(sa.storage)-1 {
			copy(sa.storage[index+1:], sa.storage[index:len(sa.storage)-1])
		}
		sa.storage[index] = value
	}
	return !found
}

func (sa *SortedArray[T]) Remove(value T) bool {
	index, found := sa.findIndexOf(value)
	if found {
		sa.storage = slices.Delete(sa.storage, index, index+1)
	}
	return found
}

func (sa *SortedArray[T]) RemoveIf(predicate func(T) bool) bool {
	altered := false
	for i := 0; i < len(sa.storage); i++ {
		if predicate(sa.storage[i]) {
			sa.storage = slices.Delete(sa.storage, i, i+1)
			altered = true
			i-- // Adjust index after deletion
		}
	}
	return altered
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
