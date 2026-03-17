# Development Guidelines

## Codebase Overview

This is a fuzzy string matching library for station/location names, designed to handle linguistic variation: diacritics, phonetic equivalences, dialect variants (Serbian ekavica/ijekavica), Cyrillic transliterations, and minor typos. The public API is:

```go
matcher := approximatematch.NewMatcher(wordList)
results := matcher.Find(query)  // []Match, sorted by Score descending
```

### Package Layout

```
approximatematch/       (root package — library)
  matcher.go            Matcher struct, Find(), matchScore()
  normalize.go          Normalize(), ConsonantSkeleton()
  runestat.go           RuneStat, buildRuneStats(), calcAbsDiffSum()

cmd/
  main.go               main() only — demo entry point

test/
  stations.go           Test fixture: StationData struct + Stations slice
  matcher_test.go       Integration tests (all linguistic test groups)

matcher_test.go         Unit tests: score bounds, benchmarks (root package)
normalize_test.go       Unit tests: Normalize() and ConsonantSkeleton() cases
```

The `test/` directory is its own package (`package test`) to keep the large station fixture and integration tests separate from the library's unit tests.

### Matching Algorithm

`NewMatcher` preprocesses each word in the list:
1. Runs `Normalize()` — Unicode NFD, diacritic stripping, lowercase, phonetic substitutions (see below)
2. Runs `ConsonantSkeleton()` — strips vowels from the normalized form
3. Calls `buildRuneStats()` — maps each rune to its frequency and all substrings starting at that position

`Find` scores every candidate against the query using `matchScore()`, which:
- Computes the **longest common substring (LCS)** length between normalized forms
- Computes `calcAbsDiffSum()` — sum of absolute character-frequency deltas
- Combines them: `score = lcs/maxLen * (1 - delta/totalChars)`
- Runs the same computation on the **consonant skeletons** (weighted by `skeletonMatchWeight = 0.90`)
- Takes the max of the two scores

Results with score > 0 are returned sorted descending.

### Normalization Pipeline (`normalize.go`)

`Normalize()` applies these steps in order:
1. Unicode NFD decomposition + remove category-`Mn` nonspacing marks (handles most diacritics)
2. Remove spaces, lowercase everything
3. Explicit single-char replacements: `đ → d`
4. Multi-char phonetic substitutions (applied in order — longer patterns first):
   - Slavic digraphs: `ije → e`, `lj → l`, `nj → n`, `dj → d`
   - German/English clusters: `sch → s`, `sh → s`, `ch → c`, `ph → f`, `th → t`, `ck → k`
   - Vowel collapses: `ee → i`, `oo → u`, `ou → u`
   - Double consonant reduction: `bb → b`, `cc → c`, … `zz → z`

`ConsonantSkeleton()` calls `Normalize()` then strips all vowels (`a e i o u`).

---

## Variable and Function Naming

- Use full, descriptive names. Single-letter and abbreviated names (`ss`, `ns`, `wi`, `iw`, `sk`, `b`) are not acceptable.
- Name variables after what they represent, not their type (`charFreqDelta`, not `m`; `longestCommonSubstr`, not `max`).

## Code Organization

- Do not use separator comments (e.g. `// ---- helpers ----`) to divide sections of a file. Create separate files instead.
- Do not keep dead code or compatibility shims. Delete unused functions, types, and variables outright.
- `cmd/main.go` must contain only the `main` function. All other declarations belong in their own files.

## Testing

- Use `github.com/stretchr/testify` for all assertions:
  - `assert` — non-fatal checks; the test continues and reports all failures at once.
  - `require` — fatal checks where continuing would panic or produce meaningless results (e.g. checking `results[0]`
    after verifying `len(results) > 0`), or where the failure represents an invariant that must never be violated (e.g.
    `require.Empty` for false-positive checks).
- Use `t.Run()` for every logical group of sub-cases (exact matches, typos, phonetic variants, etc.).
- Use table-driven tests (`[]struct{ ... }` + a loop) instead of repeating the same assertion call with different
  arguments. This applies even when the struct has extra fields like `mustNotNames []string` alongside
  `wantName string`.
- When sub-tests share identical loop logic but have different inputs, define a local `run` helper *inside* the parent
  test function (not at package level) and pass the cases slice to it:

```go
func TestDisambiguation(t *testing.T) {
matcher := ...

run := func (t *testing.T, cases []struct {
query        string
wantName     string
mustNotNames []string
}) {
t.Helper()
for _, tc := range cases {
t.Run(tc.query, func (t *testing.T) { ... })
}
}

t.Run("BelgradeFamily", func (t *testing.T) {
run(t, []struct{ ... }{
{"belgrade", "Beograd Centar", []string{"Novi Sad"}},
...
})
})
}
```

- Do not add `t.Skip` to paper over data gaps.
- Integration tests live in `test/matcher_test.go`; unit tests for each source file live alongside it in the root package.

### Test Data and `github.com/samber/lo`

Use `github.com/samber/lo` for all common data reshaping — it is far more readable than manual loops:

```go
lo.Keys(m) // map → key slice (replaces: for k := range m { keys = append(keys, k) })
lo.Filter(s, fn) // keep elements matching a predicate
lo.Compact(s) // remove zero values (empty strings, 0, nil, …)
lo.Uniq(s)    // deduplicate a slice
```

Shared test fixtures (lookup maps, derived slices) that are used across multiple tests should be **package-level
variables initialized with an IIFE** (immediately-invoked function expression), not helper functions:

```go
// Good — data lives in a var; the lambda that builds it has no name and can't be called elsewhere
var nameToStationName = func () map[string]string {
result := make(map[string]string)
for _, station := range integration.Stations {
allNames := lo.Compact(append([]string{station.Name, station.NameEn, station.NameCyr}, station.ProductionAliases...))
for _, name := range allNames {
result[name] = station.Name
}
}
return result
}()
```

### Break complex expressions into named local variables

Avoid chaining operations or nesting calls when the result has a meaningful name. Each intermediate variable should
express one clear thought; the reader should be able to understand the code by reading it like a sentence.

```go
// Bad — one dense line that requires parsing from inside out
for _, station := range lo.Filter(integration.Stations, func (s integration.StationData, _ int) bool { return !s.Blacklisted }) {
for _, name := range lo.Uniq(lo.Compact([]string{station.Name, station.NameEn, station.NameCyr})) {

// Good — each step has a name that explains what it is
nonBlacklistedStations := lo.Filter(integration.Stations, func (s integration.StationData, _ int) bool {
return !s.Blacklisted
})
for _, station := range nonBlacklistedStations {
officialNames := lo.Uniq(lo.Compact([]string{station.Name, station.NameEn, station.NameCyr}))
for _, name := range officialNames {
```

This applies everywhere, not just in tests: if a sub-expression has a name, give it one.

Do **not** wrap simple, one-liner constructor calls in named helper functions. Construct objects inline where they are
used:

```go
// Good
matcher := approxmatch.NewMatcher(lo.Keys(nameToStationName), nil)

// Bad — unnecessary indirection that obscures what is actually constructed
func newUnifiedMatcher() (*approxmatch.Matcher, map[string]string) { ... }
matcher, stationToName := newUnifiedMatcher()
```

If a custom threshold is needed in one test, pass it inline with a local variable:

```go
strictThreshold := 0.6
matcher := approxmatch.NewMatcher(lo.Keys(officialNameToStationName), &strictThreshold)
```

## Diacritics and Unicode Normalization

- Handle diacritics via Unicode NFD decomposition + removal of nonspacing marks (category `Mn`) using `golang.org/x/text`. This covers hundreds of characters automatically.
- Only add explicit replacement rules for characters that do not decompose under NFD (e.g. `đ`).

## Development Workflow

```
make lint    # golangci-lint on all packages
make test    # go test ./...
make build   # builds binary to build/matcher
make run     # build + execute
```

Linting is configured in `.golangci.yml`. Enabled linter groups: correctness (`errcheck`, `govet`, `staticcheck`, `ineffassign`, `bodyclose`, `noctx`, `errorlint`), style (`revive`, `gocritic`, `misspell`, `gosec`, `exhaustive`), and formatting (`goimports`). The `exported` comment rule from `revive` and the `ifElseChain` rule from `gocritic` are disabled.

## Dependencies

- `golang.org/x/text` — Unicode normalization (`norm`, `runes`, `transform`).
- `github.com/stretchr/testify` — test assertions (`assert`, `require`).
- `github.com/samber/lo` — generic slice/map utilities (`lo.Keys`, `lo.Filter`, `lo.Without`, `lo.Uniq`, …).
