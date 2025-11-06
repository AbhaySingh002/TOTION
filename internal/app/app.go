package app

// models of BubbleTea
import (
	"context"
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
	"github.com/charmbracelet/lipgloss"
	"github.com/joho/godotenv"
	"google.golang.org/genai"
)

var NotesDir string

func init() {
	homedir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal("Error getting home dir", err)
	}
	NotesDir = fmt.Sprintf("%s/.totion", homedir)
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning : error in getting the api-key : %v", err)
	}
}

type Model struct {
	NewFileInput           textinput.Model
	CreateFileInputVisible bool
	CurrentNote            *os.File
	NoteContent            textarea.Model
	List                   list.Model
	ListVisible            bool
	ErrMsg                 string
	Ctx                    context.Context
	Client                 *genai.Client
	Suggestion             string
	AutoCompleteEnabled    bool
	Width                  int
	Height                 int
}

func (m Model) Init() tea.Cmd {
	return nil
}

type suggestionMsg struct {
	suggestion string
	err        error
}

func (m *Model) generateSuggestionCmd() tea.Cmd {
	return func() tea.Msg {
		if m.Client == nil {
			return suggestionMsg{"", fmt.Errorf("AI client not available")}
		}
		temp := float32(0.9)
		prompt := fmt.Sprintf(SystemPrompt, m.NoteContent.Value())
		resp, err := m.Client.Models.GenerateContent(m.Ctx, GenaiModel, genai.Text(prompt), &genai.GenerateContentConfig{Temperature: &temp})
		if err != nil {
			return suggestionMsg{"", err}
		}
		if len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
			return suggestionMsg{"", fmt.Errorf("no suggestion generated")}
		}
		sugg := resp.Text()
		// Truncate to a reasonable length, e.g., first 100 chars or until next period
		return suggestionMsg{sugg, nil}
	}
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

	case suggestionMsg:
		if msg.err != nil {
			m.ErrMsg = fmt.Sprintf("Suggestion error: %v", msg.err)
			m.Suggestion = ""
		} else {
			m.Suggestion = msg.suggestion
			m.ErrMsg = ""
		}
		return m, nil

	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height
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
		case "ctrl+t":
			if m.CurrentNote != nil {
				m.AutoCompleteEnabled = !m.AutoCompleteEnabled
				if !m.AutoCompleteEnabled {
					m.Suggestion = ""
				}
				m.ErrMsg = fmt.Sprintf("Autocomplete %s", map[bool]string{true: "enabled", false: "disabled"}[m.AutoCompleteEnabled])
				return m, nil
			}

		case "tab":
			if m.CurrentNote != nil && m.AutoCompleteEnabled && m.Suggestion != "" {
				current := m.NoteContent.Value()
				m.NoteContent.SetValue(current + " " + m.Suggestion)
				m.Suggestion = ""
				return m, nil
			}

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

		case " ":
			if m.Suggestion != "" {
				break
			}
			if m.CurrentNote != nil && m.AutoCompleteEnabled {
				return m, m.generateSuggestionCmd()
			}
		case "ctrl+g":
			if m.CurrentNote != nil && m.AutoCompleteEnabled {
				return m, m.generateSuggestionCmd()
			}
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
	// Calculate available width for text wrapping
	h, _ := styles.DocStyle.GetFrameSize()
	availableWidth := m.Width - h
	if availableWidth < 40 {
		availableWidth = 40 // Minimum width
	}

	errView := ""
	if m.ErrMsg != "" {
		errStyle := styles.DocStyle.Foreground(styles.CursorStyle.GetForeground()).Bold(true).Margin(0, 0, 1, 0).Width(availableWidth)
		errView = errStyle.Render(m.ErrMsg) + "\n"
	}

	var view string
	var help string = GeneralHelp // Default to GeneralHelp

	if m.CreateFileInputVisible {
		view = m.NewFileInput.View()
		help = GeneralHelp
	} else if m.CurrentNote != nil {
		view = m.NoteContent.View()
		if m.AutoCompleteEnabled && m.Suggestion != "" {
			suggStyle := lipgloss.NewStyle().Italic(true).Foreground(lipgloss.Color("#f1e588ff"))
			suggestionText := suggStyle.Render(m.Suggestion)
			suggestionLineStyle := lipgloss.NewStyle().Width(availableWidth)
			suggestionLine := suggestionLineStyle.Render(fmt.Sprintf("Suggestion: %s (Tab to accept)", suggestionText))
			view += "\n" + suggestionLine + "\n"
		}
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

	welcome := styles.WelcomeStyle.Width(availableWidth).Render("Welcome to the TOTION 🧠")
	asciiArt := AsciiArt // defined in the data.go
	totionView := styles.TotionLogostyle.Width(availableWidth).Render(asciiArt)
	description := styles.DescriptionStyle.Width(availableWidth).Render("Your personal note-taking companion • Create, edit, and manage your notes with ease using Terminal.")

	autocompleteStatus := ""
	if m.CurrentNote != nil {
		status := "off"
		if m.AutoCompleteEnabled {
			status = "on"
		}
		var nextSuggestion string
		if m.AutoCompleteEnabled {
			nextSuggestion = "Ctrl+G: Get the next suggestion"
		} else {
			nextSuggestion = ""
		}
		statusText := fmt.Sprintf("[Autocomplete: %s - Ctrl+T to toggle]", status)
		if nextSuggestion != "" {
			statusText += "\n" + nextSuggestion
		}
		autocompleteStatusStyle := lipgloss.NewStyle().Width(availableWidth).Align(lipgloss.Right)
		autocompleteStatus = "\n" + autocompleteStatusStyle.Render(statusText)
	}

	return fmt.Sprintf("%s\n%s%s%s\n%s\n\n%s\n\n%s", welcome, errView, totionView, autocompleteStatus, description, view, help)
}

func InitialModel() Model {

	ti := tui.NewTextInput()
	nt := tui.NewTextArea()
	noteList := file.NotesFiles(NotesDir)
	finallist := list.New(noteList, list.NewDefaultDelegate(), 0, 0)
	finallist.Title = "All Notes 📒"
	finallist.Styles.Title = styles.ListTitleStyle

	api_key := os.Getenv("GEMINI_API_KEY")
	var client *genai.Client
	if api_key != "" {
		ctx := context.Background()
		var err error
		client, err = genai.NewClient(ctx, &genai.ClientConfig{
			APIKey:  api_key,
			Backend: genai.BackendGeminiAPI,
		})
		if err != nil {
			log.Printf("Failed to initialise the gemini client: %v", err)
			client = nil
		}
	} else {
		log.Printf("Api key is not set, AI Suggestion is disabled.")
	}

	return Model{
		NewFileInput:           ti,
		CreateFileInputVisible: false,
		NoteContent:            nt,
		List:                   finallist,
		ListVisible:            false,
		ErrMsg:                 "",
		Ctx:                    context.Background(),
		Client:                 client,
		Suggestion:             "",
		AutoCompleteEnabled:    false,
		Width:                  80,
		Height:                 24,
	}
}
