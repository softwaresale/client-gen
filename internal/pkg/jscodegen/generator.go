package jscodegen

import (
	"fmt"
	"github.com/softwaresale/client-gen/v2/internal/pkg/codegen"
	"io"
)

type JSCodeGenerator struct {
	formatter      *JSCodeFormatter
	inputTypeCache *codegen.UserTypeCache
}

func NewJSCodeGenerator(writer io.Writer, inputTypeCache *codegen.UserTypeCache) (*JSCodeGenerator, error) {

	formatter := NewJSCodeFormatter(writer, JSTypeMapper{})

	return &JSCodeGenerator{
		formatter:      formatter,
		inputTypeCache: inputTypeCache,
	}, nil
}

func (gen *JSCodeGenerator) GenerateServiceInterface(definition codegen.ServiceDefinition) error {

	interfaceName := fmt.Sprintf("I%s", definition.Name)
	builder := gen.formatter.StartInterface(interfaceName)

	for _, endpoint := range definition.Endpoints {

		// find an input type for this
		inputType, exists := gen.inputTypeCache.GetEndpointInputType(endpoint)
		if !exists {
			return fmt.Errorf("endpoint %s(%s) so no input type", endpoint.Name, endpoint.Method)
		}

		// create a function for this endpoint
		endpointFunc, err := gen.createEndpointFunction(endpoint, inputType)
		if err != nil {
			return err
		}

		err = builder.Function(*endpointFunc)
		if err != nil {
			return err
		}
	}

	return builder.Complete()
}

func (gen *JSCodeGenerator) GenerateEndpoint(endpointDef codegen.APIEndpoint, inputType codegen.DynamicType) error {
	functionBuilder, err := gen.createEndpointFunction(endpointDef, inputType)
	if err != nil {
		return err
	}

	err = functionBuilder.Complete()
	if err != nil {
		return err
	}

	// write function body

	return nil
}

func (gen *JSCodeGenerator) createEndpointFunction(endpointDef codegen.APIEndpoint, inputType codegen.DynamicType) (*FunctionBuilder, error) {
	functionBuilder := gen.formatter.FunctionSignature(endpointDef.Name)
	inputParam := "input"

	err := functionBuilder.Param(inputParam, inputType)
	if err != nil {
		return nil, err
	}

	err = functionBuilder.ReturnType(endpointDef.ResponseBody.Type)
	if err != nil {
		return nil, err
	}

	return functionBuilder, nil
}

func (gen *JSCodeGenerator) GenerateInputType(endpoint codegen.APIEndpoint) (*codegen.DynamicType, error) {
	typeName := fmt.Sprintf("%sInput", endpoint.Name)

	var err error
	ifaceBuilder := gen.formatter.StartInterface(typeName)

	// add a body property
	if endpoint.RequestBody.Type.TypeID != codegen.TypeID_VOID {
		err = ifaceBuilder.Property("requestBody", endpoint.RequestBody.Type)
		if err != nil {
			return nil, err
		}
	}

	// add all parameters
	for varName, varType := range endpoint.PathVariables {
		err = ifaceBuilder.Property(varName, varType.Type)
		if err != nil {
			return nil, err
		}
	}

	err = ifaceBuilder.Complete()
	if err != nil {
		return nil, err
	}

	// create a type
	return &codegen.DynamicType{
		TypeID:    codegen.TypeID_USER,
		Reference: typeName,
	}, nil
}
