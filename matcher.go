package main

import (
	"sort"
	"strings"
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

// RuneStat holds the occurrence count of a rune in a string and all substrings
// of that string starting at each occurrence position.
type RuneStat struct {
	num        int
	substrings []string
}

// lenPrefix returns the byte length of the longest common prefix between sample
// and any of the given candidate strings.
func lenPrefix(sample string, candidates ...string) int {
	prefixEnded := make(map[int]bool, len(candidates))
	for byteOffset, letter := range sample {
		matched := false
		for candidateIndex, candidate := range candidates {
			if !prefixEnded[candidateIndex] && byteOffset < len(candidate) && letter == rune(candidate[byteOffset]) {
				matched = true
			} else {
				prefixEnded[candidateIndex] = true
			}
		}
		if !matched {
			return byteOffset
		}
	}
	return len(sample)
}

// diacriticOrSpaceSet implements runes.Set to identify runes stripped during
// Unicode normalization: nonspacing combining marks (accents, carons, etc.)
// and ASCII spaces.
type diacriticOrSpaceSet struct{}

func (diacriticOrSpaceSet) Contains(r rune) bool {
	return unicode.Is(unicode.Mn, r) || r == ' '
}

// normReplacements is applied in order during phonetic normalization.
// Order matters: longer patterns must precede their sub-patterns.
//
// Diacritics that decompose under Unicode NFD (š→s, č→c, ž→z, ć→c, etc.) are
// handled automatically by the Unicode transform; only đ is listed here because
// it has no NFD decomposition.
var normReplacements = [][2]string{
	// Slavic ijekavica → ekavica (must precede any "je" or "ie" rules)
	{"ije", "e"},
	// Slavic digraphs
	{"lj", "l"},
	{"nj", "n"},
	{"dj", "d"},
	// đ has no Unicode NFD decomposition, handle explicitly
	{"đ", "d"},
	// Foreign multi-char phonetics (sch before sh, sh before s)
	{"sch", "s"},
	{"sh", "s"},
	{"ch", "c"},
	{"zh", "z"},
	{"ph", "f"},
	{"th", "t"},
	{"ck", "k"},
	{"ee", "i"},
	{"oo", "u"},
	{"ou", "u"},
	// Germanic w → v
	{"w", "v"},
	// Reduce double consonants (duplicates from transliteration / typos)
	{"bb", "b"}, {"cc", "c"}, {"dd", "d"}, {"ff", "f"}, {"gg", "g"},
	{"kk", "k"}, {"ll", "l"}, {"mm", "m"}, {"nn", "n"}, {"pp", "p"},
	{"rr", "r"}, {"ss", "s"}, {"tt", "t"}, {"vv", "v"}, {"zz", "z"},
}

// normalize applies Unicode NFD decomposition to strip combining diacritical
// marks (covering š→s, č→c, ž→z, ć→c and many others), removes spaces,
// lowercases, then applies phonetic equivalence rules so that
// different-but-equivalent spellings converge to the same form.
func normalize(s string) string {
	t := transform.Chain(norm.NFD, runes.Remove(diacriticOrSpaceSet{}), norm.NFC)
	s, _, _ = transform.String(t, s)
	s = strings.ToLower(s)
	for _, replacement := range normReplacements {
		s = strings.ReplaceAll(s, replacement[0], replacement[1])
	}
	return s
}

// consonantSkeleton removes all vowels from an already-normalized string.
// This lets the matcher handle severe vowel confusion,
// e.g. "padgareeka" → pdgrk ≈ pdgrc ← podgorica.
func consonantSkeleton(s string) string {
	var builder strings.Builder
	builder.Grow(len(s))
	for _, char := range s {
		switch char {
		case 'a', 'e', 'i', 'o', 'u':
			// drop vowel
		default:
			builder.WriteRune(char)
		}
	}
	return builder.String()
}

// buildRuneStats returns, for every character in s, its count and all
// substrings of s starting at that character's position. Used by the
// substring-prefix matching in matchScore.
func buildRuneStats(s string) map[rune]RuneStat {
	stats := make(map[rune]RuneStat, len(s))
	for byteOffset, char := range s {
		stat := stats[char]
		stat.num++
		stat.substrings = append(stat.substrings, s[byteOffset:])
		stats[char] = stat
	}
	return stats
}

// calcAbsDiffSum sums the absolute character-frequency deltas.
func calcAbsDiffSum(charFreqDelta map[rune]int) int {
	total := 0
	for _, delta := range charFreqDelta {
		if delta > 0 {
			total += delta
		} else {
			total -= delta
		}
	}
	return total
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
	for byteOffset, char := range word {
		charFreqDelta[char]--
		sampleStat, found := sampleStats[char]
		if !found {
			continue
		}
		if prefixLen := lenPrefix(word[byteOffset:], sampleStat.substrings...); prefixLen > longestCommonSubstr {
			longestCommonSubstr = prefixLen
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

	// Penalise by how many characters are unaccounted for (relative to total).
	totalUnmatchedChars := calcAbsDiffSum(charFreqDelta)
	unmatchedRatio := float64(totalUnmatchedChars) / float64(len(sample)+len(word))

	penalty := unmatchedRatio
	if penalty > 1 {
		penalty = 1
	}
	return lcsRatio * (1.0 - penalty)
}

// skeletonMatchWeight slightly discounts skeleton-based matches vs. full-form
// matches, because a consonant skeleton discards vowel information.
const skeletonMatchWeight = 0.90

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
	words []indexedWord
}

// NewMatcher builds and returns a Matcher for the given fixed word list.
// All heavy preprocessing happens here so that each Find call is fast.
func NewMatcher(words []string) *Matcher {
	indexed := make([]indexedWord, len(words))
	for i, word := range words {
		normalized := normalize(word)
		skeleton := consonantSkeleton(normalized)
		indexed[i] = indexedWord{
			original:        word,
			normalized:      normalized,
			skeleton:        skeleton,
			normalizedStats: buildRuneStats(normalized),
			skeletonStats:   buildRuneStats(skeleton),
		}
	}
	return &Matcher{words: indexed}
}

// Find returns all entries from the search list ranked by similarity to sample,
// best first. Entries with no commonality at all are omitted.
func (m *Matcher) Find(sample string) []Match {
	normalizedSample := normalize(sample)
	skeletonSample := consonantSkeleton(normalizedSample)

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
		if score > 0 {
			results = append(results, Match{Word: entry.original, Score: score})
		}
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Score > results[j].Score
	})
	return results
}
