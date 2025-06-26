/*=================================== LICENSE =======================================

                                   Apache License
                             Version 2.0, January 2004
                          http://www.apache.org/licenses/

============================== BRIEF FILE DESCRIPTION ===============================

This files contains the implementation of a Queue as a golang channel. It is a simple
wrapper that exposes a Queue API. This wrapper is needed to provide modularity and
ease of replacement in case in future a lockfree queue implementation is needed.
===================================================================================*/

package data

import (
	"sync/atomic"
)

type ConcurrentQueue[T any] struct {
	channel chan T
	size    atomic.Int64
}

func NewConcurrentQueue[T any]() *ConcurrentQueue[T] {
	return &ConcurrentQueue[T]{
		channel: make(chan T),
		size:    atomic.Int64{},
	}
}

func (q *ConcurrentQueue[T]) Push(value T) {
	q.size.Add(1)
	q.channel <- value
}

func (q *ConcurrentQueue[T]) Pop() (T, bool) {
	q.size.Add(-1)
	value, ok := <-q.channel
	return value, ok
}

func (q *ConcurrentQueue[T]) IsEmpty() bool {
	return q.size.Load() == 0
}
