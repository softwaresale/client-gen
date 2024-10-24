package codegen

import (
	"fmt"
	mapset "github.com/deckarep/golang-set/v2"
	"reflect"
)

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

// MapTypeToDynamicType converts a Go type to a DynamicType
func MapTypeToDynamicType(goTp reflect.Type) (DynamicType, error) {
	var tp DynamicType

	switch goTp.Kind() {
	case reflect.Bool:
		tp.TypeID = TypeID_BOOLEAN

	case reflect.String:
		tp.TypeID = TypeID_STRING

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		tp.TypeID = TypeID_INTEGER

	case reflect.Float32, reflect.Float64:
		tp.TypeID = TypeID_FLOAT

	case reflect.Slice, reflect.Array:
		tp.TypeID = TypeID_ARRAY
		inner, err := MapTypeToDynamicType(goTp.Elem())
		if err != nil {
			return DynamicType{}, err
		}
		tp.Inner = []DynamicType{inner}

	case reflect.Struct:
		return DynamicType{}, fmt.Errorf("structures are not supported as types yet")

	default:
		return DynamicType{}, fmt.Errorf("unsupported type %s", goTp.Kind().String())
	}

	return tp, nil
}

// GetDynamicTypeForValue maps a go value into a DynamicType
func GetDynamicTypeForValue(value any) (DynamicType, error) {
	valueTp := reflect.TypeOf(value)
	return MapTypeToDynamicType(valueTp)
}

// ITypeMapper provides an interface for mapping dynamic types into language-specific types.
type ITypeMapper interface {
	Convert(dtype DynamicType) (string, error)
}
