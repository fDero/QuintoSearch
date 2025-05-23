package data

import (
	"testing"
)

func TestSortedArray(t *testing.T) {

	sa := NewSortedArray(
		func(a, b int) bool { return a < b },
		func(a, b int) bool { return a == b },
	)

	if sa.Size() != 0 {
		t.Errorf("Expected size to be 0, got %d", sa.Size())
	}

	if sa.Contains(5) {
		t.Error("Expected to not contain 5")
	}
}

func TestInsertOnce(t *testing.T) {

	sa := NewSortedArray(
		func(a, b int) bool { return a < b },
		func(a, b int) bool { return a == b },
	)

	sa.Insert(5)

	if sa.Size() != 1 {
		t.Errorf("Expected size to be 1, got %d", sa.Size())
	}

	if !sa.Contains(5) {
		t.Error("Expected to contain 5")
	}
}

func insertAll(t *testing.T, sa *SortedArray[int], elements []int) {
	for _, elem := range elements {
		sa.Insert(elem)
	}
}

func ensureContains(t *testing.T, sa *SortedArray[int], elements []int) {
	for _, elem := range elements {
		if !sa.Contains(elem) {
			t.Errorf("Expected to contain %d", elem)
		}
	}
}

func TestInsertMultiple(t *testing.T) {

	sa := NewSortedArray(
		func(a, b int) bool { return a < b },
		func(a, b int) bool { return a == b },
	)

	elements := []int{5, 3, 8, 1, 4, 10, 0, 7, 2, 6, 9}
	insertAll(t, sa, elements)

	if sa.Size() != len(elements) {
		t.Errorf("Expected size to be %d, got %d", len(elements), sa.Size())
	}

	ensureContains(t, sa, elements)
}

func TestEnsureSorted(t *testing.T) {

	sa := NewSortedArray(
		func(a, b int) bool { return a < b },
		func(a, b int) bool { return a == b },
	)

	elements := []int{5, 3, 8, 1, 4, 10, 0, 7, 2, 6, 9}
	insertAll(t, sa, elements)

	previous := -1
	for elem := range sa.Iterate() {
		if elem < previous {
			t.Errorf("Expected sorted order, found %d after %d", elem, previous)
		}
	}
}

func TestEnsureSortedEvenAfterRemove(t *testing.T) {

	sa := NewSortedArray(
		func(a, b int) bool { return a < b },
		func(a, b int) bool { return a == b },
	)

	elements := []int{5, 3, 8, 1, 4, 10, 0, 7, 2, 6, 9}
	for _, elem := range elements {
		sa.Insert(elem)
	}

	for _, elem := range elements {
		if !sa.Contains(elem) {
			t.Errorf("Expected to contain %d", elem)
		}
	}

	previous := -1
	for elem := range sa.Iterate() {
		if elem < previous {
			t.Errorf("Expected sorted order, found %d after %d", elem, previous)
		}
	}

	sa.Remove(10)
	sa.Remove(3)
	sa.Remove(2)

	previous = -1
	for elem := range sa.Iterate() {
		if elem < previous {
			t.Errorf("Expected sorted order, found %d after %d", elem, previous)
		}
	}

	if sa.Size() != len(elements)-3 {
		t.Errorf("Expected size to be %d, got %d", len(elements)-3, sa.Size())
	}
}
