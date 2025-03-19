# Development Guide

This guide provides comprehensive instructions and best practices for developing new features in the Sova CLI project.

## Table of Contents
1. [Development Setup](#development-setup)
2. [Project Structure](#project-structure)
3. [Feature Development Process](#feature-development-process)
4. [Testing Guidelines](#testing-guidelines)
5. [Code Quality Standards](#code-quality-standards)
6. [Documentation Requirements](#documentation-requirements)
7. [Release Process](#release-process)
8. [Common Pitfalls and Solutions](#common-pitfalls-and-solutions)

## Development Setup

### Prerequisites
- Go 1.21 or later
- Git
- GitHub account
- Development tools:
  ```bash
  # Install development tools
  go install golang.org/x/tools/cmd/godoc@latest
  go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
  go install github.com/securego/gosec/v2/cmd/gosec@latest
  ```

### Initial Setup
1. Fork the repository
2. Clone your fork:
   ```bash
   git clone https://github.com/YOUR_USERNAME/sova-cli.git
   cd sova-cli
   ```
3. Add upstream remote:
   ```bash
   git remote add upstream https://github.com/go-sova/sova-cli.git
   ```
4. Create a new feature branch:
   ```bash
   git checkout -b feature/your-feature-name
   ```

## Project Structure

```
.
├── cmd/          # Command implementations
│   └── your-feature/
│       └── main.go
├── internal/     # Private application code
│   └── your-feature/
│       └── feature.go
├── pkg/         # Public libraries
│   └── your-feature/
│       └── feature.go
├── docs/        # Documentation
├── tests/       # Test files
└── scripts/     # Build and utility scripts
```

### Directory Purposes
- `cmd/`: Contains the main applications for the project
- `internal/`: Private application and library code
- `pkg/`: Library code that's ok to use by external applications
- `docs/`: Project documentation
- `tests/`: Additional external test applications and test data
- `scripts/`: Various build, install, analysis, and other utility scripts

## Feature Development Process

### 1. Planning Phase
1. Create a feature proposal in GitHub Issues
2. Define requirements:
   - User interface (CLI commands and flags)
   - Expected behavior
   - Error handling
   - Success/failure conditions
3. Document use cases and examples
4. Identify potential edge cases

### 2. Implementation Phase
1. Create necessary directories and files
2. Write tests first (TDD approach)
3. Implement the feature
4. Add command registration
5. Update documentation

### 3. Testing Phase
1. Write unit tests
2. Write integration tests
3. Test edge cases
4. Run all tests:
   ```bash
   go test -v -race -cover ./...
   ```
5. Run linting:
   ```bash
   golangci-lint run
   ```
6. Run security checks:
   ```bash
   gosec ./...
   ```

## Testing Guidelines

### Unit Testing
```go
func TestYourFeature(t *testing.T) {
    tests := []struct {
        name     string
        input    string
        expected string
    }{
        {
            name:     "basic case",
            input:    "test",
            expected: "result",
        },
        // Add more test cases
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := YourFunction(tt.input)
            if result != tt.expected {
                t.Errorf("got %v, want %v", result, tt.expected)
            }
        })
    }
}
```

### Integration Testing
```go
func TestIntegration(t *testing.T) {
    // Setup
    cleanup := setupTestEnv(t)
    defer cleanup()

    // Test implementation
    result := RunFullFeature()
    assertExpectedResult(t, result)
}
```

### Test Coverage Requirements
- Minimum 80% code coverage for new features
- 100% coverage for critical paths
- Test all error conditions
- Test edge cases

## Code Quality Standards

### Go Best Practices
1. Follow standard Go project layout
2. Use `gofmt` for formatting
3. Follow Go naming conventions
4. Keep functions small and focused
5. Use proper error handling
6. Add logging where appropriate

### Error Handling
```go
// Good
if err != nil {
    return fmt.Errorf("failed to process input: %w", err)
}

// Bad
if err != nil {
    return err
}
```

### Logging
```go
// Good
log.Printf("Processing request: %s", requestID)

// Bad
fmt.Printf("Processing request: %s", requestID)
```

## Documentation Requirements

### Code Documentation
```go
// Package feature provides functionality for X
package feature

// ProcessInput handles the input processing with the following steps:
// 1. Validates input
// 2. Transforms data
// 3. Returns result
func ProcessInput(input string) (string, error) {
    // Implementation
}
```

### README Updates
- Add feature description
- Include usage examples
- Document new flags/options
- Update installation instructions if needed

### CHANGELOG Updates
```markdown
## [Unreleased]
### Added
- New feature X with support for Y and Z
- Additional configuration options

### Changed
- Updated existing feature to support new use case

### Fixed
- Bug in feature X when handling edge case
```

## Release Process

### Version Bumping
1. Update version in relevant files
2. Update CHANGELOG.md
3. Create git tag
4. Create GitHub release

### Release Checklist
- [ ] All tests pass
- [ ] Documentation is updated
- [ ] CHANGELOG.md is updated
- [ ] Version numbers are updated
- [ ] Release notes are prepared
- [ ] Installation instructions are updated

## Common Pitfalls and Solutions

### 1. Race Conditions
**Problem**: Concurrent access to shared resources
**Solution**: Use mutexes or channels appropriately
```go
type SafeCounter struct {
    mu    sync.Mutex
    count int
}

func (c *SafeCounter) Increment() {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.count++
}
```

### 2. Resource Leaks
**Problem**: Unclosed resources (files, connections)
**Solution**: Use defer or context cancellation
```go
func processFile(filename string) error {
    file, err := os.Open(filename)
    if err != nil {
        return err
    }
    defer file.Close()
    // Process file
}
```

### 3. Error Handling
**Problem**: Lost error context
**Solution**: Use error wrapping
```go
if err != nil {
    return fmt.Errorf("failed to process %s: %w", filename, err)
}
```

## Getting Help

- Check existing issues and PRs
- Ask in GitHub Discussions
- Contact maintainers
- Review documentation

## Contributing Guidelines

1. Fork the repository
2. Create your feature branch
3. Commit your changes
4. Push to the branch
5. Create a Pull Request

### Commit Message Format
```
type(scope): subject

body

footer
```

Types:
- feat: New feature
- fix: Bug fix
- docs: Documentation
- style: Formatting
- refactor: Code restructuring
- test: Adding tests
- chore: Maintenance

## Additional Resources

- [Go Documentation](https://golang.org/doc/)
- [Effective Go](https://golang.org/doc/effective_go)
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- [Go Testing](https://golang.org/pkg/testing/) 