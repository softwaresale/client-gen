package codegen

import (
	"fmt"
	"regexp"
	"strings"
)

// PathVariableMapper maps a template variable into some replacement string
type PathVariableMapper func(string) (string, error)

type URITemplate struct {
	Template  string             // The URI template
	VarMapper PathVariableMapper // Strategy for mapping path variables to replacements
	Prefix    string             // a prefix string to append to beginning of resolved string
}

// FormatTemplate takes an API endpoint URI template and expands it into a template string
func FormatTemplate(template URITemplate) (string, error) {
	parser := regexp.MustCompile(`\{\{\s*([a-zA-Z_][a-zA-Z0-9_]*)\s*}}`)

	for _, templateVar := range parser.FindAllStringSubmatch(template.Template, -1) {
		wholeMatch := templateVar[0]
		variableName := templateVar[1]

		// map our variable into a replacement
		replacement, err := template.VarMapper(variableName)
		if err != nil {
			return "", fmt.Errorf("error while expanding variable '%s': %w", variableName, err)
		}

		// replace the wholeMatch with the replacement
		template.Template = strings.Replace(template.Template, wholeMatch, replacement, -1)
	}

	// if there is a prefix, make sure that it's there
	if len(template.Prefix) > 0 {
		if template.Template[0] == '/' && template.Prefix[len(template.Prefix)-1] == '/' {
			// if both have slash, remove a slash
			template.Template = fmt.Sprintf("%s%s", template.Prefix, template.Template[1:])
		} else if template.Template[0] != '/' && template.Prefix[len(template.Prefix)-1] != '/' {
			// if neither have slash, add
			template.Template = fmt.Sprintf("%s/%s", template.Prefix, template.Template)
		} else {
			// one of them has a slash, so don't include
			template.Template = fmt.Sprintf("%s%s", template.Prefix, template.Template)
		}
	}

	return template.Template, nil
}
