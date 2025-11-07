# Testing Guide for Totion

This document provides information on how to write and run tests for the Totion project.

## Running Tests

### Run All Tests
```bash
go test ./...
```

### Run Tests with Verbose Output
```bash
go test ./... -v
```

### Run Tests for a Specific Package
```bash
# Test file operations
go test ./internal/file/... -v

# Test app logic
go test ./internal/app/... -v

# Test UI components
go test ./internal/tui/... -v
```

### Run Tests with Coverage
```bash
go test ./... -cover
```

### Generate Coverage Report
```bash
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

## Test Structure

### Test Files
Test files follow the Go naming convention: `*_test.go` alongside the source files they test.

- `internal/file/file_test.go` - Tests for file operations
- `internal/app/app_test.go` - Tests for core app logic
- `internal/tui/components_test.go` - Tests for UI components
- `internal/testhelpers/testhelpers.go` - Shared test utilities

## Writing Tests

### Basic Test Structure

```go
func TestFunctionName(t *testing.T) {
    // Arrange - set up test data
    // Act - execute the function
    // Assert - verify the results
}
```

### Table-Driven Tests

For testing multiple scenarios:

```go
func TestFunctionName(t *testing.T) {
    tests := []struct {
        name     string
        input    string
        expected string
    }{
        {"test case 1", "input1", "expected1"},
        {"test case 2", "input2", "expected2"},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := FunctionName(tt.input)
            if result != tt.expected {
                t.Errorf("got %s, want %s", result, tt.expected)
            }
        })
    }
}
```

### Test Helpers

Use the test helpers in `internal/testhelpers/testhelpers.go`:

```go
import "github.com/AbhaySingh002/Totion/internal/testhelpers"

func TestSomething(t *testing.T) {
    tmpDir := testhelpers.SetupTestEnv(t)
    defer os.RemoveAll(tmpDir) // cleanup handled by SetupTestEnv
    
    filePath := testhelpers.CreateTestNoteFile(t, tmpDir, "test", "content")
    // ... test code
}
```

## Test Coverage

### Current Test Coverage Areas

1. **File Operations** (`internal/file/file_test.go`)
   - ✅ Listing notes from directory
   - ✅ Filtering markdown files
   - ✅ Ignoring non-markdown files and directories
   - ✅ Note model methods

2. **App Logic** (`internal/app/app_test.go`)
   - ✅ Model initialization
   - ✅ Opening and creating files
   - ✅ Saving notes
   - ✅ Handling UI messages (window size, suggestions)
   - ✅ View rendering

3. **UI Components** (`internal/tui/components_test.go`)
   - ✅ Text input configuration
   - ✅ Text area configuration


## Best Practices

### 1. Use Temporary Directories
Always use temporary directories for file system tests:

```go
tmpDir := setupTestNotesDir(t)
defer os.RemoveAll(tmpDir)
```

### 2. Clean Up Resources
Always clean up resources (files, directories, etc.) after tests:

```go
t.Cleanup(func() {
    os.RemoveAll(tmpDir)
})
```

### 3. Test Edge Cases
Test both happy paths and error conditions:

- Empty inputs
- Invalid inputs
- Missing files
- Permission errors

### 4. Use Subtests
Use `t.Run()` to organize related test cases:

```go
func TestFunction(t *testing.T) {
    t.Run("success case", func(t *testing.T) { ... })
    t.Run("error case", func(t *testing.T) { ... })
}
```

### 5. Avoid Testing Implementation Details
Focus on testing behavior, not internal implementation:

- Test public APIs
- Test observable behavior
- Avoid testing private functions directly

## Mocking

### Mocking External Dependencies

For testing components that depend on external services (like the AI client), consider:

1. **Interface-based design**: Define interfaces for external dependencies
2. **Dependency injection**: Pass dependencies as parameters
3. **Mock implementations**: Create test implementations of interfaces

Example:

```go
type AIClient interface {
    GenerateSuggestion(ctx context.Context, prompt string) (string, error)
}

// In tests
type mockAIClient struct {
    suggestion string
    err        error
}

func (m *mockAIClient) GenerateSuggestion(ctx context.Context, prompt string) (string, error) {
    return m.suggestion, m.err
}
```

## Continuous Integration

To run tests in CI/CD:

```bash
# Run tests and fail if coverage is below threshold
go test ./... -cover -coverprofile=coverage.out
go tool cover -func=coverage.out
```

## Debugging Tests

### Run Specific Test
```bash
go test -v -run TestSpecificFunction ./internal/file/...
```

### Run Tests with Race Detector
```bash
go test -race ./...
```

### Verbose Test Output
```bash
go test -v ./...
```

## Common Issues

### Issue: Tests fail due to file permissions
**Solution**: Ensure test files are created with proper permissions and cleaned up after tests.

### Issue: Tests interfere with each other
**Solution**: Use isolated temporary directories for each test case.

### Issue: Tests depend on external services
**Solution**: Mock external dependencies or skip tests that require external services.

## Additional Resources

- [Go Testing Package Documentation](https://pkg.go.dev/testing)
- [Go Testing Best Practices](https://golang.org/doc/effective_go#testing)
- [Table-Driven Tests](https://github.com/golang/go/wiki/TableDrivenTests)

