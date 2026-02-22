# Repository Guidelines

## Project Structure & Module Organization
- `cmd/fod/main.go`: CLI entrypoint for the `fod` binary.
- `*.go` at repo root: core library and UI/selection logic.
- `*_test.go`: Go unit tests (currently `define_test.go`).
- `pkg-config-files/`: build/config artifacts used by the project.
- `Makefile`, `go.mod`, `go.sum`: build orchestration and dependencies.

## Build, Test, and Development Commands
- `make test`: run all Go tests (`go test ./...`).
- `make lint`: run `golangci-lint` on all packages.
- `make fmt`: format with `goimports`.
- `make deps`: download modules (`go mod download`).
- `go run ./cmd/fod`: run the CLI locally.
- `go build ./cmd/fod`: build the CLI binary.

## Coding Style & Naming Conventions
- Go standard formatting; use `make fmt` (Goimports).
- Follow Go naming: exported identifiers in `CamelCase`, unexported in `camelCase`.
- Keep files organized by feature/concern; tests live next to the code they cover.

## Testing Guidelines
- Framework: Go `testing` package.
- Test files must be named `*_test.go`; test functions `TestXxx`.
- Run all tests with `make test` or `go test ./...`.

## Commit & Pull Request Guidelines
- Commit messages follow a simple type prefix pattern: `fix: ...`, `refactor: ...`, `test: ...`.
- PRs should include:
  - A concise description of behavior changes.
  - Tests run (or a note if not run).
  - Linked issue if applicable.

## Tooling & Configuration
- Required tools: Go, `golangci-lint`, `make` (see `README.md`).
- Optional tooling managed via `mise.toml` (`golangci-lint`, `node`).
