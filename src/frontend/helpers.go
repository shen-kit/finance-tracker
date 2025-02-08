package frontend

import "regexp"

func isPartialDate(s string, r rune) bool {
	regex0 := regexp.MustCompile(`^\d{0,4}$`)
	regex1 := regexp.MustCompile(`^\d{4}-\d{0,2}$`)
	regex2 := regexp.MustCompile(`^\d{4}-\d{2}-\d{0,2}$`)
	return regex0.MatchString(s) || regex1.MatchString(s) || regex2.MatchString(s)
}
