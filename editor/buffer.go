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

func (b buffer) stringRange(yi, yj, xi, xj int) string {
	var lines [][]byte
	if b.lines() < yj {
		lines = make([][]byte, len(b.buf[yi:]))
		copy(lines, b.buf[yi:])
	} else {
		lines = make([][]byte, len(b.buf[yi:yj+1]))
		copy(lines, b.buf[yi:yj+1])
	}

	for i, v := range lines {
		if len(v) < xi {
			lines[i] = []byte{}
			continue
		}
		if len(v) < xj {
			lines[i] = v[xi:]
			continue
		}
		lines[i] = v[xi:xj]
	}
	return string(bytes.Join(lines, []byte("\n")))
}
