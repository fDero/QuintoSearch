package data

import (
	"strconv"
	"sync"
	"sync/atomic"
	"testing"
)

func TestConcurrentListInsertions(t *testing.T) {
	const numItems = 1000
	list := NewLinkedList[string]()
	var wg sync.WaitGroup

	entries := make([]ConcurrentListEntry[string], numItems)
	wg.Add(numItems)
	for i := range numItems / 2 {
		go func() {
			entry := list.InsertFront("item-" + strconv.Itoa(i))
			entries[i] = entry
			wg.Done()
			entry = list.InsertFront("item-" + strconv.Itoa(i))
			entries[2*i] = entry
			wg.Done()
		}()
	}
	wg.Wait()

	aliveNodes := CountIterations(list.IterateForward())
	if aliveNodes != numItems {
		t.Fatalf("Expected %d items, got %d", numItems, aliveNodes)
	}

	if size := list.Size(); size != numItems {
		t.Fatalf("Expected list size %d, got %d", numItems, size)
	}
}

func TestConcurrentListInsertionsAndDeletions(t *testing.T) {
	const numItemsStep1 = 500
	const numItemsStep2 = 500
	var expectedSize atomic.Int64

	list := NewLinkedList[string]()
	var wg sync.WaitGroup

	entries := make([]ConcurrentListEntry[string], numItemsStep1)
	wg.Add(numItemsStep1)
	for i := range numItemsStep1 {
		go func() {
			defer wg.Done()
			defer expectedSize.Add(1)
			entry := list.InsertFront("[S1]-item-" + strconv.Itoa(i))
			entries[i] = entry
		}()
	}
	wg.Wait()

	entries2 := make([]ConcurrentListEntry[string], numItemsStep1)
	wg.Add(numItemsStep2)
	for i := range numItemsStep2 {
		go func() {
			defer wg.Done()
			defer expectedSize.Add(1)
			entry := list.InsertFront("[S2]-item-" + strconv.Itoa(i))
			entries2[i] = entry

			if i%2 == 0 {
				expectedSize.Add(-1)
				entries[i].Remove()
			}
		}()
	}
	wg.Wait()

	expectedSizeInt := int(expectedSize.Load())

	aliveNodes := CountIterations(list.IterateForward())
	if aliveNodes != expectedSizeInt {
		t.Fatalf("Expected %d items, got %d", expectedSizeInt, aliveNodes)
	}

	if size := list.Size(); size != expectedSizeInt {
		t.Fatalf("Expected list size %d, got %d", expectedSizeInt, size)
	}
}

func TestConcurrentListInsertionsAndDeletionsAndIterations(t *testing.T) {
	const numItemsStep1 = 500
	const numItemsStep2 = 500
	var expectedSize atomic.Int64

	list := NewLinkedList[string]()
	var wg sync.WaitGroup

	entries := make([]ConcurrentListEntry[string], numItemsStep1)
	wg.Add(numItemsStep1)
	for i := range numItemsStep1 {
		go func() {
			defer wg.Done()
			defer expectedSize.Add(1)
			entry := list.InsertFront("[S1]-item-" + strconv.Itoa(i))
			entries[i] = entry
		}()
	}
	wg.Wait()

	entries2 := make([]ConcurrentListEntry[string], numItemsStep1)
	wg.Add(numItemsStep2)
	for i := range numItemsStep2 {
		go func() {
			defer wg.Done()
			defer expectedSize.Add(1)
			entry := list.InsertFront("[S2]-item-" + strconv.Itoa(i))
			entries2[i] = entry

			if i%2 == 0 {
				expectedSize.Add(-1)
				entries[i].Remove()
			} else {
				for entry := range list.IterateForward() {
					_ = entry.Value()
				}
			}
		}()
	}
	wg.Wait()

	expectedSizeInt := int(expectedSize.Load())

	aliveNodes := CountIterations(list.IterateForward())
	if aliveNodes != expectedSizeInt {
		t.Fatalf("Expected %d items, got %d", expectedSizeInt, aliveNodes)
	}

	if size := list.Size(); size != expectedSizeInt {
		t.Fatalf("Expected list size %d, got %d", expectedSizeInt, size)
	}
}

func TestConcurrentListInsertionsAndDeletionsAndBackwardsIterations(t *testing.T) {
	const numItemsStep1 = 500
	const numItemsStep2 = 500
	var expectedSize atomic.Int64

	list := NewLinkedList[string]()
	var wg sync.WaitGroup

	entries := make([]ConcurrentListEntry[string], numItemsStep1)
	wg.Add(numItemsStep1)
	for i := range numItemsStep1 {
		go func() {
			defer wg.Done()
			defer expectedSize.Add(1)
			entry := list.InsertFront("[S1]-item-" + strconv.Itoa(i))
			entries[i] = entry
		}()
	}
	wg.Wait()

	entries2 := make([]ConcurrentListEntry[string], numItemsStep1)
	wg.Add(numItemsStep2)
	for i := range numItemsStep2 {
		go func() {
			defer wg.Done()
			defer expectedSize.Add(1)
			entry := list.InsertFront("[S2]-item-" + strconv.Itoa(i))
			entries2[i] = entry

			if i%2 == 0 {
				expectedSize.Add(-1)
				entries[i].Remove()
			} else {
				for entry := range list.IterateBackwards() {
					_ = entry.Value()
				}
			}
		}()
	}
	wg.Wait()

	expectedSizeInt := int(expectedSize.Load())

	aliveNodes := CountIterations(list.IterateBackwards())
	if aliveNodes != expectedSizeInt {
		t.Fatalf("Expected %d items, got %d", expectedSizeInt, aliveNodes)
	}

	if size := list.Size(); size != expectedSizeInt {
		t.Fatalf("Expected list size %d, got %d", expectedSizeInt, size)
	}
}
