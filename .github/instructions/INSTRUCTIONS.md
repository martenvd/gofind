# INSTRUCTIONS

These instructions define how to implement changes in this repository safely and consistently.

## Scope and Style

- Keep changes minimal and targeted to the requested behavior.
- Use the existing package boundaries (`main`, `internal/core`, `internal/utils`).
- Avoid introducing new dependencies unless clearly justified.
- Keep CLI behavior stable unless explicitly requested to change it.

## Architecture Guidance

- Startup and flag parsing belong in `main.go`.
- Interactive UX belongs in `internal/core/finder.go` or `internal/core/prompt.go`.
- Filesystem and filtering helpers belong in `internal/utils/utils.go`.
- Cache persistence belongs in `internal/core/directories.go`.

## CLI Compatibility Requirements

- Preserve `gfdir.sh` integration: selected path must be plain stdout output.
- Do not print extra non-error text during successful path selection.
- If user-facing messages are needed, prefer stderr or strictly controlled output paths.

## Config and Cache Expectations

- Cache path: `~/.gofind/dirs.txt`
- Config path: `~/.gofind/config.json`
- Config schema currently expects:

```json
{
  "path": "/absolute/root/path"
}
```

- If config is missing or invalid, fall back gracefully.

## Flags and Arguments

- Supported update flags should remain aligned across the codebase.
- If changing flags, update:
  - `main.go` (`flag` parsing)
  - `internal/utils/utils.go` (`IsFlag`)
  - any docs mentioning usage

## Safety and Error Handling

- Never panic for ordinary user errors if a graceful message is possible.
- Return errors where possible instead of swallowing them.
- Keep terminal UI responsive and avoid blocking operations in event handlers.

## Verification Checklist

Before finalizing a change:

1. Run `go build ./...`
2. Validate `gofind` no-arg mode
3. Validate query mode (`gofind <query>`)
4. Validate cache refresh mode (`gofind -u`)
5. Confirm output still works with shell helper logic in `gfdir.sh`
