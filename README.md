# Totion 🧠

A beautiful terminal-based note-taking application built with Go and [Bubble Tea](https://github.com/charmbracelet/bubbletea). Totion provides an intuitive, keyboard-driven interface for creating, editing, and managing your notes directly from the terminal.

![Totion Welcome Screen](totion.png)

## Features

- 📝 **Create and edit notes** - Write and edit markdown notes effortlessly
- 📋 **List all notes** - Browse all your notes with a beautiful terminal UI
- 🗑️ **Delete notes** - Remove notes you no longer need
- 💾 **Auto-save** - Notes are automatically saved when you close them
- 🎨 **Beautiful UI** - Modern terminal interface with styled components
- 🔍 **Search** - Filter through your notes using the built-in search functionality

## Installation

### Prerequisites

- Go 1.25.3 or later
- A terminal with support for ANSI colors

### Building from Source

1. Clone the repository:
```bash
git clone https://github.com/AbhaySingh002/Totion.git
cd Totion
```
• set the api_key in the ./internal/app/data.go
2. Build the application:
```bash
make build
```

Or manually:
```bash
go build -o Totion ./cmd/totion
```

3. Run the application on mac/linux:
```bash
make run
```

Or directly:
```bash
./Totion
```

### Building for Windows

To build for Windows:
```bash
make windows
```

This will create `Totion.exe` in the current directory.

## Usage

When you first run Totion, it will create a `.totion` directory in your home directory where all your notes will be stored as `.md` files.

### Keyboard Shortcuts

#### General Navigation
- `Ctrl+N` - Create a new note
- `Ctrl+L` - List all notes
- `Esc` - Return to home screen / Cancel current action
- `Ctrl+C` - Quit Totion

#### When Editing a Note
- `Ctrl+S` - Save the current note
- `Esc` - Save and close the current note
- `Ctrl+N` - Save current note and create a new one
- `Ctrl+L` - Save current note and open the notes list

#### When Viewing Notes List
- `↑/↓` - Navigate through notes
- `Enter` - Open selected note
- `Delete/Backspace` - Delete selected note
- `/` or start typing - Filter/search notes

## Project Structure

```
Totion/
├── cmd/
│   └── totion/
│       └── main.go          # Application entry point
├── internal/
│   ├── app/
│   │   ├── app.go           # Main application logic and Bubble Tea model
│   │   └── data.go          # Constants and help text
│   ├── file/
│   │   └── file.go          # File operations and note listing
│   ├── styles/
│   │   └── styles.go       # UI styling and colors
│   └── tui/
│       └── components.go    # TUI components (text input, textarea)
├── go.mod                   # Go module dependencies
├── makefile                 # Build commands
└── README.md               # This file
```

## Notes Storage

All notes are stored in `~/.totion/` directory as Markdown (`.md`) files. You can:
- Access your notes directly from the file system
- Edit them with any text editor
- Sync the directory with cloud storage services
- Backup the entire directory

## Dependencies

- [Bubble Tea](https://github.com/charmbracelet/bubbletea) - Terminal UI framework
- [Bubbles](https://github.com/charmbracelet/bubbles) - Bubble Tea components
- [Lip Gloss](https://github.com/charmbracelet/lipgloss) - Styling library

## License

See [LICENSE](LICENSE) file for details.

## Contributing

Contributions are welcome! Feel free to open issues or submit pull requests.

---

Made with ❤️ using Go and Bubble Tea

