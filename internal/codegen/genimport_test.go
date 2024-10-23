package codegen

import (
	mapset "github.com/deckarep/golang-set/v2"
	"github.com/stretchr/testify/assert"
	"testing"
)

type ExampleImport struct {
	ProviderName string
	Provides     []string
}

func (imp *ExampleImport) ProvidedEntities() []string {
	return imp.Provides
}

func (imp *ExampleImport) Provider() string {
	return imp.ProviderName
}

func CombineExampleImports(imports []GenericImport) GenericImport {
	unique := mapset.NewSet[string]()
	for _, imp := range imports {
		unique.Append(imp.ProvidedEntities()...)
	}

	return &ExampleImport{
		ProviderName: imports[0].Provider(),
		Provides:     unique.ToSlice(),
	}
}

func TestUnionImports_SuccessfullyUnions(t *testing.T) {
	importSet1 := []ExampleImport{
		{
			ProviderName: "pkg1",
			Provides:     []string{"e1", "e2"},
		},
		{
			ProviderName: "pkg3",
			Provides:     []string{"e7", "e5"},
		},
	}

	var importSet1View []GenericImport
	for _, imp := range importSet1 {
		importSet1View = append(importSet1View, &imp)
	}

	importSet2 := []ExampleImport{
		{
			ProviderName: "pkg1",
			Provides:     []string{"e1", "e3"},
		},
		{
			ProviderName: "pkg3",
			Provides:     []string{"e7", "e5"},
		},
		{
			ProviderName: "pkg2",
			Provides:     []string{"e10"},
		},
	}
	var importSet2View []GenericImport
	for _, imp := range importSet2 {
		importSet2View = append(importSet2View, &imp)
	}

	unionedImports := UnionImports(CombineExampleImports, importSet1View, importSet2View)

	associatedImports := make(map[string]GenericImport)
	for _, imp := range unionedImports {
		associatedImports[imp.Provider()] = imp
	}

	assert.Equal(t, len(unionedImports), len(associatedImports))

	expectedMap := map[string][]string{
		"pkg1": {"e1", "e2", "e3"},
		"pkg3": {"e7", "e5"},
		"pkg2": {"e10"},
	}

	for provider, packages := range expectedMap {
		assert.Contains(t, associatedImports, provider)
		assert.Equal(t, associatedImports[provider].Provider(), provider)
		assert.Equal(t, associatedImports[provider].ProvidedEntities(), packages)
	}
}

func TestUnionImports_EmptyWithEmptyInput(t *testing.T) {
	result := UnionImports(CombineExampleImports)
	assert.Empty(t, result)
}
