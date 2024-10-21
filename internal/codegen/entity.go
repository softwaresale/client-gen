package codegen

// PropertySpec specifies an entity property
type PropertySpec struct {
	Type     DynamicType `json:"type"`     // Defines the type of this property
	Required bool        `json:"required"` // if true, this property must be specified. if false, can be an optional value
}

// EntitySpec specifies an entity model that is used
type EntitySpec struct {
	Name       string                  `json:"name"`       // name of this entity
	Properties map[string]PropertySpec `json:"properties"` // the properties that this entity defines
}

func (spec EntitySpec) IsValid() bool {
	return len(spec.Properties) > 0
}
