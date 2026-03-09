package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// topN returns the first n words from results (or all of them if fewer than n).
func topN(results []Match, n int) []string {
	out := make([]string, 0, n)
	for i, r := range results {
		if i >= n {
			break
		}
		out = append(out, r.Word)
	}
	return out
}

// assertTop1 asserts that the top result for query is want.
func assertTop1(t *testing.T, m *Matcher, query, want string) {
	t.Helper()
	results := m.Find(query)
	require.NotEmpty(t, results, "query %q: got no results, want %q as top", query, want)
	assert.Equal(t, want, results[0].Word,
		"query %q: unexpected top result (score %.3f)\n  full top-5: %v",
		query, results[0].Score, topN(results, 5))
}

// assertInTop asserts that want appears within the first n results for query.
func assertInTop(t *testing.T, m *Matcher, query, want string, n int) {
	t.Helper()
	results := m.Find(query)
	top := topN(results, n)
	assert.Contains(t, top, want, "query %q: %q not in top-%d; got %v", query, want, n, top)
}

// assertNotTop1 asserts that notWant is not the top result for query.
func assertNotTop1(t *testing.T, m *Matcher, query, notWant string) {
	t.Helper()
	results := m.Find(query)
	if len(results) > 0 {
		assert.NotEqual(t, notWant, results[0].Word,
			"query %q: top result is %q (score %.3f), which should NOT be first",
			query, results[0].Word, results[0].Score)
	}
}
