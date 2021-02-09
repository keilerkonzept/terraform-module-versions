package scan

import (
	"fmt"

	tgconfig "github.com/gruntwork-io/terragrunt/config"
	tgoptions "github.com/gruntwork-io/terragrunt/options"
	"github.com/hashicorp/terraform-config-inspect/tfconfig"
)

func ScanTerragrunt(paths []string) ([]Result, error) {
	var out []Result
	for _, path := range paths {
		opts, err := tgoptions.NewTerragruntOptions(path)
		if err != nil {
			return nil, fmt.Errorf("build terragrunt options for %q: %w", path, err)
		}
		paths, err := tgconfig.FindConfigFilesInPath(path, opts)
		if err != nil {
			return nil, fmt.Errorf("find terragrunt config files in %q: %w", path, err)
		}
		for _, path := range paths {
			cfg, err := tgconfig.PartialParseConfigFile(
				path,
				opts,
				nil,
				[]tgconfig.PartialDecodeSectionType{tgconfig.TerraformSource},
			)
			if err != nil {
				return nil, fmt.Errorf("read terragrunt config %q: %w", path, err)
			}
			if cfg.Terraform == nil {
				continue
			}
			if cfg.Terraform.Source == nil {
				continue
			}
			out = append(out, Result{
				Path: path,
				ModuleCall: tfconfig.ModuleCall{
					Name:    path,
					Source:  *cfg.Terraform.Source,
					Version: "",
					Pos: tfconfig.SourcePos{
						Filename: path,
					},
				},
			})
		}
	}
	return out, nil
}
