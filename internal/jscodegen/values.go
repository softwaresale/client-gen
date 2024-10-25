package jscodegen

import (
	"fmt"
	"github.com/softwaresale/client-gen/v2/internal/types"
	"reflect"
	"strings"
)

type JSValueMapper struct{}

func (mapper JSValueMapper) Convert(value types.StaticValue) (string, error) {
	valueTp := reflect.TypeOf(value)
	switch valueTp.Kind() {
	case reflect.String:
		return fmt.Sprintf(`'%s'`, value), nil

	case reflect.Bool:
		if value.(bool) == true {
			return "true", nil
		}

		return "false", nil

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return fmt.Sprintf("%d", value), nil

	case reflect.Float32, reflect.Float64:
		return fmt.Sprintf("%f", value), nil

	case reflect.Slice:
		anySlice := value.([]any)
		builder := strings.Builder{}
		builder.WriteByte('[')
		sliceLen := len(anySlice)
		for idx, item := range anySlice {
			itemValue, err := mapper.Convert(item)
			if err != nil {
				return "", fmt.Errorf("failed to map slice value: %w", err)
			}

			builder.WriteString(itemValue)
			if idx < sliceLen-1 {
				builder.WriteString(", ")
			}
		}
		builder.WriteByte(']')

		return builder.String(), nil

	default:
		return "", fmt.Errorf("failed to map value: %v", valueTp.Kind())
	}
}
