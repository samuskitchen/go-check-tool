package check_tool

import (
	"log"
	"regexp"
	"strings"
)

// SplitCamelCase receives a camelcase text and separates it into spaces
func SplitCamelCase(s string) string {
	for _, reStr := range []string{`([A-Z]+)([A-Z][a-z])`, `([a-z\d])([A-Z])`} {
		re := regexp.MustCompile(reStr)
		s = re.ReplaceAllString(s, "${1} ${2}")
	}

	return s
}

// StandardSpace remove unnecessary spaces between words
func StandardSpace(s string) string {
	for strings.Contains(s, "  ") {
		s = strings.ReplaceAll(s, "  ", " ")
	}

	return strings.TrimSpace(s)
}

// SplitKeyValue receives a text string `len=89`, separated by a `=` this function returns (valid, key, value)
func SplitKeyValue(s string) (bool, string, string) {
	ok, err := regexp.MatchString("^[a-zA-Z]+=(.)+$", s)
	if err != nil {
		log.Println("ERROR:", err)
		return false, "", ""
	}

	if !ok {
		return false, "", ""
	}

	values := strings.Split(s, "=")
	return true, values[0], values[1]
}
