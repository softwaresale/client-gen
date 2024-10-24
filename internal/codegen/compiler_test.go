package codegen

import (
	"github.com/softwaresale/client-gen/v2/internal/testutils"
	"github.com/stretchr/testify/assert"
	"io"
	"os"
	"testing"
)

type MockServiceGenerator struct {
	TestGenerateService func(writer io.Writer, service ServiceDefinition, outputPath string, resolver ImportManager) error
	TestGenerateEntity  func(writer io.Writer, entity EntitySpec, outputPath string, resolver ImportManager) error
}

func DefaultMockServiceGenerator() *MockServiceGenerator {
	return &MockServiceGenerator{
		TestGenerateService: func(writer io.Writer, service ServiceDefinition, outputPath string, resolver ImportManager) error {
			return nil
		},
		TestGenerateEntity: func(writer io.Writer, entity EntitySpec, outputPath string, resolver ImportManager) error {
			return nil
		},
	}
}

func (m *MockServiceGenerator) GenerateService(writer io.Writer, service ServiceDefinition, outputPath string, resolver ImportManager) error {
	return m.TestGenerateService(writer, service, outputPath, resolver)
}

func (m *MockServiceGenerator) GenerateEntity(writer io.Writer, entity EntitySpec, outputPath string, resolver ImportManager) error {
	return m.TestGenerateEntity(writer, entity, outputPath, resolver)
}

var outputDirectoryInfo *testutils.MockFileInfo
var mockServiceGenerator *MockServiceGenerator
var mockFileOps *testutils.MockFileOperations
var compiler *APICompiler
var apiDef APIDefinition

func setup(t *testing.T) {

	mockServiceGenerator = DefaultMockServiceGenerator()
	mockFileOps = testutils.DefaultMockFileOperations()

	mockFileOps.TestStat = func(name string) (os.FileInfo, error) {
		return outputDirectoryInfo, nil
	}

	compiler = &APICompiler{
		Generator:     mockServiceGenerator,
		FileOps:       mockFileOps,
		ImportManager: nil,
		OutputPath:    "",
	}

	apiDef = APIDefinition{
		Name:     "API",
		Entities: []EntitySpec{},
		Services: []ServiceDefinition{},
	}
}

func TestAPICompiler_Compile_NoErrorsOnHappyPath(t *testing.T) {
	setup(t)

	err := compiler.Compile(apiDef)
	assert.NoError(t, err)
}
