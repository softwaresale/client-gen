package codegen

// ISourceGenerator specifies an interface that all source code generators implement.
type ISourceGenerator interface {
	GenerateServiceInterface(serviceDef ServiceDefinition) error
	GenerateService(serviceDef ServiceDefinition) error
	GenerateEndpoint(endpointDef APIEndpoint) error
	GenerateInputType(endpoint APIEndpoint) error
}
