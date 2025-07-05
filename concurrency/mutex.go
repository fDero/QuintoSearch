/*=================================== LICENSE =======================================

                                   Apache License
                             Version 2.0, January 2004
                          http://www.apache.org/licenses/

============================== BRIEF FILE DESCRIPTION ===============================

This files contains the definition of "ReadWriteMutex" which describes a mutex that
behaves differently for readers and writers of a particular resource. Keep in mind
that reads can be performed concurrently, while writes are exclusive.

The "WritersFirstRWMutex" is a particular implementation of the "ReadWriteMutex"
interface that prioritizes writers over readers. It is designed to ensure that when a
write operation is pending, all read operations are blocked until the write operation
is completed. This is useful in scenarios where write operations are more critical
than read operations, and we want to minimize the time that writers have to wait
for readers to finish. Usually, writers are fewer than readers, and we want to keep
it that way by making sure that writers are not blocked by readers, hence they don't
pile up. In a databse context, we expect more reads than writes.
==================================================================================*/

package concurrency

import (
	"sync"
	"sync/atomic"
)

type ReadWriteMutex interface {
	RLock()
	RUnlock()
	Lock()
	Unlock()
}

type WritersFirstRWMutex struct {
	simpleRwMutex                 *sync.RWMutex
	pendingWriteOperationsCounter atomic.Int64
	writeOperationsCondVariable   *sync.Cond
}

func NewWritersFirstRWMutex() *WritersFirstRWMutex {
	mut := &WritersFirstRWMutex{
		simpleRwMutex: &sync.RWMutex{},
	}
	mut.writeOperationsCondVariable = sync.NewCond(&sync.Mutex{})
	mut.pendingWriteOperationsCounter.Store(0)
	return mut
}

func (mut *WritersFirstRWMutex) RLock() {
	mut.writeOperationsCondVariable.L.Lock()
	for mut.pendingWriteOperationsCounter.Load() > 0 {
		mut.writeOperationsCondVariable.Wait()
	}
	mut.writeOperationsCondVariable.L.Unlock()
	mut.simpleRwMutex.RLock()
}

func (mut *WritersFirstRWMutex) RUnlock() {
	mut.simpleRwMutex.RUnlock()
}

func (mut *WritersFirstRWMutex) Lock() {
	mut.pendingWriteOperationsCounter.Add(1)
	mut.simpleRwMutex.Lock()
}

func (mut *WritersFirstRWMutex) Unlock() {
	mut.simpleRwMutex.Unlock()
	mut.writeOperationsCondVariable.L.Lock()
	defer mut.writeOperationsCondVariable.L.Unlock()
	mut.pendingWriteOperationsCounter.Add(-1)
	mut.writeOperationsCondVariable.Broadcast()
}
