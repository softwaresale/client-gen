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
