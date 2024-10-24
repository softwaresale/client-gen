package codegen

import (
	"fmt"
	"github.com/iancoleman/strcase"
	"github.com/softwaresale/client-gen/v2/internal/utils"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// ServiceGenerator specifies a type that formats a service to an output
type ServiceGenerator interface {
	GenerateService(writer io.Writer, service ServiceDefinition, outputPath string, resolver ImportManager) error
	GenerateEntity(writer io.Writer, entity EntitySpec, outputPath string, resolver ImportManager) error
}

// APICompiler compiles a service into a set of target files
type APICompiler struct {
	Generator     ServiceGenerator             // generates a target service
	ImportManager ImportManager                // helps us manage imports
	FileOps       utils.CompilerFileOperations // Abstracts away how to perform file operations. This is mainly for testing
	OutputPath    string
}

func (compiler *APICompiler) Compile(api APIDefinition) error {
	err := compiler.setupOutputDirectory(compiler.OutputPath)
	if err != nil {
		return fmt.Errorf("failed to setup output directory: %w", err)
	}

	// maps user types to paths that they can be imported from
	compiler.registerEntities(api)

	// create all dependent entities
	for _, entitySpec := range api.Entities {
		// create an output file
		entityFile, entityFileName, _, err := compiler.createOutputFile(compiler.OutputPath, entitySpec.Name, "model")
		if err != nil {
			return fmt.Errorf("failed to create entity file: %w", err)
		}

		err = compiler.writeEntityToFile(entityFile, entitySpec, entityFileName)
		if err != nil {
			return fmt.Errorf("failed to write entity: %w", err)
		}
	}

	for _, service := range api.Services {
		implFile, _, servicePath, err := compiler.createOutputFile(compiler.OutputPath, service.Name, "service")
		if err != nil {
			return fmt.Errorf("failed to create api: %w", err)
		}

		// create the api implementation for each service in the API
		err = compiler.writeServiceToFile(implFile, service, servicePath)
		if err != nil {
			return fmt.Errorf("failed to write service: %w", err)
		}
	}

	return nil
}

// registerEntities registers all entities found in the API definition and works out which files
// __will eventually contain them__. This does not actually create any files or modify the output
// directory. This function just helps for generating imports
func (compiler *APICompiler) registerEntities(api APIDefinition) {
	for _, entity := range api.Entities {
		outputFileName := createOutputFilePath(entity.Name, "model")
		importProvider := formatProviderName(outputFileName)
		compiler.ImportManager.RegisterType(importProvider, entity.Name)
	}
}

// TODO i don't like this...
func formatProviderName(outputFileName string) string {
	outputFileName = strings.TrimSuffix(outputFileName, ".ts")
	outputFileName = fmt.Sprintf("./%s", outputFileName)
	return outputFileName
}

func (compiler *APICompiler) writeServiceToFile(file *os.File, service ServiceDefinition, servicePath string) error {
	var err error
	defer utils.SafeClose(file)
	err = compiler.Generator.GenerateService(file, service, servicePath, compiler.ImportManager)
	if err != nil {
		return fmt.Errorf("failed to generate api: %w", err)
	}

	return nil
}

func (compiler *APICompiler) writeEntityToFile(file *os.File, spec EntitySpec, specPath string) error {
	defer utils.SafeClose(file)

	err := compiler.Generator.GenerateEntity(file, spec, specPath, compiler.ImportManager)
	if err != nil {
		return fmt.Errorf("failed to generate entity: %w", err)
	}

	return nil
}

func createOutputFilePath(objectName, objectType string) string {
	return fmt.Sprintf("%s.%s.gen.ts", strcase.ToKebab(objectName), objectType)
}

func (compiler *APICompiler) createOutputFile(basePath, objectName, objectType string) (*os.File, string, string, error) {
	fileName := createOutputFilePath(objectName, objectType)
	filePath := filepath.Join(basePath, fileName)
	file, err := compiler.FileOps.Create(filePath)
	if err != nil {
		return nil, "", "", fmt.Errorf("failed to create impl file: %w", err)
	}

	return file, fileName, filePath, nil
}

func (compiler *APICompiler) setupOutputDirectory(path string) error {
	var err error
	if !filepath.IsAbs(path) {
		path, err = filepath.Abs(path)
		if err != nil {
			return fmt.Errorf("failed to make absolute path: %w", err)
		}
	}

	stat, err := compiler.FileOps.Stat(path)
	if err != nil {
		if !compiler.FileOps.IsNotExist(err) {
			return fmt.Errorf("failed to stat output path: %w", err)
		}

		// create the directory
		err = compiler.FileOps.Mkdir(path, os.ModePerm)
		if err != nil {
			return fmt.Errorf("failed to create output directory: %w", err)
		}

		// we successfully created the directory and we're done
		return nil
	}

	if !stat.IsDir() {
		return fmt.Errorf("output directory exists but is not a directory")
	}

	return nil
}
