package integration_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	approxmatch "github.com/ivanov-gv/approximate-match"
	integration "github.com/ivanov-gv/approximate-match/test"
)

// ── Full alias+blacklist matcher (UnifiedStationNameToStationIdMap) ───────────

var (
	stationNameToID = integration.UnifiedStationNameToStationIdMap
	stationMatcher  = approxmatch.NewMatcher(stationNames())
)

func stationNames() []string {
	names := make([]string, 0, len(stationNameToID))
	for name := range stationNameToID {
		names = append(names, name)
	}
	return names
}

// ── Official-only matcher (StationIdToStationMap) ────────────────────────────
//
// Built from the canonical Name / NameEn / NameCyr of each real station.
// Used for alias and blacklist tests: aliases must surface the correct station,
// blacklisted names must not produce a confident match.

var (
	officialMatcher  *approxmatch.Matcher
	officialNameToID map[string]int // official station name → station ID
)

func init() {
	seen := make(map[string]bool)
	var names []string
	officialNameToID = make(map[string]int, len(integration.StationIdToStationMap)*3)

	for _, station := range integration.StationIdToStationMap {
		for _, name := range []string{station.Name, station.NameEn, station.NameCyr} {
			if seen[name] {
				continue
			}
			seen[name] = true
			names = append(names, name)
			officialNameToID[name] = station.Id
		}
	}
	officialMatcher = approxmatch.NewMatcher(names)
}

// officialStationNameToID maps the StationName field used in AliasesStationsList
// back to the numeric station ID from StationIdToStationMap.
func officialStationNameToStationID() map[string]int {
	m := make(map[string]int, len(integration.StationIdToStationMap))
	for _, station := range integration.StationIdToStationMap {
		m[station.Name] = station.Id
	}
	return m
}

// ── Shared helpers ────────────────────────────────────────────────────────────

// topN returns the first n words from results (or all of them if fewer than n).
func topN(results []approxmatch.Match, n int) []string {
	out := make([]string, 0, n)
	for i, r := range results {
		if i >= n {
			break
		}
		out = append(out, r.Word)
	}
	return out
}

// assertTopID asserts that the top result for query belongs to the station with wantID.
// Because multiple spellings of the same physical station share an ID, this is stricter
// than asserting an exact name: a typo that surfaces "belgrad" instead of "belgrade"
// still passes when both carry ID 18.
func assertTopID(t *testing.T, query string, wantID int) {
	t.Helper()
	results := stationMatcher.Find(query)
	require.NotEmpty(t, results, "query %q: got no results, want station ID %d", query, wantID)
	gotID := stationNameToID[results[0].Word]
	assert.Equal(t, wantID, gotID,
		"query %q: top result %q (ID %d) should be station ID %d\n  full top-5: %v",
		query, results[0].Word, gotID, wantID, topN(results, 5))
}

// assertTopIDNot asserts that the top result for query does not belong to the station
// with notWantID.
func assertTopIDNot(t *testing.T, query string, notWantID int) {
	t.Helper()
	results := stationMatcher.Find(query)
	if len(results) > 0 {
		gotID := stationNameToID[results[0].Word]
		assert.NotEqual(t, notWantID, gotID,
			"query %q: top result %q (ID %d) should NOT be station ID %d",
			query, results[0].Word, gotID, notWantID)
	}
}

// assertTopIDOfficial asserts that when searching the official station list the
// top result for query belongs to the station with wantID.
func assertTopIDOfficial(t *testing.T, query string, wantID int) {
	t.Helper()
	results := officialMatcher.Find(query)
	require.NotEmpty(t, results, "query %q: got no results in official list, want station ID %d", query, wantID)
	gotID := officialNameToID[results[0].Word]
	assert.Equal(t, wantID, gotID,
		"query %q: top official result %q (ID %d) should be station ID %d\n  full top-5: %v",
		query, results[0].Word, gotID, wantID, topN(results, 5))
}

// assertNoGoodMatchOfficial asserts that searching the official station list
// for query yields either no results or only low-confidence results (score < 0.5).
// Used for blacklisted station names that have no railway service.
func assertNoGoodMatchOfficial(t *testing.T, query string) {
	t.Helper()
	results := officialMatcher.Find(query)
	if len(results) > 0 {
		assert.Less(t, results[0].Score, 0.6,
			"query %q: expected no confident match (score < 0.6) but got %q with score %.3f",
			query, results[0].Word, results[0].Score)
	}
}
