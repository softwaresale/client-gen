package types

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

func TestGoTypeToDynamicType_MapsString(t *testing.T) {
	dtype, err := GoTypeToDynamicType(reflect.TypeOf("hello"))
	assert.NoError(t, err)
	assert.Equal(t, TypeID_STRING, dtype.TypeID)
}

func TestGoTypeToDynamicType_MapsInt(t *testing.T) {
	dtype, err := GoTypeToDynamicType(reflect.TypeOf(42))
	assert.NoError(t, err)
	assert.Equal(t, TypeID_INTEGER, dtype.TypeID)
}

func TestGoTypeToDynamicType_MapsUint(t *testing.T) {
	dtype, err := GoTypeToDynamicType(reflect.TypeOf(uint(42)))
	assert.NoError(t, err)
	assert.Equal(t, TypeID_INTEGER, dtype.TypeID)
}

func TestGoTypeToDynamicType_MapsFloat(t *testing.T) {
	dtype, err := GoTypeToDynamicType(reflect.TypeOf(42.5))
	assert.NoError(t, err)
	assert.Equal(t, TypeID_FLOAT, dtype.TypeID)
}

func TestGoTypeToDynamicType_MapsBool(t *testing.T) {
	dtype, err := GoTypeToDynamicType(reflect.TypeOf(true))
	assert.NoError(t, err)
	assert.Equal(t, TypeID_BOOLEAN, dtype.TypeID)
}

func TestGoTypeToDynamicType_MapsSlice(t *testing.T) {
	dtype, err := GoTypeToDynamicType(reflect.TypeOf([]string{}))
	assert.NoError(t, err)
	assert.Equal(t, TypeID_ARRAY, dtype.TypeID)
	assert.NotEmpty(t, dtype.Inner)
	assert.Equal(t, TypeID_STRING, dtype.Inner[0].TypeID)
}
