package persistence

import (
	"fmt"
	"quinto/core"
	"quinto/data"
	"testing"
)

func UtilCreateChunks(term string, handler diskHandler, trackers_list [][]core.TermTracker) []core.TermTracker {
	ogkey := "term-" + term
	key := ogkey
	counter := 0
	for _, trackers := range trackers_list {
		chunk := newIndexChunk(key, handler)
		chunk.insertIterable(data.NewSliceIterator(trackers))
		counter++
		key = ogkey + "-" + fmt.Sprint(counter)
		chunk.nextChunkKey = key
		chunk.writeBack()
	}
	concatenated := []core.TermTracker{}
	for _, trackers := range trackers_list {
		concatenated = append(concatenated, trackers...)
	}
	return concatenated
}

func TestLoggingAllChunksWithPersistenceManager(t *testing.T) {
	handler := newMockDiskHandler()

	UtilCreateChunks("hello", handler, [][]core.TermTracker{
		{
			{DocId: 17, Position: 3},
			{DocId: 17, Position: 4},
			{DocId: 27, Position: 1},
			{DocId: 37, Position: 1},
		},
		{
			{DocId: 41, Position: 5},
			{DocId: 41, Position: 9},
			{DocId: 41, Position: 11},
		}})

	UtilCreateChunks("world", handler, [][]core.TermTracker{
		{
			{DocId: 27, Position: 2},
			{DocId: 37, Position: 2},
		}})

	manager := NewPersistenceManager(PersistenceConfig{
		MaxCachedChunks: 10,
		MaxChunkSize:    1024,
		IoHandler:       handler,
	})

	hellos := manager.IterateTerms("hello")

	for hello := range hellos {
		t.Log(hello.DocId, ":", hello.Position, " ")
	}

	t.Log("Logs emitted")
}

func TestRetrieveChunksWithPersistenceManager(t *testing.T) {
	handler := newMockDiskHandler()

	hello_trackers := UtilCreateChunks("hello", handler, [][]core.TermTracker{
		{
			{DocId: 17, Position: 3},
			{DocId: 17, Position: 4},
			{DocId: 27, Position: 1},
			{DocId: 37, Position: 1},
		},
		{
			{DocId: 41, Position: 5},
			{DocId: 41, Position: 9},
			{DocId: 41, Position: 11},
		}})

	world_trackers := UtilCreateChunks("world", handler, [][]core.TermTracker{
		{
			{DocId: 27, Position: 2},
			{DocId: 37, Position: 2},
		}})

	manager := NewPersistenceManager(PersistenceConfig{
		MaxCachedChunks: 10,
		MaxChunkSize:    1024,
		IoHandler:       handler,
	})

	hellos := manager.IterateTerms("hello")
	worlds := manager.IterateTerms("world")

	if iters := data.CountIterations(hellos); len(hello_trackers) != iters {
		t.Errorf("Expected %d hello trackers, got %d", len(hello_trackers), iters)
	}

	if iters := data.CountIterations(worlds); len(world_trackers) != iters {
		t.Errorf("Expected %d world trackers, got %d", len(world_trackers), iters)
	}
}
