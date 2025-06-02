package persistence

import (
	"bytes"
	"io"
)

type mockDiskHandler struct {
	mainBuffers map[string]*bytes.Buffer
}

func newMockDiskHandler() *mockDiskHandler {
	return &mockDiskHandler{
		mainBuffers: make(map[string]*bytes.Buffer),
	}
}

func (m *mockDiskHandler) getWriter(key string) (io.Writer, func(), error) {
	if _, exists := m.mainBuffers[key]; !exists {
		m.mainBuffers[key] = new(bytes.Buffer)
	}
	return m.mainBuffers[key], func() {}, nil
}

func (m *mockDiskHandler) getReader(key string) (io.ByteReader, bool) {
	mainBuffer, exists := m.mainBuffers[key]
	return bytes.NewReader(mainBuffer.Bytes()), exists
}
