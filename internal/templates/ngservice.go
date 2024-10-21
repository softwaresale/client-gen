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

type HttpRequestDef struct {
	HttpClientVar    string
	HttpMethod       string
	ResponseType     string
	URITemplate      URITemplate
	RequestBodyValue string
}

func hasRequestBody(def string) bool {
	return len(def) > 0
}

type RequestMethodDef struct {
	RequestName      string
	InputVarName     string
	RequestInputType string
	ResponseType     string
	HttpRequest      HttpRequestDef
}

func (def RequestMethodDef) HasInput() bool {
	return len(def.RequestInputType) > 0
}

type RequestInputDef struct {
	Name       string
	Properties map[string]string
}

func (def RequestInputDef) IsValid() bool {
	return len(def.Properties) > 0
}

type ServiceDef struct {
	ServiceName   string
	HttpClientVar string
	InputTypes    []RequestInputDef
	Methods       []RequestMethodDef
}

func mapHttpEndpoint(method string) string {
	return strings.ToLower(method)
}

func Translate(service codegen.ServiceDefinition) (ServiceDef, error) {

	typeMapper := jscodegen.JSTypeMapper{}
	httpClientVar := "http"

	var methods []RequestMethodDef
	var inputs []RequestInputDef
	for _, endpoint := range service.Endpoints {

		inputVarName := "input"
		bodyPropertyName := "body"

		// create the input type
		inputTypeName := strcase.ToCamel(fmt.Sprintf("%sInput", endpoint.Name))
		properties := make(map[string]string)
		for prop, tp := range endpoint.PathVariables {
			tpStr, err := typeMapper.Convert(tp.Type)
			if err != nil {
				return ServiceDef{}, err
			}

			properties[prop] = tpStr
		}

		if !endpoint.RequestBody.Type.IsVoid() {
			tpStr, err := typeMapper.Convert(endpoint.RequestBody.Type)
			if err != nil {
				return ServiceDef{}, err
			}
			properties[bodyPropertyName] = tpStr
		}

		requestInputDef := RequestInputDef{
			Name:       inputTypeName,
			Properties: properties,
		}

		if !requestInputDef.IsValid() {
			inputTypeName = ""
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
		if requestInputDef.IsValid() {
			inputs = append(inputs, requestInputDef)
		}
	}

	return ServiceDef{
		ServiceName:   service.Name,
		HttpClientVar: httpClientVar,
		Methods:       methods,
		InputTypes:    inputs,
	}, nil
}

type NGServiceGenerator struct {
	ngServiceTemplate *template.Template
}

// NewNGServiceGenerator creates a new NGService generator, which can be used to generate services
func NewNGServiceGenerator() *NGServiceGenerator {
	tmpl := template.Must(template.New("NGService").Funcs(template.FuncMap{
		"HasRequestBody": hasRequestBody,
		"ParseTemplate":  FormatTemplate,
	}).Parse(templateText))

	return &NGServiceGenerator{
		ngServiceTemplate: tmpl,
	}
}

func (generator *NGServiceGenerator) Generate(writer io.Writer, def ServiceDef) error {
	return generator.ngServiceTemplate.Execute(writer, def)
}
