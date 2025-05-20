/*=================================== LICENSE =======================================

                                   Apache License
                             Version 2.0, January 2004
                          http://www.apache.org/licenses/

============================== BRIEF FILE DESCRIPTION ===============================

This file contains a simple implementation of a set data structure. It is a generic
implementation that can store any type that is comparable. The set is implemented
as a simple wrapper around a plain old map, where the keys are the values of the
set and the values are dummy values (e.g. booleans that are always true).
==================================================================================*/

package data

type Set[T comparable] struct {
	storage map[T]bool
}

func ToSet[T comparable](values []T) Set[T] {
	set := NewSet[T]()
	for _, value := range values {
		set.InsertOne(value)
	}
	return set
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
