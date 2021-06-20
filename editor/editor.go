package editor

import (
	"bytes"
	"fmt"
	"io"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

type Editor struct {
	f   *os.File
	buf *bytes.Buffer
}

func New(f *os.File, buf *bytes.Buffer) *Editor {
	return &Editor{f: f, buf: buf}
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
	buf := bytes.NewBuffer(content)

	return New(f, buf), nil
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
		case "ctrl+c":
			return e, tea.Quit
		}
	}
	return e, nil
}

func (e *Editor) View() string {
	return e.buf.String()
}
