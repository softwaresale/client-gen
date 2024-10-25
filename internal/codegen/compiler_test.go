package codegen

import (
	importsmocks "github.com/softwaresale/client-gen/v2/internal/codegen/imports/mocks"
	outputsmocks "github.com/softwaresale/client-gen/v2/internal/codegen/outputs/mocks"
	servicegenmocks "github.com/softwaresale/client-gen/v2/internal/codegen/servicegen/mocks"
	"github.com/softwaresale/client-gen/v2/internal/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

var mockServiceGen *servicegenmocks.MockServiceGenerator
var mockImportMan *importsmocks.MockImportManager
var mockOutputMan *outputsmocks.MockCompilerOutputsManager
var compiler APICompiler
var apiDef types.APIDefinition

func setup(t *testing.T) {

	tmpDir := t.TempDir()

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
	}
}

func TestAPICompiler_Compile_DoesNothing(t *testing.T) {
	setup(t)

	mockOutputMan.On("PrepareOutputDirectory", compiler.OutputPath).Return(nil).Once()

	err := compiler.Compile(apiDef)
	assert.NoError(t, err)
	mockOutputMan.AssertExpectations(t)
}
