package approxmatch_test

import (
	"testing"

	approxmatch "github.com/ivanov-gv/approximate-match"
	"github.com/stretchr/testify/assert"
)

func TestNormalize(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want string
	}{
		{"Lowercase", "Podgorica", "podgorica"},
		{"AllCaps", "BELGRADE", "belgrade"},
		{"IjekavicaAndDigraphAndSpace", "Bijelo Polje", "belopole"},
		{"IjekavicaAndDigraphLower", "bijelo polje", "belopole"},
		{"EkavicaSameResult", "belo pole", "belopole"},
		{"DiacriticNFD_S", "Šabac", "sabac"},
		{"DiacriticNFD_C", "Čačak", "cacak"},
		{"DiacriticNFD_Acute", "Niksić", "niksic"},
		{"DoubleEToI", "padgareeka", "padgarika"},
		{"NoChange", "Sutomore", "sutomore"},
		{"SpuriousJStays", "sjutamare", "sjutamare"},
		{"WToVAndSpace", "New Belgrade", "nevbelgrade"},
		{"SpaceStripped", "novi sad", "novisad"},
		{"IjekavicaAndDigraph", "Prijepolje", "prepole"},
		{"DigraphDJ", "dj", "d"},
		{"DigraphNJ", "nj", "n"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.want, approxmatch.Normalize(tc.in))
		})
	}
}

func TestNormalizeCyrillic(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want string
	}{
		// Russian soft-sign sequences converge with Serbian Cyrillic ligatures.
		{"RussianSoftL_ValEvo", "вальево", "валево"},
		{"SerbianLjLigature_ValEvo", "ваљево", "валево"},
		{"RussianSoftL_Ljesnica", "льешница", "лешница"},
		{"SerbianLjLigature_Ljesnica", "љешница", "лешница"},
		{"RussianSoftN_Susanj", "шушань", "шушан"},
		{"SerbianNjLigature_Susanj", "шушањ", "шушан"},
		{"RussianSoftN_Vranjina", "враньина", "вранина"},
		{"SerbianNjLigature_Vranjina", "врањина", "вранина"},
		// Serbian ћ → ч so Russian and Serbian Cyrillic spellings converge.
		{"SerbianTshe_Bratonozici", "братоношићи", "братоношичи"},
		// Russian substitutes и where Serbian uses ј (which NFD-strips й to и automatically).
		{"SerbianJ_Mojkovac", "мојковац", "моиковац"},
		{"SerbianJ_Priboj", "прибој", "прибои"},
		// Russian ю → у; ы → и.
		{"RussianYu_Ljutotuk", "лютотук", "лутотук"},
		{"RussianY_Golubovcy", "голубовцы", "голубовци"},
		// Cyrillic spaces stripped same as Latin.
		{"SpaceStripped", "нови сад", "новисад"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.want, approxmatch.Normalize(tc.in))
		})
	}
}

func TestConsonantSkeleton(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want string
	}{
		{"Podgorica", "podgorica", "pdgrc"},
		{"PadgareekaPodgoricaConverge", "padgarika", "pdgrk"},
		{"NoviSad", "novisad", "nvsd"},
		{"Beograd", "beograd", "bgrd"},
		{"Sutomore", "sutomore", "stmr"},
		{"SpuriousJ", "sjutamare", "sjtmr"},
		{"AlreadyAllConsonants", "stmr", "stmr"},
		{"Empty", "", ""},
		// Cyrillic vowels are stripped too.
		{"CyrillicPodgorica", "подгорица", "пдгрц"},
		// After normalization, валево has Cyrillic vowels stripped.
		{"CyrillicValjevo", "валево", "влв"},
		// Both Russian and Serbian forms converge on the same skeleton.
		{"CyrillicLjesnicaSkeleton", "лешница", "лшнц"},
		{"CyrillicSusanjSkeleton", "шушан", "шшн"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.want, approxmatch.ConsonantSkeleton(tc.in))
		})
	}
}
