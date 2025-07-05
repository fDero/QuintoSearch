/*=================================== LICENSE =======================================

                                   Apache License
                             Version 2.0, January 2004
                          http://www.apache.org/licenses/

============================== BRIEF FILE DESCRIPTION ===============================

This files contains an implementation of a concurrency-ready a linked list. Upon
insertion, it returns a reference object that can be used to remove it from the list,
later, in a concurrently-safe manner.
===================================================================================*/

package data

import (
	"iter"
	"quinto/concurrency"
)

type ConcurrentList[T any] struct {
	head  *concurrentListNode[T]
	tail  *concurrentListNode[T]
	mutex concurrency.ReadWriteMutex
}

type concurrentListNode[T any] struct {
	value T
	prev  *concurrentListNode[T]
	next  *concurrentListNode[T]
}

func NewLinkedList[T any]() *ConcurrentList[T] {
	return &ConcurrentList[T]{
		head:  nil,
		tail:  nil,
		mutex: concurrency.NewWritersFirstRWMutex(),
	}
}

func (dll *ConcurrentList[T]) Iterate() iter.Seq[ConcurrentListEntry[T]] {
	return func(yield func(ConcurrentListEntry[T]) bool) {
		dll.mutex.RLock()
		defer dll.mutex.RUnlock()
		concurrentListNode := dll.head
		for concurrentListNode != nil {
			entry := ConcurrentListEntry[T]{list: dll, ptr: concurrentListNode}
			if !yield(entry) {
				break
			}
			concurrentListNode = concurrentListNode.next
		}
	}
}

func (dll *ConcurrentList[T]) InsertBack(value T) ConcurrentListEntry[T] {
	newNode := &concurrentListNode[T]{value: value}
	dll.mutex.Lock()
	defer dll.mutex.Unlock()
	if dll.tail == nil {
		dll.head = newNode
		dll.tail = newNode
	} else {
		newNode.prev = dll.tail
		dll.tail.next = newNode
		dll.tail = newNode
	}
	return ConcurrentListEntry[T]{list: dll, ptr: newNode}
}

func (dll *ConcurrentList[T]) InsertFront(value T) ConcurrentListEntry[T] {
	newNode := &concurrentListNode[T]{value: value}
	dll.mutex.Lock()
	defer dll.mutex.Unlock()
	if dll.head == nil {
		dll.head = newNode
		dll.tail = newNode
	} else {
		newNode.next = dll.head
		dll.head.prev = newNode
		dll.head = newNode
	}
	return ConcurrentListEntry[T]{list: dll, ptr: newNode}
}

func (dll *ConcurrentList[T]) removeNode(listNode *concurrentListNode[T]) {
	dll.mutex.Lock()
	defer dll.mutex.Unlock()

	if listNode == nil {
		return
	}

	if listNode.prev != nil {
		listNode.prev.next = listNode.next
	} else {
		dll.head = listNode.next
	}

	if listNode.next != nil {
		listNode.next.prev = listNode.prev
	} else {
		dll.tail = listNode.prev
	}
}

type ConcurrentListEntry[T any] struct {
	list *ConcurrentList[T]
	ptr  *concurrentListNode[T]
}

func (em ConcurrentListEntry[T]) Value() T {
	if em.ptr != nil && em.list != nil {
		return em.ptr.value
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
