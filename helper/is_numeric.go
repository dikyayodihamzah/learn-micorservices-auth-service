package helper

import "regexp"

func IsNumeric(word string) bool {
	for _, c := range word {
		num := regexp.MustCompile(`\d`).MatchString(string(c))
		if !num {
			return num
		}
	}
	return true
}
