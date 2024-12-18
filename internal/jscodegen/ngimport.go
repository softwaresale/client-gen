package jscodegen

import (
	"fmt"
	mapset "github.com/deckarep/golang-set/v2"
	"github.com/softwaresale/client-gen/v2/internal/codegen/imports"
	"github.com/softwaresale/client-gen/v2/internal/types"
)

type TSImport struct {
	File          string
	ProvidedTypes []string
}

func (imp *TSImport) Provider() string {
	return imp.File
}

func (imp *TSImport) ProvidedEntities() []string {
	return imp.ProvidedTypes
}

func CombineTSImports(genericImport []imports.GenericImport) imports.GenericImport {
	if len(genericImport) == 0 {
		return nil
	}

	name := genericImport[0].Provider()

	uniqueProvidedEntities := mapset.NewSet[string]()
	for _, imp := range genericImport {
		uniqueProvidedEntities.Append(imp.ProvidedEntities()...)
	}

	return &TSImport{
		File:          name,
		ProvidedTypes: uniqueProvidedEntities.ToSlice(),
	}
}

type TSImportManager struct {
	typeFiles map[string]string             // type name -> file that provides
	providers map[string]mapset.Set[string] // file -> types it provides
}

func NewTSImportManager() TSImportManager {
	return TSImportManager{
		providers: make(map[string]mapset.Set[string]),
		typeFiles: make(map[string]string),
	}
}

func (importManager *TSImportManager) RegisterProvider(providerName string) {
	_, exists := importManager.providers[providerName]
	if !exists {
		importManager.providers[providerName] = mapset.NewSet[string]()
	}
}

func (importManager *TSImportManager) RegisterType(providerName, typeName string) {
	provider, exists := importManager.providers[providerName]
	if exists {
		provider.Add(typeName)
		return
	}

	importManager.providers[providerName] = mapset.NewSet[string](typeName)

	importManager.typeFiles[typeName] = providerName
}

func (importManager *TSImportManager) GetEntityImports(entities ...types.EntitySpec) []imports.GenericImport {
	// get unique entities referenced in the entity implementation
	referencedEntities := mapset.NewSet[string]()
	for _, entity := range entities {
		for _, propSpec := range entity.Properties {
			referencedEntities.Append(propSpec.Type.TypeReferences()...)
		}
	}

	return importManager.createImportsForReferencedTypes(referencedEntities)
}

func (importManager *TSImportManager) GetServiceImports(service types.ServiceDefinition) []imports.GenericImport {
	referencedEntities := mapset.NewSet[string]()
	for _, endpoint := range service.Endpoints {
		referencedEntities.Append(endpoint.ResponseBody.Type.TypeReferences()...)
		referencedEntities.Append(endpoint.RequestBody.Type.TypeReferences()...)
		for _, prop := range endpoint.PathVariables {
			referencedEntities.Append(prop.Type.TypeReferences()...)
		}

		for _, prop := range endpoint.QueryVariables {
			referencedEntities.Append(prop.Type.TypeReferences()...)
		}
	}

	return importManager.createImportsForReferencedTypes(referencedEntities)
}

func (importManager *TSImportManager) GetImportForType(typeName string) (imports.GenericImport, error) {
	providingFile, exists := importManager.typeFiles[typeName]
	if !exists {
		return nil, fmt.Errorf("type '%s' is not registered", typeName)
	}

	return &TSImport{
		File:          providingFile,
		ProvidedTypes: []string{typeName},
	}, nil
}

func (importManager *TSImportManager) createImportsForReferencedTypes(referencedEntities mapset.Set[string]) []imports.GenericImport {
	var imports []imports.GenericImport

	// get unique files
	usedFiles := mapset.NewSet[string]()
	for _, uniqueEntity := range referencedEntities.ToSlice() {
		providingFile, exists := importManager.typeFiles[uniqueEntity]
		if exists {
			usedFiles.Add(providingFile)
		}
	}

	// turn into imports
	for _, providerFile := range usedFiles.ToSlice() {
		usedEntities := importManager.providers[providerFile].Intersect(referencedEntities)
		if usedEntities.IsEmpty() {
			continue
		}

		imp := TSImport{
			File:          providerFile,
			ProvidedTypes: usedEntities.ToSlice(),
		}

		imports = append(imports, &imp)
	}

	return imports
}
