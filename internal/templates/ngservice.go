package templates

import (
	_ "embed"
	"fmt"
	"github.com/iancoleman/strcase"
	"github.com/softwaresale/client-gen/v2/internal/codegen"
	"github.com/softwaresale/client-gen/v2/internal/jscodegen"
	"io"
	"strings"
	"text/template"
)

//go:embed ng-service.tmpl
var templateText string

//go:embed ng-entity.tmpl
var entityTemplateText string

type HttpRequestDef struct {
	HttpClientVar    string      // the name of the variable that defines the HTTP client in use
	HttpMethod       string      // The HTTP method used by this request
	ResponseType     string      // the type string of our response
	URITemplate      URITemplate // our URI template. This gets mapped into a uri string
	RequestBodyValue string      // The value to read the body type
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

type ServiceDef struct {
	ServiceName   string
	HttpClientVar string
	InputTypes    []codegen.EntitySpec
	Methods       []RequestMethodDef
}

func mapHttpEndpoint(method string) string {
	return strings.ToLower(method)
}

func translate(service codegen.ServiceDefinition) (ServiceDef, error) {

	typeMapper := jscodegen.JSTypeMapper{}
	httpClientVar := "http"

	var methods []RequestMethodDef
	var inputs []codegen.EntitySpec
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
				URITemplate: URITemplate{
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

	return ServiceDef{
		ServiceName:   service.Name,
		HttpClientVar: httpClientVar,
		Methods:       methods,
		InputTypes:    inputs,
	}, nil
}

func createInputType(endpoint codegen.APIEndpoint, bodyPropertyName string) (*codegen.EntitySpec, error) {
	inputTypeName := strcase.ToCamel(fmt.Sprintf("%sInput", endpoint.Name))
	properties := make(map[string]codegen.PropertySpec)
	for prop, tp := range endpoint.PathVariables {
		properties[prop] = codegen.PropertySpec{
			Type:     tp.Type,
			Required: tp.Required,
		}
	}

	if !endpoint.RequestBody.Type.IsVoid() {
		properties[bodyPropertyName] = codegen.PropertySpec{
			Type:     endpoint.RequestBody.Type,
			Required: endpoint.RequestBody.Required,
		}
	}

	return &codegen.EntitySpec{
		Name:       inputTypeName,
		Properties: properties,
	}, nil
}

type NGServiceGenerator struct {
	ngServiceTemplate *template.Template
}

// NewNGServiceGenerator creates a new NGService generator, which can be used to generate services
func NewNGServiceGenerator() *NGServiceGenerator {

	mapper := jscodegen.JSTypeMapper{}

	tmpl := template.Must(template.New("NGService").Funcs(template.FuncMap{
		"HasRequestBody": hasRequestBody,
		"ParseTemplate":  FormatTemplate,
		"ConvertType":    mapper.Convert,
	}).Parse(templateText))

	tmpl = template.Must(tmpl.Parse(entityTemplateText))

	return &NGServiceGenerator{
		ngServiceTemplate: tmpl,
	}
}

func (generator *NGServiceGenerator) Generate(writer io.Writer, def codegen.ServiceDefinition) error {

	translatedDef, err := translate(def)
	if err != nil {
		return fmt.Errorf("failed to translate service definition: %w", err)
	}

	return generator.ngServiceTemplate.Execute(writer, translatedDef)
}
