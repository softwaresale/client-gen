package codegen

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFormatTemplate_Works(t *testing.T) {
	tmpl := URITemplate{
		Template: "/hello/{{world}}",
		VarMapper: func(variable string) (string, error) {
			return variable, nil
		},
	}

	expected := "/hello/world"
	formatted, err := FormatTemplate(tmpl)
	assert.NoError(t, err)

	assert.Equal(t, expected, formatted)
}

func TestFormatTemplate_ShouldFailWhenMapperFails(t *testing.T) {
	tmpl := URITemplate{
		Template: "/hello/{{world}}",
		VarMapper: func(variable string) (string, error) {
			return "", fmt.Errorf("error")
		},
	}

	_, err := FormatTemplate(tmpl)
	assert.Error(t, err)
}

func TestFormatTemplate_PrependPrefixWhenPresent(t *testing.T) {
	tmpl := URITemplate{
		Template:   "/hello/{{world}}",
		VarMapper:  func(variable string) (string, error) { return variable, nil },
		PathPrefix: "prefix",
	}

	actual, err := FormatTemplate(tmpl)
	assert.NoError(t, err)

	assert.Equal(t, "prefix/hello/world", actual)
}

func TestFormatTemplate_PrependPrefixWhenPresent_WithNoLeadingSlash(t *testing.T) {
	tmpl := URITemplate{
		Template:   "hello/{{world}}",
		VarMapper:  func(variable string) (string, error) { return variable, nil },
		PathPrefix: "prefix",
	}

	actual, err := FormatTemplate(tmpl)
	assert.NoError(t, err)

	assert.Equal(t, "prefix/hello/world", actual)
}

func TestFormatTemplate_PrependPrefixWhenPresent_WithPrefixSlash(t *testing.T) {
	tmpl := URITemplate{
		Template:   "hello/{{world}}",
		VarMapper:  func(variable string) (string, error) { return variable, nil },
		PathPrefix: "prefix/",
	}

	actual, err := FormatTemplate(tmpl)
	assert.NoError(t, err)

	assert.Equal(t, "prefix/hello/world", actual)
}
