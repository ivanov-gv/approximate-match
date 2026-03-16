package approxmatch_test

import (
	"testing"

	approxmatch "github.com/ivanov-gv/approximate-match"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestScoreBounds(t *testing.T) {
	m := approxmatch.NewMatcher([]string{"podgorica", "belgrade", "novisad"})

	t.Run("ExactMatchScoresHigh", func(t *testing.T) {
		results := m.Find("podgorica")
		require.NotEmpty(t, results, "no results for exact query")
		assert.GreaterOrEqual(t, results[0].Score, 0.8,
			"exact match score too low: %.3f", results[0].Score)
	})

	t.Run("AllScoresInRange", func(t *testing.T) {
		for _, query := range []string{"podgorica", "belgrade", "xyz", "novi sad"} {
			for _, r := range m.Find(query) {
				assert.GreaterOrEqual(t, r.Score, 0.0,
					"query %q: score %.3f below 0 for %q", query, r.Score, r.Word)
				assert.LessOrEqual(t, r.Score, 1.0,
					"query %q: score %.3f above 1 for %q", query, r.Score, r.Word)
			}
		}
	})
}

func BenchmarkNewMatcher(b *testing.B) {
	words := []string{"podgorica", "belgrade", "novisad", "sarajevo", "tirana"}
	for i := 0; i < b.N; i++ {
		approxmatch.NewMatcher(words)
	}
}

func BenchmarkFind(b *testing.B) {
	m := approxmatch.NewMatcher([]string{"podgorica", "belgrade", "novisad", "sarajevo", "tirana"})
	queries := []string{"belgrade", "podgorica", "padgareeka", "sjutamare", "belo pole"}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.Find(queries[i%len(queries)])
	}
}
