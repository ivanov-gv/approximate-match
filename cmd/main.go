package main

import (
	"fmt"

	approxmatch "github.com/ivanov-gv/approximate-match"
)

func main() {
	// Example from README: find the closest match for "aaaaa".
	fmt.Println("Sample: aaaaa")
	matcher := approxmatch.NewMatcher([]string{"aaaab", "aaaac", "abcde", "bbbbb"}, nil)
	for _, match := range matcher.Find("aaaaa") {
		fmt.Printf("  %-8s score: %.3f\n", match.Word, match.Score)
	}

	// Example from README: find the closest match for "Ella" among names.
	fmt.Println("\nSample: Ella")
	matcher = approxmatch.NewMatcher([]string{"Adele", "Elaine", "Elizabeth", "Harriet", "Ingrid", "Michelle", "Ella"}, nil)
	for _, match := range matcher.Find("Ella") {
		fmt.Printf("  %-12s score: %.3f\n", match.Word, match.Score)
	}
}
