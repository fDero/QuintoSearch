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

func (m *mockDiskHandler) getWriter(key string) (writer io.Writer, finalize func(), err error) {
	tmpBuffer := new(bytes.Buffer)
	finalize = func() {
		m.mainBuffers[key] = tmpBuffer
	}
	return tmpBuffer, finalize, nil
}

func (m *mockDiskHandler) getReader(key string) (reader io.ByteReader, exists bool) {
	mainBuffer, ok := m.mainBuffers[key]
	if !ok {
		m.mainBuffers[key] = new(bytes.Buffer)
		mainBuffer = m.mainBuffers[key]
	}
	return bytes.NewReader(mainBuffer.Bytes()), ok
}
