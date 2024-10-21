package codegen

// PropertySpec specifies an entity property
type PropertySpec struct {
	Type     DynamicType // Defines the type of this property
	Required bool        // if true, this property must be specified. if false, can be an optional value
}

// EntitySpec specifies an entity model that is used
type EntitySpec struct {
	Name       string                  // name of this entity
	Properties map[string]PropertySpec // the properties that this entity defines
}

func (spec EntitySpec) IsValid() bool {
	return len(spec.Properties) > 0
}
