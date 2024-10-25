package jscodegen

import (
	"fmt"
	"github.com/softwaresale/client-gen/v2/internal/types"
	"strings"
)

type JSTypeMapper struct {
}

func (mapper JSTypeMapper) Convert(dtype types.DynamicType) (string, error) {
	var typeStr string
	switch dtype.TypeID {
	case types.TypeID_VOID:
		typeStr = "void"
	case types.TypeID_STRING:
		typeStr = "string"
	case types.TypeID_INTEGER:
		fallthrough
	case types.TypeID_FLOAT:
		typeStr = "number"
	case types.TypeID_BOOLEAN:
		typeStr = "boolean"
	case types.TypeID_USER:
		typeStr = dtype.Reference
	case types.TypeID_TIMESTAMP:
		typeStr = "Date"
	case types.TypeID_ANY:
		typeStr = "any"
	case types.TypeID_ARRAY:
		// get the inner type
		innerTypeStr, err := mapper.Convert(dtype.ArrayElementTp())
		if err != nil {
			return "", fmt.Errorf("failed to map array inner type: %w", err)
		}
		typeStr = fmt.Sprintf("%s[]", innerTypeStr)

	case types.TypeID_GENERIC:
		var genericParams []string
		for genericIdx, inner := range dtype.Inner {
			innerTypeStr, err := mapper.Convert(inner)
			if err != nil {
				return "", fmt.Errorf("failed to map generic inner type at index %d: %w", genericIdx, err)
			}

			genericParams = append(genericParams, innerTypeStr)
		}
		typeStr = fmt.Sprintf("%s<%s>", dtype.Reference, strings.Join(genericParams, ", "))

	default:
		return "", fmt.Errorf("unknown type ID %s", dtype.TypeID)
	}

	return typeStr, nil
}
