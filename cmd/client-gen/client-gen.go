package main

import (
	"github.com/softwaresale/client-gen/v2/internal/pkg/codegen"
	"github.com/softwaresale/client-gen/v2/internal/pkg/jscodegen"
	"net/http"
	"os"
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

	typeCache := codegen.NewUserTypeCache()
	generator, _ := jscodegen.NewJSCodeGenerator(os.Stdout, typeCache)

	inputType, err := generator.GenerateInputType(endpoint)
	if err != nil {
		panic(err)
	}

	typeCache.CacheEndpointInputType(endpoint, *inputType)

	err = generator.GenerateServiceInterface(service)
	if err != nil {
		panic(err)
	}
}
