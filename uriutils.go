package discovery

import (
	"regexp"
	"strings"
)

func ParseURIs(content string) []string {
	re := regexp.MustCompile(`(http(s)?://)?localhost.*"`)
	matches := re.FindAllString(content, -1)

	var result []string = []string{}
	for _, match := range matches {
		trimmed := strings.Trim(match, `"`)
		result = append(result, trimmed)
	}

	return result
}
