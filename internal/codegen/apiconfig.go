package codegen

// APIConfig represents configuration data for the API. This map serves two purpose:
//  1. It provides literal values that can be used for base configuration
//  2. a type is created to represent the configuration data, so additional configurations can be provided inside
//     the application or from external sources
type APIConfig struct {
	BaseURL string `json:"baseURL"`
}

// CreateConfigEntity introspects our config object with reflection and creates an entity that store it
func (config APIConfig) CreateConfigEntity() (EntitySpec, error) {
	entitySpec := EntitySpec{
		Name: "APIConfig",
		Properties: map[string]PropertySpec{
			"baseURL": {
				Type:     DynamicType{TypeID: TypeID_STRING},
				Required: true,
			},
		},
	}

	// create an initializer

	return entitySpec, nil
}
