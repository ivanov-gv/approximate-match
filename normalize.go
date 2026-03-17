package approxmatch

import (
	"strings"
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

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

	// Cyrillic: multi-char sequences first (must precede their sub-patterns).
	// Russian writes the Serbian lj/nj sounds as "ль"/"нь" (base + soft sign);
	// Serbian Cyrillic uses the precomposed ligatures љ/њ.
	// Both forms collapse to the base consonant so queries match either spelling.
	{"ль", "л"}, // Russian soft-L sequence → л (e.g. вальево ≈ ваљево)
	{"нь", "н"}, // Russian soft-N sequence → н (e.g. шушань ≈ шушањ)
	// Strip remaining Cyrillic soft/hard signs after multi-char rules.
	{"ь", ""},
	{"ъ", ""},
	// Russian "ю" (/ju/) corresponds to Serbian "у" after a soft consonant
	// (e.g. лютотук ≈ љутотук); "ы" corresponds to Serbian "и" (e.g. голубовцы ≈ голубовци).
	{"ю", "у"},
	{"ы", "и"},
	// Serbian Cyrillic ligatures and letters without NFD decomposition.
	{"љ", "л"}, // Serbian lj-ligature (e.g. љешница ≈ льешница after above rules)
	{"њ", "н"}, // Serbian nj-ligature
	{"ћ", "ч"}, // Serbian tshe → ч (Russian uses ч for both ч and ћ, e.g. братоношићи ≈ братоношичи)
	{"ђ", "д"}, // Serbian dje → д
	{"ј", "и"}, // Serbian j-letter → и (Russian substitutes и/й; й NFD-strips to и automatically)
}

// Normalize applies Unicode NFD decomposition to strip combining diacritical
// marks (covering š→s, č→c, ž→z, ć→c and many others), removes spaces,
// lowercases, then applies phonetic equivalence rules so that
// different-but-equivalent spellings converge to the same form.
func Normalize(s string) string {
	t := transform.Chain(norm.NFD, runes.Remove(diacriticOrSpaceSet{}), norm.NFC)
	s, _, _ = transform.String(t, s)
	s = strings.ToLower(s)
	for _, replacement := range normReplacements {
		s = strings.ReplaceAll(s, replacement[0], replacement[1])
	}
	return s
}

// ConsonantSkeleton removes all vowels from an already-normalized string.
// This lets the matcher handle severe vowel confusion,
// e.g. "padgareeka" → pdgrk ≈ pdgrc ← podgorica.
func ConsonantSkeleton(s string) string {
	var builder strings.Builder
	builder.Grow(len(s))
	for _, char := range s {
		switch char {
		case 'a', 'e', 'i', 'o', 'u',
			'а', 'е', 'и', 'о', 'у', 'я': // Cyrillic vowels (ю→у and ы→и already handled by Normalize)
			// drop vowel
		default:
			builder.WriteRune(char)
		}
	}
	return builder.String()
}
