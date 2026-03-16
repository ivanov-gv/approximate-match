package main

import (
	"fmt"

	approxmatch "github.com/ivanov-gv/approximate-match"
)

func main() {
	// Example from README: find the closest match for "aaaaa".
	fmt.Println("Sample: aaaaa")
	m := approxmatch.NewMatcher([]string{"aaaab", "aaaac", "abcde", "bbbbb"})
	for _, match := range m.Find("aaaaa") {
		fmt.Printf("  %-8s score: %.3f\n", match.Word, match.Score)
	}

	// Example from README: find the closest match for "Ella" among names.
	fmt.Println("\nSample: Ella")
	m = approxmatch.NewMatcher([]string{"Adele", "Elaine", "Elizabeth", "Harriet", "Ingrid", "Michelle", "Ella"})
	for _, match := range m.Find("Ella") {
		fmt.Printf("  %-12s score: %.3f\n", match.Word, match.Score)
	}
}
