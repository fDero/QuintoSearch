package persistence

import (
	"os"
	"quinto/misc"
	"testing"
)

func initTemporaryDirectory(t *testing.T) string {
	tempDir, err := os.MkdirTemp("", "test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	t.Cleanup(func() { os.RemoveAll(tempDir) })
	return tempDir
}

func TestNewPersistenceManager(t *testing.T) {
	tempDir := initTemporaryDirectory(t)

	pm := NewPersistenceManager(tempDir)
	if pm == nil {
		t.Fatal("Expected non-nil PersistenceManager")
	}

	if pm.dbDirectory != tempDir {
		t.Errorf("Expected dbDirectory to be %s, got %s", tempDir, pm.dbDirectory)
	}

	if pm.segments == nil {
		t.Fatal("Expected non-nil segments cache")
	}
}

func TestStoreNewDocument(t *testing.T) {
	tempDir := initTemporaryDirectory(t)
	pm := NewPersistenceManager(tempDir)
	pm.StoreNewDocument(1, func(yield func(misc.Token) bool) {
		_ = yield(misc.Token{StemmedText: "hello", Position: 1}) &&
			yield(misc.Token{StemmedText: "world", Position: 2}) &&
			yield(misc.Token{StemmedText: "hello", Position: 3})
	})
}

func TestStoreAndRetrieve(t *testing.T) {
	tempDir := initTemporaryDirectory(t)
	pm := NewPersistenceManager(tempDir)
	pm.StoreNewDocument(1, func(yield func(misc.Token) bool) {
		_ = yield(misc.Token{StemmedText: "hello", Position: 1}) &&
			yield(misc.Token{StemmedText: "world", Position: 2}) &&
			yield(misc.Token{StemmedText: "hello", Position: 3})
	})
	pm.StoreNewDocument(2, func(yield func(misc.Token) bool) {
		_ = yield(misc.Token{StemmedText: "hello", Position: 1}) &&
			yield(misc.Token{StemmedText: "there", Position: 2})
	})
	expected := []misc.TermTracker{
		{DocumentId: 1, Position: 1},
		{DocumentId: 1, Position: 3},
		{DocumentId: 2, Position: 1},
	}
	counter := 0
	for term := range pm.GetInvertedListIterator("hello") {
		if counter >= len(expected) {
			t.Errorf("Expected %d elements, but got more", len(expected))
			break
		}
		if term != expected[counter] {
			t.Errorf("Expected %v, but got %v", expected[counter], term)
		}
		counter++
	}
}
