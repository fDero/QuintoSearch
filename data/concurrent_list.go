/*=================================== LICENSE =======================================

                                   Apache License
                             Version 2.0, January 2004
                          http://www.apache.org/licenses/

============================== BRIEF FILE DESCRIPTION ===============================

This file contains the implementation of a lockfree concurrent doubly linked list
implementation. It works because the two links (pointers) are stored in a single
object, and it's address it's used atomically using the `atomic.Pointer` primitive.
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
	size atomic.Int64

	pruningMutex         sync.Mutex
	markedForDeletion    chan *concurrentListNode[T]
	markedNodesThreshold float64
}

type concurrentListNode[T any] struct {
	mark atomic.Bool
	link atomic.Pointer[concurrentListLink[T]]
	item T
}

type concurrentListLink[T any] struct {
	next *concurrentListNode[T]
	prev *concurrentListNode[T]
}

func NewLinkedList[T any]() *ConcurrentList[T] {
	emptyList := &ConcurrentList[T]{}
	emptyList.head.Store(nil)
	emptyList.tail.Store(nil)
	emptyList.size.Store(0)
	return emptyList
}

func (list *ConcurrentList[T]) Size() int {
	return int(list.size.Load())
}

func (list *ConcurrentList[T]) InsertFront(value T) ConcurrentListEntry[T] {
	newNode := &concurrentListNode[T]{item: value}
	defer list.size.Add(1)
	for {
		currentHead := list.head.Load()
		newNode.link.Store(&concurrentListLink[T]{
			next: list.head.Load(),
			prev: nil,
		})
		if list.head.CompareAndSwap(currentHead, newNode) {
			if currentHead != nil {
				oldHeadLink := currentHead.link.Load()
				oldHeadLink.prev = newNode
			}
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
				return
			}
			link := cursor.link.Load()
			cursor = link.next
		}
	}
}

func (list *ConcurrentList[T]) IterateBackwards() iter.Seq[ConcurrentListEntry[T]] {
	return func(yield func(ConcurrentListEntry[T]) bool) {
		cursor := list.tail.Load()
		for cursor != nil {
			entry := ConcurrentListEntry[T]{list: list, ptr: cursor}
			if !cursor.mark.Load() && !yield(entry) {
				return
			}
			link := cursor.link.Load()
			cursor = link.prev
		}
	}
}

func (list *ConcurrentList[T]) rebindAndSkipLeft(listNode *concurrentListNode[T]) {
	for {
		currentLinkCenter := listNode.link.Load()
		currentLinkLeft := currentLinkCenter.prev.link.Load()
		newLinkLeft := &concurrentListLink[T]{
			next: currentLinkCenter.next,
			prev: currentLinkLeft.prev,
		}
		if currentLinkCenter.prev.link.CompareAndSwap(currentLinkLeft, newLinkLeft) {
			break
		}
	}
}

func (list *ConcurrentList[T]) rebindAndSkipRight(listNode *concurrentListNode[T]) {
	for {
		currentLinkCenter := listNode.link.Load()
		currentLinkRight := currentLinkCenter.next.link.Load()
		newLinkRight := &concurrentListLink[T]{
			next: currentLinkRight.next,
			prev: currentLinkCenter.prev,
		}
		if currentLinkCenter.next.link.CompareAndSwap(currentLinkRight, newLinkRight) {
			break
		}
	}
}

func (list *ConcurrentList[T]) tryPrune() {
	toPruneCount := len(list.markedForDeletion)
	sizeFloat := float64(list.size.Load())
	countFloat := float64(toPruneCount)
	if toPruneCount == 0 || sizeFloat/countFloat <= list.markedNodesThreshold {
		return
	}
	if !list.pruningMutex.TryLock() {
		return
	}
	defer list.pruningMutex.Unlock()
	for range len(list.markedForDeletion) {
		toPrune := <-list.markedForDeletion
		list.removeNode(toPrune)
	}
}

func (list *ConcurrentList[T]) removeNode(listNode *concurrentListNode[T]) {
	if listNode == nil {
		return
	}
	if !listNode.mark.Load() {
		listNode.mark.Store(true)
		list.size.Add(-1)
	}
	link := listNode.link.Load()
	if link.prev == nil {
		list.markedForDeletion <- listNode
		return
	}
	if link.next == nil {
		return
	}
	list.rebindAndSkipLeft(listNode)
	list.rebindAndSkipRight(listNode)
	list.tryPrune()
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
