package codegen

import "github.com/softwaresale/client-gen/v2/internal/types"

// GenericImport provides an interface for generalized imports. An import accesses a number of
// entities from a given provider
type GenericImport interface {
	ProvidedEntities() []string // ProvidedEntities gets the entities provided by this import
	Provider() string           // Provider gets the name of this import provider
}

// ImportManager provides an interface for 1) registering imports and types they provide and 2) figuring out
// which types need to be imported for the given type
type ImportManager interface {
	RegisterProvider(providerName string)                              // RegisterProvider creates a new empty provider
	RegisterType(providerName, typeName string)                        // RegisterType adds a type to the given provider
	GetEntityImports(entity ...types.EntitySpec) []GenericImport       // GetEntityImports gets a list of imports needed by this collection of entities
	GetServiceImports(service types.ServiceDefinition) []GenericImport // GetServiceImports get all entities needed for the given service
}

// ImportCombiner combines multiple imports with the same provider into a single import. Implementation
// is target-specific
type ImportCombiner func([]GenericImport) GenericImport

// UnionImports combines imports that use the same provider
func UnionImports(combiner ImportCombiner, importSets ...[]GenericImport) []GenericImport {

	// find unique map of providers
	providerMap := make(map[string][]GenericImport)
	for _, importSet := range importSets {
		for _, importSpec := range importSet {
			duplicateProviders, exists := providerMap[importSpec.Provider()]
			if !exists {
				providerMap[importSpec.Provider()] = []GenericImport{importSpec}
				continue
			}

			providerMap[importSpec.Provider()] = append(duplicateProviders, importSpec)
		}
	}

	// now that we have duplicates, combine down
	var finalizedImports []GenericImport
	for _, relatedImports := range providerMap {
		combined := combiner(relatedImports)
		finalizedImports = append(finalizedImports, combined)
	}

	return finalizedImports
}
