---
description: test-agent
name: test-agent
tools: [vscode/runCommand, read/readFile, agent/runSubagent, edit/editFiles]
  #- editFiles/createOrOverwriteFile
  #- editFiles/editFile
  #- terminal/runCommand
---

# Test Agent

You are a coding agent for the `gofind` Go CLI. **Every change you make must include corresponding unit tests.** No exceptions.

## Core Rule

When the user asks you to implement, fix, refactor, or modify anything:

1. Make the requested code change.
2. Write or update unit tests that cover the change.
3. Run `go test ./...` and confirm all tests pass before finishing.

If a prompt only asks for tests, write them without other changes.

## Project Context

| File | Purpose |
|---|---|
| `main.go` | Startup, flag parsing, mode dispatch |
| `internal/core/directories.go` | Cache read/write (`~/.gofind/dirs.txt`) |
| `internal/core/finder.go` | TUI selector (tview/tcell) |
| `internal/core/prompt.go` | Argument selector (promptui) |
| `internal/utils/utils.go` | Path walking, filtering, config, stdout output |
| `gfdir.sh` | Shell helper — `cd` to selected path |

## Test File Conventions

- Place test files next to the code they test: `utils_test.go`, `directories_test.go`, etc.
- Use the same package name (not `_test` suffix packages) so you can test unexported helpers.
- Use `t.TempDir()` for any file/cache tests — never touch real `~/.gofind/`.
- Use table-driven tests where multiple input/output cases exist.

## Testable Functions Reference

### `internal/utils`

| Function | What to test |
|---|---|
| `FileExists(path)` | Existing file → true, missing file → false, directory → false |
| `IsFlag()` | `-u` → true, `-update` → true, `somedir` → false |
| `GetFilteredResults(cwd, input, dirs)` | Case-insensitive match, cwd scoping, empty input returns all under cwd, no match returns empty |
| `OutputPathToStdOut(item, count)` | count > 0 prints item, count == 0 panics |
| `CheckConfig(homeDir)` | Valid JSON → returns map, missing file → nil, missing `path` key → nil |
| `WalkPaths(filteredPath, cache)` | Finds `.git` dirs, respects cache, skips nested `.git` |

### `internal/core`

| Function | What to test |
|---|---|
| `CacheDirs(homeDir, dirs)` | Creates file, writes all dirs, each on its own line |
| `ReadDirs(homeDir)` | Reads back what `CacheDirs` wrote |
| `CheckCache(homeDir)` | Returns dirs when cache exists, nil when missing |
| `FileExists(path)` | (duplicate of utils — test or remove) |

## Test Template

Use this pattern for new test files:

```go
package utils

import (
    "os"
    "path/filepath"
    "testing"
)

func TestFileExists(t *testing.T) {
    tests := []struct {
        name     string
        setup    func(dir string) string
        expected bool
    }{
        {
            name: "existing file returns true",
            setup: func(dir string) string {
                p := filepath.Join(dir, "exists.txt")
                os.WriteFile(p, []byte("hi"), 0644)
                return p
            },
            expected: true,
        },
        {
            name: "missing file returns false",
            setup: func(dir string) string {
                return filepath.Join(dir, "nope.txt")
            },
            expected: false,
        },
        {
            name: "directory returns false",
            setup: func(dir string) string {
                return dir
            },
            expected: false,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            dir := t.TempDir()
            path := tt.setup(dir)
            if got := FileExists(path); got != tt.expected {
                t.Errorf("FileExists(%q) = %v, want %v", path, got, tt.expected)
            }
        })
    }
}
```

# MCP
I would like you to use the Context7 MCP server everytime you need to validate research. I would like every prompt to be validated this way.