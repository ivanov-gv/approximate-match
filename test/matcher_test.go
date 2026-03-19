package integration_test

import (
	"testing"

	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	approxmatch "github.com/ivanov-gv/approximate-match"
	integration "github.com/ivanov-gv/approximate-match/test"
)

// nameToStationName maps every indexed name (Name, NameEn, NameCyr, and
// ProductionAliases) of every station to the station's canonical Name field.
var nameToStationName = func() map[string]string {
	result := make(map[string]string)
	for _, station := range integration.Stations {
		allNames := lo.Compact(append([]string{station.Name, station.NameEn, station.NameCyr}, station.ProductionAliases...))
		for _, name := range allNames {
			result[name] = station.Name
		}
	}
	return result
}()

// officialNameToStationName maps all names of non-blacklisted stations
// (Name, NameEn, NameCyr, and ProductionAliases) to the station's canonical
// Name field. ProductionAliases are included so that every alias resolves
// correctly when searched in the official index.
var officialNameToStationName = func() map[string]string {
	result := make(map[string]string)
	nonBlacklistedStations := lo.Filter(integration.Stations, func(s integration.StationData, _ int) bool {
		return !s.Blacklisted
	})
	for _, station := range nonBlacklistedStations {
		allNames := lo.Uniq(lo.Compact(append(
			[]string{station.Name, station.NameEn, station.NameCyr},
			station.ProductionAliases...,
		)))
		for _, name := range allNames {
			result[name] = station.Name
		}
	}
	return result
}()

func TestPositiveCases(t *testing.T) {
	matcher := approxmatch.NewMatcher(lo.Keys(nameToStationName), nil)

	t.Run("ExactMatches", func(t *testing.T) {
		cases := []struct {
			query    string
			wantName string
		}{
			{"podgorica", "Podgorica"},
			{"Podgorica", "Podgorica"},
			{"PODGORICA", "Podgorica"},
			{"bar", "Bar"},
			{"kotor", "Kotor"},
			{"budva", "Budva"},
			{"tivat", "Tivat"},
			{"tirana", "Tirana"},
			{"novisad", "Novi Sad"},
			{"niksic", "Nikšić"},
			{"sarajevo", "Sarajevo"},
			{"subotica", "Subotica"},
		}
		for _, tc := range cases {
			t.Run(tc.query, func(t *testing.T) {
				results := matcher.Find(tc.query)
				require.NotEmpty(t, results, "query %q: got no results, want %q", tc.query, tc.wantName)
				gotName := nameToStationName[results[0].Word]
				assert.Equal(t, tc.wantName, gotName,
					"query %q: top result %q maps to %q, want %q",
					tc.query, results[0].Word, gotName, tc.wantName)
			})
		}
	})

	t.Run("SpacedNames", func(t *testing.T) {
		cases := []struct {
			query    string
			wantName string
		}{
			{"novi sad", "Novi Sad"},
			{"nova pazova", "Nova Pazova"},
			{"stara pazova", "Stara Pazova"},
			{"bijelo polje", "Bijelo Polje"},
			{"herceg novi", "Herceg Novi"},
			{"beograd centar", "Beograd Centar"},
		}
		for _, tc := range cases {
			t.Run(tc.query, func(t *testing.T) {
				results := matcher.Find(tc.query)
				require.NotEmpty(t, results, "query %q: got no results, want %q", tc.query, tc.wantName)
				gotName := nameToStationName[results[0].Word]
				assert.Equal(t, tc.wantName, gotName,
					"query %q: top result %q maps to %q, want %q",
					tc.query, results[0].Word, gotName, tc.wantName)
			})
		}
	})

	t.Run("MinorTypos", func(t *testing.T) {
		cases := []struct {
			query    string
			wantName string
		}{
			{"belgarde", "Beograd Centar"}, // transposition
			{"belgade", "Beograd Centar"},  // missing r
			{"belgrate", "Beograd Centar"}, // transposition r/t
			{"podgorcia", "Podgorica"},     // transposition c/i
			{"sutmore", "Sutomore"},        // missing o
			{"kolasin", "Kolašin"},
			{"sutomore", "Sutomore"},
			{"mojkovac", "Mojkovac"},
			{"bijelopolje", "Bijelo Polje"},
		}
		for _, tc := range cases {
			t.Run(tc.query, func(t *testing.T) {
				results := matcher.Find(tc.query)
				require.NotEmpty(t, results, "query %q: got no results, want %q", tc.query, tc.wantName)
				gotName := nameToStationName[results[0].Word]
				assert.Equal(t, tc.wantName, gotName,
					"query %q: top result %q maps to %q, want %q",
					tc.query, results[0].Word, gotName, tc.wantName)
			})
		}
	})
}

func TestPhoneticCases(t *testing.T) {
	matcher := approxmatch.NewMatcher(lo.Keys(nameToStationName), nil)

	t.Run("VowelShifts", func(t *testing.T) {
		cases := []struct {
			query    string
			wantName string
		}{
			{"padgareeka", "Podgorica"}, // "pod-go-REE-ka" heard as "padgareeka"
			{"podgoriika", "Podgorica"}, // doubled vowel
			{"podgoorica", "Podgorica"}, // "oo" → u normalisation
			{"sjutamare", "Sutomore"},   // vowel shifts + spurious j
			{"sutomare", "Sutomore"},    // o→a vowel confusion
		}
		for _, tc := range cases {
			t.Run(tc.query, func(t *testing.T) {
				results := matcher.Find(tc.query)
				require.NotEmpty(t, results, "query %q: got no results, want %q", tc.query, tc.wantName)
				gotName := nameToStationName[results[0].Word]
				assert.Equal(t, tc.wantName, gotName,
					"query %q: top result %q maps to %q, want %q",
					tc.query, results[0].Word, gotName, tc.wantName)
			})
		}
	})

	t.Run("EkavicaIjekavica", func(t *testing.T) {
		cases := []struct {
			query    string
			wantName string
		}{
			{"belo pole", "Bijelo Polje"},
			{"belo polje", "Bijelo Polje"},
			{"bijelo polje", "Bijelo Polje"},
		}
		for _, tc := range cases {
			t.Run(tc.query, func(t *testing.T) {
				results := matcher.Find(tc.query)
				require.NotEmpty(t, results, "query %q: got no results, want %q", tc.query, tc.wantName)
				gotName := nameToStationName[results[0].Word]
				assert.Equal(t, tc.wantName, gotName,
					"query %q: top result %q maps to %q, want %q",
					tc.query, results[0].Word, gotName, tc.wantName)
			})
		}
	})

	t.Run("TransliterationVariants", func(t *testing.T) {
		cases := []struct {
			query    string
			wantName string
		}{
			{"niksic", "Nikšić"},
			{"nickshicsh", "Nikšić"},
			{"priboj", "Priboj"},
			{"belgrade", "Beograd Centar"},
			{"Belgrade", "Beograd Centar"},
			{"novi sad", "Novi Sad"},
			{"Novi Sad", "Novi Sad"},
			{"shkoder", "Shkoder"},
			{"shushann", "Šušanj"},
		}
		for _, tc := range cases {
			t.Run(tc.query, func(t *testing.T) {
				results := matcher.Find(tc.query)
				require.NotEmpty(t, results, "query %q: got no results, want %q", tc.query, tc.wantName)
				gotName := nameToStationName[results[0].Word]
				assert.Equal(t, tc.wantName, gotName,
					"query %q: top result %q maps to %q, want %q",
					tc.query, results[0].Word, gotName, tc.wantName)
			})
		}
	})

	t.Run("RussianTsForC", func(t *testing.T) {
		// "ts" has no explicit rule; consonant skeleton bridges pdgrtsk → pdgrc.
		results := matcher.Find("podgoritsa")
		require.NotEmpty(t, results, "query 'podgoritsa': got no results, want Podgorica")
		gotName := nameToStationName[results[0].Word]
		assert.Equal(t, "Podgorica", gotName,
			"query 'podgoritsa': top result %q maps to %q, want Podgorica",
			results[0].Word, gotName)
	})
}

func TestDisambiguation(t *testing.T) {
	matcher := approxmatch.NewMatcher(lo.Keys(nameToStationName), nil)

	run := func(t *testing.T, cases []struct {
		query        string
		wantName     string
		mustNotNames []string
	}) {
		t.Helper()
		for _, tc := range cases {
			t.Run(tc.query, func(t *testing.T) {
				results := matcher.Find(tc.query)
				require.NotEmpty(t, results, "query %q: got no results, want %q", tc.query, tc.wantName)
				gotName := nameToStationName[results[0].Word]
				assert.Equal(t, tc.wantName, gotName,
					"query %q: top result %q maps to %q, want %q",
					tc.query, results[0].Word, gotName, tc.wantName)
				for _, mustNot := range tc.mustNotNames {
					assert.NotEqual(t, mustNot, gotName,
						"query %q: must not surface %q", tc.query, mustNot)
				}
			})
		}
	}

	t.Run("BelgradeFamily", func(t *testing.T) {
		run(t, []struct {
			query        string
			wantName     string
			mustNotNames []string
		}{
			{"belgrade", "Beograd Centar", []string{"Novi Sad"}},
			{"beograd", "Beograd Centar", []string{"Novi Sad"}},
			{"beograd centar", "Beograd Centar", nil},
		})
	})

	t.Run("NoviSadVsPazova", func(t *testing.T) {
		run(t, []struct {
			query        string
			wantName     string
			mustNotNames []string
		}{
			{"novi sad", "Novi Sad", []string{"Nova Pazova", "Stara Pazova"}},
			{"nova pazova", "Nova Pazova", []string{"Stara Pazova", "Novi Sad"}},
			{"stara pazova", "Stara Pazova", []string{"Nova Pazova"}},
		})
	})

	t.Run("ShortNames", func(t *testing.T) {
		run(t, []struct {
			query        string
			wantName     string
			mustNotNames []string
		}{
			{"bar", "Bar", nil},
			{"kotor", "Kotor", []string{"Kolašin", "Kosjerić"}},
		})
	})

	t.Run("CompoundNames", func(t *testing.T) {
		run(t, []struct {
			query        string
			wantName     string
			mustNotNames []string
		}{
			{"prijepolje", "Prijepolje", nil},
			{"tirana", "Tirana", nil},
		})
	})

	// These four stations share tokens ("novi", "beograd") that cause frequent
	// cross-matches. Each query must surface its own station at the top.
	t.Run("NoviBeogradHercegNoviFamily", func(t *testing.T) {
		run(t, []struct {
			query        string
			wantName     string
			mustNotNames []string
		}{
			// Novi Beograd must not surface as Beograd Centar or Novi Sad.
			{"novi beograd", "Novi Beograd", []string{"Beograd Centar", "Novi Sad", "Herceg Novi"}},
			{"novibeograd", "Novi Beograd", []string{"Beograd Centar", "Novi Sad"}},
			{"novi belgrado", "Novi Beograd", []string{"Beograd Centar", "Novi Sad"}},
			{"novi beograde", "Novi Beograd", []string{"Beograd Centar", "Novi Sad"}}, // minor typo
			{"новый белград", "Novi Beograd", []string{"Beograd Centar", "Novi Sad"}}, // minor typo
			{"нови белград", "Novi Beograd", []string{"Beograd Centar", "Novi Sad"}},  // minor typo
			{"нови београд", "Novi Beograd", []string{"Beograd Centar", "Novi Sad"}},  // minor typo

			// Herceg Novi must not surface as Novi Sad or Novi Beograd.
			{"herceg novi", "Herceg Novi", []string{"Novi Sad", "Novi Beograd"}},
			{"hercegnovi", "Herceg Novi", []string{"Novi Sad", "Novi Beograd"}},
			{"hertzeg novi", "Herceg Novi", []string{"Novi Sad", "Novi Beograd"}}, // ch/tz transliteration
			{"херцег нови", "Herceg Novi", []string{"Novi Sad", "Novi Beograd"}},  // ch/tz transliteration
			{"герцег нови", "Herceg Novi", []string{"Novi Sad", "Novi Beograd"}},  // ch/tz transliteration

			// Beograd Centar must win over Novi Beograd for bare "beograd".
			{"beograd", "Beograd Centar", []string{"Novi Beograd"}},
			{"belgrade", "Beograd Centar", []string{"Novi Beograd"}},
			{"белград", "Beograd Centar", []string{"Novi Beograd"}},

			// Novi Sad must win over Novi Beograd for bare "novi sad".
			{"novi sad", "Novi Sad", []string{"Novi Beograd", "Herceg Novi"}},
			{"нови сад", "Novi Sad", []string{"Novi Beograd", "Herceg Novi"}},
			{"новый сад", "Novi Sad", []string{"Novi Beograd", "Herceg Novi"}},
		})
	})
}

func TestFalsePositives(t *testing.T) {
	matcher := approxmatch.NewMatcher(lo.Keys(nameToStationName), nil)

	for _, query := range []string{"london", "chicago", "berlin"} {
		t.Run(query, func(t *testing.T) {
			results := matcher.Find(query)
			require.Empty(t, results,
				"query %q: expected no results above score threshold", query)
		})
	}
}

func TestFalseNegatives(t *testing.T) {
	matcher := approxmatch.NewMatcher(lo.Keys(nameToStationName), nil)

	cases := []struct {
		query    string
		wantName string
	}{
		{"padgareeka", "Podgorica"},
		{"podgoritsa", "Podgorica"},
		{"bar", "Bar"},
		{"kos", "Kos"},
		{"bijelo polje", "Bijelo Polje"},
		{"beograd centar", "Beograd Centar"},
		{"herceg novi", "Herceg Novi"},
	}
	for _, tc := range cases {
		t.Run(tc.query, func(t *testing.T) {
			results := matcher.Find(tc.query)
			require.NotEmpty(t, results, "query %q: got no results, want %q", tc.query, tc.wantName)
			gotName := nameToStationName[results[0].Word]
			assert.Equal(t, tc.wantName, gotName,
				"query %q: top result %q maps to %q, want %q",
				tc.query, results[0].Word, gotName, tc.wantName)
		})
	}
}

func TestCyrillicCases(t *testing.T) {
	matcher := approxmatch.NewMatcher(lo.Keys(nameToStationName), nil)

	t.Run("ExactMatches", func(t *testing.T) {
		cases := []struct {
			query    string
			wantName string
		}{
			{"подгорица", "Podgorica"},
			{"бар", "Bar"},
			{"сутоморе", "Sutomore"},
			{"нови сад", "Novi Sad"},
			{"никшић", "Nikšić"},
			{"тирана", "Tirana"},
			{"мојковац", "Mojkovac"},
		}
		for _, tc := range cases {
			t.Run(tc.query, func(t *testing.T) {
				results := matcher.Find(tc.query)
				require.NotEmpty(t, results, "query %q: got no results, want %q", tc.query, tc.wantName)
				gotName := nameToStationName[results[0].Word]
				assert.Equal(t, tc.wantName, gotName,
					"query %q: top result %q maps to %q, want %q",
					tc.query, results[0].Word, gotName, tc.wantName)
			})
		}
	})

	t.Run("SpacedNames", func(t *testing.T) {
		cases := []struct {
			query    string
			wantName string
		}{
			{"нови сад", "Novi Sad"},      // space stripped → новисад
			{"бело поле", "Bijelo Polje"}, // space stripped → белополе
			{"бијело поље", "Bijelo Polje"},
			{"херцег нови", "Herceg Novi"},
			{"београд центар", "Beograd Centar"},
		}
		for _, tc := range cases {
			t.Run(tc.query, func(t *testing.T) {
				results := matcher.Find(tc.query)
				require.NotEmpty(t, results, "query %q: got no results, want %q", tc.query, tc.wantName)
				gotName := nameToStationName[results[0].Word]
				assert.Equal(t, tc.wantName, gotName,
					"query %q: top result %q maps to %q, want %q",
					tc.query, results[0].Word, gotName, tc.wantName)
			})
		}
	})

	t.Run("RussianSoftSign", func(t *testing.T) {
		// Russian writes Serbian lj/nj sounds as "ль"/"нь"; normalization collapses
		// both to the base consonant so they match the Serbian Cyrillic ligature form.
		cases := []struct {
			query    string
			wantName string
		}{
			{"льешница", "Lješnica"},     // Russian ль ≈ Serbian љ
			{"шушань", "Šušanj"},         // Russian нь ≈ Serbian њ
			{"вальево", "Valjevo"},       // Russian ль ≈ Serbian љ
			{"враньина", "Vranjina"},     // Russian нь ≈ Serbian њ
			{"требальево", "Trebaljevo"}, // Russian ль ≈ Serbian љ
		}
		for _, tc := range cases {
			t.Run(tc.query, func(t *testing.T) {
				results := matcher.Find(tc.query)
				require.NotEmpty(t, results, "query %q: got no results, want %q", tc.query, tc.wantName)
				gotName := nameToStationName[results[0].Word]
				assert.Equal(t, tc.wantName, gotName,
					"query %q: top result %q maps to %q, want %q",
					tc.query, results[0].Word, gotName, tc.wantName)
			})
		}
	})

	t.Run("RussianVowelVariants", func(t *testing.T) {
		// Russian ю → у and ы → и let Russian and Serbian spellings converge.
		cases := []struct {
			query    string
			wantName string
		}{
			{"лютотук", "Ljutotuk"},    // Russian ю ≈ Serbian у after soft consonant (лю ≈ љу)
			{"голубовцы", "Golubovci"}, // Russian ы ≈ Serbian и at word end
		}
		for _, tc := range cases {
			t.Run(tc.query, func(t *testing.T) {
				results := matcher.Find(tc.query)
				require.NotEmpty(t, results, "query %q: got no results, want %q", tc.query, tc.wantName)
				gotName := nameToStationName[results[0].Word]
				assert.Equal(t, tc.wantName, gotName,
					"query %q: top result %q maps to %q, want %q",
					tc.query, results[0].Word, gotName, tc.wantName)
			})
		}
	})

	t.Run("SerbianVsRussianSpelling", func(t *testing.T) {
		// Russian transliteration of Serbian station names uses different letters
		// for Serbian ћ (→ ч) and ј (→ и/й).
		cases := []struct {
			query    string
			wantName string
		}{
			{"братоношичи", "Bratonožići"}, // Russian ч ≈ Serbian ћ in братоношићи
			{"никшич", "Nikšić"},           // Russian ч ≈ Serbian ћ in никшић
			{"моиковац", "Mojkovac"},       // Russian и ≈ Serbian ј in мојковац
			{"прибои", "Priboj"},           // Russian и ≈ Serbian ј in прибој
			{"лаиковац", "Lajkovac"},       // Russian и ≈ Serbian ј in лајковац
		}
		for _, tc := range cases {
			t.Run(tc.query, func(t *testing.T) {
				results := matcher.Find(tc.query)
				require.NotEmpty(t, results, "query %q: got no results, want %q", tc.query, tc.wantName)
				gotName := nameToStationName[results[0].Word]
				assert.Equal(t, tc.wantName, gotName,
					"query %q: top result %q maps to %q, want %q",
					tc.query, results[0].Word, gotName, tc.wantName)
			})
		}
	})
}

func TestScoreBounds(t *testing.T) {
	matcher := approxmatch.NewMatcher(lo.Keys(nameToStationName), nil)

	t.Run("ExactMatchScoresHigh", func(t *testing.T) {
		results := matcher.Find("podgorica")
		require.NotEmpty(t, results, "no results for exact query 'podgorica'")
		assert.Greater(t, results[0].Score, approxmatch.DefaultScoreThreshold,
			"exact match score too low: %.3f", results[0].Score)
	})

	t.Run("AllScoresInRange", func(t *testing.T) {
		for _, query := range []string{"podgorica", "belgrade", "padgareeka", "xyz", "novi sad"} {
			t.Run(query, func(t *testing.T) {
				for _, result := range matcher.Find(query) {
					assert.GreaterOrEqual(t, result.Score, approxmatch.DefaultScoreThreshold,
						"query %q: score %.3f below threshold for %q", query, result.Score, result.Word)
					assert.LessOrEqual(t, result.Score, 1.0,
						"query %q: score %.3f above 1 for %q", query, result.Score, result.Word)
				}
			})
		}
	})
}

func TestOfficialStationNames(t *testing.T) {
	officialMatcher := approxmatch.NewMatcher(lo.Keys(officialNameToStationName), nil)

	nonBlacklistedStations := lo.Filter(integration.Stations, func(station integration.StationData, _ int) bool {
		return !station.Blacklisted
	})
	for _, station := range nonBlacklistedStations {
		t.Run(station.Name, func(t *testing.T) {
			results := officialMatcher.Find(station.Name)
			require.NotEmpty(t, results, "station %q: no results for Name", station.Name)
			assert.Equal(t, station.Name, officialNameToStationName[results[0].Word],
				"station %q: top result %q should map to %q",
				station.Name, results[0].Word, station.Name)

			if station.NameEn != station.Name && station.NameEn != "" {
				results = officialMatcher.Find(station.NameEn)
				require.NotEmpty(t, results, "station %q: no results for NameEn %q", station.Name, station.NameEn)
				assert.Equal(t, station.Name, officialNameToStationName[results[0].Word],
					"station %q: NameEn %q top result %q should map to %q",
					station.Name, station.NameEn, results[0].Word, station.Name)
			}

			results = officialMatcher.Find(station.NameCyr)
			require.NotEmpty(t, results, "station %q: no results for NameCyr %q", station.Name, station.NameCyr)
			assert.Equal(t, station.Name, officialNameToStationName[results[0].Word],
				"station %q: NameCyr %q top result %q should map to %q",
				station.Name, station.NameCyr, results[0].Word, station.Name)
		})
	}
}

func TestAliasesMatchOfficialStations(t *testing.T) {
	officialMatcher := approxmatch.NewMatcher(lo.Keys(officialNameToStationName), nil)

	stationsWithAliases := lo.Filter(integration.Stations, func(station integration.StationData, _ int) bool {
		return !station.Blacklisted && len(station.ProductionAliases) > 0
	})
	for _, station := range stationsWithAliases {
		t.Run(station.Name, func(t *testing.T) {
			for _, alias := range station.ProductionAliases {
				t.Run(alias, func(t *testing.T) {
					results := officialMatcher.Find(alias)
					require.NotEmpty(t, results, "alias %q → station %q: no results in official list", alias, station.Name)
					gotName := officialNameToStationName[results[0].Word]
					assert.Equal(t, station.Name, gotName,
						"alias %q → station %q: top result %q maps to %q, want %q",
						alias, station.Name, results[0].Word, gotName, station.Name)
				})
			}
		})
	}
}

func TestBlacklistedStationsNoMatch(t *testing.T) {
	// Borderline skeleton matches for short Cyrillic queries (e.g. Тиват→Лутово)
	// reach scores just above DefaultScoreThreshold. A slightly stricter threshold
	// eliminates them while keeping all legitimate station matches. This test
	// demonstrates passing a custom threshold to NewMatcher.
	strictThreshold := 0.6
	officialMatcher := approxmatch.NewMatcher(lo.Keys(officialNameToStationName), &strictThreshold)

	for _, station := range integration.Stations {
		if !station.Blacklisted {
			continue
		}
		for _, name := range lo.Compact([]string{station.Name, station.NameCyr}) {
			t.Run(name, func(t *testing.T) {
				results := officialMatcher.Find(name)
				var resultWord string
				var resultScore float64
				if len(results) > 0 {
					resultWord = results[0].Word
					resultScore = results[0].Score
				}
				assert.Empty(t, results,
					"query %q: expected no results at strict threshold, got top result %q (score %.3f)",
					name, resultWord, resultScore)
			})
		}
	}
}
