package misc

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
