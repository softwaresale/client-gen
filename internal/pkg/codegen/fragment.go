package codegen

// TypeDecl is a fragment that represents a type specifier
type TypeDecl struct {
	Type     DynamicType
	Required bool
}

// VariableDecl is a variable, which has a name and a type
type VariableDecl struct {
	Name string
	Type TypeDecl
}

// Record is a simple type that holds key-value properties
type Record struct {
	Name      string
	Variables map[string]VariableDecl
}

// FunctionSignature describes a general function/method signature
type FunctionSignature struct {
	Name string
	Ret  DynamicType
	Args []VariableDecl
}

// FunctionImpl is a function with an actual implementation
type FunctionImpl struct {
	Signature FunctionSignature
	// TODO body
}

// Interface is a functional type that has functions that another type inherits
type Interface struct {
	Name      string
	Functions []FunctionSignature
}

// Class is a specification of a class with actual implementation
type Class struct {
	Name                string
	Properties          []VariableDecl
	InjectionProperties []VariableDecl
	Methods             []FunctionImpl
}
