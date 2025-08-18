package utils

import "strings"

func GenerateNameAlias(name string) string {
	words := strings.Fields(name) // split by space, auto skip multiple spaces
	for i, w := range words {
		if len(w) > 1 {
			words[i] = string(w[0]) + strings.Repeat("*", len(w)-1)
		} else {
			words[i] = w
		}
	}
	return strings.Join(words, " ")
}
