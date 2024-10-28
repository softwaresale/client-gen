package types

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

// ITypeMapper provides an interface for mapping dynamic types into language-specific types.
type ITypeMapper interface {
	Convert(dtype DynamicType) (string, error)
}

// GoValueToDynamicType gets the type of the value via reflection and finds its type using GoValueToDynamicType
func GoValueToDynamicType(goVal any) (DynamicType, error) {
	tp := reflect.TypeOf(goVal)
	return GoTypeToDynamicType(tp)
}

// GoTypeToDynamicType maps a reflected go type to a dynamic type
func GoTypeToDynamicType(goTp reflect.Type) (DynamicType, error) {

	var dtype DynamicType

	switch goTp.Kind() {
	case reflect.String:
		dtype.TypeID = TypeID_STRING

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		dtype.TypeID = TypeID_INTEGER

	case reflect.Float32, reflect.Float64:
		dtype.TypeID = TypeID_FLOAT

	case reflect.Bool:
		dtype.TypeID = TypeID_BOOLEAN

	case reflect.Slice:
		dtype.TypeID = TypeID_ARRAY
		// figure out what the inner type is
		elemType := goTp.Elem()
		innerTp, err := GoTypeToDynamicType(elemType)
		if err != nil {
			return DynamicType{}, fmt.Errorf("failed to map slice element type '%s': %w", elemType.String(), err)
		}
		dtype.Inner = append(dtype.Inner, innerTp)

	default:
		return dtype, fmt.Errorf("unsupported type '%s'", goTp.Kind().String())
	}

	return dtype, nil
}
