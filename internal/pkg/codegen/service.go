package codegen

// ServiceDefinition defines a service that consumes a controller
type ServiceDefinition struct {
	Name      string        `json:"name"`
	BaseURL   string        `json:"baseURL"`
	Endpoints []APIEndpoint `json:"endpoints"`
}
