package main

import (
	"fmt"
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
			{
				Name:     "updateModel",
				Endpoint: "/api/v1/person/{{id}}",
				Method:   http.MethodPut,
				PathVariables: map[string]codegen.RequestValue{
					"id": {
						Type: codegen.DynamicType{
							TypeID: codegen.TypeID_STRING,
						},
						Required: true,
					},
				},
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
					Required: true,
				},
			},
		},
	}

	definedTypes := map[string]string{
		"PersonModel": "lib/common",
	}

	resolver := codegen.ComposePackageResolver(
		jscodegen.WellKnownTypeResolver,
		func(typeRef string) (string, error) {
			pkg, exists := definedTypes[typeRef]
			if !exists {
				return "", fmt.Errorf("failed to resolve type '%s'", typeRef)
			}

			return pkg, nil
		},
	)

	typeManager := codegen.TypeManager{
		Filter: func(typeRef string) bool {
			if len(typeRef) == 0 {
				return false
			}

			_, err := resolver(typeRef)
			return err == nil
		},
		PkgResolver: resolver,
	}

	serviceCompiler := codegen.ServiceCompiler{
		TypeManager: typeManager,
	}

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
