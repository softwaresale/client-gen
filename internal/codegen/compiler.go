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
	GenerateService(writer io.Writer, service ServiceDefinition) error
	GenerateEntity(writer io.Writer, entity EntitySpec) error
}

// ServiceCompiler compiles a service into a set of target files
type ServiceCompiler struct {
	Generator  IServiceGenerator // generates a target service
	OutputPath string
}

func (compiler *ServiceCompiler) Compile(service ServiceDefinition) error {
	err := setupOutputDirectory(compiler.OutputPath)
	if err != nil {
		return fmt.Errorf("failed to setup output directory: %w", err)
	}

	// create all dependent entities
	for _, entitySpec := range service.Entities {
		// create an output file
		entityFileName := fmt.Sprintf("%s.model.gen.ts", strcase.ToKebab(entitySpec.Name))
		entityPath := filepath.Join(compiler.OutputPath, entityFileName)
		entityFile, err := os.Create(entityPath)
		if err != nil {
			return fmt.Errorf("failed to create entity file: %w", err)
		}

		err = compiler.writeEntityToFile(entityFile, entitySpec)
		if err != nil {
			return fmt.Errorf("failed to write entity: %w", err)
		}
	}

	// create the service implementation
	implFileName := fmt.Sprintf("%s.service.gen.ts", strcase.ToKebab(service.Name))
	implPath := filepath.Join(compiler.OutputPath, implFileName)
	implFile, err := os.Create(implPath)
	if err != nil {
		return fmt.Errorf("failed to create impl file: %w", err)
	}
	defer utils.SafeClose(implFile)

	err = compiler.Generator.GenerateService(implFile, service)
	if err != nil {
		return fmt.Errorf("failed to generate service: %w", err)
	}

	return nil
}

func (compiler *ServiceCompiler) writeEntityToFile(file *os.File, spec EntitySpec) error {
	defer utils.SafeClose(file)

	err := compiler.Generator.GenerateEntity(file, spec)
	if err != nil {
		return fmt.Errorf("failed to generate entity: %w", err)
	}

	return nil
}

func setupOutputDirectory(path string) error {
	stat, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("failed to stat output path: %w", err)
		}

		// create the directory
		err = os.Mkdir(path, os.ModePerm)
		if err != nil {
			return fmt.Errorf("failed to create output directory: %w", err)
		}
	}

	if !stat.IsDir() {
		return fmt.Errorf("output directory exists but is not a directory")
	}

	return nil
}
