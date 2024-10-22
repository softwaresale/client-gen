package codegen

// ServiceDefinition defines a service that consumes a controller
type ServiceDefinition struct {
	Name      string        `json:"name"`      // defines the name of the service. Don't include any suffixes
	Endpoints []APIEndpoint `json:"endpoints"` // defines all endpoints defined by the controller
}
