package codegen

import "fmt"

type TypeFilter func(string) bool
type PackageResolver func(string) (string, error)

type TypeManager struct {
	Filter      TypeFilter
	PkgResolver PackageResolver
}

func ComposePackageResolver(resolvers ...PackageResolver) PackageResolver {
	return func(pkg string) (string, error) {

		var errs []error

		for _, resolver := range resolvers {
			resolved, err := resolver(pkg)
			if err == nil {
				return resolved, nil
			}

			// otherwise, append the error to the list
			errs = append(errs, err)
		}

		// if we got all the way through, we failed

		return "", fmt.Errorf("failed to resolve package %s", pkg)
	}
}

// ResolveImports - Given a
func ResolveImports(service CompiledService, resolver PackageResolver, typeFilter TypeFilter) (*ImportBlock, error) {
	uniqueReferences := collectTypes(service)

	imports := make(map[string][]string)
	for reference, _ := range uniqueReferences {
		// filter out anything if needed
		if !typeFilter(reference) {
			continue
		}

		// get package provider
		pkg, err := resolver(reference)
		if err != nil {
			return nil, fmt.Errorf("failed to resolve package for type '%s': %w", reference, err)
		}

		imports[pkg] = append(imports[pkg], reference)
	}

	return &ImportBlock{
		Packages: imports,
	}, nil
}

func collectTypes(service CompiledService) map[string]bool {
	uniqueReferences := make(map[string]bool)

	for _, property := range service.InputRecords {
		for _, varType := range property.Variables {
			addAllReferences(&uniqueReferences, varType.Type.Type.TypeReferences())
		}
	}

	for _, funcSig := range service.ServiceInterface.Functions {
		addAllReferences(&uniqueReferences, funcSig.Ret.TypeReferences())

		for _, arg := range funcSig.Args {
			addAllReferences(&uniqueReferences, arg.Type.Type.TypeReferences())
		}
	}

	for _, property := range service.Implementation.InjectionProperties {
		addAllReferences(&uniqueReferences, property.Type.Type.TypeReferences())
	}

	for _, property := range service.Implementation.Properties {
		addAllReferences(&uniqueReferences, property.Type.Type.TypeReferences())
	}

	return uniqueReferences
}

func addAllReferences(uniqueReferences *map[string]bool, references []string) {
	for _, ref := range references {
		(*uniqueReferences)[ref] = true
	}
}
