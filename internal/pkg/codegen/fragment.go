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
	HttpCall  HttpRequest
}

// Interface is a functional type that has functions that another type inherits
type Interface struct {
	Name      string
	Functions []FunctionSignature
}

// Class is a specification of a class with actual implementation
type Class struct {
	Name                string
	Interfaces          []DynamicType
	Properties          []VariableDecl
	InjectionProperties []VariableDecl
	Methods             []FunctionImpl
}

// HttpRequest models an http client call
type HttpRequest struct {
	Method          string
	UrlTemplate     string
	BaseEndpointVar string
	ClientVar       string
	InputVar        string
	RequestBody     DynamicType
	ResponseBody    DynamicType
}
