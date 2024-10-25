package jscodegen

import (
	"github.com/softwaresale/client-gen/v2/internal/codegen"
	"github.com/softwaresale/client-gen/v2/internal/codegen/outputs"
)

// NewNGCompiler creates a new angular API compiler that produces Angular code
func NewNGCompiler(outputDirectory string) codegen.APICompiler {
	ngServiceGen := NewNGServiceGenerator()
	ngImportMgr := NewTSImportManager()

	return codegen.APICompiler{
		Generator:     ngServiceGen,
		ImportManager: &ngImportMgr,
		OutputsManager: &outputs.DirectoryCompilerOutputsManager{
			BasePath: outputDirectory,
		},
		OutputPath: outputDirectory,
	}
}
