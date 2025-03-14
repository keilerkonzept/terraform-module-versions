package scan

import (
	"fmt"
	"path/filepath"

	"github.com/hashicorp/terraform-config-inspect/tfconfig"
	"github.com/keilerkonzept/terraform-module-versions/pkg/modulecall"
)

type Result struct {
	ModuleCall tfconfig.ModuleCall
	Path       string
}

var foundNestedModules map[string]Result = make(map[string]Result)

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
		// transform slice of modules to map of modules to identify duplicates
		modulesByMap := make(map[string]Result)

		for _, module := range out {
			p, err := filepath.Abs(module.ModuleCall.Source)
			if err != nil {
				return nil, err
			}

			modulesByMap[p] = module
		}

		err = scanLocalModules(modulesByMap)
		if err != nil {
			return nil, err
		}

		// transform map of modules to list of modules
		out = make([]Result, 0)
		for _, module := range foundNestedModules {
			out = append(out, module)
		}
	}

	return out, nil
}

/*
	This function scans for nested terraform modules by parsing local modules recursively.
	The key for the map is the fully qualified path to the source of the module
	The key is always expanded so that relative paths are not seen a duplicates.

	Consider the following scenario:
	The naming module has the same name and path, however, they exist in 2 different directories.
	The source path for module.a.module.naming is : ./mod1/naming
	The source path for module.b.module.naming is : ./mod2/naming

	main.tf
	```hcl
	module "a" {
		source = "./mod1"
	}

	module "b" {
		source = "./mod2"
	}
	```

	./mod1:
	```hcl
	resource "..." {
	}

	...

	module "naming" {
		source = "./naming"
	}
	```

	./mod2:
	```hcl
		resource "..." {
	}

	...

	module "naming" {
		source = "./naming"
	}
	```
*/
func scanLocalModules(existingModules map[string]Result) error {
	for baseModuleSource, m := range existingModules {
		dirPath := filepath.Dir(m.Path)
		parsed, err := modulecall.Parse(m.ModuleCall)
		if err != nil {
			return err
		}

		// if the module source is local, load it and scan for more local modules
		if parsed.Source.Local != nil {
			module, err := tfconfig.LoadModule(dirPath + "/" + m.ModuleCall.Source)
			if err != nil {
				return fmt.Errorf("read terraform module %q: %w", m.ModuleCall.Source, err)
			}

			for _, call := range module.ModuleCalls {
				if call == nil {
					continue
				}
				module := Result{
					Path:       call.Pos.Filename,
					ModuleCall: *call,
				}

				// make a fully qualified path to ensure modules are distinct
				p, err := filepath.Abs(filepath.Join(baseModuleSource, call.Source))
				if err != nil {
					return err
				}

				// if the module has already been seen, skip processing it.
				if _, ok := foundNestedModules[p]; ok {
					continue
				} else {
					foundNestedModules[p] = module
				}

				singleModule := make(map[string]Result)
				singleModule[p] = module

				err = scanLocalModules(singleModule)

				if err != nil {
					return err
				}
			}
		}
	}

	for baseModuleSource, nestedModule := range existingModules {
		// make a fully qualified path to ensure modules are distinct
		p, err := filepath.Abs(filepath.Join(baseModuleSource, nestedModule.ModuleCall.Source))
		if err != nil {
			return err
		}

		if _, ok := foundNestedModules[p]; !ok {
			foundNestedModules[p] = nestedModule
		}
	}

	return nil
}
