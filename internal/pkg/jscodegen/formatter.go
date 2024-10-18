package jscodegen

import (
	"fmt"
	"github.com/softwaresale/client-gen/v2/internal/pkg/codegen"
	"io"
	"net/http"
	"strings"
)

type JSCodeFormatter struct {
	output     io.Writer
	typeMapper codegen.ITypeMapper
}

func NewJSCodeFormatter(output io.Writer, typeMapper codegen.ITypeMapper) *JSCodeFormatter {
	return &JSCodeFormatter{
		output:     output,
		typeMapper: typeMapper,
	}
}

func (formatter *JSCodeFormatter) Format(service codegen.CompiledService) error {
	// format this service
	var err error

	err = formatter.formatImports(service.Imports)
	if err != nil {
		return fmt.Errorf("failed to format imports: %w", err)
	}

	formatter.infallibleFprint("\n")

	// first, lift all inputs
	for _, record := range service.InputRecords {
		err = formatter.formatRecord(record)
		if err != nil {
			return fmt.Errorf("failed to format record '%s': %w", record.Name, err)
		}

		formatter.infallibleFprint("\n")
	}

	// next, emit the service interface
	err = formatter.formatInterface(service.ServiceInterface)
	if err != nil {
		return fmt.Errorf("failed to format interface '%s': %w", service.ServiceInterface.Name, err)
	}

	// next, format the actual implementation
	err = formatter.formatServiceImplementation(service.Implementation)
	if err != nil {
		return fmt.Errorf("failed to format service impl '%s': %w", service.Implementation.Name, err)
	}

	return nil
}

func (formatter *JSCodeFormatter) formatImports(imports codegen.ImportBlock) error {
	for pkg, types := range imports.Packages {
		typesStr := strings.Join(types, ", ")
		formatter.infallibleFprintf("import { %s } from '%s';\n", typesStr, pkg)
	}

	return nil
}

func (formatter *JSCodeFormatter) formatRecord(record codegen.Record) error {
	formatter.infallibleFprintf("export interface %s {\n", record.Name)
	for _, property := range record.Variables {
		typeStr, err := formatter.typeMapper.Convert(property.Type.Type)
		if err != nil {
			return fmt.Errorf("failed to map property type for '%s': %w", property.Name, err)
		}

		formatter.infallibleFprintf("\t%s: %s;\n", property.Name, typeStr)
	}

	formatter.infallibleFprint("}\n")

	return nil
}

func (formatter *JSCodeFormatter) formatInterface(iface codegen.Interface) error {
	formatter.infallibleFprintf("export interface %s {\n", iface.Name)
	for _, signature := range iface.Functions {
		formatter.infallibleFprintf("\t%s(", signature.Name)
		var params []string
		for _, param := range signature.Args {
			typeStr, err := formatter.typeMapper.Convert(param.Type.Type)
			if err != nil {
				return fmt.Errorf("failed to map arg type for '%s': %w", param.Name, err)
			}

			param := fmt.Sprintf("%s: %s", param.Name, typeStr)
			params = append(params, param)
		}

		// write rest of arguments
		paramsStr := strings.Join(params, ", ")
		formatter.infallibleFprintf("%s)", paramsStr)

		// if return type is non-void, then write a void
		if signature.Ret.TypeID != codegen.TypeID_VOID {
			retTypeStr, err := formatter.typeMapper.Convert(signature.Ret)
			if err != nil {
				return fmt.Errorf("failed to map return type: %w", err)
			}

			formatter.infallibleFprintf(": %s", retTypeStr)
		}

		formatter.infallibleFprintf(";\n")
	}

	formatter.infallibleFprint("}\n")

	return nil
}

func (formatter *JSCodeFormatter) formatServiceImplementation(class codegen.Class) error {
	// start the class
	formatter.infallibleFprintf("export class %s ", class.Name)

	// implement any interfaces
	if len(class.Interfaces) > 0 {
		var ifaceNames []string
		for _, iface := range class.Interfaces {
			ifaceName, err := formatter.typeMapper.Convert(iface)
			if err != nil {
				return fmt.Errorf("failed to map interface type: %w", err)
			}

			ifaceNames = append(ifaceNames, ifaceName)
		}

		interfaceList := strings.Join(ifaceNames, ", ")
		formatter.infallibleFprintf("implements %s ", interfaceList)
	}

	formatter.infallibleFprint("{\n")

	// add a constructor for injection properties
	if len(class.InjectionProperties) > 0 {

		formatter.infallibleFprint("\tconstructor(\n")

		var params []string
		for _, injectionProp := range class.InjectionProperties {
			typeStr, err := formatter.typeMapper.Convert(injectionProp.Type.Type)
			if err != nil {
				return fmt.Errorf("failed to map property type for '%s': %w", injectionProp.Name, err)
			}

			param := fmt.Sprintf("\t\tprivate readonly %s: %s", injectionProp.Name, typeStr)
			params = append(params, param)
		}

		paramsStr := strings.Join(params, ",\n")
		formatter.infallibleFprintf("%s\n\t) {}\n", paramsStr)
	}

	// add the implementation
	for _, method := range class.Methods {
		// write the function signature
		formatter.infallibleFprint("\t")
		err := formatter.formatFunctionSignature(method.Signature)
		if err != nil {
			return fmt.Errorf("failed to format function signature: %w", err)
		}

		// write the body

		// begin body
		formatter.infallibleFprint(" {\n")

		// write the http call
		formatter.infallibleFprint("\t\treturn ")

		err = formatter.formatHttpCall(method.HttpCall)
		if err != nil {
			return fmt.Errorf("failed to format http call: %w", err)
		}

		formatter.infallibleFprint(";\n")

		// close body
		formatter.infallibleFprint("\t")
		formatter.infallibleFprint("}\n")
	}

	formatter.infallibleFprint("}\n")

	return nil
}

func (formatter *JSCodeFormatter) formatFunctionSignature(signature codegen.FunctionSignature) error {

	var args []string
	for _, arg := range signature.Args {
		typeStr, err := formatter.typeMapper.Convert(arg.Type.Type)
		if err != nil {
			return fmt.Errorf("failed to map arg type for '%s': %w", arg.Name, err)
		}

		argStr := fmt.Sprintf("%s: %s", arg.Name, typeStr)
		args = append(args, argStr)
	}

	formatter.infallibleFprintf("%s(%s)", signature.Name, strings.Join(args, ", "))

	if signature.Ret.TypeID != codegen.TypeID_VOID {
		retTypeStr, err := formatter.typeMapper.Convert(signature.Ret)
		if err != nil {
			return fmt.Errorf("failed to map return type: %w", err)
		}

		formatter.infallibleFprintf(": %s", retTypeStr)
	}

	return nil
}

func (formatter *JSCodeFormatter) formatHttpCall(call codegen.HttpRequest) error {
	formatter.infallibleFprintf("this.%s.", call.ClientVar)
	// figure out which method and type variable we might need
	var requestMethod string

	switch call.Method {
	case http.MethodGet:
		requestMethod = "get"
	case http.MethodPost:
		requestMethod = "post"
	case http.MethodPut:
		requestMethod = "put"
	default:
		return fmt.Errorf("invalid method '%s'", call.Method)
	}

	formatter.infallibleFprint(requestMethod)

	// optionally provide a generic string if a response body is specified
	if call.ResponseBody.TypeID != codegen.TypeID_VOID {
		typeStr, err := formatter.typeMapper.Convert(call.ResponseBody)
		if err != nil {
			return fmt.Errorf("failed to map response body type: %w", err)
		}
		formatter.infallibleFprintf("<%s>", typeStr)
	}

	// first, make the endpoint template
	templateFormatter := func(pathVariable string) (string, error) {
		return fmt.Sprintf("${%s.%s}", call.InputVar, pathVariable), nil
	}

	expandedTemplateString, err := FormatTemplate(call.UrlTemplate, templateFormatter)
	if err != nil {
		return fmt.Errorf("failed to format template: %w", err)
	}

	formatter.infallibleFprintf("(`%s`", expandedTemplateString)

	// if there are options or further arguments, do that here
	if call.RequestBody.TypeID != codegen.TypeID_VOID {
		// TODO body is hardcoded -- BAD!!
		formatter.infallibleFprintf(", %s.body", call.InputVar)
	}

	// close the function
	formatter.infallibleFprintf(")")

	return nil
}

func (formatter *JSCodeFormatter) formatVariableDecl(decl codegen.VariableDecl) error {
	typeStr, err := formatter.typeMapper.Convert(decl.Type.Type)
	if err != nil {
		return fmt.Errorf("failed to map decl type for '%s': %w", decl.Name, err)
	}

	formatter.infallibleFprintf("%s: %s", decl.Name, typeStr)

	return nil
}

func (formatter *JSCodeFormatter) infallibleFprintf(format string, a ...interface{}) {
	_, err := fmt.Fprintf(formatter.output, format, a...)
	if err != nil {
		panic(err)
	}
}

func (formatter *JSCodeFormatter) infallibleFprint(input string) {
	_, err := fmt.Fprint(formatter.output, input)
	if err != nil {
		panic(err)
	}
}
