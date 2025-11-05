package file

import (
	"log"
	"os"

	"github.com/charmbracelet/bubbles/list"
)

type Note struct {
	title string
	desc  string
}

func (n Note) Title() string       { return n.title }
func (n Note) Description() string { return n.desc }
func (n Note) FilterValue() string { return n.title }

func NotesFiles(notesDir string) []list.Item {
	entries, err := os.ReadDir(notesDir)
	if err != nil {
		log.Fatal(err)
	}
	items := make([]list.Item, 0)
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		if len(name) > 3 && name[len(name)-3:] == ".md" {
			info, err := e.Info()
			if err != nil {
				log.Fatalf("%v", err)
			}
			modTime := info.ModTime().Format("2006-01-02 15:04")
			items = append(items, Note{
				title: name[:len(name)-3],
				desc:  modTime,
			})
		}
	}
	return items
}