# Claude Development Rules

This document outlines the development rules and guidelines for this project.

## Test-Driven Development (TDD)

All code changes must follow TDD practices:

1. **Write Tests First**: Before implementing any feature or fix, write a failing test that describes the expected behavior
2. **Minimal Implementation**: Write only enough code to make the test pass
3. **Refactor**: Clean up code while keeping all tests green
4. **No Untested Code**: Every function, method, and exported API must have corresponding tests

### Test Organization
- Unit tests: `*_test.go` files alongside implementation
- Test data: `testdata/` directory
- Integration tests: `integration_test.go` files

## Development Journal

All changes must be documented in `JOURNAL.md`:

### Required Information
- **Date**: YYYY-MM-DD format
- **Description**: What was changed and why
- **Tests**: Which tests were written/modified
- **Implementation**: Brief description of approach
- **Issues**: Any problems encountered and how they were solved
- **Decisions**: Design decisions and trade-offs considered

### Update Frequency
- Update journal before committing code
- Group related changes in a single journal entry
- Keep entries concise but informative

## Commit Message Format

All commits must follow Conventional Commits specification:

### Format
```
<type>: <description>

[optional body]

[optional footer]
```

### Types
- **feat**: New feature for the user
- **fix**: Bug fix
- **test**: Adding or modifying tests
- **refactor**: Code change that neither fixes a bug nor adds a feature
- **docs**: Documentation only changes
- **chore**: Changes to build process, dependencies, or tooling
- **style**: Code style changes (formatting, missing semicolons, etc.)
- **perf**: Performance improvements

### Examples
```
feat: add dice roll parser with regex support

test: add unit tests for spell name validation

fix: handle empty sections in markdown parser

docs: update README with usage examples

refactor: extract dice notation parsing to separate function
```

### Rules
- Use imperative mood: "add" not "added" or "adds"
- Don't capitalize first letter
- No period at the end of subject line
- Keep subject line under 50 characters
- Separate subject from body with blank line
- Wrap body at 72 characters
- Use body to explain what and why, not how

## Code Quality Standards

- **Linting**: Code must pass `go vet` and `golint`
- **Formatting**: Use `gofmt` before committing
- **Error Handling**: All errors must be properly handled or explicitly ignored
- **Documentation**: All exported functions must have Go doc comments
- **Naming**: Follow Go naming conventions (MixedCaps, not snake_case)

## Development Workflow

1. Write failing test (update JOURNAL.md)
2. Implement minimal code to pass test
3. Run tests: `go test ./...`
4. Refactor if needed
5. Run linters: `go vet ./...`
6. Format code: `gofmt -w .`
7. Update JOURNAL.md with final changes
8. Commit with conventional commit message
9. Repeat

## Prohibited Practices

- ❌ No code without tests
- ❌ No commits without journal updates
- ❌ No non-conventional commit messages
- ❌ No unformatted code
- ❌ No ignored errors without justification
- ❌ No unexported functions/types without doc comments
