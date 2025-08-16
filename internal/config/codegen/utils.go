package main

import (
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

const endpointDelimiter string = "-"

func toCamelCase(input string) string {
	titleCaser := cases.Title(language.English)
	parts := strings.FieldsFunc(input, func(r rune) bool {
		return r == '_' || r == '-' || r == '.'
	})

	var result strings.Builder
	for _, part := range parts {
		result.WriteString(titleCaser.String(part))
	}
	return result.String()
}

func toKeyCase(input string, delimiters ...string) string {
	upperCaser := cases.Upper(language.English)

	result := input
	for _, delimiter := range delimiters {
		result = strings.ReplaceAll(result, delimiter, "_")
	}

	return upperCaser.String(result)
}
