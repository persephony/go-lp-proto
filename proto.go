package lproto

// This package implements a length prefixed framing protocol.  It adds a 4 byte
// header containing the payload size to each message making the max message size ~4GB.

import (
	"encoding/binary"
	"fmt"
	"io"
)

const (
	// maximum payload size per-frame
	payloadMaxSize uint32 = (4 * 1024 * 1024 * 1024) - 1
)

var (
	errExceededPayload = fmt.Errorf("payload size exceeded")
	errIncompleteRead  = fmt.Errorf("incomplete read")
)

// FrameWriter implements writing a length prefixed protocol.
type FrameWriter struct {
	w io.Writer
}

// NewFrameWriter returns a new frame writer with the given underlying writer.
func NewFrameWriter(w io.Writer) *FrameWriter {
	return &FrameWriter{w}
}

// WriteHeader write the header byte.  This can be used to signal different message
// rpc types.
func (w *FrameWriter) WriteHeader(h byte) error {
	_, err := w.w.Write([]byte{h})
	return err
}

// WriteFrame writes the given data as a length prefixed frame.
func (w *FrameWriter) WriteFrame(data []byte) (int, error) {
	l := len(data)
	if l > int(payloadMaxSize) {
		return 0, errExceededPayload
	}

	err := binary.Write(w.w, binary.BigEndian, uint32(l))
	if err != nil {
		return 0, err
	}

	var c int
	for {
		if c == l {
			break
		}

		var n int
		if n, err = w.w.Write(data[c:]); err != nil {
			break
		}

		c += n
	}
	return c, err
}

// FrameReader implements reading a length prefixed framing protocol
type FrameReader struct {
	r io.Reader
}

// NewFrameReader returns a new frame reader given the underlying reader.
func NewFrameReader(r io.Reader) *FrameReader {
	return &FrameReader{r}
}

// ReadHeader reads a single byte header from the underlying header.  This can be
// used to differentiate messages and/or rpc's.
func (r *FrameReader) ReadHeader() (byte, error) {
	p := make([]byte, 1)
	n, err := r.r.Read(p)
	if err == nil && n != 1 {
		err = errIncompleteRead
	}

	return p[0], err
}

func (r *FrameReader) readSize() (uint32, error) {
	var sz uint32
	err := binary.Read(r.r, binary.BigEndian, &sz)
	if err != nil {
		return 0, err
	}

	return sz, err
}

// ReadFrame reads a frame from the reader using the prefixed size.
func (r *FrameReader) ReadFrame() ([]byte, error) {
	sz, err := r.readSize()
	if err != nil {
		return nil, err
	}
	s := int(sz)

	p := make([]byte, s)
	var c int
	for {
		if c == s {
			break
		}
		var n int
		if n, err = r.r.Read(p[c:]); err != nil {
			break
		}
		c += n

	}
	return p, nil
}
