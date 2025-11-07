package testhelpers

import (
	"os"
	"path/filepath"
	"testing"
)

// SetupTestEnv creates a temporary directory and returns its path
// It also sets up a cleanup function that will be called when the test completes
func SetupTestEnv(t *testing.T) string {
	tmpDir, err := os.MkdirTemp("", "totion-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}

	t.Cleanup(func() {
		if err := os.RemoveAll(tmpDir); err != nil {
			t.Errorf("Failed to cleanup temp dir: %v", err)
		}
	})

	return tmpDir
}

// CreateTestNoteFile creates a markdown note file in the given directory
func CreateTestNoteFile(t *testing.T, dir, name, content string) string {
	filePath := filepath.Join(dir, name+".md")
	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to create test note file: %v", err)
	}
	return filePath
}

// CreateTestFiles creates multiple test note files
func CreateTestFiles(t *testing.T, dir string, files map[string]string) []string {
	var paths []string
	for name, content := range files {
		path := CreateTestNoteFile(t, dir, name, content)
		paths = append(paths, path)
	}
	return paths
}

// FileExists checks if a file exists at the given path
func FileExists(t *testing.T, path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	t.Fatalf("Error checking file existence: %v", err)
	return false
}

// ReadFileContent reads and returns the content of a file
func ReadFileContent(t *testing.T, path string) string {
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}
	return string(content)
}
