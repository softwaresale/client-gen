package codegen

import (
	"fmt"
)

/// Convert a service definition into a concrete syntax tree of fragments

type CompiledService struct {
	inputRecords     []Record
	serviceInterface Interface
	implementation   Class
}

type ServiceCompiler struct {
}

func (compiler *ServiceCompiler) CompileService(serviceDef ServiceDefinition) CompiledService {

	// start an interface
	var inputTypeDefs []Record
	serviceInterface, _ := compiler.createServiceInterface(serviceDef)
	serviceClass, _ := compiler.createServiceClass(serviceDef)

	// create a number of record types
	for _, endpoint := range serviceDef.Endpoints {
		// create endpoint input types
		endpointInputRecord, endpointInputType := compiler.createEndpointInputPayload(endpoint)
		inputTypeDefs = append(inputTypeDefs, endpointInputRecord)

		// Create a signature for this endpoint
		endpointSignature := compiler.createEndpointSignature(endpoint, endpointInputType)

		// append the endpoint signature to the interface
		serviceInterface.Functions = append(serviceInterface.Functions, endpointSignature)

		// create an implementation for this endpoint
		// TODO

		// add the implementation to the class
		methodImpl := FunctionImpl{
			Signature: endpointSignature,
		}

		// add implementation
		serviceClass.Methods = append(serviceClass.Methods, methodImpl)
	}

	// create endpoint function signatures for all endpoints
	return CompiledService{
		implementation:   serviceClass,
		serviceInterface: serviceInterface,
		inputRecords:     inputTypeDefs,
	}
}

func (compiler *ServiceCompiler) createServiceInterface(definition ServiceDefinition) (Interface, DynamicType) {
	typeName := fmt.Sprintf("I%s", definition.Name)

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

func (compiler *ServiceCompiler) createServiceClass(definition ServiceDefinition) (Class, DynamicType) {
	class := Class{
		Name:                definition.Name,
		Properties:          make([]VariableDecl, 0),
		InjectionProperties: make([]VariableDecl, 0),
		Methods:             make([]FunctionImpl, 0),
	}

	dtype := DynamicType{
		TypeID:    TypeID_USER,
		Reference: definition.Name,
	}

	return class, dtype
}

func (compiler *ServiceCompiler) createEndpointInputPayload(endpoint APIEndpoint) (Record, DynamicType) {
	// TODO ensure proper casing endpoint
	typeName := fmt.Sprintf("%sInput", endpoint.Name)

	record := Record{
		Name:      typeName,
		Variables: make(map[string]VariableDecl),
	}

	// TODO add support for nullable/optional parameters

	// set the body if needed
	if endpoint.RequestBody.Type.TypeID != TypeID_VOID {
		record.Variables["body"] = VariableDecl{
			Name: "body",
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

func (compiler *ServiceCompiler) createEndpointSignature(endpoint APIEndpoint, inputType DynamicType) FunctionSignature {

	return FunctionSignature{
		Name: endpoint.Name,
		Ret:  endpoint.ResponseBody.Type,
		Args: []VariableDecl{
			{
				Name: "input",
				Type: TypeDecl{
					Type: inputType,
				},
			},
		},
	}
}
