package outputs

import (
	"github.com/softwaresale/client-gen/v2/internal/types"
	"io"
)

// CompilerOutputWriter is an output that we can write our targets to. They are writable and can provide a location
type CompilerOutputWriter interface {
	io.Writer
	CompilerOutputLocation
}

type CompilerOutputLocation interface {
	Name() string              // get the name of the output
	Location() (string, error) // Where this output is located
}

// CompilerOutputsManager helps abstract away the details of creating files to write the compiler outputs to. Instead,
// this interface produces different outputs for different input types
type CompilerOutputsManager interface {
	PrepareOutputDirectory(path string) error                                                  // Prepares the directory that these outputs will go
	CreateServiceOutput(serviceDef types.ServiceDefinition) (CompilerOutputWriter, error)      // create a writer to write this service definition to
	ComputeServiceLocation(serviceDef types.ServiceDefinition) (CompilerOutputLocation, error) // figure out where this service will be located without actually creating the output
	CreateModelOutput(model types.EntitySpec) (CompilerOutputWriter, error)                    // create a writer to write this model entity to
	ComputeModelLocation(model types.EntitySpec) (CompilerOutputLocation, error)               // figure out where this entity will be located without actually creating the output
}
