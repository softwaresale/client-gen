package jscodegen

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var mapper JSValueMapper

func setup(t *testing.T) {
	mapper = JSValueMapper{}
}

func TestJSValueMapper_Convert_MapString(t *testing.T) {
	setup(t)
	result, err := mapper.Convert("Hello World")
	assert.Nil(t, err)
	assert.Equal(t, "'Hello World'", result)
}

func TestJSValueMapper_Convert_MapInt(t *testing.T) {
	setup(t)
	result, err := mapper.Convert(1)
	assert.Nil(t, err)
	assert.Equal(t, "1", result)
}

func TestJSValueMapper_Convert_MapFloat(t *testing.T) {
	setup(t)
	result, err := mapper.Convert(1.1)
	assert.Nil(t, err)
	assert.Equal(t, "1.100000", result)
}

func TestJSValueMapper_Convert_MapBool(t *testing.T) {
	setup(t)
	result, err := mapper.Convert(true)
	assert.Nil(t, err)
	assert.Equal(t, "true", result)
}

/* array is not yet supported
func TestJSValueMapper_Convert_MapSlice(t *testing.T) {
	setup(t)
	result, err := mapper.Convert([]string{"Hello", "World"})
	assert.Nil(t, err)
	assert.Equal(t, "['Hello', 'World]", result)
}
*/
