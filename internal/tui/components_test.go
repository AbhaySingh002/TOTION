package tui

import (
	"testing"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
)

func TestNewTextInput(t *testing.T) {
	ti := NewTextInput()

	t.Run("placeholder is set", func(t *testing.T) {
		if ti.Placeholder != "What do you wanna name it?" {
			t.Errorf("Expected placeholder 'What do you wanna name it?', got '%s'", ti.Placeholder)
		}
	})

	t.Run("character limit is set", func(t *testing.T) {
		if ti.CharLimit != 30 {
			t.Errorf("Expected CharLimit 30, got %d", ti.CharLimit)
		}
	})

	t.Run("width is set", func(t *testing.T) {
		if ti.Width != 50 {
			t.Errorf("Expected Width 50, got %d", ti.Width)
		}
	})

	t.Run("focus method exists and returns command", func(t *testing.T) {
		// Focus() returns a tea.Cmd, which means the component is set up for focus
		focusCmd := ti.Focus()
		if focusCmd == nil {
			t.Error("Expected Focus() to return a command")
		}
	})

	t.Run("returns textinput.Model", func(t *testing.T) {
		var _ textinput.Model = ti
	})
}

func TestNewTextArea(t *testing.T) {
	nt := NewTextArea()

	t.Run("placeholder is set", func(t *testing.T) {
		if nt.Placeholder != "Type your notes...." {
			t.Errorf("Expected placeholder 'Type your notes....', got '%s'", nt.Placeholder)
		}
	})

	t.Run("show line numbers is false", func(t *testing.T) {
		if nt.ShowLineNumbers != false {
			t.Error("Expected ShowLineNumbers to be false")
		}
	})

	t.Run("focus method exists and returns command", func(t *testing.T) {
		// Focus() returns a tea.Cmd, which means the component is set up for focus
		focusCmd := nt.Focus()
		if focusCmd == nil {
			t.Error("Expected Focus() to return a command")
		}
	})

	t.Run("returns textarea.Model", func(t *testing.T) {
		var _ textarea.Model = nt
	})
}
