package go_ringbuffer

import (
	"io"
	"testing"
)

const (
	testBufferSize = 10
)

var (
	smallbuffer  = []byte("12")
	unevenbuffer = []byte("123")
	equalbuffer  = []byte("1234567890")
	largebuffer  = []byte("1234567890ABCDEFGHIJKLMNOPQRSTUVWXYZ")
)

func TestPeekBufferInit(t *testing.T) {
	t.Run("Init", func(t *testing.T) {
		rb := NewPeekBuffer(testBufferSize)
		// check if the buffer is empty and has capacity of testBufferSize
		if rb.Len != 0 {
			t.Errorf("Expected buffer length to be 0 but got %d", rb.Len)
		}
		if rb.Cap != testBufferSize {
			t.Errorf("Expected buffer capacity to be 10 but got %d", rb.Cap)
		}
	})
}

func TestPeekBufferWriteSmall(t *testing.T) {
	t.Run("WriteSmallBuffer", func(t *testing.T) {
		rb := NewPeekBuffer(testBufferSize)
		n, err := rb.Write(smallbuffer)
		if err != nil {
			t.Errorf("Expected no error but got %v", err)
		}
		if n != 2 {
			t.Errorf("Expected to write 2 bytes but wrote %d", n)
		}
		if rb.Len != 2 {
			t.Errorf("Expected buffer length to be 2 but got %d", rb.Len)
		}
		if rb.Pos != 2 {
			t.Errorf("Expected buffer position to be 2 but got %d", rb.Pos)
		}
	})
}

func TestPeekBufferWriteEqual(t *testing.T) {
	t.Run("WriteEqualBuffer", func(t *testing.T) {
		rb := NewPeekBuffer(testBufferSize)
		n, err := rb.Write(equalbuffer)
		if err != nil {
			t.Errorf("Expected no error but got %v", err)
		}
		if n != 10 {
			t.Errorf("Expected to write 10 bytes but wrote %d", n)
		}
		if rb.Len != 10 {
			t.Errorf("Expected buffer length to be 10 but got %d", rb.Len)
		}
		if rb.Pos != 0 {
			t.Errorf("Expected buffer position to be 0 but got %d", rb.Pos)
		}
	})
}

func TestPeekBufferWriteLarge(t *testing.T) {
	t.Run("WriteLargeBuffer", func(t *testing.T) {
		rb := NewPeekBuffer(testBufferSize)
		n, err := rb.Write(largebuffer)
		if err != nil {
			t.Errorf("Expected no error but got %v", err)
		}
		if n != len(largebuffer) {
			t.Errorf("Expected to write %d bytes but wrote %d", len(largebuffer), n)
		}
		if rb.Len != 10 {
			t.Errorf("Expected buffer length to be 10 but got %d", rb.Len)
		}
		if rb.Pos != 0 {
			t.Errorf("Expected buffer position to be 0 but got %d", rb.Pos)
		}
	})
}

func TestPeekBufferReadSmall(t *testing.T) {
	t.Run("ReadSmallBuffer", func(t *testing.T) {
		rb := NewPeekBuffer(testBufferSize)
		rb.Write(smallbuffer)
		readbuffer := make([]byte, len(smallbuffer))
		n, err := rb.Read(readbuffer)
		if err != nil {
			t.Errorf("Expected no error but got %v", err)
		}
		if n != len(smallbuffer) {
			t.Errorf("Expected to read %d bytes but read %d", len(smallbuffer), n)
		}
		if string(readbuffer) != "12" {
			t.Errorf("Expected to read '12' but got '%s'", string(readbuffer))
		}
	})
}

func TestPeekBufferWriteSmallReadLarge(t *testing.T) {
	t.Run("WriteSmallReadLargeBuffer", func(t *testing.T) {
		rb := NewPeekBuffer(testBufferSize)
		rb.Write(smallbuffer)
		readbuffer := make([]byte, len(largebuffer))
		n, err := rb.Read(readbuffer)
		if err != nil {
			t.Errorf("Expected no error but got %v", err)
		}
		if n != 2 {
			t.Errorf("Expected to read 2 bytes but read %d", n)
		}
		if string(readbuffer[:n]) != "12" {
			t.Errorf("Expected to read '12' but got '%s'", string(readbuffer[:n]))
		}
	})
}

func TestPeekBufferReadEqual(t *testing.T) {
	t.Run("ReadEqualBuffer", func(t *testing.T) {
		rb := NewPeekBuffer(testBufferSize)
		rb.Write(equalbuffer)
		readbuffer := make([]byte, len(equalbuffer))
		n, err := rb.Read(readbuffer)
		if err != nil {
			t.Errorf("Expected no error but got %v", err)
		}
		if n != len(equalbuffer) {
			t.Errorf("Expected to read %d bytes but read %d", len(equalbuffer), n)
		}
		if string(readbuffer) != "1234567890" {
			t.Errorf("Expected to read '1234567890' but got '%s'", string(readbuffer))
		}
	})
}

func TestPeekBufferReadLarge(t *testing.T) {
	t.Run("ReadLargeBuffer", func(t *testing.T) {
		rb := NewPeekBuffer(testBufferSize)
		rb.Write(largebuffer)
		readbuffer := make([]byte, len(largebuffer))
		n, err := rb.Read(readbuffer)
		if err != nil {
			t.Errorf("Expected no error but got %v", err)
		}
		if n != testBufferSize {
			t.Errorf("Expected to read %d bytes but read %d", testBufferSize, n)
		}
		if string(readbuffer[:testBufferSize]) != "QRSTUVWXYZ" {
			t.Errorf("Expected to read 'QRSTUVWXYZ' but got '%s'", string(readbuffer[:testBufferSize]))
		}
	})
}

func TestPeekBufferReadEmpty(t *testing.T) {
	t.Run("ReadEmptyBuffer", func(t *testing.T) {
		rb := NewPeekBuffer(testBufferSize)
		readbuffer := make([]byte, 1)
		n, err := rb.Read(readbuffer)
		if err == nil {
			t.Errorf("Expected an error but got nil")
		}
		if err != io.ErrUnexpectedEOF {
			t.Errorf("Expected io.ErrUnexpectedEOF but got %v", err)
		}
		if n != 0 {
			t.Errorf("Expected to read 0 bytes but read %d", n)
		}
	})
}

func TestPeekBufferFill(t *testing.T) {
	t.Run("FillBuffer", func(t *testing.T) {
		rb := NewPeekBuffer(testBufferSize)
		rb.Write(smallbuffer)
		rb.Write(smallbuffer)
		rb.Write(smallbuffer)
		rb.Write(smallbuffer)
		rb.Write(smallbuffer)
		if rb.Len != testBufferSize {
			t.Errorf("Expected buffer length to be %d but got %d", testBufferSize, rb.Len)
		}
		if rb.Pos != 0 {
			t.Errorf("Expected buffer position to be 0 but got %d", rb.Pos)
		}
	})
}

func TestPeekBufferOverfill(t *testing.T) {
	t.Run("OverfillBuffer", func(t *testing.T) {
		rb := NewPeekBuffer(testBufferSize)
		rb.Write(unevenbuffer)
		rb.Write(unevenbuffer)
		rb.Write(unevenbuffer)
		rb.Write(unevenbuffer)
		if rb.Len != testBufferSize {
			t.Errorf("Expected buffer length to be %d but got %d", testBufferSize, rb.Len)
		}
		if rb.Pos != 2 {
			t.Errorf("Expected buffer position to be 2 but got %d", rb.Pos)
		}
	})
}

func TestPeekBufferReadOverfill(t *testing.T) {
	t.Run("ReadOverfillBuffer", func(t *testing.T) {
		rb := NewPeekBuffer(testBufferSize)
		rb.Write(unevenbuffer)
		rb.Write(unevenbuffer)
		rb.Write(unevenbuffer)
		rb.Write(unevenbuffer)
		readbuffer := make([]byte, 4)
		n, err := rb.Read(readbuffer)
		if err != nil {
			t.Errorf("Expected no error but got %v", err)
		}
		if n != 4 {
			t.Errorf("Expected to read 4 bytes but read %d", n)
		}
		if string(readbuffer) != "3123" {
			t.Errorf("Expected to read '3123' but got '%s'", string(readbuffer))
		}
	})
}

func TestPeekBufferReadOverfillWrap(t *testing.T) {
	t.Run("ReadOverfillWrapBuffer", func(t *testing.T) {
		rb := NewPeekBuffer(testBufferSize)
		rb.Write(unevenbuffer)
		rb.Write(unevenbuffer)
		rb.Write(unevenbuffer)
		rb.Write(unevenbuffer)
		if rb.Pos != 2 {
			t.Errorf("Expected buffer position to be 2 but got %d", rb.Pos)
		}
		readbuffer := make([]byte, testBufferSize)
		n, err := rb.Read(readbuffer)
		if err != nil {
			t.Errorf("Expected no error but got %v", err)
		}
		if n != testBufferSize {
			t.Errorf("Expected to read %d bytes but read %d", testBufferSize, n)
		}
		if string(readbuffer) != "3123123123" {
			t.Errorf("Expected to read '3123123123' but got '%s'", string(readbuffer))
		}
	})
}

func TestPeekBufferToBuffer(t *testing.T) {
	t.Run("ToBuffer", func(t *testing.T) {
		rb := NewPeekBuffer(testBufferSize)
		rb.Write(unevenbuffer)
		if rb.Pos != 3 {
			t.Errorf("Expected buffer position to be 3 but got %d", rb.Pos)
		}
		readbuffer := rb.Bytes()
		if len(readbuffer) != 3 {
			t.Errorf("Expected to read %d bytes but read %d", 3, len(readbuffer))
		}
		if string(readbuffer) != "123" {
			t.Errorf("Expected to read '123' but got '%s'", string(readbuffer))
		}
	})
}

func TestPeekBufferToString(t *testing.T) {
	t.Run("ToString", func(t *testing.T) {
		rb := NewPeekBuffer(testBufferSize)
		rb.Write(unevenbuffer)
		if rb.Pos != 3 {
			t.Errorf("Expected buffer position to be 3 but got %d", rb.Pos)
		}
		readstring := rb.String()
		if len(readstring) != 3 {
			t.Errorf("Expected to read %d bytes but read %d", 3, len(readstring))
		}
		if readstring != "123" {
			t.Errorf("Expected to read '123' but got '%s'", readstring)
		}
	})
}

func TestPeekBufferReset(t *testing.T) {
	t.Run("Reset", func(t *testing.T) {
		rb := NewPeekBuffer(testBufferSize)
		rb.Write(unevenbuffer)
		if rb.Len != 3 {
			t.Errorf("Expected buffer length to be 3 but got %d", rb.Len)
		}
		if rb.Pos != 3 {
			t.Errorf("Expected buffer position to be 3 but got %d", rb.Pos)
		}
		rb.Reset()
		if rb.Len != 0 {
			t.Errorf("Expected buffer length to be 0 but got %d", rb.Len)
		}
		if rb.Pos != 0 {
			t.Errorf("Expected buffer position to be 0 but got %d", rb.Pos)
		}
	})
}

func BenchmarkPeekBufferWriteBlocks(b *testing.B) {
	rb := NewPeekBuffer(testBufferSize)
	for i := 0; i < b.N; i++ {
		rb.Write(smallbuffer)
	}
}

func BenchmarkPeekBufferWriteRead5to1(b *testing.B) {
	rbuf := make([]byte, len(equalbuffer))
	rb := NewPeekBuffer(testBufferSize)
	for i := 0; i < b.N; i++ {
		rb.Write(smallbuffer)
		if i%5 == 0 {
			rb.Read(rbuf)
		}
	}
}
