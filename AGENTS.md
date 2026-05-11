# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

### Build
```bash
make build          # Build binary to build/bin/gsemver (also runs go generate)
make generate       # Run go generate (regenerates mocks via mockgen)
```

### Test
```bash
make test-style                                                    # Run golangci-lint
make test-coverage                                                 # Unit tests with coverage report
make test-integration                                              # Integration tests (builds first)
go test ./pkg/version/... -run TestBumpMajor -v                   # Run a single unit test
go test ./test/integration/... -run TestConventionalCommits -v    # Run a single integration test
```

### Format
```bash
make format    # Run goimports with local prefix (github.com/arnaud-deprez/gsemver)
```

### Release
```bash
make test-release    # Dry-run release (no publish, shows computed next version + changelog)
make release         # Tag and publish via goreleaser
```

## Architecture

The codebase is split into public API (`pkg/`) and internal implementation (`internal/`), with the CLI layer in `cmd/`.

### Public API (`pkg/`)

- **`pkg/version`** — Core semver logic. The central type is `BumpStrategy`, which holds regex patterns for detecting commit severity and a list of `BumpBranchesStrategy` entries evaluated in order. `BumpStrategy.Bump()` fetches tags, resolves the last annotated tag, gets commits since that tag, and delegates to a `versionBumper` function chosen by matching the current branch against `BumpBranchesStrategy.BranchesPattern`. `NewConventionalCommitBumpStrategy` wires up the default Conventional Commits configuration.

- **`pkg/version.GitRepo`** — Interface that `BumpStrategy` depends on. Has five methods: `FetchTags`, `GetCommits`, `CountCommits`, `GetLastRelativeTag`, `GetCurrentBranch`. Mocks are generated at `pkg/version/mock/git_repo.go` via `//go:generate mockgen`.

- **`pkg/git`** — Data types only: `Commit`, `Hash`, `Tag`, `Signature`. No behavior.

### Implementation (`internal/`)

- **`internal/git`** — CLI-based `GitRepo` implementation (`gitRepoCLI`) that shells out to `git`. `factory.go` constructs it with a `commitParser` that parses `git log` output.

- **`internal/command`** — Thin wrapper around `os/exec` for running shell commands with directory context.

- **`internal/utils`** — Go template helpers (used for pre-release/build-metadata templates), regexp utilities, array helpers.

- **`internal/version`** — `BuildInfo` struct for the tool's own version (injected at build time via `-ldflags`).

- **`internal/release/main.go`** — Standalone program that calls `gsemver` itself to compute the next version during the release process. Used by `Makefile` via `go run internal/release/main.go`.

### CLI (`cmd/`)

Built with Cobra. Three subcommands: `bump` (calls `pkg/version.BumpStrategy.Bump()`), `version` (prints `BuildInfo`), `docs` (generates Markdown CLI docs into `docs/cmd/`).

Configuration is loaded by Viper from `.gsemver.yaml` (current dir) or `~/.gsemver.yaml`. The `GIT_BRANCH` environment variable overrides branch detection for CI environments running in detached HEAD state.

### Integration Tests (`test/integration/`)

Tests build a real temporary git repo at `build/git-tmp`, make commits and tags using the `git` CLI, then invoke the compiled `gsemver` binary. They verify the full end-to-end bump logic including branch strategies, conventional commits, merge scenarios, and go module tag prefixes.

## Key Conventions

- **Annotated tags only**: `gsemver` uses `git describe` which requires annotated tags (`git tag -a -m "..."`) — lightweight tags won't work correctly.
- **Mock generation**: Never edit files under `*/mock/` directly; regenerate with `make generate`.
- **`pkg/` vs `internal/`**: `pkg/` is the stable public API intended for library consumers. Logic that isn't part of the public contract belongs in `internal/`.
- **Templates**: Pre-release and build metadata strings in `BumpBranchesStrategy` are Go templates with access to `Context` (branch, commits, last tag/version). Sprig functions are available.
- **Custom errors**: When returning a custom error, use `pkg/error.NewError(format, args...)` for standalone errors or `pkg/error.NewErrorC(cause, format, args...)` to wrap an underlying error. These produce a consistent `"<message> caused by: <cause>"` format. Do not use bare `fmt.Errorf` for domain errors in `pkg/` or `internal/`.
- **One statement per line**: Do not combine a statement and a condition in a single `if`. Write the statement on its own line first, then check the result. Prefer `err := foo(); \n if err != nil {` over `if err := foo(); err != nil {`.
