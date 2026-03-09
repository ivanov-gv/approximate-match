# Development Guidelines

## Variable and Function Naming

- Use full, descriptive names. Single-letter and abbreviated names (`ss`, `ns`, `wi`, `iw`, `sk`, `b`) are not acceptable.
- Name variables after what they represent, not their type (`charFreqDelta`, not `m`; `longestCommonSubstr`, not `max`).

## Code Organization

- Do not use separator comments (e.g. `// ---- helpers ----`) to divide sections of a file. Create separate files instead.
- Do not keep dead code or compatibility shims. Delete unused functions, types, and variables outright.
- `main.go` must contain only the `main` function. All other declarations belong in their own files.

## Testing

- Use `github.com/stretchr/testify` for all assertions:
  - `assert` — non-fatal checks; the test continues and reports all failures at once.
  - `require` — fatal checks where continuing would panic or produce meaningless results (e.g. checking `results[0]` after verifying `len(results) > 0`).
- Use `t.Run()` for every logical group of sub-cases (exact matches, typos, phonetic variants, etc.).
- Use table-driven tests (`[]struct{ ... }` + a loop) instead of repeating the same assertion call with different arguments.
- Keep test helpers in `helpers_test.go`, not mixed into the main test file.

## Diacritics and Unicode Normalization

- Handle diacritics via Unicode NFD decomposition + removal of nonspacing marks (category `Mn`) using `golang.org/x/text`. This covers hundreds of characters automatically.
- Only add explicit replacement rules for characters that do not decompose under NFD (e.g. `đ`).

## Dependencies

- `golang.org/x/text` — Unicode normalization (`norm`, `runes`, `transform`).
- `github.com/stretchr/testify` — test assertions (`assert`, `require`).
