package servicegen

import (
	"github.com/softwaresale/client-gen/v2/internal/codegen/imports"
	"github.com/softwaresale/client-gen/v2/internal/types"
	"io"
)

// ServiceGenerator specifies a type that formats a service to an output
//
//go:generate mockery --name ServiceGenerator --structname MockServiceGenerator --outpkg servicegen_mocks
type ServiceGenerator interface {
	GenerateService(writer io.Writer, service types.ServiceDefinition, resolver imports.ImportManager) error
	GenerateEntity(writer io.Writer, entity types.EntitySpec, resolver imports.ImportManager) error
}
