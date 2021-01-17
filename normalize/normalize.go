package normalize

import (
	"log"
	"regexp"
	"strings"
)

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

func Number(words []string) []string {
	re, err := regexp.Compile("[۰-۹]+|[0-9]+")
	if err != nil{
		log.Fatal(err)
	}

	ans := make([]string, 0)
	for _, word := range words {
		if !re.MatchString(word){
			ans = append(ans, word)
		}
	}

	return ans
}

func Normalize(word string) []string {
	return Number(punctuation(word))
}