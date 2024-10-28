package types

// StaticValue encodes a statically-known, compile-time go value. This value can be mapped into target constants
// as well. Use this for static initializers in target code.
type StaticValue any

// EntityInitializer describes how we can initialize an entity with compile time constants
type EntityInitializer struct {
	PropertyValues map[string]StaticValue
}

// ValueMapper provides an interface for converting a StaticValue into target code
type ValueMapper interface {
	Convert(value StaticValue) (string, error)
}
