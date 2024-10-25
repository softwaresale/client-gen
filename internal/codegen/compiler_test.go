package codegen

import (
	"fmt"
	importsmocks "github.com/softwaresale/client-gen/v2/internal/codegen/imports/mocks"
	outputsmocks "github.com/softwaresale/client-gen/v2/internal/codegen/outputs/mocks"
	servicegenmocks "github.com/softwaresale/client-gen/v2/internal/codegen/servicegen/mocks"
	"github.com/softwaresale/client-gen/v2/internal/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

var mockServiceGen *servicegenmocks.MockServiceGenerator
var mockImportMan *importsmocks.MockImportManager
var mockOutputMan *outputsmocks.MockCompilerOutputsManager
var compiler APICompiler
var apiDef types.APIDefinition

func setup(t *testing.T) {
	tmpDir := "output"

	mockServiceGen = servicegenmocks.NewMockServiceGenerator(t)
	mockImportMan = importsmocks.NewMockImportManager(t)
	mockOutputMan = outputsmocks.NewMockCompilerOutputsManager(t)

	compiler = APICompiler{
		Generator:      mockServiceGen,
		ImportManager:  mockImportMan,
		OutputsManager: mockOutputMan,
		OutputPath:     tmpDir,
	}

	apiDef = types.APIDefinition{
		Name:     "api",
		Services: []types.ServiceDefinition{},
		Entities: []types.EntitySpec{},
		Config: types.APIConfig{
			BaseURL: "http://localhost:8080",
		},
	}

	// This will always be called
	mockOutputMan.On("PrepareOutputDirectory", compiler.OutputPath).Return(nil).Once()
}

func configureDefaultAPIConfig(t *testing.T) {
	outputName := "config"
	mockConfigOutput := outputsmocks.NewMockCompilerOutputWriter(t)
	mockConfigOutput.On("Name").Return(outputName).Once()
	mockConfigOutput.On("Close").Return(nil).Once()
	mockOutputMan.On("CreateConfigOutput", apiDef.Config).Return(mockConfigOutput, nil).Once()
	mockImportMan.On("RegisterType", fmt.Sprintf("./%s", outputName), "APIConfig").Return().Once()
	mockServiceGen.On("GenerateConfig", mockConfigOutput, apiDef.Config, mockImportMan).Return(nil).Once()
}

func TestAPICompiler_Compile_DoesNothing(t *testing.T) {
	setup(t)
	configureDefaultAPIConfig(t)

	err := compiler.Compile(apiDef)
	assert.NoError(t, err)
	mockOutputMan.AssertExpectations(t)
	mockServiceGen.AssertExpectations(t)
	mockImportMan.AssertExpectations(t)
}

func TestAPICompiler_Compile_GeneratesAnEntity(t *testing.T) {
	setup(t)
	configureDefaultAPIConfig(t)

	// no properties
	entity1 := types.EntitySpec{
		Name:       "entity1",
		Properties: nil,
	}
	apiDef.Entities = append(apiDef.Entities, entity1)

	mockLocation := outputsmocks.NewMockCompilerOutputLocation(t)
	mockLocation.On("Name").Return(entity1.Name).Once()
	mockOutput := outputsmocks.NewMockCompilerOutputWriter(t)
	mockOutput.On("Close").Return(nil).Once()
	mockOutputMan.On("ComputeModelLocation", entity1).Return(mockLocation, nil).Once()
	mockOutputMan.On("CreateModelOutput", entity1).Return(mockOutput, nil).Once()

	mockImportMan.On("RegisterType", mock.Anything, entity1.Name).Once()

	mockServiceGen.On("GenerateEntity", mockOutput, entity1, mockImportMan).Return(nil).Once()

	err := compiler.Compile(apiDef)
	assert.NoError(t, err)
	mockOutputMan.AssertExpectations(t)
	mockServiceGen.AssertExpectations(t)
	mockImportMan.AssertExpectations(t)
}

func TestAPICompiler_Compile_GeneratesAService(t *testing.T) {
	setup(t)
	configureDefaultAPIConfig(t)

	service1 := types.ServiceDefinition{
		Name: "service1",
	}
	apiDef.Services = append(apiDef.Services, service1)

	mockOutput := outputsmocks.NewMockCompilerOutputWriter(t)
	mockOutput.On("Close").Return(nil).Once()
	mockOutputMan.On("CreateServiceOutput", service1).Return(mockOutput, nil).Once()

	mockServiceGen.On("GenerateService", mockOutput, service1, mockImportMan).Return(nil).Once()

	err := compiler.Compile(apiDef)
	assert.NoError(t, err)
	mockOutputMan.AssertExpectations(t)
	mockServiceGen.AssertExpectations(t)
	mockImportMan.AssertExpectations(t)
}
