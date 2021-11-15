package scan

import (
	"fmt"
	"path/filepath"

	"github.com/hashicorp/terraform-config-inspect/tfconfig"
	"github.com/keilerkonzept/terraform-module-versions/pkg/modulecall"
)

// if the local module search goes more than this number, it will stop searching.
const MaxRecursionDepth int = 10

type Result struct {
	ModuleCall tfconfig.ModuleCall
	Path       string
}

func Scan(paths []string, searchLocalModules bool) ([]Result, error) {
	var out []Result

	for _, path := range paths {
		module, err := tfconfig.LoadModule(path)
		if err != nil {
			return nil, fmt.Errorf("read terraform module %q: %w", path, err)
		}
		for _, call := range module.ModuleCalls {
			if call == nil {
				continue
			}
			out = append(out, Result{
				Path:       call.Pos.Filename,
				ModuleCall: *call,
			})
		}
	}

	if searchLocalModules {
		var err error

		out, err = scanLocalModules(out, MaxRecursionDepth)
		if err != nil {
			return nil, err
		}
	}

	return out, nil
}

//	This function scans for nested terraform modules by parsing local modules recursively.
func scanLocalModules(modules []Result, depth int) ([]Result, error) {
	var out []Result

	if depth < 0 {
		fmt.Printf("[WARN] Max depth of %d was reached. Module dependencies may have a loop that causes infinite recursion.\n", MaxRecursionDepth)
		return out, nil
	}

	for _, m := range modules {
		dirPath := filepath.Dir(m.Path)
		parsed, err := modulecall.Parse(m.ModuleCall)
		if err != nil {
			return nil, err
		}

		// if the module source is local, load the nested module and scan for more local modules
		if parsed.Source.Local != nil {
			module, err := tfconfig.LoadModule(dirPath + "/" + m.ModuleCall.Source)
			if err != nil {
				return nil, fmt.Errorf("read terraform module %q: %w", m.ModuleCall.Source, err)
			}

			for _, call := range module.ModuleCalls {
				if call == nil {
					continue
				}
				mod := Result{
					Path:       call.Pos.Filename,
					ModuleCall: *call,
				}

				singleModule := make([]Result, 0)
				singleModule = append(singleModule, mod)

				nestedModules, err := scanLocalModules(singleModule, depth-1)
				if err != nil {
					return nil, err
				}

				out = append(out, nestedModules...)
			}
		}
	}

	result := append(modules, out...)

	return result, nil
}
