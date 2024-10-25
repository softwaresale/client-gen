package types

import (
	"fmt"
	"reflect"
)

// APIConfig configures additional traits about this API.
type APIConfig struct {
	BaseURL string `json:"baseURL"` // Base URL of this API. All endpoints are relative to this endpoint
}

// CreateEntitySpec creates an entity spec that can represent our API config
func (apiConfig APIConfig) CreateEntitySpec() (EntitySpec, error) {

	entity := EntitySpec{
		Name:       "APIConfig",
		Properties: make(map[string]PropertySpec),
	}

	// reflect over the fields and get everything
	fieldCount := reflect.TypeOf(apiConfig).NumField()
	for i := 0; i < fieldCount; i++ {
		field := reflect.TypeOf(apiConfig).Field(i)
		// TODO it's kinda weird that i'm just using the json tag. Figure this out later...
		name := field.Tag.Get("json")
		tp, err := GoTypeToDynamicType(field.Type)
		if err != nil {
			return EntitySpec{}, fmt.Errorf("failed to map type for property '%s': %w", name, err)
		}

		entity.Properties[name] = PropertySpec{
			Type:     tp,
			Required: true,
		}
	}

	return entity, nil
}
