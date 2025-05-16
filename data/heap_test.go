package data

import (
	"testing"
)

func TestEmptyHeap(t *testing.T) {
	heap := NewHeap(func(a, b int) bool {
		return a < b
	})

	if size := heap.Size(); size != 0 {
		t.Errorf("Expected heap size to be 0, got %d", size)
	}

	if _, exists := heap.Peek(); exists {
		t.Error("Expected Peek to return false for empty heap")
	}

	if _, exists := heap.Pop(); exists {
		t.Error("Expected Pop to return false for empty heap")
	}
}

func TestSingleElementHeap(t *testing.T) {
	heap := NewHeap(func(a, b int) bool { return a < b })

	heap.Push(42)

	if size := heap.Size(); size != 1 {
		t.Errorf("Expected heap size to be 1, got %d", size)
	}

	if value, exists := heap.Peek(); !exists || value != 42 {
		t.Errorf("Expected Peek to return 42, got %v", value)
	}

	if value, exists := heap.Pop(); !exists || value != 42 {
		t.Errorf("Expected Pop to return 42, got %v", value)
	}

	if size := heap.Size(); size != 0 {
		t.Errorf("Expected heap size to be 0 after Pop, got %d", size)
	}
}

func TestMultipleElementsHeap(t *testing.T) {
	heap := NewHeap(func(a, b int) bool { return a < b })

	elements := []int{5, 3, 8, 1, 4, 10, 0, 7, 2, 6, 9}
	for _, elem := range elements {
		heap.Push(elem)
	}

	if size := heap.Size(); size != len(elements) {
		t.Errorf("Expected heap size to be %d, got %d", len(elements), size)
	}

	lastPopped := -1
	for range elements {
		if value, exists := heap.Pop(); !exists {
			t.Error("Expected Pop to return a value, but got false")
		} else if value < lastPopped {
			t.Errorf("Expected values to be in ascending order, got %d after %d", value, lastPopped)
		}
	}

	if size := heap.Size(); size != 0 {
		t.Errorf("Expected heap size to be 0 after Pop, got %d", size)
	}
}
