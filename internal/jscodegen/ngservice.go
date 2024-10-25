package jscodegen

import (
	_ "embed"
	"fmt"
	"github.com/iancoleman/strcase"
	"github.com/softwaresale/client-gen/v2/internal/codegen"
	"github.com/softwaresale/client-gen/v2/internal/codegen/imports"
	"github.com/softwaresale/client-gen/v2/internal/types"
	"io"
	"strings"
	"text/template"
)

//go:embed ng-service.tmpl
var templateText string

//go:embed ng-entity.tmpl
var entityTemplateText string

//go:embed ng-standalone-entity.tmpl
var standaloneEntityTemplateText string

//go:embed ng-imports.tmpl
var importsTemplateText string

//go:embed ng-config.tmpl
var configTemplateText string

type HttpRequestDef struct {
	HttpClientVar    string              // the name of the variable that defines the HTTP client in use
	HttpMethod       string              // The HTTP method used by this request
	ResponseType     string              // the type string of our response
	URITemplate      codegen.URITemplate // our URI template. This gets mapped into a uri string
	RequestBodyValue string              // The value to read the body type
}

func hasRequestBody(def string) bool {
	return len(def) > 0
}

type RequestMethodDef struct {
	RequestName      string         // The name of this request
	InputVarName     string         // The variable name of the input payload type
	RequestInputType string         // the type string of the input payload
	ResponseType     string         // The type string of the response
	HttpRequest      HttpRequestDef // The http request that should be called in this endpoint
}

func (def RequestMethodDef) HasInput() bool {
	return len(def.RequestInputType) > 0
}

// EntityDef defines the template for a standalone entity file
type EntityDef struct {
	Entity  types.EntitySpec        // the entity we are generating
	Imports []imports.GenericImport // imports used by this entity
}

type ServiceDef struct {
	ServiceName   string
	HttpClientVar string
	InputTypes    []types.EntitySpec
	Methods       []RequestMethodDef
	Imports       []imports.GenericImport
}

// ConfigDef defines what we need to model for our API configuration providers
type ConfigDef struct {
	APIName      string                  // what the name of the overall API configuration is
	ConfigEntity types.EntitySpec        // The record that houses our
	ConfigInit   types.EntityInitializer // how to configure the default configuration
}

func mapHttpEndpoint(method string) string {
	return strings.ToLower(method)
}

type NGServiceGenerator struct {
	ngServiceTemplate *template.Template
	ngEntityTemplate  *template.Template
	ngConfigTemplate  *template.Template
}

// NewNGServiceGenerator creates a new NGService generator, which can be used to generate services
func NewNGServiceGenerator() *NGServiceGenerator {

	typeMapper := JSTypeMapper{}
	valueMapper := JSValueMapper{}

	funcMap := template.FuncMap{
		"HasRequestBody": hasRequestBody,
		"ParseTemplate":  codegen.FormatTemplate,
		"ConvertType":    typeMapper.Convert,
		"ConvertValue":   valueMapper.Convert,
	}

	serviceTmpl := template.Must(template.New("NGService").Funcs(funcMap).Parse(templateText))
	serviceTmpl = template.Must(serviceTmpl.Parse(importsTemplateText))
	serviceTmpl = template.Must(serviceTmpl.Parse(entityTemplateText))

	entityTmpl := template.Must(template.New("NGEntity").Funcs(funcMap).Parse(entityTemplateText))
	entityTmpl = template.Must(entityTmpl.Parse(importsTemplateText))
	entityTmpl = template.Must(entityTmpl.Parse(standaloneEntityTemplateText))

	configTmpl := template.Must(template.New("NGConfig").Funcs(funcMap).Parse(configTemplateText))
	configTmpl = template.Must(configTmpl.Parse(entityTemplateText))

	return &NGServiceGenerator{
		ngServiceTemplate: serviceTmpl,
		ngEntityTemplate:  entityTmpl,
		ngConfigTemplate:  configTmpl,
	}
}

func (generator *NGServiceGenerator) GenerateService(writer io.Writer, def types.ServiceDefinition, resolver imports.ImportManager) error {
	translatedDef, err := translateService(def, resolver)
	if err != nil {
		return fmt.Errorf("failed to translateService service definition: %w", err)
	}

	return generator.ngServiceTemplate.Execute(writer, translatedDef)
}

func translateService(service types.ServiceDefinition, importResolver imports.ImportManager) (ServiceDef, error) {

	typeMapper := JSTypeMapper{}
	httpClientVar := "http"

	var methods []RequestMethodDef
	var inputs []types.EntitySpec
	for _, endpoint := range service.Endpoints {

		inputVarName := "input"
		bodyPropertyName := "body"

		requestInputDef, err := createInputType(endpoint, bodyPropertyName)
		if err != nil {
			return ServiceDef{}, err
		}

		inputTypeName := ""
		if requestInputDef.IsValid() {
			inputTypeName = requestInputDef.Name
			inputs = append(inputs, *requestInputDef)
		}

		requestBodyValue := ""
		if !endpoint.RequestBody.Type.IsVoid() {
			requestBodyValue = fmt.Sprintf("%s.%s", inputVarName, bodyPropertyName)
		}

		responseType, err := typeMapper.Convert(endpoint.ResponseBody.Type)
		if err != nil {
			return ServiceDef{}, err
		}

		methodDef := RequestMethodDef{
			RequestName:      endpoint.Name,
			InputVarName:     inputVarName,
			RequestInputType: inputTypeName,
			ResponseType:     responseType,
			HttpRequest: HttpRequestDef{
				HttpClientVar: httpClientVar,
				HttpMethod:    mapHttpEndpoint(endpoint.Method),
				ResponseType:  responseType,
				URITemplate: codegen.URITemplate{
					Template: endpoint.Endpoint,
					VarMapper: func(pathVar string) (string, error) {
						return fmt.Sprintf("${%s.%s}", inputVarName, pathVar), nil
					},
				},
				RequestBodyValue: requestBodyValue,
			},
		}

		methods = append(methods, methodDef)
	}

	inputImportMap := importResolver.GetEntityImports(inputs...)
	serviceImportMap := importResolver.GetServiceImports(service)

	importMap := imports.UnionImports(CombineTSImports, inputImportMap, serviceImportMap)

	return ServiceDef{
		ServiceName:   service.Name,
		HttpClientVar: httpClientVar,
		Methods:       methods,
		InputTypes:    inputs,
		Imports:       importMap,
	}, nil
}

func createInputType(endpoint types.APIEndpoint, bodyPropertyName string) (*types.EntitySpec, error) {
	inputTypeName := strcase.ToCamel(fmt.Sprintf("%sInput", endpoint.Name))
	properties := make(map[string]types.PropertySpec)
	for prop, tp := range endpoint.PathVariables {
		properties[prop] = types.PropertySpec{
			Type:     tp.Type,
			Required: tp.Required,
		}
	}

	if !endpoint.RequestBody.Type.IsVoid() {
		properties[bodyPropertyName] = types.PropertySpec{
			Type:     endpoint.RequestBody.Type,
			Required: endpoint.RequestBody.Required,
		}
	}

	return &types.EntitySpec{
		Name:       inputTypeName,
		Properties: properties,
	}, nil
}

func (generator *NGServiceGenerator) GenerateEntity(writer io.Writer, def types.EntitySpec, resolver imports.ImportManager) error {
	entity := translateEntity(def, resolver)
	return generator.ngEntityTemplate.Execute(writer, entity)
}

func translateEntity(spec types.EntitySpec, importResolver imports.ImportManager) EntityDef {

	entityImports := importResolver.GetEntityImports(spec)

	return EntityDef{
		Entity:  spec,
		Imports: entityImports,
	}
}

func (generator *NGServiceGenerator) GenerateConfig(writer io.Writer, config types.APIConfig, resolver imports.ImportManager) error {
	configDef, err := generator.translateConfig(config, resolver)
	if err != nil {
		return fmt.Errorf("failed to translate config def: %w", err)
	}

	return generator.ngConfigTemplate.Execute(writer, configDef)
}

func (generator *NGServiceGenerator) translateConfig(config types.APIConfig, resolver imports.ImportManager) (*ConfigDef, error) {

	// create the type that will represent our config
	configType, err := config.CreateEntitySpec()
	if err != nil {
		return nil, fmt.Errorf("failed to create API config entity: %w", err)
	}

	configInit, err := config.ConfigEntityInitializer()
	if err != nil {
		return nil, fmt.Errorf("failed to create API config entity: %w", err)
	}

	return &ConfigDef{
		APIName:      "",
		ConfigEntity: configType,
		ConfigInit:   configInit,
	}, nil
}
