package editor

import (
	"fmt"
	"io"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/muesli/termenv"
)

type Editor struct {
	f *os.File

	buf     *buffer
	display *display
	cursor  *cursor
}

type display struct {
	window *window
	dRange *displayRange
}

type window struct {
	width, height int
}

type displayRange struct {
	top, bottom, left, right int
}

type cursor struct {
	line, column int
}

func New(f *os.File, buf *buffer) (*Editor, error) {
	return &Editor{
		f:   f,
		buf: buf,
		display: &display{
			window: &window{},
			dRange: &displayRange{},
		},
		cursor: &cursor{},
	}, nil
}

func Open(fileName string) (*Editor, error) {
	f, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return nil, fmt.Errorf("failed to open the file: %s", fileName)
	}
	content, err := io.ReadAll(f)
	if err != nil {
		return nil, fmt.Errorf("failed to read the file: %s", f.Name())
	}
	buf := newBuffer(content)

	return New(f, buf)
}

func (e *Editor) Close() error {
	return e.f.Close()
}

func (e *Editor) Start() error {
	p := tea.NewProgram(e, tea.WithAltScreen())
	return p.Start()
}

func (e *Editor) adjustDisplayRange() {
	// move to the outside of the display range
	if e.display.dRange.left > e.cursor.column {
		if e.cursor.column > e.display.window.width {
			e.display.dRange.right = e.cursor.column + e.display.window.width/2
			e.display.dRange.left = e.display.dRange.right - e.display.window.width
		} else {
			e.display.dRange.left = 0
			e.display.dRange.right = e.display.window.width
		}
	}
	// up
	if e.display.dRange.top > e.cursor.line {
		e.display.dRange.top--
		e.display.dRange.bottom--
	}
	// down
	if e.display.dRange.bottom < e.cursor.line {
		e.display.dRange.top++
		e.display.dRange.bottom++
	}
	// left
	if e.display.dRange.left > e.cursor.column {
		e.display.dRange.left--
		e.display.dRange.right--
	}
	// right
	if e.display.dRange.right < e.cursor.column+1 {
		e.display.dRange.left++
		e.display.dRange.right++
	}
}

// Elm Architecture

func (e *Editor) Init() tea.Cmd {
	return nil
}

func (e *Editor) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "k":
			if e.cursor.line > 0 {
				e.cursor.line--
				// move to short line
				if e.cursor.column > len(e.buf.line(e.cursor.line)) {
					e.cursor.column = len(e.buf.line(e.cursor.line)) - 1
					if e.cursor.column < 0 {
						e.cursor.column = 0
					}
				}
			}
		case "j":
			if e.cursor.line < len(e.buf.lines)-1 {
				e.cursor.line++
				// move to short line
				if e.cursor.column > len(e.buf.line(e.cursor.line)) {
					e.cursor.column = len(e.buf.line(e.cursor.line)) - 1
					if e.cursor.column < 0 {
						e.cursor.column = 0
					}
				}
			}
		case "h":
			if e.cursor.column > 0 {
				e.cursor.column--
			}
		case "l":
			if e.cursor.column < len(e.buf.line(e.cursor.line))-1 {
				e.cursor.column++
			}
		case "ctrl+c":
			return e, tea.Quit
		}
	case tea.WindowSizeMsg:
		e.display.window.width = msg.Width
		e.display.window.height = msg.Height
		e.display.dRange.bottom = e.display.dRange.top + msg.Height - 1
		e.display.dRange.right = e.display.dRange.left + msg.Width
	}

	e.adjustDisplayRange()
	return e, nil
}

func (e *Editor) View() string {
	top := e.display.dRange.top
	bottom := e.display.dRange.bottom
	left := e.display.dRange.left
	right := e.display.dRange.right

	// secure space for cursor
	buf := e.buf.copy()
	for i := range buf.lines {
		if len(buf.lines[i]) == 0 {
			buf.lines[i] = append(buf.lines[i], ' ')
		}
	}

	// style cursor
	line := buf.line(e.cursor.line)
	style := termenv.String(string(line[e.cursor.column]))
	p := termenv.ColorProfile()
	style = style.Foreground(p.Color("#000000")).Background(p.Color("#eeeeee"))
	tmp := string(line[:e.cursor.column]) + style.String() + string(line[e.cursor.column+1:])
	buf.lines[e.cursor.line] = []byte(tmp)

	return buf.stringRange(top, bottom, left, right+len(style.String())-1)
}
