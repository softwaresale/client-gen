package jscodegen

import (
	"fmt"
	"regexp"
	"strings"
)

// PathVariableMapper maps a template variable into some replacement string
type PathVariableMapper func(string) (string, error)

// FormatTemplate takes an API endpoint URI template and expands it into a template string
func FormatTemplate(template string, mapper PathVariableMapper) (string, error) {
	parser := regexp.MustCompile(`\{\{\s*([a-zA-Z_][a-zA-Z0-9_]*)\s*}}`)

	for _, templateVar := range parser.FindAllStringSubmatch(template, -1) {
		wholeMatch := templateVar[0]
		variableName := templateVar[1]

		// map our variable into a replacement
		replacement, err := mapper(variableName)
		if err != nil {
			return "", fmt.Errorf("error while expanding variable '%s': %w", variableName, err)
		}

		// replace the wholeMatch with the replacement
		template = strings.Replace(template, wholeMatch, replacement, -1)
	}

	return template, nil
}
