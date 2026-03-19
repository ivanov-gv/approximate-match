package integration_test

import (
	"testing"

	"github.com/samber/lo"

	approxmatch "github.com/ivanov-gv/approximate-match"
)

func BenchmarkNewMatcherSmallList(b *testing.B) {
	words := []string{"podgorica", "belgrade", "novisad", "sarajevo", "tirana"}
	for i := 0; i < b.N; i++ {
		approxmatch.NewMatcher(words, nil)
	}
}

// BenchmarkFindSmallList — 5-word synthetic list, AMD Ryzen 7 5700U:
// 9339 ns/op, 9545 B/op, 24 allocs/op
func BenchmarkFindSmallList(b *testing.B) {
	matcher := approxmatch.NewMatcher([]string{"podgorica", "belgrade", "novisad", "sarajevo", "tirana"}, nil)
	queries := []string{"belgrade", "podgorica", "padgareeka", "sjutamare", "belo pole"}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		matcher.Find(queries[i%len(queries)])
	}
}

func BenchmarkNewMatcher(b *testing.B) {
	names := lo.Keys(nameToStationName)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		approxmatch.NewMatcher(names, nil)
	}
}

// BenchmarkFind — full station list (~370 entries), AMD Ryzen 7 5700U:
// 112498 ns/op, 12560 B/op, 24 allocs/op
func BenchmarkFind(b *testing.B) {
	matcher := approxmatch.NewMatcher(lo.Keys(nameToStationName), nil)
	queries := []string{
		"belgrade", "podgorica", "padgareeka", "sjutamare", "belo pole",
		"novi sad", "stara pazova",
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		matcher.Find(queries[i%len(queries)])
	}
}
