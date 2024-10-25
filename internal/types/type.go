package types

import mapset "github.com/deckarep/golang-set/v2"

const (
	TypeID_VOID      = "VOID"
	TypeID_STRING    = "STRING"
	TypeID_INTEGER   = "INTEGER"
	TypeID_FLOAT     = "FLOAT"
	TypeID_BOOLEAN   = "BOOLEAN"
	TypeID_USER      = "USER"
	TypeID_ARRAY     = "ARRAY"
	TypeID_GENERIC   = "GENERIC"
	TypeID_TIMESTAMP = "TIMESTAMP"
	TypeID_ANY       = "ANY"
)

// DynamicType specifies a dynamic type that is specified by the user
type DynamicType struct {
	TypeID    string        `json:"typeID"`    // An identifier for this type. Comes from predefined enum
	Reference string        `json:"reference"` // Used by different types, references to this entity
	Inner     []DynamicType `json:"nested"`    // Related types, used by generics
}

func (tp DynamicType) IsVoid() bool {
	return tp.TypeID == TypeID_VOID
}

func (tp DynamicType) ArrayElementTp() DynamicType {
	if tp.TypeID != TypeID_ARRAY {
		panic("type is not an array")
	}

	if len(tp.Inner) < 1 {
		panic("array type does not have at least one inner type")
	}

	return tp.Inner[0]
}

func (tp DynamicType) TypeReferences() []string {
	references := mapset.NewSet[string](tp.Reference)
	for _, inner := range tp.Inner {
		inner := inner.TypeReferences()
		references.Append(inner...)
	}

	// remove the empty string -- we do not care about this.
	references.Remove("")

	return references.ToSlice()
}

// ITypeMapper provides an interface for mapping dynamic types into language-specific types.
type ITypeMapper interface {
	Convert(dtype DynamicType) (string, error)
}
