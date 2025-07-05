/*=================================== LICENSE =======================================

                                   Apache License
                             Version 2.0, January 2004
                          http://www.apache.org/licenses/

============================== BRIEF FILE DESCRIPTION ===============================

This files contains the implementation of a R. Kent Treiber lock-free stack. It is
a concurrent data structure that allows multiple goroutines to push and pop elements
from the stack without blocking each other. The stack is implemented using atomic
operations to ensure thread safety without relying on traditional mutex-like locks.
===================================================================================*/

package data

import (
	"sync/atomic"
)

type ConcurrentStack[T any] struct {
	head atomic.Pointer[concurrentStackNode[T]]
}

type concurrentStackNode[T any] struct {
	next  atomic.Pointer[concurrentStackNode[T]]
	value T
}

func NewConcurrentStack[T any]() *ConcurrentStack[T] {
	return &ConcurrentStack[T]{}
}

func (s *ConcurrentStack[T]) Push(value T) {
	newNode := &concurrentStackNode[T]{value: value}
	for {
		currentHead := s.head.Load()
		newNode.next.Store(currentHead)
		if s.head.CompareAndSwap(currentHead, newNode) {
			return
		}
	}
}

func (s *ConcurrentStack[T]) Pop() (T, bool) {
	var zeroValue T
	for {
		currentHead := s.head.Load()
		if currentHead == nil {
			return zeroValue, false
		}
		nextNode := currentHead.next.Load()
		if s.head.CompareAndSwap(currentHead, nextNode) {
			return currentHead.value, true
		}
	}
}

func (s *ConcurrentStack[T]) IsEmpty() bool {
	return s.head.Load() == nil
}
