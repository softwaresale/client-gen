package main

import (
	codegen2 "github.com/softwaresale/client-gen/v2/internal/codegen"
	"github.com/softwaresale/client-gen/v2/internal/templates"
	"net/http"
	"os"
)

func main() {
	service := codegen2.ServiceDefinition{
		Name: "ModelsService",
		Endpoints: []codegen2.APIEndpoint{
			{
				Name:     "getModels",
				Endpoint: "/api/v1/person/{{age}}/name/{{name}}/somethingelse",
				Method:   http.MethodGet,
				PathVariables: map[string]codegen2.RequestValue{
					"age": {
						Type: codegen2.DynamicType{
							TypeID: codegen2.TypeID_STRING,
						},
						Required: false,
					},
					"name": {
						Type: codegen2.DynamicType{
							TypeID: codegen2.TypeID_STRING,
						},
						Required: true,
					},
				},
				RequestBody: codegen2.RequestValue{
					Type: codegen2.DynamicType{
						TypeID: codegen2.TypeID_VOID,
					},
					Required: false,
				},
				ResponseBody: codegen2.RequestValue{
					Type: codegen2.DynamicType{
						TypeID:    codegen2.TypeID_ARRAY,
						Reference: "",
						Inner: []codegen2.DynamicType{
							{
								TypeID:    codegen2.TypeID_USER,
								Reference: "PersonModel",
							},
						},
					},
				},
			},
			{
				Name:     "createModel",
				Endpoint: "/api/v1/person",
				Method:   http.MethodPost,
				RequestBody: codegen2.RequestValue{
					Type: codegen2.DynamicType{
						TypeID:    codegen2.TypeID_USER,
						Reference: "PersonModel",
					},
					Required: true,
				},
				ResponseBody: codegen2.RequestValue{
					Type: codegen2.DynamicType{
						TypeID:    codegen2.TypeID_USER,
						Reference: "PersonModel",
					},
				},
			},
			{
				Name:     "updateModel",
				Endpoint: "/api/v1/person/{{id}}",
				Method:   http.MethodPut,
				PathVariables: map[string]codegen2.RequestValue{
					"id": {
						Type: codegen2.DynamicType{
							TypeID: codegen2.TypeID_STRING,
						},
						Required: true,
					},
				},
				RequestBody: codegen2.RequestValue{
					Type: codegen2.DynamicType{
						TypeID:    codegen2.TypeID_USER,
						Reference: "PersonModel",
					},
					Required: true,
				},
				ResponseBody: codegen2.RequestValue{
					Type: codegen2.DynamicType{
						TypeID:    codegen2.TypeID_USER,
						Reference: "PersonModel",
					},
					Required: true,
				},
			},
		},
	}

	ngServiceGen := templates.NewNGServiceGenerator()
	err := ngServiceGen.Generate(os.Stdout, service)
	if err != nil {
		panic(err)
	}
}
