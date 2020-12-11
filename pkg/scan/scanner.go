package scan

import (
	"fmt"

	"github.com/hashicorp/terraform-config-inspect/tfconfig"
)

type Result struct {
	ModuleCall tfconfig.ModuleCall
	Path       string
}

func Scan(paths []string) ([]Result, error) {
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
	return out, nil
}
