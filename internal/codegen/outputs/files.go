package outputs

import (
	"fmt"
	"github.com/iancoleman/strcase"
	"github.com/softwaresale/client-gen/v2/internal/types"
	"os"
	"path/filepath"
)

// TODO add checking to output location type to ensure that this is what we want

type FileCompilerOutputLocation string

func (path FileCompilerOutputLocation) Name() string {
	pathStr := string(path)
	return filepath.Base(pathStr)
}

func (path FileCompilerOutputLocation) Location() string {
	return string(path)
}

type FileCompilerOutput struct {
	file    *os.File
	absPath FileCompilerOutputLocation
}

func (f FileCompilerOutput) Write(p []byte) (int, error) {
	return f.file.Write(p)
}

func (f FileCompilerOutput) Location() string {
	return f.absPath.Location()
}

func (f FileCompilerOutput) Name() string {
	return f.absPath.Name()
}

func (f FileCompilerOutput) Close() error {
	return f.file.Close()
}

const (
	OutputType_SERVICE = "service"
	OutputType_MODEL   = "model"
	OutputType_CONFIG  = "config"
)

// DirectoryCompilerOutputsManager outputs our files in a directory
type DirectoryCompilerOutputsManager struct {
	BasePath string // path that all outputs are relative to
}

func (outputs *DirectoryCompilerOutputsManager) PrepareOutputDirectory(path string) error {
	return setupOutputDirectory(path)
}

func (outputs *DirectoryCompilerOutputsManager) ComputeServiceLocation(serviceDef types.ServiceDefinition) (CompilerOutputLocation, error) {
	outputAbsPath, err := outputs.createOutputFilePath(serviceDef.Name, OutputType_SERVICE)
	if err != nil {
		return nil, fmt.Errorf("unable to compute service location: %w", err)
	}

	return FileCompilerOutputLocation(outputAbsPath), nil
}

func (outputs *DirectoryCompilerOutputsManager) CreateServiceOutput(serviceDef types.ServiceDefinition) (CompilerOutputWriter, error) {
	outputFile, outputAbsPath, err := outputs.createOutputFile(serviceDef.Name, OutputType_SERVICE)
	if err != nil {
		return nil, fmt.Errorf("failed to create output file: %w", err)
	}

	return &FileCompilerOutput{
		file:    outputFile,
		absPath: FileCompilerOutputLocation(outputAbsPath),
	}, nil
}

func (outputs *DirectoryCompilerOutputsManager) ComputeModelLocation(model types.EntitySpec) (CompilerOutputLocation, error) {
	outputAbsPath, err := outputs.createOutputFilePath(model.Name, OutputType_MODEL)
	if err != nil {
		return nil, fmt.Errorf("unable to compute service location: %w", err)
	}

	return FileCompilerOutputLocation(outputAbsPath), nil
}

func (outputs *DirectoryCompilerOutputsManager) CreateModelOutput(model types.EntitySpec) (CompilerOutputWriter, error) {
	outputFile, outputAbsPath, err := outputs.createOutputFile(model.Name, OutputType_MODEL)
	if err != nil {
		return nil, fmt.Errorf("failed to create output file: %w", err)
	}

	return &FileCompilerOutput{
		file:    outputFile,
		absPath: FileCompilerOutputLocation(outputAbsPath),
	}, nil
}

func (outputs *DirectoryCompilerOutputsManager) CreateConfigOutput(config types.APIConfig) (CompilerOutputWriter, error) {
	outputFile, outputAbsPath, err := outputs.createOutputFile("APIConfig", OutputType_CONFIG)
	if err != nil {
		return nil, fmt.Errorf("failed to create output file: %w", err)
	}

	return &FileCompilerOutput{
		file:    outputFile,
		absPath: FileCompilerOutputLocation(outputAbsPath),
	}, nil
}

func (outputs *DirectoryCompilerOutputsManager) ComputeConfigLocation(config types.APIConfig) (CompilerOutputLocation, error) {
	outputAbsPath, err := outputs.createOutputFilePath("APIConfig", OutputType_CONFIG)
	if err != nil {
		return nil, fmt.Errorf("failed to create output file: %w", err)
	}

	return FileCompilerOutputLocation(outputAbsPath), nil
}

func createOutputFileName(objectName, objectType string) string {
	return fmt.Sprintf("%s.%s.gen.ts", strcase.ToKebab(objectName), objectType)
}

func (outputs *DirectoryCompilerOutputsManager) createOutputFilePath(objectName, objectType string) (string, error) {
	outputFileName := createOutputFileName(objectName, objectType)
	outputPath := filepath.Join(outputs.BasePath, outputFileName)
	outputAbsPath, err := filepath.Abs(outputPath)
	if err != nil {
		return "", fmt.Errorf("failed to get absolute path of output file: %w", err)
	}

	return outputAbsPath, nil
}

func (outputs *DirectoryCompilerOutputsManager) createOutputFile(name, outputType string) (*os.File, string, error) {
	outputAbsPath, err := outputs.createOutputFilePath(name, outputType)
	if err != nil {
		return nil, "", fmt.Errorf("failed to get absolute path of output file: %w", err)
	}

	outputFile, err := os.Create(outputAbsPath)
	if err != nil {
		return nil, "", fmt.Errorf("failed to create output file: %w", err)
	}

	return outputFile, outputAbsPath, nil
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
		err = os.MkdirAll(path, os.ModePerm)
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
