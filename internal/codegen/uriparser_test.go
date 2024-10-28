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

func TestFormatTemplate_ShouldPrependPrefix(t *testing.T) {
	tmpl := URITemplate{
		Template:  "/hello/{{world}}",
		VarMapper: func(variable string) (string, error) { return variable, nil },
		Prefix:    "/prefix",
	}

	formatted, err := FormatTemplate(tmpl)
	assert.NoError(t, err)

	assert.Equal(t, "/prefix/hello/world", formatted)
}

func TestFormatTemplate_ShouldPrependPrefix_AccountingForSlash1(t *testing.T) {
	tmpl := URITemplate{
		Template:  "hello/{{world}}",
		VarMapper: func(variable string) (string, error) { return variable, nil },
		Prefix:    "/prefix",
	}

	formatted, err := FormatTemplate(tmpl)
	assert.NoError(t, err)

	assert.Equal(t, "/prefix/hello/world", formatted)
}

func TestFormatTemplate_ShouldPrependPrefix_AccountingForSlash2(t *testing.T) {
	tmpl := URITemplate{
		Template:  "hello/{{world}}",
		VarMapper: func(variable string) (string, error) { return variable, nil },
		Prefix:    "/prefix/",
	}

	formatted, err := FormatTemplate(tmpl)
	assert.NoError(t, err)

	assert.Equal(t, "/prefix/hello/world", formatted)
}

func TestFormatTemplate_ShouldPrependPrefix_AccountingForSlash3(t *testing.T) {
	tmpl := URITemplate{
		Template:  "/hello/{{world}}",
		VarMapper: func(variable string) (string, error) { return variable, nil },
		Prefix:    "/prefix/",
	}

	formatted, err := FormatTemplate(tmpl)
	assert.NoError(t, err)

	assert.Equal(t, "/prefix/hello/world", formatted)
}
