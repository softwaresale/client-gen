package types

import (
	"github.com/stretchr/testify/assert"
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
