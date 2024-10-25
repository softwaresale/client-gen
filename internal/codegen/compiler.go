package codegen

import (
	"fmt"
	"github.com/softwaresale/client-gen/v2/internal/codegen/imports"
	"github.com/softwaresale/client-gen/v2/internal/codegen/outputs"
	"github.com/softwaresale/client-gen/v2/internal/codegen/servicegen"
	"github.com/softwaresale/client-gen/v2/internal/types"
	"github.com/softwaresale/client-gen/v2/internal/utils"
	"strings"
)

// APICompiler compiles a service into a set of target files
type APICompiler struct {
	Generator      servicegen.ServiceGenerator    // generates a target service
	ImportManager  imports.ImportManager          // helps us manage imports
	OutputsManager outputs.CompilerOutputsManager // facilitates writing compiler outputs
	OutputPath     string
}

func (compiler *APICompiler) Compile(api types.APIDefinition) error {

	var err error

	err = compiler.OutputsManager.PrepareOutputDirectory(compiler.OutputPath)
	if err != nil {
		return fmt.Errorf("failed to setup output directory: %w", err)
	}

	// maps user types to paths that they can be imported from
	err = compiler.registerEntities(api)
	if err != nil {
		return fmt.Errorf("failed to register entities: %w", err)
	}

	// create all dependent entities
	for _, entitySpec := range api.Entities {
		err = compiler.compileEntity(entitySpec)
		if err != nil {
			return fmt.Errorf("failed to compile entity '%s': %w", entitySpec.Name, err)
		}
	}

	for _, service := range api.Services {
		err = compiler.compileService(service)
		if err != nil {
			return fmt.Errorf("failed to compile service '%s': %w", service.Name, err)
		}
	}

	return nil
}

func (compiler *APICompiler) compileEntity(entitySpec types.EntitySpec) error {
	entityWriter, err := compiler.OutputsManager.CreateModelOutput(entitySpec)
	if err != nil {
		return fmt.Errorf("failed to create model output: %w", err)
	}

	defer utils.SafeClose(entityWriter)

	err = compiler.Generator.GenerateEntity(entityWriter, entitySpec, compiler.ImportManager)
	if err != nil {
		return fmt.Errorf("failed to write entity: %w", err)
	}

	return nil
}

func (compiler *APICompiler) compileService(service types.ServiceDefinition) error {
	implWriter, err := compiler.OutputsManager.CreateServiceOutput(service)
	if err != nil {
		return fmt.Errorf("failed to create service output: %w", err)
	}
	defer utils.SafeClose(implWriter)

	// create the api implementation for each service in the API
	err = compiler.Generator.GenerateService(implWriter, service, compiler.ImportManager)
	if err != nil {
		return fmt.Errorf("failed to write service: %w", err)
	}

	return nil
}

// registerEntities registers all entities found in the API definition and works out which files
// __will eventually contain them__. This does not actually create any files or modify the output
// directory. This function just helps for generating imports
func (compiler *APICompiler) registerEntities(api types.APIDefinition) error {
	for _, entity := range api.Entities {
		output, err := compiler.OutputsManager.ComputeModelLocation(entity)
		if err != nil {
			return fmt.Errorf("failed to compute model location: %w", err)
		}

		importProvider := formatProviderName(output.Name())
		compiler.ImportManager.RegisterType(importProvider, entity.Name)
	}

	return nil
}

// TODO i don't like this...
func formatProviderName(outputFileName string) string {
	outputFileName = strings.TrimSuffix(outputFileName, ".ts")
	outputFileName = fmt.Sprintf("./%s", outputFileName)
	return outputFileName
}
