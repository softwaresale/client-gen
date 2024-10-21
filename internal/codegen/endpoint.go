package codegen

// RequestValue specifies a value that is passed in an API endpoint. Each value is typed and has optional
// metadata
type RequestValue struct {
	Type     DynamicType `json:"type"`
	Required bool        `json:"required"`
}

// APIEndpoint is an endpoint to call
type APIEndpoint struct {
	Name           string                  `json:"name"`
	Endpoint       string                  `json:"endpoint"`
	Method         string                  `json:"method"`
	PathVariables  map[string]RequestValue `json:"pathVariables"`
	RequestBody    RequestValue            `json:"requestBody"`
	ResponseBody   RequestValue            `json:"responseBody"`
	QueryVariables map[string]RequestValue `json:"queryVariables"`
}
