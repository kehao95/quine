package tools

import (
	"bytes"
	"sync"
)

// SafeBuffer is a thread-safe byte buffer.
// Pump goroutines write to it concurrently; Peek/TakeAll read from it.
type SafeBuffer struct {
	buf bytes.Buffer
	mu  sync.Mutex
}

// Write appends data to the buffer. Safe for concurrent use.
func (b *SafeBuffer) Write(p []byte) (int, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.buf.Write(p)
}

// Peek returns the buffer contents without consuming them.
func (b *SafeBuffer) Peek() string {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.buf.String()
}

// TakeAll returns the buffer contents and resets it (consume pattern).
func (b *SafeBuffer) TakeAll() string {
	b.mu.Lock()
	defer b.mu.Unlock()
	s := b.buf.String()
	b.buf.Reset()
	return s
}

// Len returns the current byte count in the buffer.
func (b *SafeBuffer) Len() int {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.buf.Len()
}
