package main

import (
	"fmt"
	codegen2 "github.com/softwaresale/client-gen/v2/internal/codegen"
	jscodegen2 "github.com/softwaresale/client-gen/v2/internal/jscodegen"
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
						Required: true,
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
						TypeID:    codegen2.TypeID_STRING,
						Reference: "",
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

	definedTypes := map[string]string{
		"PersonModel": "lib/common",
	}

	resolver := codegen2.ComposePackageResolver(
		jscodegen2.WellKnownTypeResolver,
		func(typeRef string) (string, error) {
			pkg, exists := definedTypes[typeRef]
			if !exists {
				return "", fmt.Errorf("failed to resolve type '%s'", typeRef)
			}

			return pkg, nil
		},
	)

	typeManager := codegen2.TypeManager{
		Filter: func(typeRef string) bool {
			if len(typeRef) == 0 {
				return false
			}

			_, err := resolver(typeRef)
			return err == nil
		},
		PkgResolver: resolver,
	}

	serviceCompiler := codegen2.ServiceCompiler{
		TypeManager: typeManager,
	}

	compiledService, err := serviceCompiler.CompileService(service)
	if err != nil {
		panic(err)
	}

	formatter := jscodegen2.NewJSCodeFormatter(os.Stdout, jscodegen2.JSTypeMapper{})

	err = formatter.Format(*compiledService)
	if err != nil {
		panic(err)
	}
}
