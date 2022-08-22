package scan

import (
	"fmt"
	"github.com/hashicorp/terraform-config-inspect/tfconfig"
	"io/fs"
	"path/filepath"
)

type Result struct {
	ModuleCall tfconfig.ModuleCall
	Path       string
}

func Scan(paths []string, recursive bool) ([]Result, error) {
	// recurse through the provided paths in a non-overlapping way
	if recursive {
		recursivePaths, err := loadSubPaths(paths)
		if err != nil {
			return nil, err
		}

		paths = recursivePaths
	}

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

func loadSubPaths(paths []string) ([]string, error) {
	combinedPaths := make([]string, 0)

	for _, path := range paths {
		err := filepath.WalkDir(path, func(path string, d fs.DirEntry, err error) error {
			// don't worry about anything that isn't a directory
			if !d.IsDir() {
				return nil
			}

			combinedPaths = append(combinedPaths, path)
			return nil
		})

		if err != nil {
			return paths, err
		}
	}

	return combinedPaths, nil
}
