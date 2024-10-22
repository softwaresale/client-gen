package codegen

import (
	"fmt"
	"strings"
)

// UserTypeResolver resolves imports for user defined types
type UserTypeResolver struct {
	userDefinedTypes    map[string]EntitySpec      // unique user identified types
	typePackageProvider map[string]string          // map of (userTypeReference -> package that provides it)
	packageCache        map[string]map[string]bool // map of (package -> types it provides)
}

func NewUserTypeResolver() UserTypeResolver {
	return UserTypeResolver{
		userDefinedTypes:    make(map[string]EntitySpec),
		typePackageProvider: make(map[string]string),
		packageCache:        make(map[string]map[string]bool),
	}
}

// RegisterType registers a type and where it is defined
func (resolver *UserTypeResolver) RegisterType(entity EntitySpec, entityPath string) bool {

	_, alreadyRegistered := resolver.userDefinedTypes[entity.Name]
	if alreadyRegistered {
		return false
	}

	resolver.userDefinedTypes[entity.Name] = entity
	resolver.typePackageProvider[entity.Name] = entityPath
	existingMultiset, ok := resolver.packageCache[entityPath]
	if !ok {
		resolver.packageCache[entityPath] = make(map[string]bool)
		existingMultiset = resolver.packageCache[entityPath]
	}

	existingMultiset[entity.Name] = true

	return true
}

func (resolver *UserTypeResolver) CreateImportMap(neededEntities []string, servicePath string) (map[string][]string, error) {

	if len(neededEntities) == 0 {
		return resolver.createDefaultImportMap(servicePath)
	}

	// unique entities
	uniqueEntities := make(map[string]bool)
	for _, entity := range neededEntities {
		uniqueEntities[entity] = true
	}

	// map of package and unique entities it provides
	imports := make(map[string]map[string]bool)

	for entity, _ := range uniqueEntities {
		// get the package that provides this entity
		providingPackage, exists := resolver.typePackageProvider[entity]
		if !exists {
			return nil, fmt.Errorf("entity type '%s' was not registered", entity)
		}

		// uniquely add that entity to the import multimap
		existingSequence, exists := imports[providingPackage]
		if !exists {
			imports[providingPackage] = make(map[string]bool)
			existingSequence = imports[providingPackage]
		}

		existingSequence[entity] = true
	}

	// make all entity paths relative to our service path
	relativeImports := make(map[string][]string)
	for providingPackage, uniqueEntitySet := range imports {
		relativePath, err := makeRelativePath(servicePath, providingPackage)
		if err != nil {
			return nil, err
		}

		for entity, _ := range uniqueEntitySet {
			relativeImports[relativePath] = append(relativeImports[relativePath], entity)
		}
	}

	return relativeImports, nil
}

func (resolver *UserTypeResolver) createDefaultImportMap(servicePath string) (map[string][]string, error) {
	importMap := make(map[string][]string)
	for providingPackage, uniqueEntities := range resolver.packageCache {
		relativeImport, err := makeRelativePath(servicePath, providingPackage)
		if err != nil {
			return nil, err
		}

		// ditch file extensions
		relativeImport = strings.TrimSuffix(relativeImport, ".ts")

		var entitySlice []string
		for entity := range uniqueEntities {
			entitySlice = append(entitySlice, entity)
		}

		importMap[relativeImport] = entitySlice
	}

	return importMap, nil
}

func makeRelativePath(servicePath, packagePath string) (string, error) {
	relativePath := fmt.Sprintf("./%s", packagePath)
	/*
		relativePath, err := filepath.Rel(servicePath, providingPackage)
		if err != nil {
			return nil, fmt.Errorf("failed to relativize path: %w", err)
		}
	*/

	relativePath = strings.TrimSuffix(relativePath, ".ts")

	return relativePath, nil
}
