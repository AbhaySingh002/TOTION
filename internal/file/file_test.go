package file

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/charmbracelet/bubbles/list"
)

// setupTestDir creates a temporary directory for testing
func setupTestDir(t *testing.T) string {
	tmpDir, err := os.MkdirTemp("", "totion-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	return tmpDir
}

// cleanupTestDir removes the temporary directory
func cleanupTestDir(t *testing.T, dir string) {
	if err := os.RemoveAll(dir); err != nil {
		t.Errorf("Failed to cleanup temp dir: %v", err)
	}
}

// createTestNote creates a test note file in the given directory
func createTestNote(t *testing.T, dir, name string, content string) string {
	filePath := filepath.Join(dir, name+".md")
	err := os.WriteFile(filePath, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to create test note: %v", err)
	}
	// Set a known mod time for consistency in tests
	targetTime := time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC)
	if err := os.Chtimes(filePath, targetTime, targetTime); err != nil {
		t.Logf("Warning: could not set mod time: %v", err)
	}
	return filePath
}

func TestNotesFiles(t *testing.T) {
	t.Run("empty directory", func(t *testing.T) {
		tmpDir := setupTestDir(t)
		defer cleanupTestDir(t, tmpDir)

		items := NotesFiles(tmpDir)
		if len(items) != 0 {
			t.Errorf("Expected 0 items, got %d", len(items))
		}
	})

	t.Run("directory with markdown files", func(t *testing.T) {
		tmpDir := setupTestDir(t)
		defer cleanupTestDir(t, tmpDir)

		// Create test notes
		createTestNote(t, tmpDir, "note1", "# Note 1\nContent 1")
		createTestNote(t, tmpDir, "note2", "# Note 2\nContent 2")
		createTestNote(t, tmpDir, "note3", "# Note 3\nContent 3")

		items := NotesFiles(tmpDir)
		if len(items) != 3 {
			t.Errorf("Expected 3 items, got %d", len(items))
		}

		// Verify note titles
		titles := make(map[string]bool)
		for _, item := range items {
			note, ok := item.(Note)
			if !ok {
				t.Errorf("Item is not a Note type")
				continue
			}
			titles[note.Title()] = true
		}

		expectedTitles := map[string]bool{"note1": true, "note2": true, "note3": true}
		for title := range expectedTitles {
			if !titles[title] {
				t.Errorf("Expected note %s not found", title)
			}
		}
	})

	t.Run("ignores non-markdown files", func(t *testing.T) {
		tmpDir := setupTestDir(t)
		defer cleanupTestDir(t, tmpDir)

		// Create markdown files
		createTestNote(t, tmpDir, "note1", "content")

		// Create non-markdown files
		nonMdFile := filepath.Join(tmpDir, "text.txt")
		os.WriteFile(nonMdFile, []byte("text"), 0644)

		items := NotesFiles(tmpDir)
		if len(items) != 1 {
			t.Errorf("Expected 1 item (only .md files), got %d", len(items))
		}
	})

	t.Run("ignores directories", func(t *testing.T) {
		tmpDir := setupTestDir(t)
		defer cleanupTestDir(t, tmpDir)

		createTestNote(t, tmpDir, "note1", "content")

		// Create a subdirectory
		subDir := filepath.Join(tmpDir, "subdir")
		os.Mkdir(subDir, 0755)
		createTestNote(t, subDir, "subnote", "content")

		items := NotesFiles(tmpDir)
		if len(items) != 1 {
			t.Errorf("Expected 1 item (ignoring subdirectory), got %d", len(items))
		}
	})

	t.Run("filters correctly", func(t *testing.T) {
		tmpDir := setupTestDir(t)
		defer cleanupTestDir(t, tmpDir)

		createTestNote(t, tmpDir, "test-note", "content")

		items := NotesFiles(tmpDir)
		if len(items) != 1 {
			t.Fatalf("Expected 1 item, got %d", len(items))
		}

		note := items[0].(Note)
		if note.FilterValue() != "test-note" {
			t.Errorf("Expected FilterValue 'test-note', got '%s'", note.FilterValue())
		}
	})
}

func TestNote(t *testing.T) {
	t.Run("note methods", func(t *testing.T) {
		note := Note{
			title: "test-title",
			desc:  "2024-01-15 10:30",
		}

		if note.Title() != "test-title" {
			t.Errorf("Expected Title 'test-title', got '%s'", note.Title())
		}

		if note.Description() != "2024-01-15 10:30" {
			t.Errorf("Expected Description '2024-01-15 10:30', got '%s'", note.Description())
		}

		if note.FilterValue() != "test-title" {
			t.Errorf("Expected FilterValue 'test-title', got '%s'", note.FilterValue())
		}
	})

	t.Run("implements list.Item interface", func(t *testing.T) {
		var _ list.Item = Note{}
	})
}

func TestNotesFilesErrorHandling(t *testing.T) {
	t.Run("non-existent directory", func(t *testing.T) {
		// This test expects NotesFiles to handle errors gracefully
		// Currently it calls log.Fatal which will exit the program
		// In a real scenario, you might want to return an error instead
		// Note: This will cause log.Fatal, so we can't easily test it
		// This is a code smell - the function should return an error
		// For now, we'll skip this test or document the behavior
		t.Skip("NotesFiles uses log.Fatal which exits the program - consider refactoring to return error")
	})
}
