package codegen

// ServiceDefinition defines a service that consumes a controller
type ServiceDefinition struct {
	Name      string        `json:"name"`
	Endpoints []APIEndpoint `json:"endpoints"`
}
