package lproto

import (
	"bytes"
	"testing"
)

func Test_FrameWriter(t *testing.T) {
	bw := bytes.NewBuffer(nil)
	wr := &FrameWriter{bw}
	n, err := wr.WriteFrame(make([]byte, payloadMaxSize))
	if err != nil {
		t.Fatal(err)
	}

	if n != int(payloadMaxSize) {
		t.Fatal("payload may not be completely written")
	}

}

func Test_FrameReader(t *testing.T) {
	bw := bytes.NewBuffer(nil)

	wr := &FrameWriter{bw}

	if _, err := wr.WriteFrame(make([]byte, int64(payloadMaxSize)+1)); err != errExceededPayload {
		t.Fatal("should fail with size exceeded")
	}

	n, err := wr.WriteFrame([]byte("foobar"))
	if err != nil {
		t.Fatal(err)
	}

	if n != 6 {
		t.Fatal("did not write full data")
	}

	wr.WriteFrame([]byte("foobar1"))
	wr.WriteFrame([]byte("foobar2"))

	rd := &FrameReader{bw}
	f, err := rd.ReadFrame()
	if err != nil {
		t.Fatal(err)
	}
	if string(f) != "foobar" {
		t.Fatal("mismatch")
	}

	f1, _ := rd.ReadFrame()
	f2, _ := rd.ReadFrame()
	if string(f1) != "foobar1" {
		t.Fatal("mismatch")
	}
	if string(f2) != "foobar2" {
		t.Fatal("mismatch")
	}

}
