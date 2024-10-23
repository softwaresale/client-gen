package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/softwaresale/client-gen/v2/internal/codegen"
	"github.com/softwaresale/client-gen/v2/internal/jscodegen"
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

	apiDef, err := readAPIDefinition(args.InputSpec)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	ngServiceGen := templates.NewNGServiceGenerator()

	outputDirectory := fmt.Sprintf("./output/%s", apiDef.Name)
	ngImportMgr := jscodegen.NewTSImportManager()

	compiler := codegen.APICompiler{
		Generator:     ngServiceGen,
		ImportManager: &ngImportMgr,
		OutputPath:    outputDirectory,
	}

	err = compiler.Compile(apiDef)
	if err != nil {
		panic(err)
	}
}

func readAPIDefinition(path string) (codegen.APIDefinition, error) {
	serviceFileContents, err := os.ReadFile(path)
	if err != nil {
		return codegen.APIDefinition{}, fmt.Errorf("failed to open API definition file: %w", err)
	}

	var apiDef codegen.APIDefinition
	err = json.Unmarshal(serviceFileContents, &apiDef)
	if err != nil {
		return codegen.APIDefinition{}, fmt.Errorf("failed to parse API definition file: %w", err)
	}

	return apiDef, nil
}
