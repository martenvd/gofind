# SKILLS

This file describes practical skills for implementing changes in `gofind`.

## Skill: CLI Flow Changes

Use when changing startup behavior, argument parsing, or mode switching.

- Primary file: `main.go`
- Related files: `internal/core/prompt.go`, `internal/core/finder.go`, `internal/utils/utils.go`
- Guardrails:
  - Keep existing mode split (no args -> TUI, args -> prompt)
  - Keep update flow working with cache rebuild

## Skill: Directory Discovery Logic

Use when changing search scope, filtering, or repository detection.

- Primary file: `internal/utils/utils.go` (`WalkPaths`, `GetFilteredResults`)
- Guardrails:
  - Preserve `.git` detection intent
  - Keep cache compatibility (`[]string` of absolute paths)
  - Avoid regressions in large directory trees

## Skill: Cache Persistence

Use when changing cache storage format or lifecycle.

- Primary file: `internal/core/directories.go`
- Storage:
  - `~/.gofind/dirs.txt` contains newline-separated paths
- Guardrails:
  - Maintain backward compatibility unless migration is explicit
  - Keep read/write behavior deterministic

## Skill: Interactive Selector UX

Use when changing keyboard behavior or selection interactions.

- TUI mode: `internal/core/finder.go` (tview/tcell)
- Prompt mode: `internal/core/prompt.go` (promptui)
- Guardrails:
  - Enter key should resolve to a selected path output
  - Escape/normal-mode behavior should remain predictable

## Skill: Shell Integration

Use when output contracts or install logic change.

- Shell helper: `gfdir.sh`
- Build/install: `makefile`
- Guardrails:
  - `gfdir.sh` expects raw selected path from `gofind`
  - Avoid printing extra output that breaks `cd` behavior

## Skill: Config Evolution

Use when extending runtime config options.

- Config reader: `internal/utils/utils.go` (`CheckConfig`)
- Config file: `~/.gofind/config.json`
- Current key: `path`
- Guardrails:
  - Missing config must not block normal usage
  - Unknown keys should be tolerated

## Skill: Release Readiness

Use before finalizing any non-trivial change.

Checklist:

1. `go build ./...`
2. Verify no-args TUI launch
3. Verify query argument mode
4. Verify update flag behavior
5. Confirm selected path output remains stdout-only
