package main

import (
	"fmt"
	"github.com/softwaresale/client-gen/v2/internal/pkg/codegen"
	"net/http"
)

func main() {
	endpoint := codegen.APIEndpoint{
		Name:     "GetModels",
		Endpoint: "/api/v1/person/{{age}}/name/{{name}}/somethingelse",
		Method:   http.MethodGet,
		PathVariables: map[string]codegen.RequestValue{
			"age": {
				Type: codegen.DynamicType{
					TypeID: codegen.TypeID_STRING,
				},
				Required: true,
			},
			"name": {
				Type: codegen.DynamicType{
					TypeID: codegen.TypeID_STRING,
				},
				Required: true,
			},
		},
		RequestBody: codegen.RequestValue{
			Type: codegen.DynamicType{
				TypeID: codegen.TypeID_VOID,
			},
			Required: false,
		},
		ResponseBody: codegen.RequestValue{
			Type: codegen.DynamicType{
				TypeID:    codegen.TypeID_STRING,
				Reference: "",
			},
		},
		QueryVariables: nil,
	}

	service := codegen.ServiceDefinition{
		Name:    "ModelsService",
		BaseURL: "/api/v1/models",
		Endpoints: []codegen.APIEndpoint{
			endpoint,
		},
	}

	serviceCompiler := codegen.ServiceCompiler{}

	compiledService := serviceCompiler.CompileService(service)

	fmt.Printf("%v\n", compiledService)
}
