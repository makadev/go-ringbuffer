package go_ringbuffer

import (
	"io"
	"sync"
)

// PeekBuffer is a thread-safe ring buffer that allows typical append-only (stream) write access.
// The difference to a stream ringbuffer is that it also allows random full reads of the buffer contents ("peeking").
type PeekBuffer struct {
	// Lock is the mutex for the buffer.
	Lock sync.RWMutex
	// Buffer is the underlying byte slice.
	Buffer []byte
	// Cap is the capacity of the buffer.
	Cap int
	// Len is the length of the buffer (number of bytes written).
	// If Len == Cap, the buffer is full and the next write will overwrite the oldest data.
	Len int
	// Pos is the current position in the buffer (next write position).
	// If Pos == Cap, the next write will wrap around to the beginning of the buffer.
	//
	// Note: If Len < Cap, Pos is the next write position but not the first byte to read.
	Pos int
}

// Create a new buffer with given size.
func NewPeekBuffer(size int) *PeekBuffer {
	buf := make([]byte, size)
	return &PeekBuffer{
		Buffer: buf,
		Cap:    size,
		Pos:    0,
	}
}

// Write writes p to the buffer.
// If the buffer is full, it will overwrite the oldest data.
// It returns the number of bytes written from p (0 <= n <= len(p)) and an error, if any.
// Implements the io.Writer interface.
func (r *PeekBuffer) Write(p []byte) (n int, err error) {
	var written int
	r.Lock.Lock()
	defer r.Lock.Unlock()
	if len(p) <= r.Cap-r.Pos {
		// case 1.: p is smaller than the buffer and smaller than (or equal) the remaining space
		written = copy(r.Buffer[r.Pos:], p)
		r.Pos += written
	} else if len(p) >= r.Cap {
		// case 2.: p is filling the buffer, we only need to write the last buffer.cap bytes
		// note: This is effectively the same as creating a new buffer with the last buffer.cap bytes of p hence
		// any read pointer into the buffer will be invalid after this operation
		written = copy(r.Buffer, p[len(p)-r.Cap:])
		r.Pos = 0
	} else {
		// case 3.: p is smaller than the buffer but larger than the remaining space so it wraps around.
		written = copy(r.Buffer[r.Pos:], p)
		remaining := r.Cap - r.Pos
		written += copy(r.Buffer, p[remaining:])
		r.Pos = len(p) - remaining
	}

	if r.Pos >= r.Cap {
		// always wrap around the write position if it exceeds the capacity
		r.Pos = 0
	}

	if r.Len < r.Cap {
		// update the length of the buffer if it is not using its full capacity
		r.Len += written
		if r.Len > r.Cap {
			r.Len = r.Cap
		}
	}
	// always return len(p)
	return len(p), nil
}

// Read reads up to len(p) bytes from the buffer into p.
// It returns the number of bytes read (0 <= n <= len(p)) and an error, if any.
// If the buffer is empty, Read returns io.ErrUnexpectedEOF.
// Implements the io.Reader interface.
func (r *PeekBuffer) Read(p []byte) (n int, err error) {
	r.Lock.RLock()
	defer r.Lock.RUnlock()
	// get the number of bytes to read
	var readlen = len(p)
	if len(p) > r.Len {
		readlen = r.Len
		// if the buffer is empty, exit early
		if readlen == 0 {
			return 0, io.ErrUnexpectedEOF
		}
	}

	var readPos = r.Pos
	if r.Pos >= r.Len {
		// if the read position is equal to the length the buffer was not fully written the first time
		// so we need to read from the beginning
		readPos = 0
	}

	if readPos+readlen <= r.Len {
		// case 1.: readlen bytes can be read without wrapping around
		n = copy(p, r.Buffer[readPos:readPos+readlen])
	} else {
		// case 2.: readlen bytes need to be read with wrapping around
		remaining := r.Cap - r.Pos
		n = copy(p, r.Buffer[r.Pos:])
		n += copy(p[n:], r.Buffer[:readlen-remaining])
	}
	return n, nil
}

func (r *PeekBuffer) Reset() {
	r.Lock.Lock()
	defer r.Lock.Unlock()
	r.Len = 0
	r.Pos = 0
}

func (r *PeekBuffer) Bytes() []byte {
	r.Lock.RLock()
	defer r.Lock.RUnlock()
	var res []byte = make([]byte, r.Len)
	if r.Len > 0 {
		r.Read(res)
	}
	return res
}

func (r *PeekBuffer) String() string {
	return string(r.Bytes())
}
