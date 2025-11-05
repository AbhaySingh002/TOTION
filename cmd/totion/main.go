package main

import (
	"fmt"
	"os"

	"github.com/AbhaySingh002/Totion/internal/app"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	if err := os.MkdirAll(app.NotesDir, 0750); err != nil {
		fmt.Printf("could not create notes directory: %v", err)
		os.Exit(1)
	}

	p := tea.NewProgram(app.InitialModel())

	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}