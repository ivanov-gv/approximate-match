package integration_test

import (
	approxmatch "github.com/ivanov-gv/approximate-match"
	integration "github.com/ivanov-gv/approximate-match/test"
)

// newUnifiedMatcher builds a Matcher from UnifiedStationNameToStationIdMap using
// the default score threshold, and returns it together with the map so callers
// can look up station IDs by name.
func newUnifiedMatcher() (*approxmatch.Matcher, map[string]int) {
	stationNameToID := integration.UnifiedStationNameToStationIdMap
	names := make([]string, 0, len(stationNameToID))
	for name := range stationNameToID {
		names = append(names, name)
	}
	return approxmatch.NewMatcher(names, nil), stationNameToID
}

// newUnifiedMatcherWithThreshold is like newUnifiedMatcher but applies a custom
// score threshold, useful for demonstrating threshold configuration in tests.
func newUnifiedMatcherWithThreshold(threshold float64) (*approxmatch.Matcher, map[string]int) {
	stationNameToID := integration.UnifiedStationNameToStationIdMap
	names := make([]string, 0, len(stationNameToID))
	for name := range stationNameToID {
		names = append(names, name)
	}
	return approxmatch.NewMatcher(names, &threshold), stationNameToID
}

// newOfficialMatcher builds a Matcher from the canonical Name / NameEn / NameCyr
// fields of StationIdToStationMap using the default score threshold, and returns
// it together with a map from each name to its station ID.
func newOfficialMatcher() (*approxmatch.Matcher, map[string]int) {
	return newOfficialMatcherWithThreshold(approxmatch.DefaultScoreThreshold)
}

// newOfficialMatcherWithThreshold is like newOfficialMatcher but applies a custom
// score threshold, useful for demonstrating threshold configuration in tests.
func newOfficialMatcherWithThreshold(threshold float64) (*approxmatch.Matcher, map[string]int) {
	seen := make(map[string]bool)
	var names []string
	officialNameToID := make(map[string]int, len(integration.StationIdToStationMap)*3)
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
	return approxmatch.NewMatcher(names, &threshold), officialNameToID
}

// officialStationNameToStationID maps the StationName field used in
// AliasesStationsList back to the numeric station ID from StationIdToStationMap.
func officialStationNameToStationID() map[string]int {
	m := make(map[string]int, len(integration.StationIdToStationMap))
	for _, station := range integration.StationIdToStationMap {
		m[station.Name] = station.Id
	}
	return m
}
