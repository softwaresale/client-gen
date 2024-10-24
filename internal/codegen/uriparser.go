package codegen

import (
	"fmt"
	"regexp"
	"strings"
)

// PathVariableMapper maps a template variable into some replacement string
type PathVariableMapper func(string) (string, error)

type URITemplate struct {
	Template   string             // The actual URI template
	VarMapper  PathVariableMapper // Determines how to swap out template variables
	PathPrefix string             // some text to prepend to the template URI
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

	// prepend the prefix if it's present
	if len(template.PathPrefix) > 0 {
		printTmpl := "%s/%s"

		// figure out if we need the slash
		if template.Template[0] == '/' || template.PathPrefix[len(template.PathPrefix)-1] == '/' {
			printTmpl = "%s%s"
		}

		template.Template = fmt.Sprintf(printTmpl, template.PathPrefix, template.Template)
	}

	return template.Template, nil
}
