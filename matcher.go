package main

import (
	"sort"
	"strings"
)

// normReplacements is applied in order during phonetic normalization.
// Order matters: longer patterns before their sub-patterns.
var normReplacements = [][2]string{
	// Slavic ijekavica → ekavica (must precede any "je" or "ie" rules)
	{"ije", "e"},
	// Slavic digraphs
	{"lj", "l"},
	{"nj", "n"},
	{"dj", "d"},
	// Diacritics
	{"đ", "d"},
	{"š", "s"},
	{"č", "c"},
	{"ž", "z"},
	{"ć", "c"},
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

// normalize strips spaces, lowercases, and applies phonetic rules so that
// different-but-equivalent spellings converge to the same form.
func normalize(s string) string {
	s = strings.ToLower(s)
	s = strings.ReplaceAll(s, " ", "")
	for _, r := range normReplacements {
		s = strings.ReplaceAll(s, r[0], r[1])
	}
	return s
}

// consonantSkeleton removes all vowels from a (already normalized) string.
// This lets the matcher handle severe vowel confusion, e.g. "padgareeka" → pdgrk ≈ pdgrc ← podgorica.
func consonantSkeleton(s string) string {
	var b strings.Builder
	b.Grow(len(s))
	for _, r := range s {
		switch r {
		case 'a', 'e', 'i', 'o', 'u':
			// drop vowel
		default:
			b.WriteRune(r)
		}
	}
	return b.String()
}

// buildRuneStats returns, for every character in s, its count and all
// substrings of s starting at that character's position.  These are used by
// the substring-prefix matching in matchScore.
func buildRuneStats(s string) map[rune]RuneStat {
	stats := make(map[rune]RuneStat, len(s))
	for i, ch := range s {
		st := stats[ch]
		st.num++
		st.substrings = append(st.substrings, s[i:])
		stats[ch] = st
	}
	return stats
}

// calcDiff sums the absolute character-frequency deltas.
func calcDiff(wordDiff map[rune]int) int {
	d := 0
	for _, v := range wordDiff {
		if v > 0 {
			d += v
		} else {
			d -= v
		}
	}
	return d
}

// matchScore returns a value in [0, 1] measuring how closely word matches
// sample (with precomputed sampleStats).  0 means no shared characters at all.
func matchScore(sample string, sampleStats map[rune]RuneStat, word string) float64 {
	if len(sample) == 0 || len(word) == 0 {
		return 0
	}

	// Start with frequency counts from sample; decrement as we process word.
	wordDiff := make(map[rune]int, len(sampleStats))
	for k, v := range sampleStats {
		wordDiff[k] = v.num
	}

	var maxSubstrLen int
	for i, ch := range word {
		wordDiff[ch]--
		stat, found := sampleStats[ch]
		if !found {
			continue
		}
		if n := lenPrefix(word[i:], stat.substrings...); n > maxSubstrLen {
			maxSubstrLen = n
		}
	}

	if maxSubstrLen == 0 {
		return 0
	}

	// Normalise LCS length against the longer of the two strings.
	maxLen := len(sample)
	if len(word) > maxLen {
		maxLen = len(word)
	}
	lcsRatio := float64(maxSubstrLen) / float64(maxLen)

	// Penalise by how many characters are unaccounted for (relative to total).
	totalDiff := calcDiff(wordDiff)
	diffRatio := float64(totalDiff) / float64(len(sample)+len(word))

	penalty := diffRatio
	if penalty > 1 {
		penalty = 1
	}
	return lcsRatio * (1.0 - penalty)
}

// skelWeight slightly discounts skeleton-based matches vs. full-form matches,
// because a consonant skeleton discards information.
const skelWeight = 0.90

// indexedWord holds all precomputed representations of one entry in the search list.
type indexedWord struct {
	original  string
	norm      string
	skel      string
	normStats map[rune]RuneStat
	skelStats map[rune]RuneStat
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
	iw := make([]indexedWord, len(words))
	for i, w := range words {
		n := normalize(w)
		sk := consonantSkeleton(n)
		iw[i] = indexedWord{
			original:  w,
			norm:      n,
			skel:      sk,
			normStats: buildRuneStats(n),
			skelStats: buildRuneStats(sk),
		}
	}
	return &Matcher{words: iw}
}

// Find returns all entries from the search list ranked by similarity to sample,
// best first.  Entries with no commonality at all are omitted.
func (m *Matcher) Find(sample string) []Match {
	normSample := normalize(sample)
	skelSample := consonantSkeleton(normSample)

	normSampleStats := buildRuneStats(normSample)
	skelSampleStats := buildRuneStats(skelSample)

	results := make([]Match, 0, len(m.words)/2)

	for _, iw := range m.words {
		ns := matchScore(normSample, normSampleStats, iw.norm)
		ss := matchScore(skelSample, skelSampleStats, iw.skel) * skelWeight

		score := ns
		if ss > score {
			score = ss
		}
		if score > 0 {
			results = append(results, Match{Word: iw.original, Score: score})
		}
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Score > results[j].Score
	})
	return results
}
