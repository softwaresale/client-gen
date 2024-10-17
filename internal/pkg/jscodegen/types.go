package jscodegen

import (
	"fmt"
	"github.com/softwaresale/client-gen/v2/internal/pkg/codegen"
)

type JSTypeMapper struct {
}

func (mapper JSTypeMapper) Convert(dtype codegen.DynamicType) (string, error) {
	var typeStr string
	switch dtype.TypeID {
	case codegen.TypeID_VOID:
		typeStr = "void"
	case codegen.TypeID_STRING:
		typeStr = "string"
	case codegen.TypeID_INTEGER:
		fallthrough
	case codegen.TypeID_FLOAT:
		typeStr = "number"
	case codegen.TypeID_BOOLEAN:
		typeStr = "boolean"
	case codegen.TypeID_USER:
		typeStr = dtype.Reference
	default:
		return "", fmt.Errorf("unknown type ID %s", dtype.TypeID)
	}

	return typeStr, nil
}
