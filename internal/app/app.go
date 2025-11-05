package app

// models of BubbleTea
import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/AbhaySingh002/Totion/internal/file"
	"github.com/AbhaySingh002/Totion/internal/styles"
	"github.com/AbhaySingh002/Totion/internal/tui"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

var NotesDir string

func init() {
	homedir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal("Error getting home dir", err)
	}
	NotesDir = fmt.Sprintf("%s/.totion", homedir)
}

type Model struct {
	NewFileInput           textinput.Model
	CreateFileInputVisible bool
	CurrentNote            *os.File
	NoteContent            textarea.Model
	List                   list.Model
	ListVisible            bool
	ErrMsg                 string
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m *Model) OpenOrCreateFile(filePath string) error {
	var f *os.File
	var err error
	_, statErr := os.Stat(filePath)
	if statErr == nil {
		f, err = os.OpenFile(filePath, os.O_RDWR, 0644)
	} else if os.IsNotExist(statErr) {
		f, err = os.Create(filePath)
	} else {
		return statErr
	}
	if err != nil {
		return err
	}

	content, err := io.ReadAll(f)
	if err != nil {
		f.Close()
		return err
	}
	if _, err := f.Seek(0, 0); err != nil {
		f.Close()
		return err
	}

	m.CurrentNote = f
	m.NoteContent.SetValue(string(content))
	m.ErrMsg = ""
	return nil
}

func (m *Model) SaveNote() {
	if m.CurrentNote == nil {
		return
	}
	if err := m.CurrentNote.Truncate(0); err != nil {
		m.ErrMsg = fmt.Sprintf("Truncate error: %v", err)
		return
	}
	if _, err := m.CurrentNote.Seek(0, 0); err != nil {
		m.ErrMsg = fmt.Sprintf("Seek error: %v", err)
		return
	}
	if _, err := m.CurrentNote.WriteString(m.NoteContent.Value()); err != nil {
		m.ErrMsg = fmt.Sprintf("Write error: %v", err)
		return
	}
	if err := m.CurrentNote.Close(); err != nil {
		m.ErrMsg = fmt.Sprintf("Close error: %v", err)
		return
	}
	m.CurrentNote = nil
	m.NoteContent.SetValue("")
	m.ErrMsg = ""
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		h, v := styles.DocStyle.GetFrameSize()
		contentWidth := msg.Width - h
		contentHeight := msg.Height - v - 10
		m.List.SetSize(contentWidth, contentHeight)
		m.NoteContent.SetWidth(contentWidth)
		m.NoteContent.SetHeight(contentHeight)
		m.NewFileInput.Width = contentWidth
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+l":
			if !m.ListVisible {
				m.ListVisible = true
				m.CreateFileInputVisible = false
				if m.CurrentNote != nil {
					m.SaveNote()
				}
				m.List.SetItems(file.NotesFiles(NotesDir))
				m.ErrMsg = ""
				return m, nil
			}

		case "delete", "backspace":
			if m.ListVisible {
				item, ok := m.List.SelectedItem().(file.Note)
				if ok {
					filePath := fmt.Sprintf("%s/%s.md", NotesDir, item.Title())
					if err := os.Remove(filePath); err != nil {
						m.ErrMsg = fmt.Sprintf("Error deleting file: %v", err)
					} else {
						m.ErrMsg = ""
						m.List.SetItems(file.NotesFiles(NotesDir))
					}
				} else {
					m.ErrMsg = "No item selected. Use arrow keys to select a note."
				}
				return m, nil
			}

		case "esc":
			if m.CreateFileInputVisible {
				m.CreateFileInputVisible = false
				m.ErrMsg = ""
			}
			if m.CurrentNote != nil {
				m.SaveNote()
				m.CurrentNote = nil
				m.ErrMsg = ""
			}
			if m.ListVisible {
				if m.List.FilterState() == list.Filtering {
					break
				}
				m.ListVisible = false
				m.ErrMsg = ""
			}
			return m, nil
		case "ctrl+c":
			if m.CurrentNote != nil {
				m.SaveNote()
			}
			return m, tea.Quit
		case "ctrl+n":
			if m.CurrentNote != nil {
				m.SaveNote()
			}
			m.ListVisible = false
			m.CreateFileInputVisible = true
			m.ErrMsg = ""
			return m, nil
		case "enter":
			if m.CurrentNote != nil {
				break
			}
			if m.ListVisible {
				item, ok := m.List.SelectedItem().(file.Note)
				if ok {
					filePath := fmt.Sprintf("%s/%s.md", NotesDir, item.Title())
					if err := m.OpenOrCreateFile(filePath); err != nil {
						m.ErrMsg = fmt.Sprintf("Error opening file: %v", err)
					} else {
						m.ListVisible = false
					}
				} else {
					m.ErrMsg = "No item selected. Use arrow keys to select a note."
				}
				return m, nil
			}
			fileName := strings.TrimSpace(m.NewFileInput.Value())
			if fileName != "" {
				filePath := fmt.Sprintf("%s/%s.md", NotesDir, fileName)
				if err := m.OpenOrCreateFile(filePath); err != nil {
					m.ErrMsg = fmt.Sprintf("Error creating/opening file: %v", err)
				} else {
					m.CreateFileInputVisible = false
					m.NewFileInput.SetValue("")
				}
			}
			return m, nil

		case "ctrl+s":
			m.SaveNote()
			return m, nil
		}
	}
	if m.ListVisible {
		m.List, cmd = m.List.Update(msg)
		return m, cmd
	}

	if m.CreateFileInputVisible {
		m.NewFileInput, cmd = m.NewFileInput.Update(msg)
	}

	if m.CurrentNote != nil {
		m.NoteContent, cmd = m.NoteContent.Update(msg)
	}
	return m, cmd
}

func (m Model) View() string {
	errView := ""
	if m.ErrMsg != "" {
		errStyle := styles.DocStyle.Foreground(styles.CursorStyle.GetForeground()).Bold(true).Margin(0, 0, 1, 0)
		errView = errStyle.Render(m.ErrMsg) + "\n"
	}

	var view string
	var help string = GeneralHelp // Default to GeneralHelp

	if m.CreateFileInputVisible {
		view = m.NewFileInput.View()
		help = GeneralHelp
	} else if m.CurrentNote != nil {
		view = m.NoteContent.View()
		help = SaveHelp
	} else if m.ListVisible {
		if len(m.List.Items()) == 0 {
			view = "All Notes 📒\n\nNo notes yet. Press Ctrl+N to create one."
		} else {
			view = m.List.View()
		}
		help = ListHelp
	} else {
		view = "No note open. Press Ctrl+N to create one or Ctrl+L to list existing notes."
	}

	welcome := styles.WelcomeStyle.Render("Welcome to the TOTION 🧠")
	asciiArt := AsciiArt // defined in the data.go
	totionView := styles.TotionLogostyle.Render(asciiArt)
	description := styles.DescriptionStyle.Render("Your personal note-taking companion • Create, edit, and manage your notes with ease using Terminal.")

	return fmt.Sprintf("%s\n%s%s\n%s\n\n%s\n\n%s", welcome, errView, totionView, description, view, help)
}

func InitialModel() Model {

	ti := tui.NewTextInput()
	nt := tui.NewTextArea()
	noteList := file.NotesFiles(NotesDir)
	finallist := list.New(noteList, list.NewDefaultDelegate(), 0, 0)
	finallist.Title = "All Notes 📒"
	finallist.Styles.Title = styles.ListTitleStyle

	return Model{
		NewFileInput:           ti,
		CreateFileInputVisible: false,
		NoteContent:            nt,
		List:                   finallist,
		ListVisible:            false,
		ErrMsg:                 "",
	}
}
