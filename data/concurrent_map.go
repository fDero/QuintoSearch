/*=================================== LICENSE =======================================

                                   Apache License
                             Version 2.0, January 2004
                          http://www.apache.org/licenses/

============================== BRIEF FILE DESCRIPTION ===============================

This files contains a wrapper around a `sync.Map`, which provides a concurrent map
implementation. Wrapping is needed since the choice of `sync.Map` is not
necessarily the best one for all use cases, and might be replaced in the future.
===================================================================================*/

package data

import (
	"sync"
)

type ConcurrentMap[K comparable, V any] struct {
	storage sync.Map
}

func NewConcurrentMap[K comparable, V any]() *ConcurrentMap[K, V] {
	return &ConcurrentMap[K, V]{}
}

func (m *ConcurrentMap[K, V]) Set(key K, value V) {
	m.storage.Store(key, value)
}

func (m *ConcurrentMap[K, V]) Get(key K) (V, bool) {
	value, exists := m.storage.Load(key)
	if !exists {
		var zeroValue V
		return zeroValue, false
	}
	return value.(V), true
}

func (m *ConcurrentMap[K, V]) Delete(key K) {
	m.storage.Delete(key)
}
