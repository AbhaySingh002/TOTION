package tui

import (
	"github.com/AbhaySingh002/Totion/internal/styles"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
)

func NewTextInput() textinput.Model {
	ti := textinput.New()
	ti.Placeholder = "What do you wanna name it?"
	ti.Focus()
	ti.CharLimit = 30
	ti.Width = 50
	ti.Cursor.Style = styles.CursorStyle
	return ti
}

func NewTextArea() textarea.Model {
	nt := textarea.New()
	nt.Placeholder = "write the content here"
	nt.Focus()
	nt.ShowLineNumbers = false
	nt.Placeholder = "Type your notes...."
	nt.Cursor.Style = styles.CursorStyle
	return nt
}