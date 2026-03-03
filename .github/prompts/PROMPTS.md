# PROMPTS

Reusable prompts for working on `gofind`.

## 1) Add a Feature

```
Implement this feature in gofind: <feature description>.

Constraints:
- Keep changes minimal and consistent with existing package structure.
- Preserve gfdir.sh compatibility (selected path printed to stdout only).
- Update docs if CLI behavior changes.

Validation:
- Run go build ./...
- List manual test commands and expected behavior.
```

## 2) Fix a Bug

```
Investigate and fix this gofind bug: <bug description>.

Please:
- Identify root cause first.
- Implement the smallest safe fix.
- Avoid unrelated refactors.
- Include a concise explanation of why the fix works.

Validation:
- Run go build ./...
- Provide exact steps to reproduce before/after.
```

## 3) Improve UX in TUI Mode

```
Improve TUI behavior in internal/core/finder.go for this request: <ux change>.

Requirements:
- Keep existing key behavior unless explicitly changing it.
- Preserve Enter-to-output path behavior.
- Do not add heavy dependencies.

Validation:
- Run go build ./...
- Document keyboard interactions that changed.
```

## 4) Add/Adjust CLI Flags

```
Add or modify CLI flags for gofind: <flag change>.

Requirements:
- Keep main.go and utils.IsFlag aligned.
- Maintain backward compatibility unless otherwise specified.
- Update relevant markdown docs.

Validation:
- Run go build ./...
- Show sample command invocations and expected output.
```

## 5) Performance Pass (Directory Discovery)

```
Optimize path discovery performance in internal/utils/utils.go.

Requirements:
- Preserve .git-based repository detection semantics.
- Preserve cache behavior and output contract.
- Focus on measurable improvements and readability.

Validation:
- Run go build ./...
- Summarize complexity and practical impact.
```

## 6) Documentation Sync

```
Update repository docs to match current gofind behavior.

Include:
- Usage modes (no args, query arg, update flags)
- Cache/config locations and schema
- Install and shell integration notes

Keep docs concise and command-oriented.
```
