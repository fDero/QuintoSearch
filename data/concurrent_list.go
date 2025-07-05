/*=================================== LICENSE =======================================

                                   Apache License
                             Version 2.0, January 2004
                          http://www.apache.org/licenses/

============================== BRIEF FILE DESCRIPTION ===============================

This files contains an implementation of a concurrency-ready a linked list. Upon
insertion, it returns a reference object that can be used to remove it from the list,
later, in a concurrently-safe manner. The current implementation uses the following
optimizations:
	1. Insertions are lockfree whenever possible (when the list is not empty).
	2. Deletions are lazy, and are actually performed only when necessary
	3. Is it allowed for updates to be visible after some time

Real eager deletion is performed once the rateo of nodes marked for deletion
over the total number of nodes exceeds a threshold, which is set to 0.4. Of course,
only one thread can perform the deletion at a time, if multiple threads want to
perform a deletion, the first one to acquire the lock will perform it, while
the others will just exit. If somehow the deletion will not bring the rateo below
the threshold, the deletion will be performed again, on the next call to removeNode.
===================================================================================*/

package data

import (
	"iter"
	"sync"
	"sync/atomic"
)

type ConcurrentList[T any] struct {
	head atomic.Pointer[concurrentListNode[T]]
	tail atomic.Pointer[concurrentListNode[T]]

	structure_mutex sync.Mutex
	node_count      atomic.Int64
	marked_count    atomic.Int64
	rm_threshold    float64
}

type concurrentListNode[T any] struct {
	next *concurrentListNode[T]
	prev *concurrentListNode[T]
	mark atomic.Bool
	item T
}

func NewLinkedList[T any]() *ConcurrentList[T] {
	return &ConcurrentList[T]{}
}

func (list *ConcurrentList[T]) InsertFront(value T) ConcurrentListEntry[T] {
	newNode := &concurrentListNode[T]{item: value}
	if list.tail.Load() == nil {
		list.structure_mutex.Lock()
		defer list.structure_mutex.Unlock()
		if list.tail.Load() != nil {
			return list.InsertFront(value)
		}
		list.head.Store(newNode)
		list.tail.Store(newNode)
		return ConcurrentListEntry[T]{list: list, ptr: newNode}
	}
	for {
		currentHead := list.head.Load()
		newNode.next = currentHead
		if list.head.CompareAndSwap(currentHead, newNode) {
			currentHead.prev = newNode
			return ConcurrentListEntry[T]{list: list, ptr: newNode}
		}
	}
}

func (list *ConcurrentList[T]) IterateForward() iter.Seq[ConcurrentListEntry[T]] {
	return func(yield func(ConcurrentListEntry[T]) bool) {
		cursor := list.head.Load()
		for cursor != nil {
			entry := ConcurrentListEntry[T]{list: list, ptr: cursor}
			if !cursor.mark.Load() && !yield(entry) {
				break
			}
			cursor = cursor.next
		}
	}
}

func (list *ConcurrentList[T]) IterateBackwards() iter.Seq[ConcurrentListEntry[T]] {
	return func(yield func(ConcurrentListEntry[T]) bool) {
		cursor := list.tail.Load()
		for cursor != nil {
			entry := ConcurrentListEntry[T]{list: list, ptr: cursor}
			if !cursor.mark.Load() && !yield(entry) {
				break
			}
			cursor = cursor.prev
		}
	}
}

func (list *ConcurrentList[T]) mustPrune() bool {
	currentNodeCount := list.node_count.Load()
	currentMarkedCount := list.marked_count.Load()
	currentMarkedRateo := float64(currentMarkedCount) / float64(currentNodeCount)
	return currentMarkedRateo > list.rm_threshold
}

func (list *ConcurrentList[T]) attemptPrune() {
	if !list.structure_mutex.TryLock() {
		return
	}
	defer list.structure_mutex.Unlock()
	cursor := list.head.Load()
	for cursor != nil {
		if cursor.mark.Load() && cursor.prev != nil {
			cursor.prev.next = cursor.next
		}
		if cursor.mark.Load() && cursor.next != nil {
			cursor.next.prev = cursor.prev
		}
		if cursor.mark.Load() {
			list.node_count.Add(-1)
			list.marked_count.Add(-1)
		}
		cursor = cursor.next
	}
}

func (list *ConcurrentList[T]) removeNode(listNode *concurrentListNode[T]) {
	if listNode == nil {
		return
	}
	listNode.mark.Store(true)
	list.marked_count.Add(1)
	if list.mustPrune() {
		list.attemptPrune()
	}
}

type ConcurrentListEntry[T any] struct {
	list *ConcurrentList[T]
	ptr  *concurrentListNode[T]
}

func (em ConcurrentListEntry[T]) Value() T {
	if em.ptr != nil && em.list != nil {
		return em.ptr.item
	}
	var zeroValue T
	return zeroValue
}

func (em *ConcurrentListEntry[T]) Remove() {
	if em.ptr != nil && em.list != nil {
		em.list.removeNode(em.ptr)
		em.ptr = nil
		em.list = nil
	}
}
