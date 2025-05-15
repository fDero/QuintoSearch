/*=================================== LICENSE =======================================

                                   Apache License
                             Version 2.0, January 2004
                          http://www.apache.org/licenses/

============================== BRIEF FILE DESCRIPTION ===============================

This files contains some utilities that are often needed in the codebase. Such
utilities are not specific to a particular package, but are used in multiple packages.
==================================================================================*/

package misc

import (
	"iter"
)

type Set[T comparable] struct {
	storage map[T]bool
}

func NewSet[T comparable]() Set[T] {
	return Set[T]{
		storage: make(map[T]bool),
	}
}

func (s *Set[T]) InsertOne(value T) {
	s.storage[value] = true
}

func (s *Set[T]) InsertAll(other *Set[T]) {
	for value := range other.storage {
		s.storage[value] = true
	}
}

func (s *Set[T]) Contains(value T) bool {
	flag, exists := s.storage[value]
	return exists && flag
}

func (s *Set[T]) Size() int {
	return len(s.storage)
}

func NewSliceIterator[T any](slice []T) iter.Seq[T] {
	return func(yield func(T) bool) {
		for _, value := range slice {
			if !yield(value) {
				break
			}
		}
	}
}
