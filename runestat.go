package approxmatch

import "unicode/utf8"

// RuneStat holds the occurrence count of a rune in a string and all substrings
// of that string starting at each occurrence position.
type RuneStat struct {
	num        int
	substrings []string
}

// lenPrefix returns the byte length of the longest common prefix between sample
// and any of the given candidate strings.
// Each candidate tracks its own byte offset so multi-byte UTF-8 characters
// (e.g. Cyrillic, 2 bytes each) are compared correctly by rune, not by byte.
func lenPrefix(sample string, candidates ...string) int {
	candidateByteOffsets := make([]int, len(candidates))
	candidateEnded := make([]bool, len(candidates))

	for sampleByteOffset, sampleRune := range sample {
		matched := false
		for i, candidate := range candidates {
			if candidateEnded[i] {
				continue
			}
			if candidateByteOffsets[i] >= len(candidate) {
				candidateEnded[i] = true
				continue
			}
			candidateRune, runeSize := utf8.DecodeRuneInString(candidate[candidateByteOffsets[i]:])
			if sampleRune == candidateRune {
				matched = true
				candidateByteOffsets[i] += runeSize
			} else {
				candidateEnded[i] = true
			}
		}
		if !matched {
			return sampleByteOffset
		}
	}
	return len(sample)
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
