# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

`devctl` is a pluggable CLI tool designed to reduce developer friction by codifying rote logic execution. Built in Go, it focus on a plugins and "universal" local developer needs capabilities. 

## Common Development Commands

### Building
- **Local build**: `make build` - Builds for current OS/architecture to `bin/`
- **Cross-platform build**: `GOOS=linux GOARCH=amd64 make build` - Cross-compile for specific target
- **Build all platforms**: `make build-all` - Builds for Linux and macOS on both amd64/arm64
- **Run locally without building**: `go run main.go [command]` - Execute CLI directly from source

### Testing
- **Run all tests**: `go test ./...`
- **Run specific package tests**: `go test ./pkg/secrets/` or `go test ./cmd/`
- **Run with verbose output**: `go test -v ./...`

### Cleanup
- **Clean build artifacts**: `make clean` - Removes `bin/` directory

## Architecture

### Core Structure
- **`main.go`**: Entry point with upgrade checking and build hash injection
- **`cmd/`**: Cobra command definitions and CLI logic
  - `root.go`: Root command setup and plugin registration
  - `jwt.go`: JWT token management commands
  - `jwt_*.go`: Individual JWT subcommands (generate, configure, list, delete)
- **`pkg/`**: Reusable packages organized by functionality

### Key Design Patterns

#### Interface-Based Architecture
The codebase uses interfaces extensively for testability:
- `pkg/secrets/SecretsProvider` interface allows mocking keychain operations
- `DefaultSecretsProvider` variable can be swapped in tests
- Commands use the interface rather than concrete implementations

#### Plugin System
- Scans PATH for `devctl-*` executables and registers them as subcommands
- Plugins are external executables that follow naming convention
- Plugin registration happens in `pkg/plugin/plugin.go:RegisterPlugins()`

#### Build System
- Uses Makefile with Go build with ldflags injection for git hash versioning
- Supports cross-compilation for multiple OS/ARCH combinations
- Build artifacts include git hash in filename for version tracking

### Package Organization
- **`pkg/secrets/`**: OAuth client storage/retrieval using system keychain
- **`pkg/plugin/`**: Dynamic plugin discovery and registration
- **`pkg/update/`**: Self-update functionality via GitHub API

### Dependencies
- **Cobra**: CLI framework for commands and argument parsing
- **golang-jwt/jwt**: JWT token generation and validation
- **golang.org/x/oauth2**: OAuth2 client implementation

## Testing Strategy
- Unit tests alongside source files (`*_test.go`)
- Mock implementations in `testutil/mocks/`
- Tests use interface injection for dependencies
- Coverage includes both success and error scenarios