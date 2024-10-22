package codegen

// RequestValue specifies a value that is passed in an API endpoint. Each value is typed and has optional
// metadata
type RequestValue struct {
	Type     DynamicType `json:"type"`
	Required bool        `json:"required"`
}

// APIEndpoint is an endpoint to call
type APIEndpoint struct {
	Name           string                  `json:"name"`           // the name of the endpoint
	Endpoint       string                  `json:"endpoint"`       // the URI endpoint that this request is located at
	Method         string                  `json:"method"`         // the HTTP method that this endpoint consumes
	PathVariables  map[string]RequestValue `json:"pathVariables"`  // a map of variables that are contained in the URI
	RequestBody    RequestValue            `json:"requestBody"`    // the request attached to the body
	ResponseBody   RequestValue            `json:"responseBody"`   // the type of the response body
	QueryVariables map[string]RequestValue `json:"queryVariables"` // additional query variables append to URI
}
