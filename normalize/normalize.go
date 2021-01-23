package normalize

import (
	"log"
	"regexp"
)

func number(words []string) []string {
	re, err := regexp.Compile("[۰-۹]+|[0-9]+")
	if err != nil {
		log.Fatal(err)
	}

	ans := make([]string, 0)
	for _, word := range words {
		if !re.MatchString(word) {
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

// ZWNJ and ZWSP
func zeroWidth(words []string) []string {
	ans := make([]string, 0)
	for _, word := range words {
		normalizedWord := make([]rune, 0)
		for _, r := range []rune(word) {
			if r != '‌' && r != '​' && r != '­' {
				normalizedWord = append(normalizedWord, r)
			}
		}
		ans = append(ans, string(normalizedWord))
	}

	return ans
}

func ha(words []string) []string {
	ans := make([]string, 0)
	for _, word := range words {
		// ها
		runes := []rune(word)
		if string(runes[len(runes)-2:]) == "ها" {
			if string(runes[:len(runes)-2]) != "تن" {
				word = string(runes[:len(runes)-2])
			}
		}
		if word != "" {
			ans = append(ans, word)
		}
	}
	return ans
}

func tarin(words []string) []string {
	ans := make([]string, 0)
	for _, word := range words {
		// ترین
		runes := []rune(word)
		if len(runes) < 4 {
			continue
		}
		if string(runes[len(runes)-4:]) == "ترین" {
			word = string(runes[:len(runes)-4])

		}
		if word != "" {
			ans = append(ans, word)
		}
	}
	//	ends = ['ات', 'ان', 'ترین', 'تر', 'م', 'ت', 'ش', 'یی', 'ی', 'ها', 'ٔ', '‌ا',]
	return ans
}

func Normalize(word string) []string {
	return tarin(ha(singleChar(number(zeroWidth(punctuation(word))))))
}
