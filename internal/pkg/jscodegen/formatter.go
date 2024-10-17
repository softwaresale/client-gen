package jscodegen

import (
	"fmt"
	"github.com/softwaresale/client-gen/v2/internal/pkg/codegen"
	"io"
	"regexp"
	"strings"
)

type JSCodeFormatter struct {
	output     io.Writer
	typeMapper codegen.ITypeMapper
}

func NewJSCodeFormatter(output io.Writer, mapper codegen.ITypeMapper) *JSCodeFormatter {
	return &JSCodeFormatter{
		output:     output,
		typeMapper: mapper,
	}
}

type InterfaceBuilder struct {
	parent     *JSCodeFormatter
	name       string
	properties map[string]string
	functions  []FunctionBuilder
}

func NewInterfaceBuilder(parent *JSCodeFormatter, name string) *InterfaceBuilder {
	return &InterfaceBuilder{
		parent:     parent,
		name:       name,
		properties: make(map[string]string),
		functions:  make([]FunctionBuilder, 0),
	}
}

func (formatter *InterfaceBuilder) Property(name string, propType codegen.DynamicType) error {
	typeStr, err := formatter.parent.typeMapper.Convert(propType)
	if err != nil {
		return fmt.Errorf("failed to map property type: %w", err)
	}

	formatter.properties[name] = typeStr
	return nil
}

func (formatter *InterfaceBuilder) Function(function FunctionBuilder) error {
	formatter.functions = append(formatter.functions, function)
	return nil
}

func (formatter *InterfaceBuilder) Complete() error {

	_, err := fmt.Fprintf(formatter.parent.output, "export interface %s {\n", formatter.name)
	if err != nil {
		return err
	}

	for prop, typeStr := range formatter.properties {
		_, err = fmt.Fprintf(formatter.parent.output, "\t%s: %s;\n", prop, typeStr)
		if err != nil {
			return err
		}
	}

	for _, function := range formatter.functions {
		_, err = fmt.Fprint(formatter.parent.output, "\t")
		if err != nil {
			return err
		}

		err = function.Complete()
		if err != nil {
			return err
		}

		_, err = fmt.Fprint(formatter.parent.output, ";\n")
		if err != nil {
			return err
		}
	}

	_, err = fmt.Fprintf(formatter.parent.output, "}\n")
	if err != nil {
		return fmt.Errorf("failed to write interface closer: %w", err)
	}

	return nil
}

func (formatter *JSCodeFormatter) StartInterface(name string) *InterfaceBuilder {
	return NewInterfaceBuilder(formatter, name)
}

func (formatter *JSCodeFormatter) FunctionSignature(name string) *FunctionBuilder {
	return NewFunctionBuilder(name, formatter)
}

func (formatter *JSCodeFormatter) WriteURIString(template string, inputVarName string) error {
	var err error
	if _, err = fmt.Fprint(formatter.output, "`"); err != nil {
		return err
	}

	// TODO typecheck that all variables are included

	parser := regexp.MustCompile(`\{\{\s*(?P<pathVariable>[a-zA-Z_][a-zA-Z0-9_]*)\s*}}`)
	expansionTemplate := func(pathVariable string) string {
		return fmt.Sprintf("${%s.%s}", inputVarName, pathVariable)
	}

	pathVariableGroupIdx := parser.SubexpIndex("pathVariable")
	if pathVariableGroupIdx == -1 {
		panic("pathVariable group should be defined")
	}

	for {
		matchRange := parser.FindStringSubmatch(template)
		if matchRange == nil {
			break
		}

		whole := matchRange[0]
		pathVariable := matchRange[pathVariableGroupIdx]

		template = strings.Replace(template, whole, expansionTemplate(pathVariable), 1)
	}

	_, err = fmt.Fprintf(formatter.output, "%s`", template)
	return err
}

type FunctionBuilder struct {
	parent  *JSCodeFormatter
	name    string
	params  map[string]string
	retType string
}

func NewFunctionBuilder(name string, parent *JSCodeFormatter) *FunctionBuilder {
	return &FunctionBuilder{
		parent:  parent,
		name:    name,
		params:  make(map[string]string),
		retType: "",
	}
}

func (builder *FunctionBuilder) Param(name string, tp codegen.DynamicType) error {
	var err error
	builder.params[name], err = builder.parent.typeMapper.Convert(tp)
	if err != nil {
		delete(builder.params, name)
		return fmt.Errorf("failed to map parameter type for %s: %w", name, err)
	}

	return nil
}

func (builder *FunctionBuilder) ReturnType(tp codegen.DynamicType) error {
	var err error
	builder.retType, err = builder.parent.typeMapper.Convert(tp)
	return err
}

func (builder *FunctionBuilder) Complete() error {
	_, err := fmt.Fprintf(builder.parent.output, "%s(", builder.name)
	if err != nil {
		return fmt.Errorf("failed to format function name: %w", err)
	}

	idx := 0
	for param, typeVal := range builder.params {
		_, err = fmt.Fprintf(builder.parent.output, "%s: %s", param, typeVal)
		if err != nil {
			return fmt.Errorf("failed to format function param %s: %w", param, err)
		}

		if idx < len(builder.params)-1 {
			_, err = fmt.Fprintf(builder.parent.output, ", ")
			if err != nil {
				return fmt.Errorf("failed to format function param %s: %w", param, err)
			}
		}

		idx++
	}

	_, err = fmt.Fprintf(builder.parent.output, "): %s", builder.retType)
	if err != nil {
		return fmt.Errorf("failed to format function return: %w", err)
	}

	return nil
}
