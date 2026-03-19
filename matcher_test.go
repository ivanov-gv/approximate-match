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
}
