package editor

import (
	"bytes"
)

type buffer struct {
	lines [][]byte
}

func newBuffer(b []byte) *buffer {
	return &buffer{
		lines: bytes.Split(b, []byte("\n")),
	}
}

func (b *buffer) line(i int) []byte {
	if len(b.lines) < i {
		return nil
	}
	return b.lines[i]
}

func (b *buffer) stringRange(yi, yj, xi, xj int) string {
	var lines [][]byte
	if len(b.lines) < yj {
		lines = make([][]byte, len(b.lines[yi:]))
		copy(lines, b.lines[yi:])
	} else {
		lines = make([][]byte, len(b.lines[yi:yj+1]))
		copy(lines, b.lines[yi:yj+1])
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

func (b *buffer) copy() *buffer {
	c := &buffer{
		lines: make([][]byte, len(b.lines)),
	}
	for i := range c.lines {
		c.lines[i] = make([]byte, len(b.line(i)))
		copy(c.lines[i], b.line(i))
	}
	return c
}
