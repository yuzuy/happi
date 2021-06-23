package editor

import (
	"fmt"
	"io"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

type Editor struct {
	f *os.File

	buf     *buffer
	display *display
	cursor  *cursor
}

type display struct {
	window       *window
	displayRange *displayRange
}

type window struct {
	width  int
	height int
}

type displayRange struct {
	top, bottom, left, right int
}

func (d *displayRange) up() {
	d.top--
	d.bottom--
}

func (d *displayRange) down() {
	d.top++
	d.bottom++
}

func (d *displayRange) moveRight() {
	d.left++
	d.right++
}

func (d *displayRange) moveLeft() {
	d.left--
	d.right--
}

type cursor struct {
	line   int
	column int
}

func New(f *os.File, buf *buffer) (*Editor, error) {
	return &Editor{
		f:   f,
		buf: buf,
		display: &display{
			window:       &window{},
			displayRange: &displayRange{},
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
				if e.display.displayRange.top > e.cursor.line {
					e.display.displayRange.up()
				}
			}
		case "j":
			if e.cursor.line < e.buf.lines()-1 {
				e.cursor.line++
				if e.display.displayRange.bottom < e.cursor.line {
					e.display.displayRange.down()
				}
			}
		case "h":
			if e.cursor.column > 0 {
				e.cursor.column--
				if e.display.displayRange.left > e.cursor.column {
					e.display.displayRange.moveLeft()
				}
			}
		case "l":
			if e.cursor.column < len(e.buf.line(e.cursor.line)) {
				e.cursor.column++
				if e.display.displayRange.right < e.cursor.column {
					e.display.displayRange.moveRight()
				}
			}
		case "ctrl+c":
			return e, tea.Quit
		}
	case tea.WindowSizeMsg:
		e.display.displayRange.bottom = e.display.displayRange.top + msg.Height - 1
		e.display.displayRange.right = e.display.displayRange.left + msg.Width
	}
	return e, nil
}

func (e *Editor) View() string {
	top := e.display.displayRange.top
	bottom := e.display.displayRange.bottom
	left := e.display.displayRange.left
	right := e.display.displayRange.right
	return e.buf.stringRange(top, bottom, left, right)
}
