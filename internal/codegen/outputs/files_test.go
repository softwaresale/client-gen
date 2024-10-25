package outputs

import (
	"github.com/softwaresale/client-gen/v2/internal/types"
	"github.com/stretchr/testify/assert"
	"path/filepath"
	"testing"
)

func TestFileCompilerOutputLocation_Name_GetsFileName(t *testing.T) {
	path := "/some/absolute/path/to/filename"
	location := FileCompilerOutputLocation(path)
	assert.Equal(t, "filename", location.Name())
}

func TestFileCompilerOutputLocation_Location_IsJustPath(t *testing.T) {
	path := "/some/absolute/path/to/filename"
	location := FileCompilerOutputLocation(path)
	assert.Equal(t, path, location.Location())
}

var directoryCompilerOutput *DirectoryCompilerOutputsManager

func setup(t *testing.T) {
	tempDir := t.TempDir()
	directoryCompilerOutput = &DirectoryCompilerOutputsManager{
		BasePath: tempDir,
	}
}

func TestDirectoryCompilerOutputsManager_ComputeModelLocation_CorrectlyGeneratesLocation(t *testing.T) {
	setup(t)

	entityName := "SomeEntity"
	model := types.EntitySpec{
		Name:       entityName,
		Properties: nil,
	}

	location, err := directoryCompilerOutput.ComputeModelLocation(model)
	assert.NoError(t, err)

	// TODO this might be cheating
	expectedName := createOutputFileName(entityName, OutputType_MODEL)
	expectedLocation := filepath.Join(directoryCompilerOutput.BasePath, expectedName)

	assert.Equal(t, expectedName, location.Name())
	assert.Equal(t, expectedLocation, location.Location())
}

func TestDirectoryCompilerOutputsManager_CreateModelOutput_CorrectlyGeneratesLocation(t *testing.T) {
	setup(t)

	entityName := "SomeEntity"
	model := types.EntitySpec{
		Name:       entityName,
		Properties: nil,
	}

	output, err := directoryCompilerOutput.CreateModelOutput(model)
	assert.NoError(t, err)

	// TODO this might be cheating
	expectedName := createOutputFileName(entityName, OutputType_MODEL)
	expectedLocation := filepath.Join(directoryCompilerOutput.BasePath, expectedName)

	assert.Equal(t, expectedName, output.Name())
	assert.Equal(t, expectedLocation, output.Location())

	// also check that we can write to the location
	bytesWritten, err := output.Write([]byte("Hello World"))
	assert.NoError(t, err)
	assert.Greater(t, bytesWritten, 0)
}

func TestDirectoryCompilerOutputsManager_ComputeServiceLocation_CorrectlyGeneratesLocation(t *testing.T) {
	setup(t)

	serviceName := "SomeEntity"

	service := types.ServiceDefinition{
		Name:      serviceName,
		Endpoints: nil,
	}

	location, err := directoryCompilerOutput.ComputeServiceLocation(service)
	assert.NoError(t, err)

	// TODO this might be cheating
	expectedName := createOutputFileName(serviceName, OutputType_SERVICE)
	expectedLocation := filepath.Join(directoryCompilerOutput.BasePath, expectedName)

	assert.Equal(t, expectedName, location.Name())
	assert.Equal(t, expectedLocation, location.Location())
}

func TestDirectoryCompilerOutputsManager_CreateServiceOutput_CorrectlyGeneratesLocation(t *testing.T) {
	setup(t)

	serviceName := "SomeEntity"

	service := types.ServiceDefinition{
		Name:      serviceName,
		Endpoints: nil,
	}

	output, err := directoryCompilerOutput.CreateServiceOutput(service)
	assert.NoError(t, err)

	// TODO this might be cheating
	expectedName := createOutputFileName(serviceName, OutputType_SERVICE)
	expectedLocation := filepath.Join(directoryCompilerOutput.BasePath, expectedName)

	assert.Equal(t, expectedName, output.Name())
	assert.Equal(t, expectedLocation, output.Location())

	// also check that we can write to the output
	bytesWritten, err := output.Write([]byte("Hello World"))
	assert.NoError(t, err)
	assert.Greater(t, bytesWritten, 0)
}
