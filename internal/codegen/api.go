package codegen

// APISpec specifies an entire API, which consists of multiple services
type APISpec struct {
	Name     string              `json:"name"`     // overall API name
	Entities []EntitySpec        `json:"entities"` // the entities needed to consume this API
	Services []ServiceDefinition `json:"services"` // the services provided by this API
}
