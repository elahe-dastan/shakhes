package normalize

import (
	"strings"
)

var punctuations = []string{".", "،", ":", "؟", "!", "«", "»", "؛", "-", "…", "[", "]", "(", ")", "/", "=", "٪"}

// assume that the word can contain only one punctuation
func punctuation(word string) []string {
	var words []string
	for _, p := range punctuations {
		words = strings.Split(word, p)
		if len(words) > 1 {
			break
		}
	}

	ans := make([]string, 0)
	for _, term := range words{
		if term != ""{
			ans = append(ans, term)
		}
	}

	return ans
}

func Normalize(word string) []string {
	return punctuation(word)
}