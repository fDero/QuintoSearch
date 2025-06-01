package persistence

import (
	"bytes"
	"io"
)

type mockDiskHandler struct {
	mainBuffers      map[string]*bytes.Buffer
	temporaryBuffers map[string]*bytes.Buffer
	initFlags        map[string]bool
}

func getBufferOrCreateIt(buffers map[string]*bytes.Buffer, key string) *bytes.Buffer {
	if buffer, exists := buffers[key]; exists {
		return buffer
	}
	newBuffer := &bytes.Buffer{}
	buffers[key] = newBuffer
	return newBuffer
}

func newMockDiskHandler() *mockDiskHandler {
	return &mockDiskHandler{
		mainBuffers:      make(map[string]*bytes.Buffer),
		temporaryBuffers: make(map[string]*bytes.Buffer),
		initFlags:        make(map[string]bool),
	}
}

func (m *mockDiskHandler) getWriter(key string) (io.Writer, func(), error) {
	temporaryBuffer := getBufferOrCreateIt(m.temporaryBuffers, key)
	finalize := func() {
		m.mainBuffers[key] = temporaryBuffer
		m.temporaryBuffers[key] = &bytes.Buffer{}
		m.initFlags[key] = true
	}
	return temporaryBuffer, finalize, nil
}

func (m *mockDiskHandler) getReader(key string) (io.ByteReader, bool) {
	mainBuffer := getBufferOrCreateIt(m.mainBuffers, key)
	return bytes.NewReader(mainBuffer.Bytes()), m.initFlags[key]
}
