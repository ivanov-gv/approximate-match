package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// stationList is the fixed railway/location list used in all tests.
var stationList = []string{
	"aerodrom",              // 26
	"albania",               // -36
	"bar",                   // 1
	"baresumanovica",        // 51
	"becici",                // -28
	"belgrad",               // 18
	"belgrade",              // 18
	"beograd",               // 18
	"beogradcentar",         // 18
	"bijelopolje",           // 7
	"bioce",                 // 28
	"bosniaandherzegovina",  // -42
	"bratonozici",           // 29
	"budva",                 // -10
	"cetinje",               // -16
	"crmnica",               // 21
	"dabovici",              // 54
	"danilovgrad",           // 49
	"durmitor",              // -20
	"golubovci",             // 3
	"hercegnovi",            // -30
	"indjija",               // -2
	"kolasin",               // 5
	"kos",                   // 34
	"kosjeric",              // 13
	"kotor",                 // -14
	"krusevackipotok",       // 31
	"krusevo",               // 44
	"lajkovac",              // 15
	"lazarevac",             // 16
	"ljesnica",              // 45
	"ljutotuk",              // 48
	"lutovo",                // 30
	"matesevo",              // 35
	"mijatovokolo",          // 41
	"mojkovac",              // 6
	"moraca",                // 25
	"niksic",                // 56
	"novapazova",            // -6
	"novisad",               // 0
	"oblutak",               // 37
	"ostrog",                // 53
	"padez",                 // 36
	"perast",                // -18
	"petrovac",              // -22
	"podgorica",             // 4
	"pozega",                // 12
	"priboj",                // 10
	"pricelje",              // 46
	"prijepolje",            // 9
	"prijepoljeteretna",     // 8
	"rakovica",              // 17
	"ravnarijeka",           // 43
	"sarajevo",              // -44
	"savnik",                // -32
	"seliste",               // 33
	"shkoder",               // -40
	"slap",                  // 50
	"slijepacmost",          // 42
	"sobajici",              // 52
	"spuz",                  // 47
	"starapazova",           // -4
	"stitarickarijeka",      // 39
	"stubica",               // 55
	"subotica",              // -8
	"susanj",                // 20
	"sutomore",              // 2
	"svetistefan",           // -26
	"tirana",                // -38
	"tivat",                 // -12
	"trebaljevo",            // 38
	"trebesica",             // 32
	"ulcinj",                // -24
	"uzice",                 // 11
	"valjevo",               // 14
	"virpazar",              // 22
	"vranjina",              // 23
	"zabljak",               // -34
	"zari",                  // 40
	"zemun",                 // 19
	"zeta",                  // 24
	"zlatica",               // 27
}

func TestPositiveCases(t *testing.T) {
	m := NewMatcher(stationList)

	t.Run("ExactMatches", func(t *testing.T) {
		cases := []struct{ query, want string }{
			{"podgorica", "podgorica"},
			{"Podgorica", "podgorica"},
			{"PODGORICA", "podgorica"},
			{"bar", "bar"},
			{"kotor", "kotor"},
			{"budva", "budva"},
			{"tivat", "tivat"},
			{"tirana", "tirana"},
			{"novisad", "novisad"},
			{"niksic", "niksic"},
			{"sarajevo", "sarajevo"},
			{"subotica", "subotica"},
		}
		for _, tc := range cases {
			tc := tc
			t.Run(tc.query, func(t *testing.T) {
				assertTop1(t, m, tc.query, tc.want)
			})
		}
	})

	t.Run("SpacedNames", func(t *testing.T) {
		cases := []struct{ query, want string }{
			{"novi sad", "novisad"},
			{"nova pazova", "novapazova"},
			{"stara pazova", "starapazova"},
			{"bijelo polje", "bijelopolje"},
			{"herceg novi", "hercegnovi"},
			{"beograd centar", "beogradcentar"},
		}
		for _, tc := range cases {
			tc := tc
			t.Run(tc.query, func(t *testing.T) {
				assertTop1(t, m, tc.query, tc.want)
			})
		}
	})

	t.Run("MinorTypos", func(t *testing.T) {
		cases := []struct {
			query   string
			want    string
			withinN int
		}{
			{"belgarde", "belgrade", 2},  // transposition
			{"belgade", "belgrade", 3},   // missing r
			{"belgrate", "belgrade", 2},  // transposition r/t
			{"podgorcia", "podgorica", 1}, // transposition c/i
			{"sutmore", "sutomore", 1},   // missing o
		}
		for _, tc := range cases {
			tc := tc
			t.Run(tc.query, func(t *testing.T) {
				assertInTop(t, m, tc.query, tc.want, tc.withinN)
			})
		}
		assertTop1(t, m, "kolasin", "kolasin")
		assertTop1(t, m, "sutomore", "sutomore")
		assertTop1(t, m, "mojkovac", "mojkovac")
		assertTop1(t, m, "bijelopolje", "bijelopolje")
	})
}

func TestPhoneticCases(t *testing.T) {
	m := NewMatcher(stationList)

	t.Run("VowelShifts", func(t *testing.T) {
		cases := []struct{ query, want string }{
			{"padgareeka", "podgorica"},  // "pod-go-REE-ka" heard as "padgareeka"
			{"podgoriika", "podgorica"},  // doubled vowel
			{"podgoorica", "podgorica"},  // "oo" → u normalisation
			{"sjutamare", "sutomore"},   // vowel shifts + spurious j
			{"sutomare", "sutomore"},    // o→a vowel confusion
		}
		for _, tc := range cases {
			tc := tc
			t.Run(tc.query, func(t *testing.T) {
				assertTop1(t, m, tc.query, tc.want)
			})
		}
	})

	t.Run("EkavicaIjekavica", func(t *testing.T) {
		cases := []struct{ query, want string }{
			{"belo pole", "bijelopolje"},
			{"belo polje", "bijelopolje"},
			{"bijelo polje", "bijelopolje"},
		}
		for _, tc := range cases {
			tc := tc
			t.Run(tc.query, func(t *testing.T) {
				assertTop1(t, m, tc.query, tc.want)
			})
		}
	})

	t.Run("TransliterationVariants", func(t *testing.T) {
		cases := []struct{ query, want string }{
			{"niksic", "niksic"},
			{"priboj", "priboj"},
			{"belgrade", "belgrade"},
			{"Belgrade", "belgrade"},
			{"novi sad", "novisad"},
			{"Novi Sad", "novisad"},
			{"shkoder", "shkoder"},
		}
		for _, tc := range cases {
			tc := tc
			t.Run(tc.query, func(t *testing.T) {
				assertTop1(t, m, tc.query, tc.want)
			})
		}
	})

	t.Run("RussianTsForC", func(t *testing.T) {
		// "ts" isn't a special rule but consonant skeleton helps: pdgrtsk ≈ pdgrc
		assertInTop(t, m, "podgoritsa", "podgorica", 3)
	})
}

func TestDisambiguation(t *testing.T) {
	m := NewMatcher(stationList)

	t.Run("BelgradeFamily", func(t *testing.T) {
		assertNotTop1(t, m, "belgrade", "novisad")
		assertNotTop1(t, m, "beograd", "novisad")
		assertInTop(t, m, "belgrade", "belgrade", 1)
		assertInTop(t, m, "beograd centar", "beogradcentar", 1)
	})

	t.Run("NoviSadVsPazova", func(t *testing.T) {
		assertTop1(t, m, "novi sad", "novisad")
		assertNotTop1(t, m, "novi sad", "novapazova")
		assertNotTop1(t, m, "novi sad", "starapazova")

		assertTop1(t, m, "nova pazova", "novapazova")
		assertNotTop1(t, m, "nova pazova", "starapazova")
		assertNotTop1(t, m, "nova pazova", "novisad")

		assertTop1(t, m, "stara pazova", "starapazova")
		assertNotTop1(t, m, "stara pazova", "novapazova")
	})

	t.Run("ShortNames", func(t *testing.T) {
		assertTop1(t, m, "bar", "bar")
		assertTop1(t, m, "kotor", "kotor")
		assertNotTop1(t, m, "kotor", "kolasin")
		assertNotTop1(t, m, "kotor", "kosjeric")
	})

	t.Run("CompoundNames", func(t *testing.T) {
		assertTop1(t, m, "prijepolje", "prijepolje")
		assertTop1(t, m, "tirana", "tirana")
	})
}

func TestFalsePositives(t *testing.T) {
	m := NewMatcher(stationList)

	t.Run("UnrelatedInputsHaveLowScore", func(t *testing.T) {
		for _, query := range []string{"london", "chicago"} {
			query := query
			t.Run(query, func(t *testing.T) {
				results := m.Find(query)
				if len(results) > 0 {
					assert.Less(t, results[0].Score, 0.5,
						"query %q: top result %q has unexpectedly high score", query, results[0].Word)
				}
			})
		}
	})

	t.Run("BerlinDoesNotMatchUnrelated", func(t *testing.T) {
		assertNotTop1(t, m, "berlin", "podgorica")
		assertNotTop1(t, m, "berlin", "tirana")
	})
}

func TestFalseNegatives(t *testing.T) {
	m := NewMatcher(stationList)

	cases := []struct {
		query   string
		want    string
		withinN int
	}{
		{"padgareeka", "podgorica", 3},
		{"podgoritsa", "podgorica", 3},
		{"bar", "bar", 1},
		{"kos", "kos", 1},
		{"bijelo polje", "bijelopolje", 1},
		{"beograd centar", "beogradcentar", 1},
		{"herceg novi", "hercegnovi", 1},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.query, func(t *testing.T) {
			assertInTop(t, m, tc.query, tc.want, tc.withinN)
		})
	}
}

func TestNormalize(t *testing.T) {
	cases := []struct{ in, want string }{
		{"Podgorica", "podgorica"},
		{"BELGRADE", "belgrade"},
		// ije→e gives "belopolje", then lj→l gives "belopole"; space stripped
		{"Bijelo Polje", "belopole"},
		{"bijelo polje", "belopole"},
		{"belo pole", "belopole"},       // user input ekavica form — same result ✓
		{"Šabac", "sabac"},              // diacritic via Unicode NFD
		{"Čačak", "cacak"},
		{"Niksić", "niksic"},
		{"padgareeka", "padgarika"},     // ee→i
		{"Sutomore", "sutomore"},
		{"sjutamare", "sjutamare"},      // j stays (removed by skeleton, not normalize)
		{"New Belgrade", "nevbelgrade"}, // w→v, space stripped
		{"novi sad", "novisad"},
		{"Prijepolje", "prepole"},       // ije→e → "prepolje", then lj→l → "prepole"
		{"dj", "d"},
		{"nj", "n"},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.in, func(t *testing.T) {
			assert.Equal(t, tc.want, normalize(tc.in))
		})
	}
}

func TestConsonantSkeleton(t *testing.T) {
	cases := []struct{ in, want string }{
		{"podgorica", "pdgrc"},
		{"padgarika", "pdgrk"}, // proves padgareeka and podgorica converge
		{"novisad", "nvsd"},
		{"beograd", "bgrd"},
		{"sutomore", "stmr"},
		{"sjutamare", "sjtmr"},
		{"stmr", "stmr"}, // already all consonants
		{"", ""},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.in, func(t *testing.T) {
			assert.Equal(t, tc.want, consonantSkeleton(tc.in))
		})
	}
}

func TestScoreBounds(t *testing.T) {
	m := NewMatcher(stationList)

	t.Run("ExactMatchScoresHigh", func(t *testing.T) {
		results := m.Find("podgorica")
		require.NotEmpty(t, results, "no results for exact query")
		assert.GreaterOrEqual(t, results[0].Score, 0.8,
			"exact match score too low: %.3f", results[0].Score)
	})

	t.Run("AllScoresInRange", func(t *testing.T) {
		for _, query := range []string{"podgorica", "belgrade", "padgareeka", "xyz", "novi sad"} {
			for _, r := range m.Find(query) {
				assert.GreaterOrEqual(t, r.Score, 0.0,
					"query %q: score %.3f below 0 for %q", query, r.Score, r.Word)
				assert.LessOrEqual(t, r.Score, 1.0,
					"query %q: score %.3f above 1 for %q", query, r.Score, r.Word)
			}
		}
	})
}

func TestBenchmark(t *testing.T) {
	m := NewMatcher(stationList)
	queries := []string{
		"belgrade", "podgorica", "padgareeka", "sjutamare", "belo pole",
		"novi sad", "stara pazova", "nova pazova", "kotor", "bar",
	}
	for _, query := range queries {
		query := query
		t.Run(query, func(t *testing.T) {
			result := m.Find(query)
			assert.NotNil(t, result)
		})
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
