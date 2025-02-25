package main

import (
	"regexp"
	"strings"
)

var (
	step1 = regexp.MustCompile(`\s+`)
	step2 = regexp.MustCompile(`([a-z])([A-Z])`)
	step3 = regexp.MustCompile(`([A-Z])([A-Z][a-z])`)
)

func toSnakeCase(input string) string {
	// Step 1: Replace spaces with underscores
	spacesReplaced := step1.ReplaceAllString(input, "_")

	// Step 2: Add underscore before uppercase letters that follow lowercase letters
	// This handles camelCase -> camel_Case
	withUnderscores := step2.ReplaceAllString(spacesReplaced, "${1}_${2}")

	// Step 3: Handle consecutive uppercase letters properly
	// This handles ABCWord -> ABC_Word
	final := step3.ReplaceAllString(withUnderscores, "${1}_${2}")

	// Step 4: Convert to lowercase
	return strings.ToLower(final)
}
