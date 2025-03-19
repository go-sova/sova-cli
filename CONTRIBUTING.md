# Contributing to Sova CLI

We love your input! We want to make contributing to Sova CLI as easy and transparent as possible, whether it's:

- Reporting a bug
- Discussing the current state of the code
- Submitting a fix
- Proposing new features
- Becoming a maintainer

## Development Process
We use GitHub to host code, to track issues and feature requests, as well as accept pull requests.

1. Fork the repo and create your branch from `main`.
2. If you've added code that should be tested, add tests.
3. If you've changed APIs, update the documentation.
4. Ensure the test suite passes.
5. Make sure your code lints.
6. Issue that pull request!

## Pull Request Process

1. Update the README.md with details of changes to the interface, if applicable.
2. Update the CHANGELOG.md with a note describing your changes.
3. The PR will be merged once you have the sign-off of at least one other developer.

## Any contributions you make will be under the MIT Software License
In short, when you submit code changes, your submissions are understood to be under the same [MIT License](http://choosealicense.com/licenses/mit/) that covers the project. Feel free to contact the maintainers if that's a concern.

## Report bugs using GitHub's [issue tracker](https://github.com/yourusername/sova-cli/issues)
We use GitHub issues to track public bugs. Report a bug by [opening a new issue](https://github.com/yourusername/sova-cli/issues/new); it's that easy!

## Write bug reports with detail, background, and sample code

**Great Bug Reports** tend to have:

- A quick summary and/or background
- Steps to reproduce
  - Be specific!
  - Give sample code if you can.
- What you expected would happen
- What actually happens
- Notes (possibly including why you think this might be happening, or stuff you tried that didn't work)

## License
By contributing, you agree that your contributions will be licensed under its MIT License.

## Code of Conduct

This project and everyone participating in it is governed by our Code of Conduct. By participating, you are expected to uphold this code.

## How Can I Contribute?

### Reporting Bugs

Before creating bug reports, please check the issue list as you might find out that you don't need to create one. When you are creating a bug report, please include as many details as possible:

* Use a clear and descriptive title
* Describe the exact steps which reproduce the problem
* Provide specific examples to demonstrate the steps
* Describe the behavior you observed after following the steps
* Explain which behavior you expected to see instead and why
* Include screenshots if possible

### Suggesting Enhancements

Enhancement suggestions are tracked as GitHub issues. When creating an enhancement suggestion, please include:

* Use a clear and descriptive title
* Provide a step-by-step description of the suggested enhancement
* Provide specific examples to demonstrate the steps
* Describe the current behavior and explain which behavior you expected to see instead
* Explain why this enhancement would be useful

## Development Setup

1. Fork and clone the repository
   ```bash
   git clone https://github.com/go-sova/sova-cli.git
   ```

2. Install dependencies
   ```bash
   go mod download
   ```

3. Run tests
   ```bash
   go test ./...
   ```

4. Build the project
   ```bash
   go build
   ```

## Project Structure

```
.
├── cmd/          # Command implementations
├── internal/     # Private application code
├── pkg/         # Public libraries
├── docs/        # Documentation
└── tests/       # Test files
```

## Coding Style

* Follow standard Go project layout
* Use `gofmt` for formatting
* Follow Go naming conventions
* Write descriptive commit messages
* Add tests for new features

## Testing

* Write unit tests for new features
* Ensure all tests pass before submitting PR
* Include integration tests when needed
* Test edge cases and error conditions

## Documentation

* Update README.md if needed
* Add godoc comments to public functions
* Update wiki pages if needed
* Include examples for new features

## Commit Messages

* Use the present tense ("Add feature" not "Added feature")
* Use the imperative mood ("Move cursor to..." not "Moves cursor to...")
* Limit the first line to 72 characters or less
* Reference issues and pull requests liberally after the first line

## Release Process

1. Update version number in relevant files
2. Update CHANGELOG.md
3. Create a new GitHub release
4. Tag the release with version number
5. Update installation instructions if needed

## Questions?

Feel free to open an issue with your question or contact the maintainers directly.

Thank you for contributing to Sova CLI! 🚀 