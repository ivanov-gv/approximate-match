package approxmatch

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBuildRuneStats(t *testing.T) {
	t.Run("ASCIIRuneCount", func(t *testing.T) {
		_, runeCount := buildRuneStats("abc")
		assert.Equal(t, 3, runeCount)
	})

	t.Run("CyrillicRuneCount", func(t *testing.T) {
		// "бар" is 3 runes but 6 bytes; count must be rune-based.
		_, runeCount := buildRuneStats("бар")
		assert.Equal(t, 3, runeCount)
	})

	t.Run("Empty", func(t *testing.T) {
		stats, runeCount := buildRuneStats("")
		assert.Empty(t, stats)
		assert.Equal(t, 0, runeCount)
	})

	t.Run("RepeatedRune", func(t *testing.T) {
		stats, runeCount := buildRuneStats("aab")
		assert.Equal(t, 3, runeCount)
		require.Contains(t, stats, 'a')
		assert.Equal(t, 2, stats['a'].count)
		require.Contains(t, stats, 'b')
		assert.Equal(t, 1, stats['b'].count)
	})

	t.Run("SubstringContent", func(t *testing.T) {
		stats, _ := buildRuneStats("abc")
		// Substrings starting at 'a' should include the full string.
		require.Contains(t, stats, 'a')
		assert.Contains(t, stats['a'].substrings, "abc")
		require.Contains(t, stats, 'b')
		assert.Contains(t, stats['b'].substrings, "bc")
	})
}

func TestLenPrefix(t *testing.T) {
	t.Run("ExactMatch", func(t *testing.T) {
		assert.Equal(t, len("abc"), lenPrefix("abc", "abc"))
	})

	t.Run("PartialMatch", func(t *testing.T) {
		assert.Equal(t, 2, lenPrefix("abc", "abx"))
	})

	t.Run("NoMatch", func(t *testing.T) {
		assert.Equal(t, 0, lenPrefix("abc", "xyz"))
	})

	t.Run("EmptyCandidate", func(t *testing.T) {
		assert.Equal(t, 0, lenPrefix("abc", ""))
	})

	t.Run("EmptySample", func(t *testing.T) {
		assert.Equal(t, 0, lenPrefix("", "abc"))
	})

	t.Run("BestOfMultipleCandidates", func(t *testing.T) {
		// "abx" matches 2 runes, "abc" matches all 3.
		assert.Equal(t, len("abc"), lenPrefix("abc", "abx", "abc"))
	})

	t.Run("CyrillicExactMatch", func(t *testing.T) {
		// Result is byte length of the Cyrillic string.
		assert.Equal(t, len("бар"), lenPrefix("бар", "бар"))
	})

	t.Run("CyrillicPartialMatch", func(t *testing.T) {
		// "ба" is 4 bytes; "бар" vs "бах" share the first 2 runes.
		assert.Equal(t, len("ба"), lenPrefix("бар", "бах"))
	})
}
