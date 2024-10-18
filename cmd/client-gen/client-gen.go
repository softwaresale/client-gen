package main

import (
	"github.com/softwaresale/client-gen/v2/internal/pkg/codegen"
	"github.com/softwaresale/client-gen/v2/internal/pkg/jscodegen"
	"net/http"
	"os"
)

func main() {
	service := codegen.ServiceDefinition{
		Name: "ModelsService",
		Endpoints: []codegen.APIEndpoint{
			{
				Name:     "getModels",
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
			},
			{
				Name:     "createModel",
				Endpoint: "/api/v1/person",
				Method:   http.MethodPost,
				RequestBody: codegen.RequestValue{
					Type: codegen.DynamicType{
						TypeID:    codegen.TypeID_USER,
						Reference: "PersonModel",
					},
					Required: true,
				},
				ResponseBody: codegen.RequestValue{
					Type: codegen.DynamicType{
						TypeID:    codegen.TypeID_USER,
						Reference: "PersonModel",
					},
				},
			},
		},
	}

	serviceCompiler := codegen.ServiceCompiler{}

	compiledService, err := serviceCompiler.CompileService(service)
	if err != nil {
		panic(err)
	}

	formatter := jscodegen.NewJSCodeFormatter(os.Stdout, jscodegen.JSTypeMapper{})

	err = formatter.Format(*compiledService)
	if err != nil {
		panic(err)
	}
}
