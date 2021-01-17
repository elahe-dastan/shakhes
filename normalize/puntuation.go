package normalize

import "strings"

var punctuations = []string{".", "،", ":", "؟", "!", "«", "»", "؛", "-", "…", "[", "]", "(", ")", "/", "=", "٪", "\"", "'", "+"}

// assume that the word can contain only one punctuation
func punctuation(word string) []string {
	words := split(word, punctuations)
	ans := make([]string, 0)
	for _, term := range words{
		if term != ""{
			ans = append(ans, term)
		}
	}

	return ans
}

func split(word string, puncts []string) []string{
	if len(puncts) == 0 {
		return nil
	}

	var words []string
	for i, p := range puncts {
		terms := strings.Split(word, p)
		if len(terms) > 1 {
			for _, term := range terms {
				if term != ""{
					words = append(words, split(term, puncts[i+1:])...)
				}
			}
			return words
		}
	}

	return []string{word}
}