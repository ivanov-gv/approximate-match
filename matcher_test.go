package approxmatch_test

import (
	"testing"

	approxmatch "github.com/ivanov-gv/approximate-match"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestScoreBounds(t *testing.T) {
	matcher := approxmatch.NewMatcher([]string{"podgorica", "belgrade", "novisad"}, nil)

	t.Run("ExactMatchScoresHigh", func(t *testing.T) {
		results := matcher.Find("podgorica")
		require.NotEmpty(t, results, "no results for exact query")
		assert.Greater(t, results[0].Score, approxmatch.DefaultScoreThreshold,
			"exact match score too low: %.3f", results[0].Score)
	})

	t.Run("AllScoresInRange", func(t *testing.T) {
		for _, query := range []string{"podgorica", "belgrade", "xyz", "novi sad"} {
			t.Run(query, func(t *testing.T) {
				for _, result := range matcher.Find(query) {
					assert.GreaterOrEqual(t, result.Score, approxmatch.DefaultScoreThreshold,
						"score %.3f below threshold for %q", result.Score, result.Word)
					assert.LessOrEqual(t, result.Score, 1.0,
						"score %.3f above 1 for %q", result.Score, result.Word)
				}
			})
		}
	})

	t.Run("IdenticalStringScoresOne", func(t *testing.T) {
		exactMatcher := approxmatch.NewMatcher([]string{"podgorica"}, nil)
		results := exactMatcher.Find("podgorica")
		require.Len(t, results, 1)
		assert.InDelta(t, 1.0, results[0].Score, 0.001, "identical strings must score ~1.0")
	})

	t.Run("CyrillicScoreConsistentWithLatin", func(t *testing.T) {
		// Rune count must be used (not byte count) so that equivalent Latin and
		// Cyrillic structures produce equivalent scores.
		latinMatcher := approxmatch.NewMatcher([]string{"bar"}, nil)
		cyrillicMatcher := approxmatch.NewMatcher([]string{"бар"}, nil)
		latinResults := latinMatcher.Find("bra")       // anagram of "bar"
		cyrillicResults := cyrillicMatcher.Find("бра") // anagram of "бар"
		require.Len(t, latinResults, 1)
		require.Len(t, cyrillicResults, 1)
		assert.InDelta(t, latinResults[0].Score, cyrillicResults[0].Score, 0.01,
			"equivalent Latin/Cyrillic structures should score equally")
	})
}

func TestEdgeCases(t *testing.T) {
	t.Run("EmptyQuery", func(t *testing.T) {
		matcher := approxmatch.NewMatcher([]string{"podgorica"}, nil)
		assert.Empty(t, matcher.Find(""), "empty query must return no results")
	})

	t.Run("EmptyWordList", func(t *testing.T) {
		matcher := approxmatch.NewMatcher([]string{}, nil)
		assert.Empty(t, matcher.Find("podgorica"), "empty word list must return no results")
	})
}

func TestLeadingPrefixBeatsInterior(t *testing.T) {
	// "beograd" is a leading prefix of "beogradcentar" but an interior match in
	// "novibeograd". The leading-prefix bonus must surface "beogradcentar" first.
	matcher := approxmatch.NewMatcher([]string{"beogradcentar", "novibeograd"}, nil)
	results := matcher.Find("beograd")
	require.NotEmpty(t, results, "no results for 'beograd'")
	assert.Equal(t, "beogradcentar", results[0].Word,
		"leading prefix must outscore interior match")
}
