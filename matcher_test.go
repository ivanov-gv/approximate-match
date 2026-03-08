package main

import (
	"testing"
)

// stationList is the fixed railway/location list used in all tests.
// Keys are lowercase, space-free station identifiers; values are station IDs
// (kept here as comments to show intent).
var stationList = []string{
	"aerodrom",         // 26
	"albania",          // -36
	"bar",              // 1
	"baresumanovica",   // 51
	"becici",           // -28
	"belgrad",          // 18
	"belgrade",         // 18
	"beograd",          // 18
	"beogradcentar",    // 18
	"bijelopolje",      // 7
	"bioce",            // 28
	"bosniaandherzegovina", // -42
	"bratonozici",      // 29
	"budva",            // -10
	"cetinje",          // -16
	"crmnica",          // 21
	"dabovici",         // 54
	"danilovgrad",      // 49
	"durmitor",         // -20
	"golubovci",        // 3
	"hercegnovi",       // -30
	"indjija",          // -2
	"kolasin",          // 5
	"kos",              // 34
	"kosjeric",         // 13
	"kotor",            // -14
	"krusevackipotok",  // 31
	"krusevo",          // 44
	"lajkovac",         // 15
	"lazarevac",        // 16
	"ljesnica",         // 45
	"ljutotuk",         // 48
	"lutovo",           // 30
	"matesevo",         // 35
	"mijatovokolo",     // 41
	"mojkovac",         // 6
	"moraca",           // 25
	"niksic",           // 56
	"novapazova",       // -6
	"novisad",          // 0
	"oblutak",          // 37
	"ostrog",           // 53
	"padez",            // 36
	"perast",           // -18
	"petrovac",         // -22
	"podgorica",        // 4
	"pozega",           // 12
	"priboj",           // 10
	"pricelje",         // 46
	"prijepolje",       // 9
	"prijepoljeteretna", // 8
	"rakovica",         // 17
	"ravnarijeka",      // 43
	"sarajevo",         // -44
	"savnik",           // -32
	"seliste",          // 33
	"shkoder",          // -40
	"slap",             // 50
	"slijepacmost",     // 42
	"sobajici",         // 52
	"spuz",             // 47
	"starapazova",      // -4
	"stitarickarijeka", // 39
	"stubica",          // 55
	"subotica",         // -8
	"susanj",           // 20
	"sutomore",         // 2
	"svetistefan",      // -26
	"tirana",           // -38
	"tivat",            // -12
	"trebaljevo",       // 38
	"trebesica",        // 32
	"ulcinj",           // -24
	"uzice",            // 11
	"valjevo",          // 14
	"virpazar",         // 22
	"vranjina",         // 23
	"zabljak",          // -34
	"zari",             // 40
	"zemun",            // 19
	"zeta",             // 24
	"zlatica",          // 27
}

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

// contains returns true if word appears anywhere in results.
func contains(results []Match, word string) bool {
	for _, r := range results {
		if r.Word == word {
			return true
		}
	}
	return false
}

// ---- helpers for test assertions -----------------------------------------------

func assertTop1(t *testing.T, m *Matcher, query, want string) {
	t.Helper()
	results := m.Find(query)
	if len(results) == 0 {
		t.Errorf("query %q: got no results, want %q as top", query, want)
		return
	}
	if results[0].Word != want {
		t.Errorf("query %q: top result = %q (score %.3f), want %q\n  full top-5: %v",
			query, results[0].Word, results[0].Score, want, topN(results, 5))
	}
}

func assertTopN(t *testing.T, m *Matcher, query string, want []string) {
	t.Helper()
	results := m.Find(query)
	top := topN(results, len(want))
	for i, w := range want {
		if i >= len(top) || top[i] != w {
			t.Errorf("query %q: top-%d = %v, want %v", query, len(want), top, want)
			return
		}
	}
}

func assertInTop(t *testing.T, m *Matcher, query, want string, n int) {
	t.Helper()
	results := m.Find(query)
	top := topN(results, n)
	for _, w := range top {
		if w == want {
			return
		}
	}
	t.Errorf("query %q: %q not in top-%d; got %v", query, want, n, top)
}

func assertNotTop1(t *testing.T, m *Matcher, query, notWant string) {
	t.Helper()
	results := m.Find(query)
	if len(results) > 0 && results[0].Word == notWant {
		t.Errorf("query %q: top result is %q (score %.3f), which should NOT be first",
			query, notWant, results[0].Score)
	}
}

func assertNoResults(t *testing.T, m *Matcher, query string) {
	t.Helper()
	results := m.Find(query)
	if len(results) > 0 {
		t.Errorf("query %q: expected no results, got %v", query, topN(results, 3))
	}
}

// ---- tests ---------------------------------------------------------------------

func TestPositiveCases(t *testing.T) {
	m := NewMatcher(stationList)

	// Exact matches (case-insensitive)
	assertTop1(t, m, "podgorica", "podgorica")
	assertTop1(t, m, "Podgorica", "podgorica")
	assertTop1(t, m, "PODGORICA", "podgorica")
	assertTop1(t, m, "bar", "bar")
	assertTop1(t, m, "kotor", "kotor")
	assertTop1(t, m, "budva", "budva")
	assertTop1(t, m, "tivat", "tivat")
	assertTop1(t, m, "tirana", "tirana")
	assertTop1(t, m, "novisad", "novisad")
	assertTop1(t, m, "niksic", "niksic")
	assertTop1(t, m, "sarajevo", "sarajevo")
	assertTop1(t, m, "subotica", "subotica")

	// Station names with spaces (user input typically includes spaces)
	assertTop1(t, m, "novi sad", "novisad")
	assertTop1(t, m, "nova pazova", "novapazova")
	assertTop1(t, m, "stara pazova", "starapazova")
	assertTop1(t, m, "bijelo polje", "bijelopolje")
	assertTop1(t, m, "herceg novi", "hercegnovi")
	assertTop1(t, m, "beograd centar", "beogradcentar")

	// Minor typos (1-2 character errors) — any belgrade-family key is correct
	// since belgrad/belgrade/beograd all map to the same station ID.
	assertInTop(t, m, "belgarde", "belgrade", 2)  // transposition
	assertInTop(t, m, "belgade", "belgrade", 3)   // missing r
	assertInTop(t, m, "belgrate", "belgrade", 2)  // transposition r/t
	assertTop1(t, m, "kolasin", "kolasin")
	assertTop1(t, m, "podgorcia", "podgorica") // transposition c/i
	assertTop1(t, m, "sutomore", "sutomore")
	assertTop1(t, m, "sutmore", "sutomore")    // missing o
	assertTop1(t, m, "mojkovac", "mojkovac")
	assertTop1(t, m, "bijelopolje", "bijelopolje")
}

func TestPhoneticCases(t *testing.T) {
	m := NewMatcher(stationList)

	// Vowel-heavy foreign mispronunciations — consonant skeleton rescues these.
	// "pod-go-REE-ka" heard as "padgareeka":
	assertTop1(t, m, "padgareeka", "podgorica")
	// Same with doubled vowels spelled out differently:
	assertTop1(t, m, "podgoriika", "podgorica")
	// "oo" → u normalisation:
	assertTop1(t, m, "podgoorica", "podgorica")

	// Mishearing of Sutomore as "sjutamare" (vowel shifts + spurious j):
	assertTop1(t, m, "sjutamare", "sutomore")
	// Similar: "sutomare" (o→a vowel confusion):
	assertTop1(t, m, "sutomare", "sutomore")

	// Ekavica / ijekavica equivalence:
	// "Bijelo Polje" (ijekavica) ↔ "Belo Pole" (ekavica / simplified)
	assertTop1(t, m, "belo pole", "bijelopolje")
	assertTop1(t, m, "belo polje", "bijelopolje")
	assertTop1(t, m, "bijelo polje", "bijelopolje")

	// Diacritic-free spellings:
	assertTop1(t, m, "niksic", "niksic")   // already ASCII in list
	assertTop1(t, m, "priboj", "priboj")

	// English-like transliterations:
	assertTop1(t, m, "belgrade", "belgrade")
	assertTop1(t, m, "Belgrade", "belgrade")
	assertTop1(t, m, "novi sad", "novisad")
	assertTop1(t, m, "Novi Sad", "novisad")

	// sh/ch phonetics:
	assertTop1(t, m, "shkoder", "shkoder")  // already in list with sh

	// Russian-style "ts" → c overlap — podgoritsa:
	// "ts" isn't a special rule but consonant skeleton helps: pdgrtsk ≈ pdgrc
	assertInTop(t, m, "podgoritsa", "podgorica", 3)
}

func TestDisambiguation(t *testing.T) {
	m := NewMatcher(stationList)

	// "belgrade" / "beograd" must point to beograd-family, NOT novisad.
	assertNotTop1(t, m, "belgrade", "novisad")
	assertNotTop1(t, m, "beograd", "novisad")
	assertInTop(t, m, "belgrade", "belgrade", 1)
	assertInTop(t, m, "beograd centar", "beogradcentar", 1)

	// Novi Sad vs Nova Pazova vs Stara Pazova.
	// "novi sad" → novisad wins over novapazova.
	assertTop1(t, m, "novi sad", "novisad")
	assertNotTop1(t, m, "novi sad", "novapazova")
	assertNotTop1(t, m, "novi sad", "starapazova")

	// "nova pazova" → novapazova wins over starapazova and novisad.
	assertTop1(t, m, "nova pazova", "novapazova")
	assertNotTop1(t, m, "nova pazova", "starapazova")
	assertNotTop1(t, m, "nova pazova", "novisad")

	// "stara pazova" → starapazova wins.
	assertTop1(t, m, "stara pazova", "starapazova")
	assertNotTop1(t, m, "stara pazova", "novapazova")

	// Bar is short but must not be beaten by unrelated words.
	assertTop1(t, m, "bar", "bar")

	// "kotor" must not be beaten by "kolasin" or "kosjeric".
	assertTop1(t, m, "kotor", "kotor")
	assertNotTop1(t, m, "kotor", "kolasin")
	assertNotTop1(t, m, "kotor", "kosjeric")

	// Prijepolje vs Prijepoljeteretna: plain query → shorter/plainer match first.
	assertTop1(t, m, "prijepolje", "prijepolje")

	// Tirana must not bubble up to "virpazar" or "ravnarijeka" etc.
	assertTop1(t, m, "tirana", "tirana")
}

func TestFalsePositives(t *testing.T) {
	m := NewMatcher(stationList)

	// Completely unrelated inputs should not confidently return a station.
	// We don't require zero results (some weak character overlap is OK), but the
	// score should be low.  We check that obviously wrong top results don't appear.

	// "london" shares only scattered chars with station names.
	results := m.Find("london")
	if len(results) > 0 && results[0].Score > 0.5 {
		t.Errorf("query \"london\": top result %q has unexpectedly high score %.3f",
			results[0].Word, results[0].Score)
	}

	// "chicago" — no station should dominate.
	results = m.Find("chicago")
	if len(results) > 0 && results[0].Score > 0.5 {
		t.Errorf("query \"chicago\": top result %q score %.3f too high",
			results[0].Word, results[0].Score)
	}

	// "berlin" — weak overlap with "belgrade" family, but should not beat them.
	assertNotTop1(t, m, "berlin", "podgorica")
	assertNotTop1(t, m, "berlin", "tirana")
}

func TestFalseNegatives(t *testing.T) {
	m := NewMatcher(stationList)

	// These are cases that may be hard for the algorithm; we want to confirm
	// the right answer IS present in the top results even if not always #1.

	// Heavy vowel shift — "padgareeka" should at least appear near the top.
	assertInTop(t, m, "padgareeka", "podgorica", 3)

	// Misheard ending — "podgoritsa" (Russian "-tsa" for "-ca"):
	assertInTop(t, m, "podgoritsa", "podgorica", 3)

	// Very short query — single real station name:
	assertInTop(t, m, "bar", "bar", 1)
	assertInTop(t, m, "kos", "kos", 1)

	// Compound names with spaces where the station key has no space:
	assertInTop(t, m, "bijelo polje", "bijelopolje", 1)
	assertInTop(t, m, "beograd centar", "beogradcentar", 1)
	assertInTop(t, m, "herceg novi", "hercegnovi", 1)
}

func TestNormalize(t *testing.T) {
	cases := []struct{ in, want string }{
		{"Podgorica", "podgorica"},
		{"BELGRADE", "belgrade"},
		// ije→e gives "belopolje", then lj→l gives "belopole"; space stripped
		{"Bijelo Polje", "belopole"},
		{"bijelo polje", "belopole"},
		{"belo pole", "belopole"},      // user input ekavica form — same result ✓
		{"Šabac", "sabac"},             // diacritic
		{"Čačak", "cacak"},
		{"Niksić", "niksic"},
		{"padgareeka", "padgarika"},    // ee→i
		{"Sutomore", "sutomore"},
		{"sjutamare", "sjutamare"},     // j stays (removed by skeleton, not normalize)
		{"New Belgrade", "nevbelgrade"}, // w→v, space stripped
		{"novi sad", "novisad"},
		{"Prijepolje", "prepole"},      // ije→e → "prepolje", then lj→l → "prepole"
		{"dj", "d"},                    // dj→d
		{"nj", "n"},
	}
	for _, tc := range cases {
		got := normalize(tc.in)
		if got != tc.want {
			t.Errorf("normalize(%q) = %q, want %q", tc.in, got, tc.want)
		}
	}
}

func TestConsonantSkeleton(t *testing.T) {
	cases := []struct{ in, want string }{
		{"podgorica", "pdgrc"},
		{"padgarika", "pdgrk"},  // proves padgareeka and podgorica converge
		{"novisad", "nvsd"},
		{"beograd", "bgrd"},
		{"sutomore", "stmr"},
		{"sjutamare", "sjtmr"},
		{"stmr", "stmr"},       // already all consonants
		{"", ""},
	}
	for _, tc := range cases {
		got := consonantSkeleton(tc.in)
		if got != tc.want {
			t.Errorf("consonantSkeleton(%q) = %q, want %q", tc.in, got, tc.want)
		}
	}
}

func TestScoreBounds(t *testing.T) {
	m := NewMatcher(stationList)

	// Perfect match should score very high (close to 1).
	results := m.Find("podgorica")
	if len(results) == 0 {
		t.Fatal("no results for exact query")
	}
	if results[0].Score < 0.8 {
		t.Errorf("exact match score = %.3f, want >= 0.8", results[0].Score)
	}

	// All scores must be in [0, 1].
	for _, query := range []string{"podgorica", "belgrade", "padgareeka", "xyz", "novi sad"} {
		for _, r := range m.Find(query) {
			if r.Score < 0 || r.Score > 1 {
				t.Errorf("query %q: score %.3f out of [0,1] for %q", query, r.Score, r.Word)
			}
		}
	}
}

func TestBenchmark(t *testing.T) {
	// Sanity-check that NewMatcher + a batch of finds completes quickly
	// (go test -v shows timing; this just makes sure it doesn't panic/hang).
	m := NewMatcher(stationList)
	queries := []string{
		"belgrade", "podgorica", "padgareeka", "sjutamare", "belo pole",
		"novi sad", "stara pazova", "nova pazova", "kotor", "bar",
	}
	for _, q := range queries {
		r := m.Find(q)
		if r == nil {
			t.Errorf("nil result for %q", q)
		}
	}
}

func BenchmarkNewMatcher(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NewMatcher(stationList)
	}
}

func BenchmarkFind(b *testing.B) {
	m := NewMatcher(stationList)
	queries := []string{
		"belgrade", "podgorica", "padgareeka", "sjutamare", "belo pole",
		"novi sad", "stara pazova",
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.Find(queries[i%len(queries)])
	}
}
