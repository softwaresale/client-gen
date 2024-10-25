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

func mapHttpEndpoint(method string) string {
	return strings.ToLower(method)
}

type NGServiceGenerator struct {
	ngServiceTemplate *template.Template
	ngEntityTemplate  *template.Template
}

// NewNGServiceGenerator creates a new NGService generator, which can be used to generate services
func NewNGServiceGenerator() *NGServiceGenerator {

	mapper := JSTypeMapper{}

	funcMap := template.FuncMap{
		"HasRequestBody": hasRequestBody,
		"ParseTemplate":  codegen.FormatTemplate,
		"ConvertType":    mapper.Convert,
	}

	serviceTmpl := template.Must(template.New("NGService").Funcs(funcMap).Parse(templateText))
	serviceTmpl = template.Must(serviceTmpl.Parse(importsTemplateText))
	serviceTmpl = template.Must(serviceTmpl.Parse(entityTemplateText))

	entityTmpl := template.Must(template.New("NGEntity").Funcs(funcMap).Parse(entityTemplateText))
	entityTmpl = template.Must(entityTmpl.Parse(importsTemplateText))
	entityTmpl = template.Must(entityTmpl.Parse(standaloneEntityTemplateText))

	return &NGServiceGenerator{
		ngServiceTemplate: serviceTmpl,
		ngEntityTemplate:  entityTmpl,
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

	imports := importResolver.GetEntityImports(spec)

	return EntityDef{
		Entity:  spec,
		Imports: imports,
	}
}
