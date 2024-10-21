package codegen

import (
	"fmt"
	"github.com/iancoleman/strcase"
)

/// Convert a service definition into a concrete syntax tree of fragments

type CompiledService struct {
	InputRecords     []Record
	ServiceInterface Interface
	Implementation   Class
	Imports          ImportBlock
}

type ServiceCompiler struct {
	TypeManager TypeManager
}

func (compiler *ServiceCompiler) CompileService(serviceDef ServiceDefinition) (*CompiledService, error) {

	// start an interface
	var err error
	var inputTypeDefs []Record
	serviceInterface, serviceInterfaceType := compiler.createServiceInterface(serviceDef)

	// make a service class
	serviceClass, _ := compiler.createServiceClass(serviceDef, serviceInterfaceType)

	// create a number of record types
	for _, endpoint := range serviceDef.Endpoints {
		// create endpoint input types
		endpointInputRecord, endpointInputType := compiler.createEndpointInputPayload(endpoint)
		inputTypeDefs = append(inputTypeDefs, endpointInputRecord)

		// Create a signature for this endpoint
		inputVarName := "input"
		endpointSignature := compiler.createEndpointSignature(endpoint, endpointInputType, inputVarName)

		// append the endpoint signature to the interface
		serviceInterface.Functions = append(serviceInterface.Functions, endpointSignature)

		// create an implementation for this endpoint
		// TODO break hardcoding of given
		request := HttpRequest{
			Method:          endpoint.Method,
			UrlTemplate:     endpoint.Endpoint,
			BaseEndpointVar: "",
			ClientVar:       "http",
			InputVar:        inputVarName,
			RequestBody:     endpoint.RequestBody.Type,
			ResponseBody:    endpoint.ResponseBody.Type,
		}

		// add the implementation to the class
		methodImpl := FunctionImpl{
			Signature: endpointSignature,
			HttpCall:  request,
		}

		// add implementation
		serviceClass.Methods = append(serviceClass.Methods, methodImpl)
	}

	// create endpoint function signatures for all endpoints
	compiledService := CompiledService{
		Implementation:   serviceClass,
		ServiceInterface: serviceInterface,
		InputRecords:     inputTypeDefs,
	}

	importBlock, err := ResolveImports(compiledService, compiler.TypeManager.PkgResolver, compiler.TypeManager.Filter)
	if err != nil {
		return nil, fmt.Errorf("error while resolving imports: %w", err)
	}

	compiledService.Imports = *importBlock

	return &compiledService, nil
}

func (compiler *ServiceCompiler) createServiceInterface(definition ServiceDefinition) (Interface, DynamicType) {
	typeName := fmt.Sprintf("I%s", strcase.ToCamel(definition.Name))

	// create an empty interface
	iface := Interface{
		Name:      typeName,
		Functions: make([]FunctionSignature, 0),
	}

	dtype := DynamicType{
		TypeID:    TypeID_USER,
		Reference: typeName,
	}

	return iface, dtype
}

func (compiler *ServiceCompiler) createServiceClass(definition ServiceDefinition, serviceIFace DynamicType) (Class, DynamicType) {

	typeName := strcase.ToCamel(definition.Name)

	class := Class{
		Name:       typeName,
		Interfaces: []DynamicType{serviceIFace},
		Properties: make([]VariableDecl, 0),
		InjectionProperties: []VariableDecl{
			{
				Name: "http",
				Type: TypeDecl{
					Type: DynamicType{
						TypeID:    TypeID_USER,
						Reference: "HttpClient",
					},
					Required: true,
				},
			},
		},
		Methods: make([]FunctionImpl, 0),
	}

	dtype := DynamicType{
		TypeID:    TypeID_USER,
		Reference: typeName,
	}

	return class, dtype
}

func (compiler *ServiceCompiler) createEndpointInputPayload(endpoint APIEndpoint) (Record, DynamicType) {
	// TODO ensure proper casing endpoint
	typeName := fmt.Sprintf("%sInput", strcase.ToCamel(endpoint.Name))

	bodyPropName := "body"
	record := Record{
		Name:      typeName,
		Variables: make(map[string]VariableDecl),
	}

	// TODO add support for nullable/optional parameters

	// set the body if needed
	if endpoint.RequestBody.Type.TypeID != TypeID_VOID {
		record.Variables[bodyPropName] = VariableDecl{
			Name: bodyPropName,
			Type: TypeDecl{
				Type: endpoint.RequestBody.Type,
			},
		}
	}

	// add parameters
	for pathVarName, pathVarType := range endpoint.PathVariables {
		record.Variables[pathVarName] = VariableDecl{
			Name: pathVarName,
			Type: TypeDecl{
				Type: pathVarType.Type,
			},
		}
	}

	// create a type
	inputType := DynamicType{
		TypeID:    TypeID_USER,
		Reference: typeName,
	}

	return record, inputType
}

func (compiler *ServiceCompiler) createEndpointSignature(endpoint APIEndpoint, inputType DynamicType, inputVarName string) FunctionSignature {

	return FunctionSignature{
		Name: endpoint.Name,
		Ret: DynamicType{
			TypeID:    TypeID_GENERIC,
			Reference: "Observable",
			Inner:     []DynamicType{endpoint.ResponseBody.Type},
		},
		Args: []VariableDecl{
			{
				Name: inputVarName,
				Type: TypeDecl{
					Type: inputType,
				},
			},
		},
	}
}
