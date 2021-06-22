package editor

import (
	"bytes"
)

type buffer struct {
	buf [][]byte
}

func newBuffer(b []byte) *buffer {
	return &buffer{
		buf: bytes.Split(b, []byte("\n")),
	}
}

func (b buffer) line(i int) []byte {
	if len(b.buf) < i {
		return nil
	}
	return b.buf[i]
}

func (b buffer) lines() int {
	return len(b.buf)
}

func (b buffer) stringRange(i, j int) string {
	if len(b.buf) < j {
		return string(bytes.Join(b.buf[i:len(b.buf)], []byte("\n")))
	}
	return string(bytes.Join(b.buf[i:j], []byte("\n")))
}
