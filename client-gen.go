package main

import (
	"encoding/json"
	"flag"
	"fmt"
	codegen2 "github.com/softwaresale/client-gen/v2/internal/codegen"
	"github.com/softwaresale/client-gen/v2/internal/templates"
	"os"
)

// ProgramArgs specifies the arguments passed to this binary
type ProgramArgs struct {
	InputSpec string
	OutputDir string
	Target    TargetLanguage
}

var args ProgramArgs

const (
	TargetAngular = "angular"
	TargetSpring  = "spring"
)

type TargetLanguage string

func (t TargetLanguage) String() string { return string(t) }

func (t *TargetLanguage) Set(value string) error {
	switch value {
	case TargetAngular:
		*t = TargetAngular
	case TargetSpring:
		*t = TargetSpring
	default:
		return fmt.Errorf("unknown target language: %s", value)
	}

	return nil
}

func init() {
	flag.StringVar(&args.InputSpec, "input", "", "Path to input specification")
	flag.StringVar(&args.OutputDir, "output-dir", "", "The path to write this output to")
	flag.Var(&args.Target, "target", "The target language. Options are ['angular' (default), 'spring']")
}

func main() {

	// parse the input
	flag.Parse()

	if len(args.InputSpec) == 0 {
		fmt.Println("input specification path is required")
		return
	}

	if len(args.OutputDir) == 0 {
		fmt.Println("output dir path is required")
		return
	}

	service, err := readServiceDefinition(args.InputSpec)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	/*
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
	*/

	ngServiceGen := templates.NewNGServiceGenerator()

	compiler := codegen2.ServiceCompiler{
		Generator:  ngServiceGen,
		OutputPath: "./output/",
	}

	err = compiler.Compile(service)
	if err != nil {
		panic(err)
	}
}

func readServiceDefinition(path string) (codegen2.ServiceDefinition, error) {
	serviceFileContents, err := os.ReadFile(path)
	if err != nil {
		return codegen2.ServiceDefinition{}, fmt.Errorf("failed to open service definition file: %w", err)
	}

	var serviceDef codegen2.ServiceDefinition
	err = json.Unmarshal(serviceFileContents, &serviceDef)
	if err != nil {
		return codegen2.ServiceDefinition{}, fmt.Errorf("failed to parse service definition file: %w", err)
	}

	return serviceDef, nil
}
