package jscodegen

import (
	"fmt"
	"github.com/softwaresale/client-gen/v2/internal/pkg/codegen"
	"strings"
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
	case codegen.TypeID_ARRAY:
		// get the inner type
		innerTypeStr, err := mapper.Convert(dtype.ArrayElementTp())
		if err != nil {
			return "", fmt.Errorf("failed to map array inner type: %w", err)
		}
		typeStr = fmt.Sprintf("%s[]", innerTypeStr)

	case codegen.TypeID_GENERIC:
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

func WellKnownTypeResolver(typeRef string) (string, error) {
	switch typeRef {
	case "Observable":
		return "rxjs", nil
	case "HttpClient":
		return "@angular/common/http", nil
	default:
		return "", fmt.Errorf("type %s is not well-known", typeRef)
	}
}
