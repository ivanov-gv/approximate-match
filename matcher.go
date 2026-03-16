package approxmatch

import "sort"

// skeletonMatchWeight slightly discounts skeleton-based matches vs. full-form
// matches, because a consonant skeleton discards vowel information.
const skeletonMatchWeight = 0.90

// DefaultScoreThreshold is the minimum score a candidate must reach to be
// included in Find results when no explicit threshold is provided to NewMatcher.
// It is set just above the typical score for common transliteration noise
// (e.g. "london" → 0.40 against Balkan station names), while staying below
// the score produced by single-character-deletion typos (e.g. "belgade" → 0.48).
const DefaultScoreThreshold = 0.45

// indexedWord holds all precomputed representations of one entry in the search list.
type indexedWord struct {
	original        string
	normalized      string
	skeleton        string
	normalizedStats map[rune]RuneStat
	skeletonStats   map[rune]RuneStat
}

// Match is a single result from Matcher.Find.
type Match struct {
	Word  string
	Score float64
}

// Matcher holds a fixed search list with precomputed statistics.
// Construct once with NewMatcher, then call Find for every user query.
type Matcher struct {
	words          []indexedWord
	scoreThreshold float64
}

// NewMatcher builds and returns a Matcher for the given fixed word list.
// All heavy preprocessing happens here so that each Find call is fast.
// threshold sets the minimum score for a result to be returned by Find;
// pass nil to use DefaultScoreThreshold.
func NewMatcher(words []string, threshold *float64) *Matcher {
	scoreThreshold := DefaultScoreThreshold
	if threshold != nil {
		scoreThreshold = *threshold
	}
	indexed := make([]indexedWord, len(words))
	for i, word := range words {
		normalized := Normalize(word)
		skeleton := ConsonantSkeleton(normalized)
		indexed[i] = indexedWord{
			original:        word,
			normalized:      normalized,
			skeleton:        skeleton,
			normalizedStats: buildRuneStats(normalized),
			skeletonStats:   buildRuneStats(skeleton),
		}
	}
	return &Matcher{words: indexed, scoreThreshold: scoreThreshold}
}

// Find returns all entries from the search list ranked by similarity to sample,
// best first. Entries with no commonality at all are omitted.
func (m *Matcher) Find(sample string) []Match {
	normalizedSample := Normalize(sample)
	skeletonSample := ConsonantSkeleton(normalizedSample)

	normalizedSampleStats := buildRuneStats(normalizedSample)
	skeletonSampleStats := buildRuneStats(skeletonSample)

	results := make([]Match, 0, len(m.words)/2)

	for _, entry := range m.words {
		normalizedScore := matchScore(normalizedSample, normalizedSampleStats, entry.normalized)
		skeletonScore := matchScore(skeletonSample, skeletonSampleStats, entry.skeleton) * skeletonMatchWeight

		score := normalizedScore
		if skeletonScore > score {
			score = skeletonScore
		}
		if score >= m.scoreThreshold {
			results = append(results, Match{Word: entry.original, Score: score})
		}
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Score > results[j].Score
	})
	return results
}

// matchScore returns a value in [0, 1] measuring how closely word matches
// sample (with precomputed sampleStats). 0 means no shared characters at all.
func matchScore(sample string, sampleStats map[rune]RuneStat, word string) float64 {
	if len(sample) == 0 || len(word) == 0 {
		return 0
	}

	// Start with frequency counts from sample; decrement as we process word.
	charFreqDelta := make(map[rune]int, len(sampleStats))
	for char, stat := range sampleStats {
		charFreqDelta[char] = stat.num
	}

	var longestCommonSubstr int
	longestCommonSubstrIsLeading := false
	for byteOffset, char := range word {
		charFreqDelta[char]--
		sampleStat, found := sampleStats[char]
		if !found {
			continue
		}
		if prefixLen := lenPrefix(word[byteOffset:], sampleStat.substrings...); prefixLen > longestCommonSubstr {
			longestCommonSubstr = prefixLen
			longestCommonSubstrIsLeading = byteOffset == 0
		}
	}

	if longestCommonSubstr == 0 {
		return 0
	}

	// Normalise LCS length against the longer of the two strings.
	longerLen := len(sample)
	if len(word) > longerLen {
		longerLen = len(word)
	}
	lcsRatio := float64(longestCommonSubstr) / float64(longerLen)

	// When the sample is entirely a leading prefix of the word (both strings start
	// with the same sequence and the query is fully consumed), normalise against the
	// sample length rather than the longer length. This prevents a shorter word that
	// merely contains the query as an interior substring from outscoring a longer word
	// that starts with the full query (e.g. "beograd" should prefer "beogradcentar"
	// over "novibeograd").
	if longestCommonSubstrIsLeading && longestCommonSubstr == len(sample) {
		leadingRatio := float64(longestCommonSubstr) / float64(len(sample))
		if leadingRatio > lcsRatio {
			lcsRatio = leadingRatio
		}
	}

	// Penalise by how many characters are unaccounted for (relative to total).
	totalUnmatchedChars := calcAbsDiffSum(charFreqDelta)
	unmatchedRatio := float64(totalUnmatchedChars) / float64(len(sample)+len(word))

	penalty := unmatchedRatio
	if penalty > 1 {
		penalty = 1
	}
	return lcsRatio * (1.0 - penalty)
}
