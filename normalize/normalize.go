package normalize

import (
	"log"
	"regexp"
)

func number(words []string) []string {
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

func singleChar(words []string) []string {
	ans := make([]string, 0)
	for _, word := range words {
		if len([]rune(word)) > 1 {
			ans = append(ans, word)
		}
	}

	return ans
}

func zeroWidthNonJoiner(words []string) []string {
	ans := make([]string, 0)
	for _, word := range words {
		normalizedWord := make([]rune, 0)
		for _, r := range []rune(word){
			if r != '‌'{
				normalizedWord = append(normalizedWord, r)
			}
		}
		ans = append(ans, string(normalizedWord))
	}

	return ans
}



func Normalize(word string) []string {
	return singleChar(number(zeroWidthNonJoiner(punctuation(word))))
}