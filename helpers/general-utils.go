package helpers

import (
	"strings"
)

func FindFileNames(input string, identifiers []string) []string {
	var result []string
	words := strings.Fields(input)
	for _, word := range words {
		for _, id := range identifiers {
			if strings.Contains(word, id) {
				result = append(result, word)
				break
			}
		}
	}
	return result
}
