package codegen

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestDynamicType_IsVoid(t *testing.T) {
	tp := DynamicType{TypeID: TypeID_VOID}
	assert.True(t, tp.IsVoid())
}

func TestDynamicType_TypeReferences_GetsAllReferences(t *testing.T) {
	tp := DynamicType{
		TypeID:    TypeID_GENERIC,
		Reference: "Container",
		Inner: []DynamicType{
			{
				TypeID: TypeID_INTEGER,
			},
			{
				TypeID:    TypeID_USER,
				Reference: "Magic",
			},
			{
				TypeID:    TypeID_ARRAY,
				Reference: "",
				Inner: []DynamicType{
					{
						TypeID:    TypeID_USER,
						Reference: "World",
					},
				},
			},
		},
	}

	referencedTypes := tp.TypeReferences()

	for _, expected := range []string{"Container", "Magic", "World"} {
		assert.Contains(t, referencedTypes, expected)
	}
}

func TestDynamicType_ArrayElementTp_GetsElementType(t *testing.T) {
	expected := DynamicType{
		TypeID: TypeID_INTEGER,
	}

	tp := DynamicType{
		TypeID:    TypeID_ARRAY,
		Reference: "",
		Inner: []DynamicType{
			expected,
		},
	}

	elementTp := tp.ArrayElementTp()
	assert.Equal(t, expected, elementTp)
}

func TestDynamicType_ArrayElementTp_PanicsForNonArray(t *testing.T) {
	tp := DynamicType{
		TypeID: TypeID_INTEGER,
	}

	assert.Panics(t, func() {
		tp.ArrayElementTp()
	})
}

func TestDynamicType_ArrayElementTp_PanicsForNoInnerType(t *testing.T) {
	tp := DynamicType{
		TypeID: TypeID_ARRAY,
	}

	assert.Panics(t, func() {
		tp.ArrayElementTp()
	})
}

func TestMapTypeToDynamicType_MapBoolType(t *testing.T) {
	mappedTp, err := MapTypeToDynamicType(reflect.TypeOf(true))
	assert.NoError(t, err)
	assert.Equal(t, TypeID_BOOLEAN, mappedTp.TypeID)
}

func TestGetDynamicTypeForValue_MapBoolType(t *testing.T) {
	mappedTp, err := GetDynamicTypeForValue(true)
	assert.NoError(t, err)
	assert.Equal(t, TypeID_BOOLEAN, mappedTp.TypeID)
}

func TestMapTypeToDynamicType_MapIntType(t *testing.T) {
	var v int = 1
	mappedTp, err := MapTypeToDynamicType(reflect.TypeOf(v))
	assert.NoError(t, err)
	assert.Equal(t, TypeID_INTEGER, mappedTp.TypeID)
}

func TestGetDynamicTypeForValue_MapIntType(t *testing.T) {
	var v int = 1
	mappedTp, err := GetDynamicTypeForValue(v)
	assert.NoError(t, err)
	assert.Equal(t, TypeID_INTEGER, mappedTp.TypeID)
}

func TestMapTypeToDynamicType_MapFloatType(t *testing.T) {
	var v float32 = 1.4
	mappedTp, err := MapTypeToDynamicType(reflect.TypeOf(v))
	assert.NoError(t, err)
	assert.Equal(t, TypeID_FLOAT, mappedTp.TypeID)
}

func TestGetDynamicTypeForValue_MapFloatType(t *testing.T) {
	var v float32 = 1.4
	mappedTp, err := GetDynamicTypeForValue(v)
	assert.NoError(t, err)
	assert.Equal(t, TypeID_FLOAT, mappedTp.TypeID)
}

func TestMapTypeToDynamicType_MapStringType(t *testing.T) {
	v := "hello"
	mappedTp, err := MapTypeToDynamicType(reflect.TypeOf(v))
	assert.NoError(t, err)
	assert.Equal(t, TypeID_STRING, mappedTp.TypeID)
}

func TestGetDynamicTypeForValue_MapStringType(t *testing.T) {
	v := "hello"
	mappedTp, err := GetDynamicTypeForValue(v)
	assert.NoError(t, err)
	assert.Equal(t, TypeID_STRING, mappedTp.TypeID)
}

func TestMapTypeToDynamicType_MapStringArrayType(t *testing.T) {
	v := []string{"hello", "world"}
	mappedTp, err := MapTypeToDynamicType(reflect.TypeOf(v))
	assert.NoError(t, err)
	assert.Equal(t, TypeID_ARRAY, mappedTp.TypeID)
	assert.NotEmpty(t, mappedTp.Inner)
	assert.Equal(t, TypeID_STRING, mappedTp.Inner[0].TypeID)
}

func TestGetDynamicTypeForValue_MapStringArrayType(t *testing.T) {
	v := []string{"hello", "world"}
	mappedTp, err := GetDynamicTypeForValue(v)
	assert.NoError(t, err)
	assert.Equal(t, TypeID_ARRAY, mappedTp.TypeID)
	assert.NotEmpty(t, mappedTp.Inner)
	assert.Equal(t, TypeID_STRING, mappedTp.Inner[0].TypeID)
}
