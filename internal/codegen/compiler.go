package codegen

import (
	"fmt"
	"github.com/iancoleman/strcase"
	"github.com/softwaresale/client-gen/v2/internal/utils"
	"io"
	"os"
	"path/filepath"
)

// IServiceGenerator specifies a type that formats a service to an output
type IServiceGenerator interface {
	GenerateService(writer io.Writer, service ServiceDefinition, outputPath string, resolver UserTypeResolver) error
	GenerateEntity(writer io.Writer, entity EntitySpec, outputPath string, resolver UserTypeResolver) error
}

// APICompiler compiles a service into a set of target files
type APICompiler struct {
	Generator  IServiceGenerator // generates a target service
	OutputPath string
}

func (compiler *APICompiler) Compile(api APIDefinition) error {
	err := setupOutputDirectory(compiler.OutputPath)
	if err != nil {
		return fmt.Errorf("failed to setup output directory: %w", err)
	}

	// maps user types to paths that they can be imported from
	typeResolver := NewUserTypeResolver()

	// create all dependent entities
	for _, entitySpec := range api.Entities {
		// create an output file
		entityFile, entityFileName, _, err := createOutputFile(compiler.OutputPath, entitySpec.Name, "model")
		if err != nil {
			return fmt.Errorf("failed to create entity file: %w", err)
		}

		err = compiler.writeEntityToFile(entityFile, entitySpec, entityFileName, typeResolver)
		if err != nil {
			return fmt.Errorf("failed to write entity: %w", err)
		}

		// stash the filename so that we can write the import
		// TODO optionally check if it was actually registered
		typeResolver.RegisterType(entitySpec, entityFileName)
	}

	for _, service := range api.Services {
		implFile, _, servicePath, err := createOutputFile(compiler.OutputPath, service.Name, "service")
		if err != nil {
			return fmt.Errorf("failed to create api: %w", err)
		}

		// create the api implementation for each service in the API
		err = compiler.writeServiceToFile(implFile, service, servicePath, typeResolver)
		if err != nil {
			return fmt.Errorf("failed to write service: %w", err)
		}
	}

	return nil
}

func (compiler *APICompiler) writeServiceToFile(file *os.File, service ServiceDefinition, servicePath string, typeResolver UserTypeResolver) error {
	var err error
	defer utils.SafeClose(file)
	err = compiler.Generator.GenerateService(file, service, servicePath, typeResolver)
	if err != nil {
		return fmt.Errorf("failed to generate api: %w", err)
	}

	return nil
}

func (compiler *APICompiler) writeEntityToFile(file *os.File, spec EntitySpec, specPath string, resolver UserTypeResolver) error {
	defer utils.SafeClose(file)

	err := compiler.Generator.GenerateEntity(file, spec, specPath, resolver)
	if err != nil {
		return fmt.Errorf("failed to generate entity: %w", err)
	}

	return nil
}

func createOutputFile(basePath, objectName, objectType string) (*os.File, string, string, error) {
	fileName := fmt.Sprintf("%s.%s.gen.ts", strcase.ToKebab(objectName), objectType)
	filePath := filepath.Join(basePath, fileName)
	file, err := os.Create(filePath)
	if err != nil {
		return nil, "", "", fmt.Errorf("failed to create impl file: %w", err)
	}

	return file, fileName, filePath, nil
}

func setupOutputDirectory(path string) error {
	var err error
	if !filepath.IsAbs(path) {
		path, err = filepath.Abs(path)
		if err != nil {
			return fmt.Errorf("failed to make absolute path: %w", err)
		}
	}

	stat, err := os.Stat(path)
	if err != nil {
		if !os.IsNotExist(err) {
			return fmt.Errorf("failed to stat output path: %w", err)
		}

		// create the directory
		err = os.Mkdir(path, os.ModePerm)
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
