package jscodegen

import (
	"fmt"
	"github.com/softwaresale/client-gen/v2/internal/pkg/codegen"
)

type HttpGetModel struct {
	ClientVariable  string
	ResponseTypeStr string
	RequestStr      string
}

func CreateHttpGetModel(endpoint codegen.APIEndpoint) (*HttpGetModel, error) {
	endpointStr, err := CreateEndpointStr(endpoint.Endpoint, endpoint.PathVariables)
	if err != nil {
		return nil, fmt.Errorf("failed to create endpoint str: %w", err)
	}

	responseTypeStr := endpoint.RequestBody.Type.TypeID

	return &HttpGetModel{
		ClientVariable:  "this.http",
		ResponseTypeStr: responseTypeStr,
		RequestStr:      endpointStr,
	}, nil
}
